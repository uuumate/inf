package main

import (
	"context"
	"fmt"
	"time"

	"github.com/uuumate/inf/logging"
	"github.com/uuumate/inf/sql"
)

type DatabaseDB struct {
	Database string `json:"database" gorm:"column:Database"`
}

type Row struct {
	ID uint64 `gorm:"column:id" json:"id"`
}

func main() {
	sqlGroup, err := sql.InitSqlGroup(&sql.Config{
		Name:   "test",
		Driver: "mysql",
		Master: "root:12345678@tcp(127.0.0.1:3306)/message_contact_nvwademo?charset=utf8&parseTime=true&loc=Local",
		Slaves: []string{"root:12345678@tcp(127.0.0.1:3306)/message_contact_nvwademo?charset=utf8&parseTime=true&loc=Local"},
	})

	if err != nil {
		fmt.Printf("InitSqlGroup error: %s\n", err.Error())
		return
	}

	logging.InitLogger(&logging.LogConfig{
		LogPath:  "logs",
		LogLevel: logging.LogLevelDebug,
		Rolling:  logging.RollingFormatDay,
	})

	rows := make([]Row, 0)

	err = sqlGroup.Slave(context.TODO()).Table("contact_0").Select("id").Where("id > 1").Limit(5).Scan(&rows).Error
	if err != nil {
		fmt.Printf("crete table error: %s\n", err.Error())
	}

	fmt.Printf("result: %+v\n", rows)

	time.Sleep(time.Second * 5)
}
