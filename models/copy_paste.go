package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type CopyPaste struct {
	Id      int `orm:"pk"`
	Content string
}

func Write2CopyPaste(input string) {
	t := time.Now()
	s := fmt.Sprintf("%02d.%02d.%0.2d %0.2d:%0.2d\n\n",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())

	write2String := s + input
	o := orm.NewOrm()
	insertText := CopyPaste{Content: write2String}
	_, err := o.Insert(&insertText)
	if err != nil {
		logs.Error(err)
	}
}

func ReadFromCopyPaste() []string {
	var tmp CopyPaste
	o := orm.NewOrm()
	err := o.QueryTable("copy_paste").OrderBy("-id").One(&tmp)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return strings.Split(tmp.Content, "\n")
}
