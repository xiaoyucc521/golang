package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AccessTokenResponse 定义获取access_token的响应数据结构体
type AccessTokenResponse struct {
	ErrCode     int    `json:"errcode"`      // 错误码
	ErrMsg      string `json:"errmsg"`       // 错误信息
	AccessToken string `json:"access_token"` // access_token
	ExpiresIn   int    `json:"expires_in"`   // 过期时间
}

func main() {
	retry := 0
Do:
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=&secret="
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode == -1 && retry < 3 {
		retry++
		goto Do
	}

	var tokenResp AccessTokenResponse

	if err = json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		fmt.Println("错误", err)
	}

	fmt.Println(tokenResp.ErrCode)
	fmt.Println(tokenResp.ErrMsg)
	fmt.Println(tokenResp.AccessToken)

}
