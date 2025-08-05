package handle

import (
	"encoding/json"
	"fmt"
	"github.com/PixDevopsSre/helloWorld/config"
	logger "github.com/PixDevopsSre/helloWorld/pkg"
	"github.com/PixDevopsSre/helloWorld/util"
	cmdb_sdk "github.com/veops/cmdb-sdk-golang"
	"log"
)

const (
	cmdbURL       = "http://localhost:8000/api/v0.1"
	cmdbApiKey    = "af6164ce93784f359695563216a7127d"
	cmdbApiSecret = "2Qj8z#s~DPxLV0!Cl9uK*$gkOi^qF&BM"
)

var (
	cmdbHelper = cmdb_sdk.NewHelper(cmdbURL, cmdbApiKey, cmdbApiSecret)

	url = "http://localhost:80/cloud"
)

func MysqlHandle() error {
	ids, err := getProjectIDs()
	if err != nil {
		logger.Errorf("getProjectIDs error: %v", err)
		return err
	}

	for _, id := range ids {
		var dnsCloudData = config.CloudResource{
			CloudName: "gcp",
			CloudType: "mysql",
			CloudData: config.CloudData{
				MySqlData: config.MySqlData{
					ProjectId: id,
				},
			},
		}

		toMap, err := util.StructToMap(dnsCloudData)
		if err != nil {
			logger.Errorf("StructToMap error: %v", err)
			return err
		}

		callResp, code, err := util.ToCallHttp(url, "GET", toMap)
		if err != nil {
			logger.Errorf("ToCallHttp error: %v", err)
			return err
		}
		if code != 200 {
			continue
		}
		var apiResp config.APIResponse
		err = json.Unmarshal(callResp, &apiResp)
		if err != nil {
			logger.Errorf("json.Unmarshal error: %v", err)
			return err
		}
		err = util.ExportToExcelAppend("database_report.xlsx", "Databases", apiResp.Data)
		if err != nil {
			logger.Errorf("Failed to export DB data: %v", err)
			continue
		}
		fmt.Println("Database report created.")
	}

	return nil
}

func S3Handle() error {
	ids, err := getProjectIDs()
	if err != nil {
		logger.Errorf("getProjectIDs error: %v", err)
		return err
	}

	for _, id := range ids {
		var dnsCloudData = config.CloudResource{
			CloudName: "gcp",
			CloudType: "s3",
			CloudData: config.CloudData{
				S3Data: config.S3Data{
					ProjectId: id,
				},
			},
		}

		toMap, err := util.StructToMap(dnsCloudData)
		if err != nil {
			logger.Errorf("StructToMap error: %v", err)
			return err
		}

		callResp, code, err := util.ToCallHttp(url, "GET", toMap)
		if err != nil {
			logger.Errorf("ToCallHttp error: %v", err)
			return err
		}
		if code != 200 {
			continue
		}
		var apiResp config.APIResponse
		err = json.Unmarshal(callResp, &apiResp)
		if err != nil {
			logger.Errorf("json.Unmarshal error: %v", err)
			return err
		}

		err = util.ExportToExcelAppend("s3_report.xlsx", "gcs", apiResp.Data)
		if err != nil {
			logger.Errorf("Failed to export s3 data: %v", err)
			continue
		}
		fmt.Println("Database report created.")
	}

	return nil
}

func VmHandle() error {
	ids, err := getProjectIDs()
	if err != nil {
		logger.Errorf("getProjectIDs error: %v", err)
		return err
	}

	for _, id := range ids {
		var dnsCloudData = config.CloudResource{
			CloudName: "gcp",
			CloudType: "vm",
			CloudData: config.CloudData{
				VmData: config.VmData{
					ProjectId: id,
				},
			},
		}

		toMap, err := util.StructToMap(dnsCloudData)
		if err != nil {
			logger.Errorf("StructToMap error: %v", err)
			return err
		}

		callResp, code, err := util.ToCallHttp(url, "GET", toMap)
		if err != nil {
			logger.Errorf("ToCallHttp error: %v", err)
			return err
		}
		if code != 200 {
			continue
		}
		var apiResp config.APIResponse
		err = json.Unmarshal(callResp, &apiResp)
		if err != nil {
			logger.Errorf("json.Unmarshal error: %v", err)
			return err
		}

		err = util.ExportToExcelAppend("vm_report.xlsx", "VM", apiResp.Data)
		if err != nil {
			logger.Errorf("Failed to export Redis data: %v", err)
			continue
		}
		fmt.Println("Database report created.")
	}

	return nil
}

func RedisHandle() error {
	ids, err := getProjectIDs()
	if err != nil {
		logger.Errorf("getProjectIDs error: %v", err)
		return err
	}

	for _, id := range ids {
		var dnsCloudData = config.CloudResource{
			CloudName: "gcp",
			CloudType: "redis",
			CloudData: config.CloudData{
				RedisData: config.RedisData{
					ProjectId: id,
				},
			},
		}

		toMap, err := util.StructToMap(dnsCloudData)
		if err != nil {
			logger.Errorf("StructToMap error: %v", err)
			return err
		}

		callResp, code, err := util.ToCallHttp(url, "GET", toMap)
		if err != nil {
			logger.Errorf("ToCallHttp error: %v", err)
			return err
		}
		if code != 200 {
			continue
		}
		var apiResp config.APIResponse
		err = json.Unmarshal(callResp, &apiResp)
		if err != nil {
			logger.Errorf("json.Unmarshal error: %v", err)
			return err
		}

		err = util.ExportToExcelAppend("redis_report.xlsx", "Redis", apiResp.Data)
		if err != nil {
			logger.Errorf("Failed to export Redis data: %v", err)
			continue
		}
		fmt.Println("Database report created.")
	}

	return nil
}

func getProjectIDs() ([]string, error) {
	var projectData = &config.CloudResource{
		CloudName: "gcp",
		CloudType: "project",
	}

	toMap, err := util.StructToMap(projectData)
	if err != nil {
		logger.Errorf("StructToMap error: %v", err)
		return nil, err
	}

	callResp, _, err := util.ToCallHttp(url, "GET", toMap)
	if err != nil {
		logger.Errorf("ToCallHttp error: %v", err)
		return nil, err
	}

	var (
		resp config.Response
	)
	err = json.Unmarshal(callResp, &resp)
	if err != nil {
		logger.Errorf("json.Unmarshal error: %v", err)
		return nil, err
	}

	projectIDs := make([]string, 0, resp.Data.Count)

	for _, project := range resp.Data.Projects {
		if project.ProjectId != "" {
			projectIDs = append(projectIDs, project.ProjectId)
		} else {
			logger.Warnf("Project ID is empty for project: %s", project.Name)
		}
	}

	return projectIDs, nil
}

func upsertCI(modelName, uniqueField, uniqueValue string, attrs map[string]any) (int, error) {
	query := fmt.Sprintf(`_type:%s,%s:%s`, modelName, uniqueField, uniqueValue)

	log.Printf("正在使用正确的服务端过滤查询: %s", query)
	getCIRes, err := cmdbHelper.GetCI(query, "", "", "", 1, 0, cmdb_sdk.RetKeyDefault)
	if err != nil {
		log.Printf("Upsert - 查询CI失败 [%s]: %v", query, err)
		return 0, err
	}
	if len(getCIRes.Result) > 0 {
		// 2. 获取第一个匹配到的资产
		// getCIRes.Result 是一个列表，我们取第一个元素 [0]
		firstCI := getCIRes.Result[0]
		log.Printf("【完整CI数据】: %#v", firstCI)
		// ⭐⭐⭐【新的重要调试代码】请在这里加上这一行 ⭐⭐⭐
		idValueFromAPI := firstCI["_id"]
		log.Printf("【调试信息】收到的ID值为: %#v, 其Go语言类型为: %T", idValueFromAPI, idValueFromAPI)
		// ⭐⭐⭐【重要调试代码】请在这里加上这两行
		// 3. 从 map 中提取您需要的具体字段，比如资产的ID
		// 注意：从 map[string]any 取值需要进行类型断言
		idFloat, ok := firstCI["_id"].(float64) // JSON数字默认解析为float64
		if ok {
			// 将 float64 转换为您需要的 int64 类型
			ciID := int64(idFloat)
			fmt.Printf("找到了！资产ID是: %d\n", ciID)

			// 在这里，您就可以用这个 ciID 去调用 helper.UpdateCI(...) 了
		} else {
			fmt.Println("找到了资产，但无法解析ID字段。")
		}

		ciID := int(idFloat)
		_, err := cmdbHelper.UpdateCI(ciID, modelName, cmdb_sdk.NoAttrPolicyDefault, attrs)
		return ciID, err
	} else {
		addCIRes, err := cmdbHelper.AddCI(modelName, cmdb_sdk.NoAttrPolicyDefault, cmdb_sdk.ExistPolicyDefault, attrs)
		if err != nil {
			return 0, err
		}
		return addCIRes.CIID, nil
	}
}
