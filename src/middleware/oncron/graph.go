package oncron

import (
	"log"
	"strconv"
	"time"

	"codisgraph/src/cfg"
	"codisgraph/src/middleware/es"
	"codisgraph/src/middleware/mysql"
)

func CronGraph() {
	log.Println("start cron codisgraph")
	CodisCreate := make(map[string]string)
	CodisClose := make(map[string]string)
	endtime := time.Now().Unix()
	difftime := cfg.Get_Local("estimediff")
	esdifftime, err := strconv.ParseInt(difftime, 10, 64)
	if err != nil {
		log.Println("解析时间失败：", err)
	}
	starttime := endtime - esdifftime
	query, EsdocType, IndexName := es.EsScrollTime(starttime, endtime)
	SearchStatus, SearchResult := es.EsSrollID(query, EsdocType, IndexName)
	if SearchStatus == "NotFound" {
		log.Println("NotFound", SearchStatus, SearchResult)
	}

	CodisCreate, CodisClose = DateHandle(CodisCreate, CodisClose, SearchResult)

	scrollquery := es.EsScrollQuery(SearchResult.ScrollId)

	fornum := (SearchResult.AllHits.Total / 1000) + 1
	searchsum := 0
	for searchsum <= fornum {
		SrollStatus, SrollResult := es.EsSrollSearch(scrollquery, EsdocType, IndexName)
		if SrollStatus == "NotFound" {
			log.Println("NotFound", SrollStatus)
			break
		}
		CodisCreate, CodisClose = DateHandle(CodisCreate, CodisClose, SrollResult)
		searchsum = searchsum + 1
		if searchsum == fornum {
			break
		}
	}

	DBWrite(CodisCreate, CodisClose)

}

func DateHandle(CodisCreate, CodisClose map[string]string, esnum es.EsData) (map[string]string, map[string]string) {
	for _, val := range esnum.AllHits.SouHits {
		Clientip := val.EsSource.Clientip
		Action := val.EsSource.Action
		// Codisname := val.EsSource.Codisname
		Codisip := val.EsSource.Codisip

		if Action == "create" {
			if _, ok := CodisClose[Clientip]; ok {
				delete(CodisClose, Clientip)
			}
			CodisCreate[Clientip] = Codisip
		} else if Action == "closed" {
			if _, ok := CodisCreate[Clientip]; ok {
				delete(CodisCreate, Clientip)
				continue
			} else {
				CodisClose[Clientip] = Codisip
			}

		} else {
			continue
		}

	}
	return CodisCreate, CodisClose
}

func DBWrite(CodisCreate, CodisClose map[string]string) {
	for ckey, _ := range CodisClose {
		if mysql.DB.GraphExist(ckey) == true {
			mysql.DB.GraphDel(ckey)
		}
	}
	for ekey, eval := range CodisCreate {
		codisname, codisid := mysql.DB.CodisProxyLike(eval)
		if mysql.DB.GraphExist(ekey) == false {
			mysql.DB.GraphWrite(ekey, codisname, codisid)
		} else {
			mysql.DB.GraphUpdate(ekey, codisname, codisid)
		}

	}

}
func checkError(err error) {
	if err != nil {
		log.Fatalf("Error Elasticsearch Opration Response: %s", err)
	}
}
