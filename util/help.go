package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
)

func ToCallHttp(rowURLl, method string, queryParams map[string]interface{}) ([]byte, int, error) {
	reqUrl, err := url.Parse(rowURLl)
	if err != nil {
		return nil, 0, fmt.Errorf("parse url failed: %w", err)
	}

	var reqBody io.Reader

	method = strings.ToUpper(method)
	switch method {
	case http.MethodGet:
		// 拼接 URL 查询参数
		q := reqUrl.Query()
		for k, v := range queryParams {
			var val string
			switch vv := v.(type) {
			case string:
				val = vv
			default:
				b, err := json.Marshal(vv)
				if err != nil {
					return nil, 0, fmt.Errorf("marshal query param failed: %w", err)
				}
				val = string(b)
			}
			q.Set(k, val)
		}
		reqUrl.RawQuery = q.Encode()

	default:
		// POST ，需要把整个 queryParams 放 body
		b, err := json.Marshal(queryParams)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal body failed: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, reqUrl.String(), reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("new request failed: %w", err)
	}
	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read body failed: %w", err)
	}

	return body, resp.StatusCode, nil
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &m)
	return m, err
}

// filename: 要保存的文件名
// sheetName: 工作表名
// data: 要导出的数据，必须是 []MyStruct 或者 map[string]MyStruct 的形式
func ExportToExcelAppend(filename, sheetName string, data interface{}) error {
	// --- 1. 检查文件是否存在，然后创建或打开 ---
	var f *excelize.File
	var err error

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// 文件不存在，创建它
		f = excelize.NewFile()
	} else {
		// 文件存在，打开它
		f, err = excelize.OpenFile(filename)
		if err != nil {
			return fmt.Errorf("failed to open existing Excel file: %w", err)
		}
	}

	// --- 2. 确定工作表和起始行 ---
	// 确保工作表存在，如果不存在则创建
	index, err := f.GetSheetIndex(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get sheet index: %w", err)
	}
	if index == -1 {
		index, err := f.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("failed to create new sheet: %w", err)
		}
		f.SetActiveSheet(index)
	}

	// 获取工作表现有行，以确定下一行写入位置
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows from sheet: %w", err)
	}
	rowIndex := len(rows) + 1 // 下一个可用的行号

	// --- 3. 分析数据并动态写入 ---
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Map {
		return errors.New("data must be a slice or a map")
	}
	if val.Len() == 0 {
		return nil // 没有数据就直接返回，什么也不做
	}

	// 动态生成表头 (仅当工作表为空时)
	var headers []string
	elemType := val.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return errors.New("data element must be a struct")
	}

	// 为 map 类型额外添加 "Name" 或 "Key" 列
	if val.Kind() == reflect.Map {
		headers = append(headers, "Name")
	}
	for i := 0; i < elemType.NumField(); i++ {
		headers = append(headers, elemType.Field(i).Name)
	}

	// 如果是新表 (rowIndex == 1)，则写入表头
	if rowIndex == 1 {
		if err := f.SetSheetRow(sheetName, "A1", &headers); err != nil {
			return fmt.Errorf("failed to write headers: %w", err)
		}
		// 表头写入后，数据从下一行开始
		rowIndex++
	}

	// --- 4. 写入数据行 (追加模式) ---
	switch val.Kind() {
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			var rowData []interface{}
			elem := val.Index(i)
			for j := 0; j < elem.NumField(); j++ {
				rowData = append(rowData, elem.Field(j).Interface())
			}
			cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
			f.SetSheetRow(sheetName, cell, &rowData)
			rowIndex++ // 递增行号，为下一条数据做准备
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			var rowData []interface{}
			elem := val.MapIndex(key)
			rowData = append(rowData, key.Interface()) // 第一列是 map 的 key
			for j := 0; j < elem.NumField(); j++ {
				fieldValue := elem.Field(j).Interface()
				if reflect.TypeOf(fieldValue).Kind() == reflect.Slice {
					rowData = append(rowData, fmt.Sprintf("%v", fieldValue))
				} else {
					rowData = append(rowData, fieldValue)
				}
			}
			cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
			f.SetSheetRow(sheetName, cell, &rowData)
			rowIndex++ // 递增行号
		}
	}

	// 使用 Save() 而不是 SaveAs() 来保存对现有文件的更改
	return f.SaveAs(filename)
}
