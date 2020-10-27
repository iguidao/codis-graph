package es

import (
	"bytes"
	"codisgraph/src/cfg"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v6/esapi"
)

// 按时间聚合ES的查询条件
func EsScrollTime(starttime, endtime string) (map[string]interface{}, string, string) {
	EsdocType := cfg.Get_Elasticsearch("doctype")
	Esindex := cfg.Get_Elasticsearch("index")
	// today := time.Now().Format("2006.01.02")
	IndexName := Esindex + "-" + "*"

	starttime64, err := strconv.ParseInt(starttime, 10, 64)
	if err != nil {
		log.Println("解析时间失败：", err)
	}
	endtime64, err := strconv.ParseInt(endtime, 10, 64)
	if err != nil {
		log.Println("解析时间失败：", err)
	}

	stm := time.Unix((starttime64 - 28800), 0)
	etm := time.Unix((endtime64 - 28800), 0)
	esstart := stm.Format("2006-01-02T15:04:05.000Z")
	esend := etm.Format("2006-01-02T15:04:05.000Z")
	query := map[string]interface{}{
		"size": 1000,
		"sort": map[string]interface{}{
			"@timestamp": map[string]interface{}{
				"order": "asc",
			},
		},
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"@timestamp": map[string]interface{}{
					"gte": esstart,
					"lte": esend,
				},
			},
		},
	}
	return query, EsdocType, IndexName
}

// 聚合ES的Scroll查询条件
func EsScrollQuery(scrollid string) map[string]interface{} {
	query := map[string]interface{}{
		"scroll":    "1m",
		"scroll_id": scrollid,
	}
	return query
}

// 获取sroll的ID
func EsSrollID(query map[string]interface{}, DocType, IndexName string) (string, EsData) {
	jsonBody, _ := json.Marshal(query)
	seconds := 3
	srolles := time.Duration(seconds) * time.Minute

	var a int = 1000
	var essize *int
	essize = &a

	req := esapi.SearchRequest{
		Index:        []string{IndexName},
		DocumentType: []string{DocType},
		Scroll:       srolles,
		Size:         essize,
		Body:         bytes.NewReader(jsonBody),
	}

	res, err := req.Do(context.Background(), ES)
	checkError(err)
	defer res.Body.Close()
	// log.Println(res)
	var esdata EsData

	if res.StatusCode != 200 {
		return "NotFound", esdata
	}
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &esdata)
	if err != nil {
		log.Println("err was: ", err)
	}

	return "Found", esdata

}

// 获取sroll数据
func EsSrollSearch(query map[string]interface{}, DocType, IndexName string) (string, EsData) {
	seconds := 1
	srolles := time.Duration(seconds) * time.Minute

	jsonBody, _ := json.Marshal(query)
	req := esapi.ScrollRequest{
		Scroll: srolles,
		Body:   bytes.NewReader(jsonBody),
	}

	res, err := req.Do(context.Background(), ES)
	checkError(err)
	defer res.Body.Close()
	// log.Println(res)
	var esdata EsData

	if res.StatusCode != 200 {
		return "NotFound", esdata
	}
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &esdata)
	if err != nil {
		log.Println("err was: ", err)
	}

	return "Found", esdata

}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error Elasticsearch Opration Response: %s", err)
	}
}
