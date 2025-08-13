package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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

func GetSolscanLabel(token string) {
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
		Value: token,
	})
	// 人机验证
	req.AddCookie(&http.Cookie{
		Name: "cf_clearance",
		//Value: common.CF_clearance,
		Value: "_0xz.AVz7sD4VOk0kb4s1sasVNy5Cfmoh88VVcb_000-1755072393-1.2.1.1-GoJa8XBt3koRC_9as6wqwY2z.Tm7REwM35P8jM5YkKF7H_4SbrPp.iqR.linve1Jmr3xmrOuyp5p5FUpQuvj9urlwhbadYuNk_KkHukhpz_SZp9gfKbAv80uzBfFD9Bz4RZL6x.wIbqJOndskQ5ttUu.lAsxw1xXFPYK_l2AY3CatBVxkVT_iRYuM58v6X3t..a2TBU_L3u5HqfE0wPCqJ8.ktB2uwwEkt1.iHGGR1Y",
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
	fmt.Println("响应内容:", string(body))

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

func GetSolscanCookie() (string, []*network.Cookie, error) {
	// 连接 Chrome
	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), "http://localhost:9222")
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		log.Fatalf("启用网络失败: %v", err)
	}

	// 打开页面并等待用户登录
	loginURL := "https://solscan.io/user/signin"
	err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
		chromedp.WaitVisible(`#email`, chromedp.ByID),
	)

	// 等待 cf-turnstile-response 有效
	var cfToken string

	for {
		err = chromedp.Run(ctx,
			// 使用 name 属性直接获取 value
			chromedp.Value(`input[name="cf-turnstile-response"]`, &cfToken, chromedp.BySearch),
		)
		if err != nil {
			log.Printf("检测 cf-turnstile-response 失败: %v", err)
		} else if len(cfToken) > 100 && strings.HasPrefix(cfToken, "0.") {
			fmt.Println("Turnstile 验证完成，Token:", cfToken)
			break
		} else {
			fmt.Println("验证中，当前 token:", cfToken)
		}

		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatalf("等待验证码完成失败: %v", err)
	}
	fmt.Println("Turnstile 验证已通过")

	// 获取 Cookie
	var cookies []*network.Cookie
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().Do(ctx)
		return err
	}))
	//TODO :更新CF_clearance
	if err != nil {
		return "", nil, fmt.Errorf("获取 Cookies 失败: %w", err)
	}

	return cfToken, cookies, nil
}

// 把 []*network.Cookie 转成 HTTP 请求用的 Cookie 字符串
func cookiesToHeader(cookies []*network.Cookie) string {
	var sb strings.Builder
	for i, c := range cookies {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(c.Name)
		sb.WriteString("=")
		sb.WriteString(c.Value)
	}
	return sb.String()
}

func GetAuthToken(username, password, token string, cookies []*network.Cookie) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://api-v2.solscan.io/v2/user/login"

	// 如果接口需要发送 JSON 数据，可以按需构造
	payload := map[string]string{
		// 你可以根据API需要调整传递参数，或者为空体
		"email":    username,
		"password": password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 重要：带上 Cookie
	req.Header.Set("Cookie", cookiesToHeader(cookies))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://solscan.io")
	req.Header.Set("Referer", "https://solscan.io")
	req.Header.Set("X-Captcha-Token", token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("请求返回非200状态码: %d", resp.StatusCode)
	}

	// 读取响应内容（示例）
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	// 只定义需要的部分结构体
	type TokenResp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	var tr TokenResp
	if err := json.Unmarshal([]byte(buf.String()), &tr); err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	fmt.Println("响应内容:", buf.String())

	return token, nil
}
