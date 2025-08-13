package service

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetBaseLabel(sessionId string) {
	label := make(map[string]string)
	url := "https://basescan.org/mynotes_address"

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 通用Cookie
	req.AddCookie(&http.Cookie{
		Name:  "_ga",
		Value: "GA1.1.279061314.1754616678",
	})
	req.AddCookie(&http.Cookie{
		Name:  "_ga_TWEL8GRQ12",
		Value: "GS2.1.s1754978580$o7$g1$t1754978785$j32$l0$h0",
	})
	// 核心 Cookie
	req.AddCookie(&http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: sessionId,
	})
	req.AddCookie(&http.Cookie{
		Name:  "__cflb",
		Value: "02DiuJ1fCRi484mKRwMLZ1DrxBLfLhBdetS67mZaJxckL",
	})

	// 添加必要的请求头（模仿浏览器）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	// 输出响应状态码和内容
	//fmt.Println("响应内容:", string(body))

	// 用 strings.NewReader 将 HTML 包装为 io.Reader
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		panic(err)
	}

	// 查找所有含有 data-highlight-target 的 span 标签
	doc.Find("span[data-highlight-target]").Each(func(i int, s *goquery.Selection) {
		if val, exists := s.Attr("data-highlight-target"); exists {
			fmt.Printf("地址 %d: %s\n", i+1, val)
			selector := fmt.Sprintf("#t_%s", strings.ToLower(val))
			span := doc.Find(selector)

			// 获取 span 标签中的文本内容
			text := span.Text()
			fmt.Println("标签内容是:", text)
			label[val] = text
		}
	})
	//fmt.Printf("标签内容是: %v\n", label)
}

func GetBaseCookie(username, password string) string {
	var sessionId string
	// 连接 Chrome
	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), "http://localhost:9222")
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		log.Fatalf("启用网络失败: %v", err)
	}

	// 打开页面并等待用户登录
	loginURL := "https://basescan.org/login"
	err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
		chromedp.WaitVisible(`#ContentPlaceHolder1_txtUserName`, chromedp.ByID),

		chromedp.SendKeys(`#ContentPlaceHolder1_txtUserName`, username, chromedp.ByID),
		chromedp.SendKeys(`#ContentPlaceHolder1_txtPassword`, password, chromedp.ByID),
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

	// 点击登录按钮
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`document.getElementById("ContentPlaceHolder1_btnLogin").click()`, nil),
	)
	if err != nil {
		log.Fatalf("点击登录按钮失败: %v", err)
	}
	fmt.Println("已点击登录，等待跳转...")

	// 等待跳转，或你可以等某个元素出现
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`#showUtcLocalDate`, chromedp.ByID),
	)
	if err != nil {
		log.Fatalf("跳转钮失败: %v", err)
	}
	//time.Sleep(10 * time.Second)

	// 获取 Cookie
	var cookies []*network.Cookie
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().Do(ctx)
		return err
	}))
	if err != nil {
		log.Fatalf("获取 Cookies 失败: %v", err)
	}

	fmt.Println("登录后获取到 Cookies：")
	for _, c := range cookies {
		if c.Name == "ASP.NET_SessionId" {
			sessionId = c.Value
		}
		fmt.Printf("[Cookie] %s = %s\n", c.Name, c.Value)
	}
	return sessionId
}
