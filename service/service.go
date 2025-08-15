package service

import (
	"fmt"
	"log"
	"profile-label/common"
)

func GetAddress(loginType, username, password string) error {
	var loginURL string
	var addressURL string
	switch loginType {
	case "base":
		loginURL = common.BaseLoginUrl
		addressURL = common.BaseAddressUrl
	case "eth":
		loginURL = common.EthLoginUrl
		addressURL = common.EthAddressUrl
	case "bsc":
		loginURL = common.BscLoginUrl
		addressURL = common.BscAddressUrl
	default:
		return fmt.Errorf("暂不支持该网站")
	}
	cookies, err := GetBaseCookie(loginURL, username, password)
	if err != nil {
		log.Printf("获取 %s Cookie 失败: %v\n", loginType, err)
		return err
	}
	num := 1
	for i := 0; i < num; i++ {
		url := fmt.Sprintf(addressURL, i+1)
		num, err = GetBasePage(url, cookies)
		if err != nil {
			log.Printf("获取 %s 页码失败: %v\n", loginType, err)
			return err
		}
	}
	return nil
}
