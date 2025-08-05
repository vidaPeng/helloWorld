package main

//type APIResponse struct {
//	Data Data `json:"data"`
//}
//
//// Data matches the object inside "data".
//type Data struct {
//	Node Node `json:"node"`
//}
//
//// Node matches the directory node.
//type Node struct {
//	Nodes []LeafNode `json:"nodes"` // This holds the list of certificate nodes.
//}
//
//// LeafNode matches the final node containing certificate info.
//type LeafNode struct {
//	Value SSLInfo `json:"value"`
//}
//
//// SSLInfo represents the structure of the JSON string inside the 'value' field.
//type SSLInfo struct {
//	Cert string   `json:"cert"`
//	ID   string   `json:"id"`
//	Key  string   `json:"key"`
//	SNIs []string `json:"snis"`
//}
//
//type DD struct {
//	Name          string `json:"name"`
//	Target        string `json:"target"`
//	MonitorType   string `json:"monitor_type"`
//	ReportTypes   string `json:"report_types"`
//	Cycle         string `json:"cycle"`
//	RepeatSendGap string `json:"repeat_send_gap"`
//	Active        string `json:"active"`
//	AdvanceDay    string `json:"advance_day"`
//}
//
//func main() {
//	file, err := os.ReadFile("/Users/vida/Work/helloWorld/json/data.json")
//	if err != nil {
//		panic(err)
//	}
//	var data APIResponse
//	err = json.Unmarshal(file, &data)
//	if err != nil {
//		panic(err)
//	}
//
//	var l []DD
//
//	for _, node := range data.Data.Node.Nodes {
//		sslInfo := node.Value
//		if len(sslInfo.SNIs) == 0 {
//			continue // Skip if SNIs is empty
//		}
//
//		for _, sni := range sslInfo.SNIs {
//			var t = DD{
//				MonitorType:   "https",
//				ReportTypes:   "feishu",
//				Cycle:         "480",
//				RepeatSendGap: "2",
//				Active:        "1",
//				AdvanceDay:    "30",
//			}
//
//			t.Name = sni
//			t.Target = sni
//
//			l = append(l, t)
//		}
//	}
//
//	marshal, err := json.Marshal(l)
//	if err != nil {
//		panic(err)
//	}
//
//	err = os.WriteFile("/Users/vida/Work/helloWorld/json/test.json", marshal, 0644)
//	if err != nil {
//		panic(err)
//	}
//}
