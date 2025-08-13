package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ApiResponse struct {
	Data []NoteItem `json:"data"`
}

type NoteItem struct {
	ID        string `json:"_id"`
	Hash      string `json:"hash"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	Label     string `json:"label"`
	Note      string `json:"note"`
	Status    int    `json:"status"`
	Type      string `json:"type"`
	UpdatedAt string `json:"updatedAt"`
}

func GetSolscanLabel() {
	url := "https://api-v2.solscan.io/v2/user/label/list"

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 通用Cookie
	req.AddCookie(&http.Cookie{
		Name:  "_ga",
		Value: "GA1.1.1887025992.1754615140",
	})
	req.AddCookie(&http.Cookie{
		Name:  "_ga_PS3V7B7KV0",
		Value: "GS2.1.s1754980488$o3$g1$t1754980514$j34$l0$h0",
	})
	// 核心cookie
	req.AddCookie(&http.Cookie{
		Name:  "auth-token",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxpdXNodWFpeGluZzUyMUBnbWFpbC5jb20iLCJhY3Rpb24iOiJsb2dnZWQiLCJpYXQiOjE3NTQ5ODA1MDAsImV4cCI6MTc2NTc4MDUwMH0.vUUsacLU79Tbu3FF5emstKbMYpjbdl8661M3CXmx9Qk",
	})
	// 人机验证
	req.AddCookie(&http.Cookie{
		Name:  "cf_clearance",
		Value: "VIritQ1D89aV.mfWfzAMVdiUvnIQFpRTdUWEzC7my5c-1754980495-1.2.1.1-2Auq_QHvV2g5ZnvNdh52yN8uoCN1_nLvheQbyZ1fGCXVJJeJkbTzVyMqEeDAglHFfWUaWe4E_65wMaXEw9lZAv6fbY.PE1__qj3.CGpDh.wjHRo0ovyqWEx9iwdNM2sOui3wrqYgeWPZmlhweZ15anN3HFB0iNSqD0D2wbTE7HczAtmhKpfB6.YzYkcd1gJMSgSHQJFAK8TEG1_WyFPfxGOAeNE8k_Z6eegSqdOMj3s",
	})

	// 添加必要的请求头（模仿浏览器）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Origin", "https://solscan.io")

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	// 输出响应状态码和内容
	//fmt.Println("响应内容:", string(body))

	var result ApiResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	// 打印解析结果
	for i, item := range result.Data {
		fmt.Printf("第 %d 条记录:\n", i+1)
		fmt.Printf("  Hash: %s\n", item.Hash)
		fmt.Printf("  Label: %s\n", item.Label)
		fmt.Printf("  Type: %s\n", item.Type)
		fmt.Printf("  data: %v\n", item)
		fmt.Println("  ----")
	}
}
