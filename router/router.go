package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"profile-label/api"
)

func Run(port int) {
	r := gin.Default()

	// 注册路由
	registerRoutes(r)

	// 启动服务
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(r *gin.Engine) {
	r.GET("/solscan", api.GetSolscan)
}
