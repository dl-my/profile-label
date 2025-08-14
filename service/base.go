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
	"profile-label/common"
	"strings"
	"time"
)

func GetBaseLabel(cookies []*network.Cookie) error {
	tag := make(map[string]string)
	note := make(map[string]string)
	url := "https://basescan.org/mynotes_address"

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("创建请求失败:", err)
		return err
	}

	// 通用Cookie
	for _, c := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  c.Name,
			Value: c.Value,
		})
	}

	//req.AddCookie(&http.Cookie{
	//	Name:  "_ga",
	//	Value: "GA1.1.1082361134.1755157146",
	//})
	//req.AddCookie(&http.Cookie{
	//	Name:  "_ga_TWEL8GRQ12",
	//	Value: "GS2.1.s1755157146$o1$g1$t1755157382$j47$l0$h0",
	//})
	//// 核心 Cookie
	//req.AddCookie(&http.Cookie{
	//	Name: "ASP.NET_SessionId",
	//	Value: sessionId,
	//})
	//req.AddCookie(&http.Cookie{
	//	Name:  "__cflb",
	//	Value: "02DiuJ1fCRi484mKRwMLZ1DrxBLfLhBdexpqYQ3YhFmtx",
	//})

	// 添加必要的请求头（模仿浏览器）
	req.Header.Set("User-Agent", common.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("请求失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取响应失败:", err)
		return err
	}

	// 输出响应状态码和内容
	//fmt.Println("响应内容:", string(body))

	// 用 strings.NewReader 将 HTML 包装为 io.Reader
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Println("转化GoQuery失败: ", err)
		return err
	}

	// 查找所有含有 data-highlight-target 的 span 标签
	doc.Find("span[data-highlight-target]").Each(func(i int, s *goquery.Selection) {
		if val, exists := s.Attr("data-highlight-target"); exists {
			fmt.Printf("地址 %d: %s\n", i+1, val)
			tagSelector := fmt.Sprintf("#t_%s", strings.ToLower(val))
			tagSpan := doc.Find(tagSelector)

			// 获取 span 标签中的文本内容
			tagText := tagSpan.Text()
			fmt.Println("标签内容是:", tagText)
			tag[val] = tagText
			noteSelector := fmt.Sprintf("#n_%s", strings.ToLower(val))
			noteSpan := doc.Find(noteSelector)

			// 获取 span 标签中的文本内容
			noteText := noteSpan.Text()
			fmt.Println("笔记内容是:", noteText)
			note[val] = noteText
		}
	})
	//fmt.Printf("标签内容是: %v\n", label)
	return nil
}

func GetBaseCookie(username, password string) ([]*network.Cookie, error) {
	// 连接 Chrome
	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), "http://localhost:9222")
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	if err := chromedp.Run(ctx, network.Enable()); err != nil {
		log.Println("启用网络失败:", err)
		return nil, err
	}

	// 打开页面并等待用户登录
	loginURL := "https://basescan.org/login"
	err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
		chromedp.WaitVisible(`#ContentPlaceHolder1_txtUserName`, chromedp.ByID),

		chromedp.SendKeys(`#ContentPlaceHolder1_txtUserName`, username, chromedp.ByID),
		chromedp.SendKeys(`#ContentPlaceHolder1_txtPassword`, password, chromedp.ByID),
	)
	if err != nil {
		log.Println("用户名密码填写失败: ", err)
		return nil, err
	}

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

	// 点击登录按钮
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`document.getElementById("ContentPlaceHolder1_btnLogin").click()`, nil),
	)
	if err != nil {
		log.Fatalf("点击登录按钮失败: %v", err)
	}

	// 等待跳转，或你可以等某个元素出现
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`#showUtcLocalDate`, chromedp.ByID),
	)
	if err != nil {
		log.Println("跳转钮失败:", err)
		return nil, err
	}

	// 获取 Cookie
	var cookies []*network.Cookie
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		cookies, err = network.GetCookies().Do(ctx)
		return err
	}))
	if err != nil {
		log.Println("获取 Cookies 失败: ", err)
		return nil, err
	}

	for _, c := range cookies {
		fmt.Printf("[Cookie] %s = %s\n", c.Name, c.Value)
	}
	return cookies, nil
}
