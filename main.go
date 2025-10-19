package main

import (
	_ "CrossDevicesService/routers"

	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "genpass" {
		if len(os.Args) < 3 {
			fmt.Println("用法: go run main.go genpass <password>")
			os.Exit(1)
		}
		password := os.Args[2]
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("生成hash出错:", err)
			os.Exit(2)
		}
		fmt.Println(string(hash))
		return
	} else {
		// log
		logs.SetLogger(logs.AdapterFile, `{"filename":"logs/CrossDevicesService.log","maxdays":10}`)
		// filename and line number
		logs.EnableFuncCallDepth(true)
		logs.Async()
		logs.Async(1e3)

		beego.Run()
	}
}
