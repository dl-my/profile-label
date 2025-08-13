package utils

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func InitChrome() error {
	// 自动启动本地 Chrome
	cmd := exec.Command(
		`C:\Users\123\AppData\Local\Google\Chrome\Application\chrome.exe`,
		"--remote-debugging-port=9222",
		"--user-data-dir=E:\\chrometmp",
		"--no-first-run",
		"--no-default-browser-check",
	)
	if err := cmd.Start(); err != nil {
		log.Fatalf("启动 Chrome 失败: %v", err)
		return err
	}
	// 等待 Chrome 启动完成
	time.Sleep(1 * time.Second)
	fmt.Println("Chrome 已启动...")
	return nil
}
