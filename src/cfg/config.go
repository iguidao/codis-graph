package cfg

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func Get_Elasticsearch(get_elasticsearch string) string {
	switch get_elasticsearch {
	case "addr":
		elasticsearch_addr := viper.GetString("elasticsearch.addr")
		return elasticsearch_addr
	case "index":
		elasticsearch_index := viper.GetString("elasticsearch.index")
		return elasticsearch_index
	case "doctype":
		elasticsearch_doctype := viper.GetString("elasticsearch.doctype")
		return elasticsearch_doctype
	case "esuser":
		elasticsearch_user := viper.GetString("elasticsearch.user")
		return elasticsearch_user
	case "espassword":
		elasticsearch_password := viper.GetString("elasticsearch.password")
		return elasticsearch_password
	default:
		return "noconfig"
	}
}
func Get_Local(get_local string) string {
	switch get_local {
	case "addr":
		local_addr := viper.GetString("local.addr")
		return local_addr
	case "crongraph":
		local_crongraph := viper.GetString("local.crongraph")
		return local_crongraph
	case "logpath":
		local_logpath := viper.GetString("local.logpath")
		return local_logpath
	case "croncodis":
		local_croncodis := viper.GetString("local.croncodis")
		return local_croncodis
	case "codisurl":
		local_codisurl := viper.GetString("local.codisurl")
		return local_codisurl
	case "estimediff":
		local_estimediff := viper.GetString("local.estimediff")
		return local_estimediff
	default:
		return "noconfig"
	}
}

func Get_Info(get_type string) string {
	switch get_type {
	case "MYSQL":
		mysql_name := viper.GetString("mysql.name")
		mysql_addr := viper.GetString("mysql.addr")
		mysql_username := viper.GetString("mysql.username")
		mysql_password := viper.GetString("mysql.password")
		mysql_url := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", mysql_username, mysql_password, mysql_addr, mysql_name)
		return mysql_url
	default:
		return "noconfig"
	}
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}
	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}
	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		// 如果指定了配置文件，则解析指定的配置文件
		viper.SetConfigFile(c.Name)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		viper.AddConfigPath("yaml")
		viper.SetConfigName("config")
	}
	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")
	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// 监听配置文件是否改变,用于热更新
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed: ", e.Name)
	})
}
