package oncron

import (
	"log"

	"codisgraph/src/middleware/es"
)

func CronGraph() {

	CodisCreate := make(map[string]string)
	CodisClose := make(map[string]string)
	starttime := "1603798260"
	endtime := "1603798320"
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
	for ckey, cval := range CodisClose {

		log.Println(ckey, cval)
	}
	log.Println("this create codis")
	for ekey, eval := range CodisCreate {
		log.Println(ekey, eval)
	}
}
func checkError(err error) {
	if err != nil {
		log.Fatalf("Error Elasticsearch Opration Response: %s", err)
	}
}
