package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"simple-demo/service"
	"time"
)

func main() {
	// 设置默认时区为亚洲/上海时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	time.Local = loc

	go service.RunMessageServer()

	r := gin.Default()

	initRouter(r)

	r.Run() // 默认8080端口
}
