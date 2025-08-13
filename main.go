package main

import (
	"fmt"
	"log"
	"profile-label/service"
	"profile-label/utils"
)

func main() {
	err := utils.InitChrome()
	if err != nil {
		log.Fatalf("启动 Chrome 失败: %v", err)
	}
	token, cookie, err := service.GetSolscanCookie()
	if err != nil {
		log.Fatalf("获取 Cookie 失败: %v", err)
	}
	fmt.Println("成功获取 Cookie，开始用 Cookie 调用接口")
	authToken, err := service.GetAuthToken("liushuaixing521@gmail.com", "Lsx5211314@", token, cookie)
	if err != nil {
		log.Fatalf("调用接口失败: %v", err)
	}
	service.GetSolscanLabel(authToken)
}
