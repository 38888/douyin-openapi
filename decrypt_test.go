package douyin_openapi

import (
	"fmt"
	"testing"
	"time"
)

func TestDecryptUserInfo(t *testing.T) {
	//search_code := "xxxx" // Login 一次就会失效
	//
	//encryptedData := "xxx=="
	//iv := "xxx=="
	//rawData := `{"nickName":"xxxx","avatarUrl":"xxxx","gender":0,"city":"","province":"","country":"中国","language":""}`
	//signature := "c838a00593e7b0c51acf956ed89e17ab4af2b89a"
	//
	//appId := "ttbxxxxxxx"
	//appSecret := "xxxxxx"
	//
	//lr, err := Login(appId, appSecret, search_code, "")
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//ui, err := DecryptUserInfo(rawData, encryptedData, signature, iv, lr.Data.SessionKey)
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//t.Log("Openid", ui.Openid)
}

func TestDecryptPhoneNumber(t *testing.T) {
	//ssk := "xxx"
	//encryptedData := "xxx=="
	//iv := "xxx=="
	//
	//phone, err := DecryptPhoneNumber(ssk, encryptedData, iv)
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//t.Log("phone", phone.PhoneNumber)

	// Defining time value
	// of Since method
	now := time.Now()
	// 延迟：时间 * 秒（单位）
	time.Sleep(5 * time.Second)
	// Prints time elapse
	fmt.Println("time elapse:",
		time.Since(now))

}
