package mysql

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Base struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type CodisInfo struct {
	Base
	ClusterName string `gorm:"size:255"`
	ProxyHost   string `gorm:"size:255"`
	ServerHost  string `gorm:"size:255"`
}

type CodisGraph struct {
	Base
	ClientIp  string    `gorm:"size:255"`
	CodisInfo CodisInfo `json:"codis_info"`
}
