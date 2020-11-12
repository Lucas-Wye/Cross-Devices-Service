package main

import (
	_ "CrossDevicesService/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	// log
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/CrossDevicesService.log","maxdays":10}`)
	// filename and line number
	logs.EnableFuncCallDepth(true)
	logs.Async()
	logs.Async(1e3)

	beego.Run()
}
