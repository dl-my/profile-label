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
	sessionId := service.GetBaseCookie("dlmydlmy", "Lsx5211314@")

	fmt.Printf("sessionId: %s\n", sessionId)
	service.GetBaseLabel(sessionId)
}
