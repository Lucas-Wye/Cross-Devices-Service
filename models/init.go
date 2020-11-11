package models

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // import your used sql driver
)

func getDatabaseSource() string {
	key := os.Getenv("DATABASE_SOURCE")
	if key == "" {
		logs.Error("Env DATABASE_SOURCE is empty")
	}
	key = key + "?charset=utf8"
	return key
}

func init() {
	// connect to DB
	maxIdle := 500
	maxConn := 4000
	err := orm.RegisterDataBase(
		"default",
		"mysql",
		getDatabaseSource(),
		maxIdle,
		maxConn,
	)
	if err != nil {
		logs.Error("connect mysql err : ", err)
	} else {
		logs.Info("connect mysql success")
	}
	// register model
	orm.RegisterModel(new(CopyPaste))
	// create table
	// orm.RunSyncdb("default", false, true)
	// orm.Debug = true
}

func checkError(err error) {
	if err != nil {
		logs.Error(err)
	}
}
