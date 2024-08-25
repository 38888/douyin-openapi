package douyin_openapi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

//内容安全

const (
	securityCensorText    = "/api/v2/tags/text/antidirt"
	securityCensorImageV2 = "/api/apps/censor/image"
	securityCensorImageV3 = "/api/apps/v1/censor/image/"
)

type SecurityCensorTextParams struct {
	Tasks []Content `json:"tasks"`
}
type Content struct {
	Value string `json:"content"`
}
type SecurityCensorTextResponse struct {
	LogId string                           `json:"log_id"`
	Data  []SecurityCensorTextResponseData `json:"data"`

	//错误返回
	ErrorId   string `json:"error_id"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Exception string `json:"exception"`
}
type SecurityCensorTextResponseData struct {
	Msg      string                              `json:"msg"`
	Code     int                                 `json:"code"`
	TaskId   string                              `json:"task_id"`
	Predicts []SecurityCensorTextResponsePredict `json:"predicts"`
	DataId   string                              `json:"data_id"`
}
type SecurityCensorTextResponsePredict struct {
	Prob      int    `json:"prob"`
	Hit       bool   `json:"hit"`
	Target    string `json:"target"`
	ModelName string `json:"model_name"`
}

// SecurityCensorText 检测一段文本是否包含违法违规内容。
func (d *DouYinOpenApi) SecurityCensorText(str string) (response SecurityCensorTextResponse, err error) {
	//请求 Headers X-Token
	url := d.GetApiUrl(securityCensorText)
	token, err := d.Config.AccessToken.GetAccessToken()
	if err != nil {
		err = fmt.Errorf("AccessToken error %s", err)
		return
	}
	res, err := resty.New().R().
		SetBody(SecurityCensorTextParams{
			Tasks: []Content{{Value: str}},
		}).
		SetHeader("X-Token", token).
		SetResult(&response).
		SetError(&response).
		Post(url)
	if err != nil {
		return
	}
	if response.Code != 0 || res.StatusCode() != 200 {
		err = fmt.Errorf("SecurityText error %s", response.Message)
		return
	}
	return
}

var ModelName = map[string]string{
	"porn":                        "图片涉黄",
	"cartoon_leader":              "领导人漫画",
	"anniversary_flag":            "特殊标志",
	"sensitive_flag":              "敏感旗帜",
	"sensitive_text":              "敏感文字",
	"leader_recognition":          "敏感人物",
	"bloody":                      "图片血腥",
	"fandongtaibiao":              "未准入台标",
	"plant_ppx":                   "图片涉毒",
	"high_risk_social_event":      "社会事件",
	"high_risk_boom":              "爆炸",
	"high_risk_money":             "人民币",
	"high_risk_terrorist_uniform": "极端服饰",
	"high_risk_sensitive_map":     "敏感地图",
	"great_hall":                  "大会堂",
	"cartoon_porn":                "色情动漫",
	"party_founding_memorial":     "建党纪念",
}

type SecurityCensorImagePredict struct {
	ModelName string `json:"model_name"`
	Hit       bool   `json:"hit"`
}
type SecurityCensorImageV2Params struct {
	AppId       string `json:"app_id"`
	AccessToken string `json:"access_token"`
	Image       string `json:"image"`      //图片链接
	ImageData   string `json:"image_data"` //图片数据的 base64 格式，有 image 字段时，此字段无效
}
type SecurityCensorImageV2Response struct {
	//0 成功
	//1 参数有误
	//2 access_token 校验失败
	//3 图片下载失败
	//4 服务内部错误
	Error    int                          `json:"error"`
	Message  string                       `json:"message"`
	Predicts []SecurityCensorImagePredict `json:"predicts"`
}

// SecurityCensorImageV2 检测图片是否包含违法违规内容。
func (d *DouYinOpenApi) SecurityCensorImageV2(params SecurityCensorImageV2Params) (response SecurityCensorImageV2Response, err error) {
	url := d.GetApiUrl(securityCensorImageV2)
	token, err := d.Config.AccessToken.GetAccessToken()
	if err != nil {
		err = fmt.Errorf("AccessToken error %s", err)
		return
	}
	params.AppId = d.Config.AppId
	params.AccessToken = token
	res, err := resty.New().R().
		SetBody(params).
		SetResult(&response).
		SetError(&response).
		Post(url)
	if err != nil {
		return
	}
	if response.Error != 0 || res.StatusCode() != 200 {
		err = fmt.Errorf("SecurityCensorImageV3 error %s %d", response.Message, response.Error)
		return
	}
	return
}

type SecurityCensorImageV3Params struct {
	AppId     string `json:"app_id"`     //小程序id
	Image     string `json:"image"`      //图片链接
	ImageData string `json:"image_data"` //图片数据的 base64 格式，有 image 字段时，此字段无效
}
type SecurityCensorImageV3Response struct {
	ErrNo    int                          `json:"err_no"`
	ErrMsg   string                       `json:"err_msg"`
	LogId    string                       `json:"log_id"`
	Predicts []SecurityCensorImagePredict `json:"predicts"`
}

// SecurityCensorImageV3 检测图片是否包含违法违规内容。
func (d *DouYinOpenApi) SecurityCensorImageV3(params SecurityCensorImageV3Params) (censorImageV3Response SecurityCensorImageV3Response, err error) {
	url := d.GetApiUrl(securityCensorImageV3)

	params.AppId = d.Config.AppId

	token, err := d.Config.AccessToken.GetAccessToken()
	if err != nil {
		err = fmt.Errorf("AccessToken error %s", err)
		return
	}
	res, err := resty.New().R().
		SetBody(params).
		SetHeader("access-token", token).
		SetResult(&censorImageV3Response).
		SetError(&censorImageV3Response).
		Post(url)
	if err != nil {
		return
	}
	if censorImageV3Response.ErrNo != 0 || res.StatusCode() != 200 {
		err = fmt.Errorf("SecurityCensorImageV3 error %s %d", censorImageV3Response.ErrMsg, censorImageV3Response.ErrNo)
		return
	}
	return
}
