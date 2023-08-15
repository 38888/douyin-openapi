package douyin_openapi

import "fmt"

const (
	orderV2Push = "https://developer.toutiao.com/api/apps/order/v2/push" // 订单推送
)

// OrderV2PushParams 订单推送
type OrderV2PushParams struct {
	ClientKey   string `json:"client_key,omitempty"`   // 否 第三方在抖音开放平台申请的 ClientKey 注意：POI 订单必传 awx1334dlkfjdf
	AccessToken string `json:"access_token,omitempty"` // 是 服务端 API 调用标识，通过 getAccessToken 获取
	ExtShopId   string `json:"ext_shop_id,omitempty"`  // 否 POI 店铺同步时使用的开发者侧店铺 ID，购买店铺 ID，长度 < 256 byte 注意：POI 订单必传 ext_112233
	AppName     string `json:"app_name,omitempty"`     // 是 做订单展示的字节系 app 名称，目前为固定值“douyin” douyin
	OpenId      string `json:"open_id,omitempty"`      // 是 小程序用户的 open_id，通过 code2Session 获取 d33432323423
	OrderStatus int64  `json:"order_status,omitempty"` // 否 普通小程序订单订单状态，POI 订单可以忽略 0：待支付 1：已支付 2：已取消 4：已核销（核销状态是整单核销,即一笔订单买了 3 个券，核销是指 3 个券核销的整单） 5：退款中 6：已退款 8：退款失败 注意：普通小程序订单必传，担保支付分账依赖该状态 4
	OrderType   int64  `json:"order_type,omitempty"`   // 是 订单类型，枚举值: 0：普通小程序订单（非POI订单） 9101：团购券订单（POI 订单） 9001：景区门票订单（POI订单）0
	UpdateTime  int64  `json:"update_time,omitempty"`  // 是 订单信息变更时间，13 位毫秒级时间戳 1643189272388
	Extra       string `json:"extra,omitempty"`        // 否 自定义字段，用于关联具体业务场景下的特殊参数，长度 < 2048byte
	OrderDetail string `json:"order_detail,omitempty"` // 是 订单详情，长度 < 2048byte
}

// OrderDetailPOI9101 POI 9101团购卷类型
type OrderDetailPOI9101 struct {
	// OrderV2PushParams
	// OrderDetail
}

type OrderDetailPOI9001 struct {
	// OrderV2PushParams
	// OrderDetail
}

// OrderDetailParams 订单详情
type OrderDetailParams struct {
	OrderId    string     `json:"order_id,omitempty"`    // 是 开发者侧业务单号。用作幂等控制。该订单号是和担保支付的支付单号绑定的，也就是预下单时传入的 out_order_no 字段，长度 <= 64byte 54bb46ba
	CreateTime int64      `json:"create_time,omitempty"` // 是 订单创建的时间，13 位毫秒时间戳 1648453349123
	Status     string     `json:"status,omitempty"`      // 是 订单状态，建议采用以下枚举值： 待支付 已支付 已取消 已超时 已核销 退款中 已退款 退款失败 已支付
	Amount     int64      `json:"amount,omitempty"`      //  是 订单商品总数 2
	TotalPrice int64      `json:"total_price,omitempty"` // 是 订单总价，单位为分 8800
	DetailUrl  string     `json:"detail_url,omitempty"`  // 是 小程序订单详情页 path，长度<=1024 byte
	ItemList   []ItemList `json:"item_list,omitempty"`   // list 是 子订单商品列表，不可为空
}

// ItemList 子订单商品列表
type ItemList struct {
	ItemCode string `json:"item_code,omitempty"` // 是 开发者侧商品 ID，长度 <= 64 byte test_item_code
	Img      string `json:"img,omitempty"`       // 是 子订单商品图片 URL，长度 <= 512 byte https://xxxxxxxxxxxxxxxxxxxxxx
	Title    string `json:"title,omitempty"`     // 是 子订单商品介绍标题，长度 <= 256 byte 好日子
	SubTitle string `json:"sub_title,omitempty"` // 否 子订单商品介绍副标题，长度 <= 256 byte
	Amount   int64  `json:"amount,omitempty"`    // 否 单类商品的数目 2
	Price    int64  `json:"price,omitempty"`     // 是 单类商品的总价，单位为分 4400
}

// OrderV2PushResponse 订单推送返回
type OrderV2PushResponse struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Body    string `json:"body"`
}

// OrderV2Push 订单推送
func (d *DouYinOpenApi) OrderV2Push(normal OrderV2PushParams) (orderV2PushResponse OrderV2PushResponse, err error) {
	// normal.AppId = d.Config.AppId
	// normal.Sign = d.GenerateSign(params)
	err = d.PostJson(orderV2Push, normal, &orderV2PushResponse)
	if err != nil {
		return
	}
	if orderV2PushResponse.ErrCode != 0 {
		err = fmt.Errorf("OrderV2Push error %s %s", orderV2PushResponse.ErrMsg, orderV2PushResponse.Body)
		return
	}
	return
}
