package codisapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func GetCluster(codisurl string) []string {

	ListUrl := codisurl + "/list"
	resp, err := http.Get(ListUrl)
	if err != nil {
		log.Println("获取Codis集群错误：", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var ClusterList ClusterResult
	err = json.Unmarshal(body, &ClusterList)
	if err != nil {
		log.Println("err was : ", err)
	}

	return ClusterList

}

func GetHost(codisurl, Cname string) CodisInfo {
	TopomUrl := codisurl + "/topom/stats?forward=" + Cname
	resp, err := http.Get(TopomUrl)
	if err != nil {
		log.Println("获取Codis集群错误：", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var codisinfo CodisInfo
	err = json.Unmarshal(body, &codisinfo)
	if err != nil {
		log.Println("err was : ", err)
	}
	return codisinfo

}
