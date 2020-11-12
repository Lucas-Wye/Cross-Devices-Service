package models

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	// "github.com/astaxie/beego/orm"
)

/*
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
*/

func WriteToFile(filename string, content string) bool {
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	} else {
		f, err = os.Create(filename)
	}
	if err != nil {
		logs.Error(err)
		return false
	}

	t := time.Now()
	s := fmt.Sprintf("%02d.%02d.%0.2d %0.2d:%0.2d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())

	write2String := "\n" + s + content + "\n"
	_, err = io.WriteString(f, write2String) //写入文件(字符串)
	defer f.Close()
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

func ReadFromFile(filename string) []string {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("read fail", err)
		return []string{""}
	}
	return strings.Split(string(f), "\n")
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
