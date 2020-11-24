package main

import (
	"codisgraph/src/cfg"
	"codisgraph/src/http"
	"codisgraph/src/middleware/es"
	"codisgraph/src/middleware/mysql"
	"codisgraph/src/middleware/oncron"
	"log"

	"github.com/robfig/cron"
	// "github.com/robfig/cron"
)

func init() {
	if err := cfg.Init(""); err != nil {
		panic(err)
	}
	mysql.Connect(cfg.Get_Info("MYSQL"))
	es.Connect(cfg.Get_Elasticsearch("addr"), cfg.Get_Elasticsearch("esuser"), cfg.Get_Elasticsearch("espassword"))
}

func main() {

	log.Println("CodisGraph Server Starting")

	c := cron.New()
	//  创建定时任务

	codiscrontime := cfg.Get_Local("croncodis")
	c.AddFunc(codiscrontime, func() {
		oncron.CronCodis()
	})
	graphcrontime := cfg.Get_Local("crongraph")
	c.AddFunc(graphcrontime, func() {
		oncron.CronGraph()
	})

	c.Start()

	listen := cfg.Get_Local("addr")
	if listen == "" {
		listen = ":8080"
	}

	http.NewServer().Run(listen)

}
