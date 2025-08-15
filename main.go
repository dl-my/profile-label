package main

import (
	"fmt"
	"log"
	"profile-label/service"
	"profile-label/utils"
	"time"
)

func main1() {
	start := time.Now() // 开始计时
	err := utils.InitChrome()
	if err != nil {
		log.Fatalf("启动 Chrome 失败: %v", err)
	}
	token, cookie, err := service.GetSolscanCookie()
	if err != nil {
		log.Fatalf("获取 Cookie 失败: %v", err)
	}
	authToken, err := service.GetAuthToken("liushuaixing521@gmail.com", "Lsx5211314@", token, cookie)
	if err != nil {
		log.Fatalf("获取auth_token失败: %v", err)
	}
	fmt.Printf("成功获取 Auth Token: %s\n", authToken)
	err = service.GetSolscanLabel(authToken)
	if err != nil {
		log.Fatalf("获取标签失败: %v", err)
	}
	elapsed := time.Since(start) // 计算耗时
	fmt.Printf("请求耗时: %.3f 秒\n", elapsed.Seconds())
}

func main() {
	start := time.Now() // 开始计时
	err := utils.InitChrome()
	if err != nil {
		log.Fatalf("启动 Chrome 失败: %v", err)
	}
	// 支持base、eth、bsc
	loginType := "bsc"
	err = service.GetAddress(loginType, "dlmydlmy", "Lsx5211314@")
	if err != nil {
		log.Fatalf("获取 %s 地址 失败: %v", loginType, err)
	}
	//cookies, err := service.GetBaseCookie("eth", "dlmydlmy", "Lsx5211314@")
	//if err != nil {
	//	log.Fatalf("获取 Base Cookie 失败: %v", err)
	//}
	//num := 1
	//for i := 0; i < num; i++ {
	//	url := fmt.Sprintf(common.EthAddressUrl, i+1)
	//	num, err = service.GetBasePage(url, cookies)
	//	if err != nil {
	//		log.Fatalf("获取 Base 页码失败: %v", err)
	//	}
	//}
	elapsed := time.Since(start) // 计算耗时
	fmt.Printf("请求耗时: %.3f 秒\n", elapsed.Seconds())
}
