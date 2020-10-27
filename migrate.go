package main

import (
	"codisgraph/src/cfg"
	"codisgraph/src/middleware/mysql"
	"flag"
)

var migrate = flag.Bool("m", false, "migrate the database schemas.")

func init() {
	if err := cfg.Init(""); err != nil {
		panic(err)
	}
	mysql.Connect(cfg.Get_Info("MYSQL"))
}
func main() {
	flag.Parse()
	if *migrate {
		mysql.Migrate()
		return
	}
}
