## 介绍
大家用过redis，不知道用没用过codis，根据codis做了一个调用关系图，可以明确知道谁在使用codis

### 架构图
![](https://github.com/lightsre/codisgraph/blob/master/screenshots/codisgraph-framework.png)

### 开发环境

- running inside Mac
- go1.14.6

### 代码功能列表
1. 查询codis链接 --- 提供http的一个接口，可以查询codis的所有请求
2. 定时获取codis集群信息 --- 通过codis-fe页面地址，抓取codis的所有集群信息
3. 定时获取codis链接信息 --- 通过读取es数据，处理链接信息并存储mysql数据库

### 涉及组件
1. filebeat-7.6.1 --- 收集codis的proxy日志信息
2. kafka 2.1.1 --- 充当消息队列
3. logstash 6.8.0 --- 解析codis信息处理并存入ES
4. Elasticsearch 6.8.0 --- 存储所有的codis链接信息
5. mysql 8.0.21 --- 存储codis集群信息，以及处理后的链接信息

### 前期准备
1. 修改yaml下的config.yaml配置文件
  ```
local 下的 codisurl 配置为codis的fe页面地址，其中estimediff为获取es数据的时间间隔，单位s
mysql 下的链接配置
elasticsearch 下的链接配置
  ```

2. 所有codis的proxy节点都要配置filebeat，收集proxy信息，下面是参考配置

  ```
/etc/filebeat/filebeat.yml 文件内容如下


filebeat.config.inputs:
  enabled: true
  path: /etc/filebeat/input.conf.d/*.yml
  reload.enabled: true
  reload.period: 10s
filebeat.inputs:
- type: log
  enabled: false
  paths:
    - /var/log/*.log
filebeat.config.modules:
  path: ${path.config}/modules.d/*.yml
  reload.enabled: false
setup.template.settings:
  index.number_of_shards: 3
setup.kibana:
processors:
  - add_host_metadata:
      netinfo.enabled: true
  - drop_fields:
      fields: ["beat.hostname","prospector.type","beat.version","beat.name","agent.ephemeral_id", "agent.hostname", "agent.id", "agent.type", "agent.version", "ecs.version", "input.type", "log.offset", "version"]
output.kafka:
  hosts: ["172.30.1.2:9092"]
  topic: '%{[fields.log_topic]}'
  partition.hash:
    reachable_only: false
    hash: ["host", "log.file.path"]
  required_acks: 1
  compression: gzip
  max_message_bytes: 1000000
output.elasticsearch:
  enabled: false
max_procs: 1
  ```

  ```
/etc/filebeat/input.conf.d/codis-proxy.yml 文件



- type: log
  enabled: true
  paths:
    - /data/codis-proxy/logs/codis-proxy.log
  ignore_older: 1h
  fields:
    log_topic: codis-proxy-log
  ```

3. 使用logstash消费kafka，处理信息，并写入ES，下面是参考配置
  ```
input {
    kafka {
            bootstrap_servers => "172.30.1.2:9092"
            group_id => "logstash-codis"
            topics => ["codis-proxy-log"]
            consumer_threads => 1
            codec => "json"
        }
}
filter {
    if "create:" in [message]{
      grok {
          match => {"message" => "%{GREEDYDATA:logdate} %{WORD:cmdfile}.go:%{NUMBER:cmdnum}:\s\[%{DATA:log_level}\]\s%{WORD}\s\[%{NOTSPACE:cmdxxx}\]\s%{WORD:action}:\s%{GREEDYDATA:jsondata}"}
      }
    } else if "closed:" in [message]{
      grok {
          match => {"message" => "%{GREEDYDATA:logdate} %{WORD:cmdfile}.go:%{NUMBER:cmdnum}:\s\[%{DATA:log_level}\]\s%{WORD}\s\[%{NOTSPACE:cmdxxx}\]\s%{WORD:action}:\s%{GREEDYDATA:jsondata}, %{GREEDYDATA:reset}"}
      }
      mutate {
          remove_field => ["reset"]
      }
    }else{
      drop {}
    }
    json {
      source => "jsondata"
      #target => "codisdate"
    }
    date {
      timezone => "Asia/Shanghai"
      match => ["logdate", "yyyy/MM/dd HH:mm:ss"]
      target => "@timestamp"
    }
    mutate {
       split => ["remote", ":"]
       add_field => ["clientip", "%{remote[0]}"]
    }
    mutate {
       add_field => ["codisname", "%{[host][name]}"]
    }
    mutate {
       add_field => ["codisip", "%{[host][ip][0]}"]
    }
    mutate {
        remove_field => ["host","message","create","ops","type","lastop","remote","logdate","jsondata","log_level","cmdfile","cmdxxx","cmdnum"]
    }
}
output {
    elasticsearch {
        hosts => ['172.30.1.3:9200']
        index => "codis-proxy-%{+YYYY.MM.dd}"
        document_type => "codis-proxy"
        user => "root"
        password => "xxxxxxxxxxxxxxx"
    }
}
  ```

### 启动方法
  ```
>\# go mod init . 
>\# go mod tidy
>\# go run migrate.go -m=true
>\# go run main.go
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/health            --> codisgraph/src/http/v1.HealthCheck (3 handlers)
[GIN-debug] GET    /api/v1/cookie            --> codisgraph/src/http/v1.Cookie (3 handlers)
[GIN-debug] GET    /cgraph/v1/getall         --> codisgraph/src/http/v1.GetAll (3 handlers)
[GIN-debug] Listening and serving HTTP on 127.0.0.1:8090

等待5分钟，请求地址： http://127.0.0.1:8090/cgraph/v1/getall    
  ```

### 联系方式
mail: xiaohui920@sina.cn
