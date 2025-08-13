package utils

import (
	"github.com/bwmarrin/snowflake"
	"log"
	"sync"
)

var (
	node     *snowflake.Node
	initOnce sync.Once
)

// InitNode 初始化全局节点，只调用一次
func InitNode(nodeID int64) {
	initOnce.Do(func() {
		var err error
		node, err = snowflake.NewNode(nodeID)
		if err != nil {
			log.Fatalf("初始化 Snowflake 节点失败: %v", err)
		}
	})
}

// GenerateID 对外提供简化接口
func GenerateID() int64 {
	if node == nil {
		log.Fatal("Snowflake 节点未初始化")
	}
	return node.Generate().Int64()
}
