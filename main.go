package main

import (
	_ "upper/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	// log设置
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/FileService.log","maxdays":10}`)
	// 输出调用的文件名和文件行号
	logs.EnableFuncCallDepth(true)
	logs.Async()
	logs.Async(1e3)

	beego.Run()
}
