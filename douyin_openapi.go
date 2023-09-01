package douyin_openapi

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	accessToken "github.com/38888/douyin-openapi/access-token"
	"github.com/38888/douyin-openapi/cache"
	"github.com/38888/douyin-openapi/util"
)

const (
	code2Session = "/api/apps/v2/jscode2session" // 小程序登录地址
)

// DouYinOpenApiConfig 实例化配置
type DouYinOpenApiConfig struct {
	AppId       string
	AppSecret   string
	AccessToken accessToken.AccessToken
	Cache       cache.Cache
	IsSandbox   bool
	Token       string
	Salt        string
}

// DouYinOpenApi 基类
type DouYinOpenApi struct {
	Config  DouYinOpenApiConfig
	BaseApi string
}

// NewDouYinOpenApi 实例化一个抖音openapi实例
func NewDouYinOpenApi(config DouYinOpenApiConfig) *DouYinOpenApi {
	if config.Cache == nil {
		config.Cache = cache.NewMemory()
	}
	if config.AccessToken == nil {
		config.AccessToken = accessToken.NewDefaultAccessToken(config.AppId, config.AppSecret, config.Cache, config.IsSandbox)
	}
	BaseApi := "https://developer.toutiao.com"
	if config.IsSandbox {
		BaseApi = "https://open-sandbox.douyin.com"
	}
	return &DouYinOpenApi{
		Config:  config,
		BaseApi: BaseApi,
	}
}

// GetApiUrl 获取api地址
func (d *DouYinOpenApi) GetApiUrl(url string) string {
	return fmt.Sprintf("%s%s", d.BaseApi, url)
}

// PostJson 封装公共的请求方法
func (d *DouYinOpenApi) PostJson(api string, params interface{}, response interface{}) (err error) {
	body, err := util.PostJSON(api, params)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	return
}

// Code2SessionParams 小程序登录 所需参数
type Code2SessionParams struct {
	Appid         string `json:"appid,omitempty"`
	Secret        string `json:"secret,omitempty"`
	AnonymousCode string `json:"anonymous_code,omitempty"`
	Code          string `json:"code,omitempty"`
}

// Code2SessionResponse 小程序登录返回值
type Code2SessionResponse struct {
	ErrNo   int                      `json:"err_no,omitempty"`
	ErrTips string                   `json:"err_tips,omitempty"`
	Data    Code2SessionResponseData `json:"data,omitempty"`
}

type Code2SessionResponseData struct {
	SessionKey      string `json:"session_key,omitempty"`
	Openid          string `json:"openid,omitempty"`
	AnonymousOpenid string `json:"anonymous_openid,omitempty"`
	UnionId         string `json:"unionid,omitempty"`
}

// Code2Session 小程序登录
func (d *DouYinOpenApi) Code2Session(code, anonymousCode string) (code2SessionResponse Code2SessionResponse, err error) {
	params := Code2SessionParams{
		Appid:         d.Config.AppId,
		Secret:        d.Config.AppSecret,
		AnonymousCode: anonymousCode,
		Code:          code,
	}
	err = d.PostJson(d.GetApiUrl(code2Session), params, &code2SessionResponse)
	if err != nil {
		return
	}
	if code2SessionResponse.ErrNo != 0 {
		return code2SessionResponse, fmt.Errorf("小程序登录错误: %s %d", code2SessionResponse.ErrTips, code2SessionResponse.ErrNo)
	}
	return
}

type UserInfo struct {
	AvatarUrl string    `json:"avatarUrl"`
	NickName  string    `json:"nickName"`
	Gender    int64     `json:"gender"`
	City      string    `json:"city"`
	Province  string    `json:"province"`
	Country   string    `json:"country"`
	Language  string    `json:"language"`
	Watermark Watermark `json:"watermark"`
}
type Watermark struct {
	Appid     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

func (d *DouYinOpenApi) Decrypt(encryptedData, sessionKey, iv string) *UserInfo {
	src, _ := base64.StdEncoding.DecodeString(encryptedData)
	_key, _ := base64.StdEncoding.DecodeString(sessionKey)
	_iv, _ := base64.StdEncoding.DecodeString(iv)

	block, _ := aes.NewCipher(_key)
	mode := cipher.NewCBCDecrypter(block, _iv)
	dst := make([]byte, len(src))
	mode.CryptBlocks(dst, src)
	var p UserInfo
	if err := json.Unmarshal(dst, &p); err != nil {
		return nil
	}
	return &p
}
