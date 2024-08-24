package douyin_openapi

const (
	securityText    = "/api/v2/tags/text/antidirt"
	securityImageV2 = "/api/apps/censor/image"
	securityImageV3 = "/api/apps/v1/censor/image/"
)

type SecurityTextParams struct {
	Tasks []struct {
		Content string `json:"content"`
	} `json:"tasks"`
}
type SecurityTextResponse struct {
	LogId string `json:"log_id"`
	Data  []struct {
		Msg      string `json:"msg"`
		Code     int    `json:"code"`
		TaskId   string `json:"task_id"`
		Predicts []struct {
			Prob      int    `json:"prob"`
			Hit       bool   `json:"hit"`
			Target    string `json:"target"`
			ModelName string `json:"model_name"`
		} `json:"predicts"`
		DataId string `json:"data_id"`
	} `json:"data"`
}

func (d *DouYinOpenApi) SecurityText(str string) (securityTextResponse SecurityTextResponse, err error) {
	//请求 Headers X-Token
	//err = d.PostJson(d.GetApiUrl(securityText), normal, &orderV2PushResponse)
	//if err != nil {
	//	return
	//}
	//if orderV2PushResponse.ErrCode != 0 {
	//	err = fmt.Errorf("OrderV2Push error %s %s", orderV2PushResponse.ErrMsg, orderV2PushResponse.Body)
	//	return
	//}
	return
}
