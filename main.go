package main

import (
	"fmt"
	"log"
	"profile-label/service"
	"profile-label/utils"
)

func main1() {
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
		log.Fatalf("获取auth_token失败: %v", err)
	}
	fmt.Printf("成功获取 Auth Token: %s\n", authToken)
	err = service.GetSolscanLabel(authToken)
	if err != nil {
		log.Fatalf("获取标签失败: %v", err)
	}
}

func main() {
	err := utils.InitChrome()
	if err != nil {
		log.Fatalf("启动 Chrome 失败: %v", err)
	}
	cookies, err := service.GetBaseCookie("dlmydlmy", "Lsx5211314@")
	if err != nil {
		log.Fatalf("获取 Base Cookie 失败: %v", err)
	}
	err = service.GetBaseLabel(cookies)
	if err != nil {
		log.Fatalf("获取 Base 标签失败: %v", err)
	}
}
