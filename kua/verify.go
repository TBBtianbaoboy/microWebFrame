package kua

import (
	"github.com/mojocn/base64Captcha"
)

type V map[string]string

var store = base64Captcha.DefaultMemStore

//生成base64验证码-----------------/verify
func GenerateCaptchaHandler() H {
	drive := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	ca := base64Captcha.NewCaptcha(drive, store)
	id, b64s, err := ca.Generate()
	body := H{"code": 1, "data": b64s, "captchaId": id, "msg": "success"}
	if err != nil {
		WriteErrToLogFile(1,err)
		body = H{"code": 0, "msg": err.Error()}
	}
	return body
}

//辅助函数-----loginAndVerify
func captchaVerifyHandle(vc V) bool {

	id := vc["capId"]
	bas64 := vc["bas64"]
	if store.Verify(id, bas64, true) {
		return true
	}
	return false
}
