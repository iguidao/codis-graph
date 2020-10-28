package oncron

import (
	"codisgraph/src/cfg"
	"codisgraph/src/middleware/codisapi"
	"codisgraph/src/middleware/mysql"
	"fmt"
	"log"
	"strings"
)

func CronCodis() {
	log.Println("start cron codis")
	codisurl := cfg.Get_Local("codisurl")
	CodisCluster := codisapi.GetCluster(codisurl)
	var CodisServer []string
	var CodisProxy []string
	for _, cname := range CodisCluster {
		CodisProxy = CodisProxy[:0]
		CodisServer = CodisServer[:0]
		codisinfo := codisapi.GetHost(codisurl, cname)
		for _, Serlist := range codisinfo.Group.SModels {
			for _, servername := range Serlist.Servers {
				saddr := strings.Split(servername.Server, ":")
				CodisServer = append(CodisServer, saddr[0])
			}
		}
		for _, Prolist := range codisinfo.Proxy.PModels {
			paddr := strings.Split(Prolist.ProxyAddr, ":")
			CodisProxy = append(CodisProxy, paddr[0])
		}
		codisproxy := strings.Replace(strings.Trim(fmt.Sprint(CodisProxy), "[]"), " ", ",", -1)
		codisserver := strings.Replace(strings.Trim(fmt.Sprint(CodisServer), "[]"), " ", ",", -1)
		if mysql.DB.CodisExist(cname) {
			mysql.DB.CodisUpdate(cname, codisproxy, codisserver)
		} else {
			mysql.DB.CodisWrite(cname, codisproxy, codisserver)
		}
	}
}
