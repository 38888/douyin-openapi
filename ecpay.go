package douyin_openapi

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const (
	createOrder          = "/api/apps/ecpay/v1/create_order"             // 预下单
	queryOrder           = "/api/apps/ecpay/v1/query_order"              // 订单查询
	createRefund         = "/api/apps/ecpay/v1/create_refund"            // 退款
	queryRefund          = "/api/apps/ecpay/v1/query_refund"             // 退款结果查询
	settle               = "/api/apps/ecpay/v1/settle"                   // 结算
	querySettle          = "/api/apps/ecpay/v1/query_settle"             // 结算结果查询
	unsettleAmount       = "/api/apps/ecpay/v1/unsettle_amount"          // 可结算金额查询
	createReturn         = "/api/apps/ecpay/v1/create_return"            // 退分账
	queryReturn          = "/api/apps/ecpay/v1/query_return"             // 退分账结果查询
	queryMerchantBalance = "/api/apps/ecpay/saas/query_merchant_balance" // 商户余额查询
	merchantWithdraw     = "/api/apps/ecpay/saas/merchant_withdraw"      // 商户提现
	queryWithdrawOrder   = "/api/apps/ecpay/saas/query_withdraw_order"   // 提现结果查询
)

// CreateOrderParams 预下单接口参数
type CreateOrderParams struct {
	AppId           string          `json:"app_id,omitempty"`            // app_id string 是 64 小程序APPID tt07e3715e98c9aac0
	OutOrderNo      string          `json:"out_order_no,omitempty"`      // out_order_no string 是 64 开发者侧的订单号。 只能是数字、大小写字母_-*且在同一个app_id下唯一 7056505317450041644
	TotalAmount     int64           `json:"total_amount,omitempty"`      // total_amount number 是 取值范围： [1,10000000000] 支付价格。 单位为[分] 100，即1元
	Subject         string          `json:"subject,omitempty"`           // subject string 是 128 商品描述。 长度限制不超过 128 字节且不超过 42 字符 抖音商品XYZ
	Body            string          `json:"body,omitempty"`              // body string 是 128 商品详情 长度限制不超过 128 字节且不超过 42 字符 抖音商品XYZ
	ValidTime       int64           `json:"valid_time,omitempty"`        // valid_time number 是 取值范围： [300,172800] 订单过期时间(秒)。最小5分钟，最大2天，小于5分钟会被置为5分钟，大于2天会被置为2天 900，即15分钟
	Sign            string          `json:"sign,omitempty"`              // sign string 是 344 签名，详见签名DEMO 21fc77aeeaad725d9500062a888888a2a3d
	CpExtra         string          `json:"cp_extra,omitempty"`          // cp_extra string 否 2048 开发者自定义字段，回调原样回传。 超过最大长度会被截断 502205261403349
	NotifyUrl       string          `json:"notify_url,omitempty"`        // notify_url string 否 256 商户自定义回调地址，必须以 https 开头，支持 443 端口。 指定时，支付成功后抖音会请求该地址通知开发者 https://api.iiyyeixin.com/Notify/bytedancePay
	ThirdpartyId    string          `json:"thirdparty_id,omitempty"`     // thirdparty_id 条件选填 服务商模式接入必传 64 第三方平台服务商 id，非服务商模式留空 tt84a4f2177777e29df
	StoreUid        string          `json:"store_uid,omitempty"`         // store_uid string 条件选填 多门店模式下可传 64 可用此字段指定本单使用的收款商户号（目前为灰度功能，需要联系平台运营添加白名单，白名单添加1小时后生效；未在白名单的小程序，该字段不生效） 70084531288883795888
	DisableMsg      int             `json:"disable_msg,omitempty"`       // disable_msg number 否 是否屏蔽支付完成后推送用户抖音消息，1-屏蔽 0-非屏蔽，默认为0。 特别注意： 若接入POI, 请传1。因为POI订单体系会发消息，所以不用再接收一次担保支付推送消息， 1
	MsgPage         string          `json:"msg_page,omitempty"`          // msg_page string 否 支付完成后推送给用户的抖音消息跳转页面，开发者需要传入在app.json中定义的链接，如果不传则跳转首页。 pages/orderDetail/orderDetail?no = DYMP8218048851499944448\u0026id = 797775
	ExpandOrderInfo ExpandOrderInfo `json:"expand_order_info,omitempty"` // expand_order_info 否 - 订单拓展信息，详见下面 expand_order_info参数说明 { "original_delivery_fee":10, "actual_delivery_fee":10 }
	LimitPayWay     string          `json:"limit_pay_way,omitempty"`     // limit_pay_way string 否 64 屏蔽指定支付方式，屏蔽多个支付方式，请使用逗号","分割，枚举值： 屏蔽微信支付：LIMIT_WX 屏蔽支付宝支付：LIMIT_ALI 屏蔽抖音支付：LIMIT_DYZF 特殊说明：若之前开通了白名单，平台会保留之前屏蔽逻辑；若传入该参数，会优先以传入的为准，白名单则无效 屏蔽抖音支付和微信支付： "LIMIT_DYZF,LIMIT_WX"
}

type ExpandOrderInfo struct {
	OriginalDeliveryFee int
	ActualDeliveryFee   int
}

// CreateOrderResponse 预下单返回值
type CreateOrderResponse struct {
	ErrNo   int                     `json:"err_no,omitempty"`
	ErrTips string                  `json:"err_tips,omitempty"`
	Data    CreateOrderResponseData `json:"data,omitempty"`
}

type CreateOrderResponseData struct {
	OrderId    string `json:"order_id,omitempty"`
	OrderToken string `json:"order_token,omitempty"`
}

// CreateOrder 预下单
func (d *DouYinOpenApi) CreateOrder(params CreateOrderParams) (createOrderResponse CreateOrderResponse, err error) {
	params.AppId = d.Config.AppId
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(createOrder), params, &createOrderResponse)
	if err != nil {
		return
	}
	if createOrderResponse.ErrNo != 0 {
		return createOrderResponse, fmt.Errorf("%s %d", createOrderResponse.ErrTips, createOrderResponse.ErrNo)
	}
	return
}

// GenerateSign 生成请求签名
func (d *DouYinOpenApi) GenerateSign(params interface{}) string {
	var paramsMap map[string]interface{}
	var paramsArr []string
	j, _ := json.Marshal(&params)
	err := json.Unmarshal(j, &paramsMap)
	if err != nil {
		return ""
	}
	for k, v := range paramsMap {
		if k == "other_settle_params" || k == "app_id" || k == "thirdparty_id" || k == "sign" || k == "salt" || k == "token" {
			continue
		}
		value := strings.TrimSpace(fmt.Sprintf("%v", v))
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) > 1 {
			value = value[1 : len(value)-1]
		}
		value = strings.TrimSpace(value)
		if value == "" || value == "null" {
			continue
		}
		paramsArr = append(paramsArr, value)
	}
	paramsArr = append(paramsArr, d.Config.Salt)
	sort.Strings(paramsArr)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(paramsArr, "&"))))
}

// QueryOrderParams 订单查询接口参数
type QueryOrderParams struct {
	AppId        string `json:"app_id,omitempty"`
	OutOrderNo   string `json:"out_order_no,omitempty"`
	Sign         string `json:"sign,omitempty"`
	ThirdpartyId string `json:"thirdparty_id,omitempty"`
}

type QueryOrderResponse struct {
	ErrNo       int         `json:"err_no,omitempty"`
	ErrTips     string      `json:"err_tips,omitempty"`
	OutOrderNo  string      `json:"out_order_no,omitempty"`
	OrderId     string      `json:"order_id,omitempty"`
	PaymentInfo PaymentInfo `json:"payment_info,omitempty"`
	CpsInfo     string      `json:"cps_info,omitempty"`
	//CpsInfo     CpsInfo     `json:"cps_info,omitempty"`
}

type PaymentInfo struct {
	TotalFee    int    `json:"total_fee,omitempty"`
	OrderStatus string `json:"order_status,omitempty"` // SUCCESS：成功 TIMEOUT：超时未支付 PROCESSING：处理中 FAIL：失败
	PayTime     string `json:"pay_time,omitempty"`     // 支付时间， 格式为"yyyy-MM-dd hh:mm:ss"
	Way         int    `json:"way,omitempty"`          // 支付渠道， 1-微信支付，2-支付宝支付，10-抖音支付
	ChannelNo   string `json:"channel_no,omitempty"`
	SellerUid   string `json:"seller_uid,omitempty"`
	ItemId      string `json:"item_id,omitempty"`
	CpsInfo     string `json:"cps_info,omitempty"`
}
type CpsInfo struct {
	TotalFee string `json:"share_amount,omitempty"` //达人分佣金额，单位为分
	DouyinId string `json:"douyin_id,omitempty"`    //达人抖音号
	Nickname string `json:"nickname,omitempty"`     //达人昵称
}

// QueryOrder 支付结果查询
func (d *DouYinOpenApi) QueryOrder(outOrderNo, thirdpartyId string) (queryOrderResponse QueryOrderResponse, err error) {
	queryParams := QueryOrderParams{
		AppId:        d.Config.AppId,
		OutOrderNo:   outOrderNo,
		ThirdpartyId: thirdpartyId,
	}
	queryParams.Sign = d.GenerateSign(queryParams)
	err = d.PostJson(d.GetApiUrl(queryOrder), queryParams, &queryOrderResponse)
	if err != nil {
		return
	}
	if queryOrderResponse.ErrNo != 0 {
		return queryOrderResponse, fmt.Errorf("%s %d", queryOrderResponse.ErrTips, queryOrderResponse.ErrNo)
	}
	return
}

// PayCallbackResponse 支付回调结构体
type PayCallbackResponse struct {
	Timestamp    string                  `json:"timestamp,omitempty"`
	Nonce        string                  `json:"nonce,omitempty"`
	Msg          string                  `json:"msg,omitempty"`
	MsgStruct    PayCallbackResponseData `json:"msg_struct"`
	MsgSignature string                  `json:"msg_signature,omitempty"`
	Type         string                  `json:"type,omitempty"`
}

type PayCallbackResponseData struct {
	Appid          string `json:"appid,omitempty"`
	CpOrderNo      string `json:"cp_orderno,omitempty"`
	CpExtra        string `json:"cp_extra,omitempty"`
	Way            string `json:"way,omitempty"`
	PaymentOrderNo string `json:"payment_order_no,omitempty"`
	ChannelNo      string `json:"channel_no,omitempty"`
	TotalAmount    int    `json:"total_amount,omitempty"`
	Extra          string `json:"extra,omitempty"`
	Status         string `json:"status,omitempty"`
	ItemId         string `json:"item_id,omitempty"`
	SellerUid      string `json:"seller_uid,omitempty"`
	PaidAt         int64  `json:"paid_at,omitempty"`
	OrderId        string `json:"order_id,omitempty"`
}

// CheckResponseSign 校验回调签名
func (d *DouYinOpenApi) CheckResponseSign(oldSign string, strArr []string) error {
	sort.Strings(strArr)
	h := sha1.New()
	h.Write([]byte(strings.Join(strArr, "")))
	newSign := fmt.Sprintf("%x", h.Sum(nil))
	if newSign != oldSign {
		return fmt.Errorf("回调验签失败 newSign:%s oldSign:%s", newSign, oldSign)
	}
	return nil
}

// PayCallback 支付结果回调
func (d *DouYinOpenApi) PayCallback(body PayCallbackResponse, checkSign bool) (PayCallbackResponseData PayCallbackResponseData, err error) {
	// 判断是否需要校验签名
	if checkSign {
		sortedString := make([]string, 0)
		sortedString = append(sortedString, d.Config.Token)
		sortedString = append(sortedString, body.Timestamp)
		sortedString = append(sortedString, body.Nonce)
		sortedString = append(sortedString, body.Msg)
		err = d.CheckResponseSign(body.MsgSignature, sortedString)
		if err != nil {
			return
		}
	}
	// 解析 msg 数据到结构体
	err = json.Unmarshal([]byte(body.Msg), &PayCallbackResponseData)
	if err != nil {
		return
	}
	return
}

// CreateRefundParams 发起退款参数
type CreateRefundParams struct {
	AppId        string `json:"app_id,omitempty"`
	OutOrderNo   string `json:"out_order_no,omitempty"`
	OutRefundNo  string `json:"out_refund_no,omitempty"`
	Reason       string `json:"reason,omitempty"`
	RefundAmount int    `json:"refund_amount,omitempty"`
	Sign         string `json:"sign,omitempty"`
	CpExtra      string `json:"cp_extra,omitempty"`
	NotifyUrl    string `json:"notify_url,omitempty"`
	ThirdpartyId string `json:"thirdparty_id,omitempty"`
	DisableMsg   int    `json:"disable_msg,omitempty"`
	MsgPage      string `json:"msg_page,omitempty"`
}

type CreateRefundResponse struct {
	ErrNo    int    `json:"err_no,omitempty"`
	ErrTips  string `json:"err_tips,omitempty"`
	RefundNo string `json:"refund_no,omitempty"`
}

// CreateRefund 发起退款
func (d *DouYinOpenApi) CreateRefund(params CreateRefundParams) (createRefundResponse CreateRefundResponse, err error) {
	params.AppId = d.Config.AppId
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(createRefund), params, &createRefundResponse)
	if err != nil {
		return
	}
	if createRefundResponse.ErrNo != 0 {
		return createRefundResponse, fmt.Errorf("CreateRefund error %s %d", createRefundResponse.ErrTips, createRefundResponse.ErrNo)
	}
	return
}

// QueryRefundParams 查询退款参数
type QueryRefundParams struct {
	OutRefundNo  string `json:"out_refund_no,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	Sign         string `json:"sign,omitempty"`
	ThirdpartyId string `json:"thirdparty_id,omitempty"`
}

type QueryRefundParamsResponse struct {
	ErrNo      int    `json:"err_no,omitempty"`
	ErrTips    string `json:"err_tips,omitempty"`
	RefundInfo struct {
		RefundNo     string `json:"refund_no,omitempty"`
		RefundAmount int    `json:"refund_amount,omitempty"`
		RefundStatus string `json:"refund_status,omitempty"`
		RefundedAt   int    `json:"refunded_at,omitempty"`
		IsAllSettled bool   `json:"is_all_settled,omitempty"`
		CpExtra      string `json:"cp_extra,omitempty"`
		Msg          string `json:"msg,omitempty"`
	} `json:"refundInfo"`
}

// QueryRefund 退款结果查询
func (d *DouYinOpenApi) QueryRefund(outRefundNo, thirdpartyId string) (queryRefundParamsResponse QueryRefundParamsResponse, err error) {
	params := QueryRefundParams{
		OutRefundNo:  outRefundNo,
		AppId:        d.Config.AppId,
		ThirdpartyId: thirdpartyId,
	}
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(queryRefund), params, &queryRefundParamsResponse)
	if err != nil {
		return
	}
	if queryRefundParamsResponse.ErrNo != 0 {
		return queryRefundParamsResponse, fmt.Errorf("QueryRefund error %s %d", queryRefundParamsResponse.ErrTips, queryRefundParamsResponse.ErrNo)
	}
	return
}

// RefundCallbackResponse 退款回调结果
type RefundCallbackResponse struct {
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	Msg          string `json:"msg"`
	MsgStruct    RefundCallbackResponseMsg
	MsgSignature string `json:"msg_signature"`
	Type         string `json:"type"`
}

type RefundCallbackResponseMsg struct {
	Appid        string `json:"appid"`
	CpRefundNo   string `json:"cp_refundno"`
	CpExtra      string `json:"cp_extra"`
	Status       string `json:"status"`
	RefundAmount int    `json:"refund_amount"`
	IsAllSettled bool   `json:"is_all_settled"`
	RefundedAt   int    `json:"refunded_at"`
	Message      string `json:"message"`
	OrderId      string `json:"order_id"`
	RefundNo     string `json:"refund_no"`
}

// RefundCallback 退款结果回调
func (d *DouYinOpenApi) RefundCallback(body string, checkSign bool) (refundCallbackResponse RefundCallbackResponse, err error) {
	err = json.Unmarshal([]byte(body), &refundCallbackResponse)
	if err != nil {
		return
	}
	// 判断是否需要校验签名
	if checkSign {
		sortedString := make([]string, 0)
		sortedString = append(sortedString, d.Config.Token)
		sortedString = append(sortedString, refundCallbackResponse.Timestamp)
		sortedString = append(sortedString, refundCallbackResponse.Nonce)
		sortedString = append(sortedString, refundCallbackResponse.Msg)
		err = d.CheckResponseSign(refundCallbackResponse.MsgSignature, sortedString)
		if err != nil {
			return
		}
	}
	var msgStruct RefundCallbackResponseMsg
	// 解析 msg 数据到结构体
	err = json.Unmarshal([]byte(refundCallbackResponse.Msg), &msgStruct)
	if err != nil {
		return
	}
	refundCallbackResponse.MsgStruct = msgStruct
	return
}

// SettleParams 发起分账参数
type SettleParams struct {
	OutSettleNo  string `json:"out_settle_no,omitempty"`
	OutOrderNo   string `json:"out_order_no,omitempty"`
	SettleDesc   string `json:"settle_desc,omitempty"`
	NotifyUrl    string `json:"notify_url,omitempty"`
	CpExtra      string `json:"cp_extra,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	Sign         string `json:"sign,omitempty"`
	SettleParams string `json:"settle_params,omitempty"`
	ThirdpartyId string `json:"thirdparty_id,omitempty"`
	Finish       string `json:"finish,omitempty"`
}

// SettleParamsItem 分账方参数
type SettleParamsItem struct {
	MerchantUid string `json:"merchant_uid,omitempty"`
	Amount      int    `json:"amount,omitempty"`
}

// SettleResponse 分账结果
type SettleResponse struct {
	ErrNo    int    `json:"err_no,omitempty"`
	ErrTips  string `json:"err_tips,omitempty"`
	SettleNo string `json:"settle_no,omitempty"`
}

// Settle 发起结算及分账
func (d *DouYinOpenApi) Settle(settleParams SettleParams, settleParamsItem ...SettleParamsItem) (settleResponse SettleResponse, err error) {
	settleParams.AppId = d.Config.AppId
	settleItem, _ := json.Marshal(settleParamsItem)
	settleParams.SettleParams = string(settleItem)
	settleParams.Sign = d.GenerateSign(settleParams)
	err = d.PostJson(d.GetApiUrl(settle), settleParams, &settleResponse)
	if err != nil {
		return
	}
	if settleResponse.ErrNo != 0 {
		err = fmt.Errorf("settle error %s %d", settleResponse.ErrTips, settleResponse.ErrNo)
		return
	}
	return
}

// QuerySettleParams 结算结果查询参数
type QuerySettleParams struct {
	Sign         string `json:"sign,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	OutSettleNo  string `json:"out_settle_no,omitempty"`
	ThirdpartyId string `json:"thirdparty_id,omitempty"`
}

// QuerySettleResponse 结算结果返回值
type QuerySettleResponse struct {
	ErrNo      int    `json:"err_no"`
	ErrTips    string `json:"err_tips"`
	SettleInfo struct {
		SettleNo     string `json:"settle_no"`
		SettleAmount int    `json:"settle_amount"`
		SettleStatus string `json:"settle_status"`
		SettleDetail string `json:"settle_detail"`
		SettledAt    int    `json:"settled_at"`
		Rake         int    `json:"rake"`
		Commission   int    `json:"commission"`
		CpExtra      string `json:"cp_extra"`
		Msg          string `json:"msg"`
	} `json:"settle_info"`
}

// QuerySettle 结算结果查询 querySettle
func (d *DouYinOpenApi) QuerySettle(outSettleNo, thirdpartyId string) (querySettleResponse QuerySettleResponse, err error) {
	params := QuerySettleParams{
		AppId:        d.Config.AppId,
		OutSettleNo:  outSettleNo,
		ThirdpartyId: thirdpartyId,
	}
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(querySettle), params, &querySettleResponse)
	if err != nil {
		return
	}
	if querySettleResponse.ErrNo != 0 {
		err = fmt.Errorf("QuerySettle error %s %d", querySettleResponse.ErrTips, querySettleResponse.ErrNo)
		return
	}
	return
}

type SettleCallbackResponse struct {
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	Type         string `json:"type"`
	Msg          string `json:"msg"`
	MsgStruct    SettleCallbackResponseMsg
	MsgSignature string `json:"msg_signature"`
}

type SettleCallbackResponseMsg struct {
	AppId           string `json:"app_id"`
	CpSettleNo      string `json:"cp_settle_no"`
	CpExtra         string `json:"cp_extra"`
	Status          string `json:"status"`
	Rake            int    `json:"rake"`
	Commission      int    `json:"commission"`
	SettleDetail    string `json:"settle_detail"`
	SettledAt       int    `json:"settled_at"`
	Message         string `json:"message"`
	OrderId         string `json:"order_id"`
	ChannelSettleId string `json:"channel_settle_id"`
	SettleAmount    int    `json:"settle_amount"`
	SettleNo        string `json:"settle_no"`
	OutOrderNo      string `json:"out_order_no"`
	IsAutoSettle    bool   `json:"is_auto_settle"`
}

// SettleCallback 结算结果回调
func (d *DouYinOpenApi) SettleCallback(body string, checkSign bool) (settleCallbackResponse SettleCallbackResponse, err error) {
	err = json.Unmarshal([]byte(body), &settleCallbackResponse)
	if err != nil {
		return
	}
	// 判断是否需要校验签名
	if checkSign {
		sortedString := make([]string, 0)
		sortedString = append(sortedString, d.Config.Token)
		sortedString = append(sortedString, settleCallbackResponse.Timestamp)
		sortedString = append(sortedString, settleCallbackResponse.Nonce)
		sortedString = append(sortedString, settleCallbackResponse.Msg)
		err = d.CheckResponseSign(settleCallbackResponse.MsgSignature, sortedString)
		if err != nil {
			return
		}
	}
	var msgStruct SettleCallbackResponseMsg
	// 解析 msg 数据到结构体
	err = json.Unmarshal([]byte(settleCallbackResponse.Msg), &msgStruct)
	if err != nil {
		return
	}
	settleCallbackResponse.MsgStruct = msgStruct
	return
}

// UnsettleAmountParams 可分账余额查询
type UnsettleAmountParams struct {
	OutOrderNo     string `json:"out_order_no,omitempty"`
	AppId          string `json:"app_id,omitempty"`
	Sign           string `json:"sign,omitempty"`
	ThirdpartyId   string `json:"thirdparty_id,omitempty"`
	OutItemOrderNo string `json:"out_item_order_no,omitempty"`
}

type UnsettleAmountResponse struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
	Data    struct {
		OutOrderNo     string `json:"out_order_no"`
		UnsettleAmount int    `json:"unsettle_amount"`
		Detail         struct {
			PayInfo struct {
				OutOrderNo string `json:"out_order_no"`
				Amount     int    `json:"amount"`
			} `json:"pay_info"`
			RefundInfo []struct {
				OutRefundNo string `json:"out_refund_no"`
				Amount      int    `json:"amount"`
			} `json:"refund_info"`
			PaymentRake int `json:"payment_rake"`
			LifeRake    int `json:"life_rake"`
			Commission  int `json:"commission"`
		} `json:"detail"`
	} `json:"data"`
}

// UnsettleAmount 可分账余额查询 unsettleAmount
func (d *DouYinOpenApi) UnsettleAmount(outOrderNo, thirdpartyId, outItemOrderNo string) (unsettleAmountResponse UnsettleAmountResponse, err error) {
	params := UnsettleAmountParams{
		OutOrderNo:     outOrderNo,
		AppId:          d.Config.AppId,
		ThirdpartyId:   thirdpartyId,
		OutItemOrderNo: outItemOrderNo,
	}
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(unsettleAmount), params, &unsettleAmountResponse)
	if err != nil {
		return
	}
	if unsettleAmountResponse.ErrNo != 0 {
		err = fmt.Errorf("UnsettleAmount error %s %d", unsettleAmountResponse.ErrTips, unsettleAmountResponse.ErrNo)
		return
	}
	return
}

// CreateReturnParams 退分账参数
type CreateReturnParams struct {
	AppId        string `json:"app_id,omitempty"`
	ThirdpartyId string `json:"thirdparty_id,omitempty"`
	OutSettleNo  string `json:"out_settle_no,omitempty"`
	SettleNo     string `json:"settle_no,omitempty"`
	OutReturnNo  string `json:"out_return_no,omitempty"`
	MerchantUid  string `json:"merchant_uid,omitempty"`
	ReturnAmount int    `json:"return_amount,omitempty"`
	ReturnDesc   string `json:"return_desc,omitempty"`
	CpExtra      string `json:"cp_extra,omitempty"`
	Sign         string `json:"sign,omitempty"`
}

type CreateReturnResponse struct {
	ErrNo      int    `json:"err_no"`
	ErrTips    string `json:"err_tips"`
	ReturnInfo struct {
		AppId        string `json:"app_id"`
		ThirdpartyId string `json:"thirdparty_id"`
		SettleNo     string `json:"settle_no"`
		OutSettleNo  string `json:"out_settle_no"`
		OutReturnNo  string `json:"out_return_no"`
		MerchantUid  string `json:"merchant_uid"`
		ReturnAmount int    `json:"return_amount"`
		ReturnStatus string `json:"return_status"`
		ReturnNo     string `json:"return_no"`
		FailReason   string `json:"fail_reason"`
		FinishTime   int    `json:"finish_time"`
		CpExtra      string `json:"cp_extra"`
	} `json:"return_info"`
}

// CreateReturn 退分账 createReturn
func (d *DouYinOpenApi) CreateReturn(params CreateReturnParams) (createReturnResponse CreateReturnResponse, err error) {
	params.AppId = d.Config.AppId
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(createReturn), params, &createReturnResponse)
	if err != nil {
		return
	}
	if createReturnResponse.ErrNo != 0 {
		err = fmt.Errorf("CreateReturn error %s %d", createReturnResponse.ErrTips, createReturnResponse.ErrNo)
		return
	}
	return
}

// QueryReturnParams 退分账结果查询 参数
type QueryReturnParams struct {
	AppId        string `json:"app_id,omitempty"`
	ReturnNo     string `json:"return_no,omitempty"`
	OutReturnNo  string `json:"out_return_no,omitempty"`
	Sign         string `json:"sign,omitempty"`
	ThirdpartyId string `json:"thirdparty_id,omitempty"`
}

// QueryReturnResponse 退分账结果查询
type QueryReturnResponse struct {
	ErrNo      int    `json:"err_no"`
	ErrTips    string `json:"err_tips"`
	ReturnInfo struct {
		AppId        string `json:"app_id"`
		ThirdpartyId string `json:"thirdparty_id"`
		SettleNo     string `json:"settle_no"`
		OutSettleNo  string `json:"out_settle_no"`
		OutReturnNo  string `json:"out_return_no"`
		MerchantUid  string `json:"merchant_uid"`
		ReturnAmount int    `json:"return_amount"`
		ReturnStatus string `json:"return_status"`
		ReturnNo     string `json:"return_no"`
		FailReason   string `json:"fail_reason"`
		FinishTime   int    `json:"finish_time"`
		CpExtra      string `json:"cp_extra"`
	} `json:"return_info"`
}

// QueryReturn 退分账结果查询 queryReturn
func (d *DouYinOpenApi) QueryReturn(returnNo, outReturnNo, thirdpartyId string) (queryReturnResponse QueryReturnResponse, err error) {
	params := QueryReturnParams{
		AppId:        d.Config.AppId,
		ReturnNo:     returnNo,
		OutReturnNo:  outReturnNo,
		ThirdpartyId: thirdpartyId,
	}
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(queryReturn), params, &queryReturnResponse)
	if err != nil {
		return
	}
	if queryReturnResponse.ErrNo != 0 {
		err = fmt.Errorf("queryReturnResponse error %s %d", queryReturnResponse.ErrTips, queryReturnResponse.ErrNo)
		return
	}
	return
}

// QueryMerchantBalanceParams 可提现余额查询
type QueryMerchantBalanceParams struct {
	ThirdpartyId   string `json:"thirdparty_id,omitempty"`
	AppId          string `json:"app_id,omitempty"`
	MerchantUid    string `json:"merchant_uid,omitempty"`
	ChannelType    string `json:"channel_type,omitempty"` // alipay: 支付宝 wx: 微信	hz: 抖音支付
	Sign           string `json:"sign,omitempty"`
	MerchantEntity string `json:"merchant_entity"`
}

type QueryMerchantBalanceResponse struct {
	ErrNo       int    `json:"err_no"`
	ErrTips     string `json:"err_tips"`
	AccountInfo struct {
		OnlineBalance       int `json:"online_balance"`
		WithDrawableBalance int `json:"withdrawable_balacne"`
		FreezeBalance       int `json:"freeze_balance"`
	} `json:"account_info"`
	SettleInfo struct {
		SettleType    int    `json:"settle_type"`
		SettleAccount string `json:"settle_account"`
		BankcardNo    string `json:"bankcard_no"`
		BankName      string `json:"bank_name"`
	} `json:"settle_info"`
	MerchantEntity int `json:"merchant_entity"`
}

// QueryMerchantBalance 可提现余额查询
func (d *DouYinOpenApi) QueryMerchantBalance(params QueryMerchantBalanceParams) (queryMerchantBalanceResponse QueryMerchantBalanceResponse, err error) {
	params.AppId = d.Config.AppId
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(queryMerchantBalance), params, &queryMerchantBalanceResponse)
	if err != nil {
		return
	}
	if queryMerchantBalanceResponse.ErrNo != 0 {
		err = fmt.Errorf("queryMerchantBalanceResponse error %s %d", queryMerchantBalanceResponse.ErrTips, queryMerchantBalanceResponse.ErrNo)
		return
	}
	return
}

// MerchantWithdrawParams 商户提现参数
type MerchantWithdrawParams struct {
	ThirdpartyId   string `json:"thirdparty_id,omitempty"`
	AppId          string `json:"app_id,omitempty"`
	MerchantUid    string `json:"merchant_uid,omitempty"`
	ChannelType    string `json:"channel_type,omitempty"` // alipay: 支付宝 wx: 微信  hz: 抖音支付 yeepay: 易宝
	WithdrawAmount int    `json:"withdraw_amount,omitempty"`
	OutOrderId     string `json:"out_order_id,omitempty"`
	Sign           string `json:"sign,omitempty"`
	Callback       string `json:"callback,omitempty"`
	CpExtra        string `json:"cp_extra,omitempty"`
	MerchantEntity int    `json:"merchant_entity,omitempty"`
}

type MerchantWithdrawResponse struct {
	ErrNo          int    `json:"err_no"`
	ErrTips        string `json:"err_tips"`
	OrderId        string `json:"order_id"`
	MerchantEntity int    `json:"merchant_entity"`
}

// MerchantWithdraw 提现
func (d *DouYinOpenApi) MerchantWithdraw(params MerchantWithdrawParams) (merchantWithdrawResponse MerchantWithdrawResponse, err error) {
	params.AppId = d.Config.AppId
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(merchantWithdraw), params, &merchantWithdrawResponse)
	if err != nil {
		return
	}
	if merchantWithdrawResponse.ErrNo != 0 {
		err = fmt.Errorf("merchantWithdrawResponse error %s %d", merchantWithdrawResponse.ErrTips, merchantWithdrawResponse.ErrNo)
		return
	}
	return
}

// QueryWithdrawOrderParams 提现结果查询
type QueryWithdrawOrderParams struct {
	ThirdpartyId string `json:"thirdparty_id,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	MerchantUid  string `json:"merchant_uid,omitempty"`
	ChannelType  string `json:"channel_type,omitempty"`
	OutOrderId   string `json:"out_order_id,omitempty"`
	Sign         string `json:"sign,omitempty"`
}

type QueryWithdrawOrderResponse struct {
	ErrNo     int    `json:"err_no"`
	ErrTips   string `json:"err_tips"`
	Status    string `json:"status"` // 状态枚举值: 成功:SUCCESS 失败: FAIL 处理中: PROCESSING 退票: REEXCHANGE 注： 退票：商户的提现申请请求通过渠道（微信/支付宝/抖音支付）提交给银行处理后，银行返回结果是处理成功，渠道返回给商户提现成功，但间隔一段时间后，银行再次通知渠道处理失败并返还款项给渠道，渠道再将该笔失败款返还至商户在渠道的账户余额中
	StatusMsg string `json:"statusMsg"`
}

// QueryWithdrawOrder 提现结果查询
func (d *DouYinOpenApi) QueryWithdrawOrder(params QueryWithdrawOrderParams) (queryWithdrawOrderResponse QueryWithdrawOrderResponse, err error) {
	params.AppId = d.Config.AppId
	params.Sign = d.GenerateSign(params)
	err = d.PostJson(d.GetApiUrl(queryWithdrawOrder), params, &queryWithdrawOrderResponse)
	if err != nil {
		return
	}
	if queryWithdrawOrderResponse.ErrNo != 0 {
		err = fmt.Errorf("queryWithdrawOrderResponse error %s %d", queryWithdrawOrderResponse.ErrTips, queryWithdrawOrderResponse.ErrNo)
		return
	}
	return
}

// MerchantWithdrawCallbackResponse 提现回调返回值解析
type MerchantWithdrawCallbackResponse struct {
	MsgSignature string `json:"msg_signature"`
	Nonce        string `json:"nonce"`
	Timestamp    string `json:"timestamp"`
	Type         string `json:"type"`
	Msg          string `json:"msg"`
	MsgStruct    MerchantWithdrawCallbackResponseMsg
}

type MerchantWithdrawCallbackResponseMsg struct {
	Status     string `json:"status"`
	Extra      string `json:"extra"`
	Message    string `json:"message"`
	WithdrawAt int    `json:"withdraw_at"`
	OrderId    string `json:"order_id"`
	OutOrderId string `json:"out_order_id"`
	ChOrderId  string `json:"ch_order_id"`
}

// MerchantWithdrawCallback 提现回调
func (d *DouYinOpenApi) MerchantWithdrawCallback(body string, checkSign bool) (merchantWithdrawCallbackResponse MerchantWithdrawCallbackResponse, err error) {
	err = json.Unmarshal([]byte(body), &merchantWithdrawCallbackResponse)
	if err != nil {
		return
	}
	// 判断是否需要校验签名
	if checkSign {
		sortedString := make([]string, 0)
		sortedString = append(sortedString, d.Config.Token)
		sortedString = append(sortedString, merchantWithdrawCallbackResponse.Timestamp)
		sortedString = append(sortedString, merchantWithdrawCallbackResponse.Nonce)
		sortedString = append(sortedString, merchantWithdrawCallbackResponse.Msg)
		err = d.CheckResponseSign(merchantWithdrawCallbackResponse.MsgSignature, sortedString)
		if err != nil {
			return
		}
	}
	var msgStruct MerchantWithdrawCallbackResponseMsg
	// 解析 msg 数据到结构体
	err = json.Unmarshal([]byte(merchantWithdrawCallbackResponse.Msg), &msgStruct)
	if err != nil {
		return
	}
	merchantWithdrawCallbackResponse.MsgStruct = msgStruct
	return
}
