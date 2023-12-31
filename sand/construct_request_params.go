package sand

import (
	"ceas-go-demo/crypt"
	"ceas-go-demo/utils"
	"encoding/json"
	"fmt"
	"time"
)

const (
	MID         = "6888802117311"       // 商户编号
	SignType    = "SHA1WithRSA"         // 签名方式
	EncryptType = "AES"                 // 加密方式
	Version     = "1.0.0"               // API接口版本
	Layout      = "2006-01-02 15:04:05" // 时间戳格式
)

// Req 请求报文 公共报文头
type Req struct {
	Mid             string `json:"mid"`             // 商户号
	Sign            string `json:"sign"`            // 签名 使用商户私钥对data参数进行RSA签名(SHA1WithRSA)，签名结果采用base64编码
	Timestamp       string `json:"timestamp"`       // 格式 时间戳 2021-02-21 20:28:10
	Version         string `json:"version"`         // 版本号 1.0.0
	CustomerOrderNo string `json:"customerOrderNo"` // 商户订单号
	SignType        string `json:"signType"`        // 签名方式 SHA1WithRSA
	EncryptType     string `json:"encryptType"`     // 加密方式 AES
	EncryptKey      string `json:"encryptKey"`      // 加密 Key 使用杉德公钥对16位随机数进行RSA加密(RSA/ECB/PKCS1Padding) 加密结果采用base64编码
	Data            string `json:"data"`            // 使用16位随机数对明文参数进行AES加密(AES/ECB/PKCS5Padding) 加密结果采用base64编码
}

// CancelAccountData 销户请求数据
type CancelAccountData struct {
	Mid             string `json:"mid"`             // 商户号
	CustomerOrderNo string `json:"customerOrderNo"` // 商户订单号
	BizUserNo       string `json:"bizUserNo"`       // 会员编号 需要注销的会员
	BizType         string `json:"bizType"`         // 操作类型 "CLOSE" 销户
	NotifyUrl       string `json:"notifyUrl"`       // 异步通知地址
}

// ConfirmData 确认销户请求数据
type ConfirmData struct {
	Mid                string `json:"mid"`                // 商户号
	CustomerOrderNo    string `json:"customerOrderNo"`    // 商户订单号
	BizUserNo          string `json:"bizUserNo"`          // 会员编号 需要注销的会员
	OriCustomerOrderNo string `json:"oriCustomerOrderNo"` // 原交易订单号
	SmsCode            string `json:"smsCode"`            // 注销短信验证码
}

//https://api.dev.zijinwenchuang.com/callback/account_callback

// ConstructCancelAccountRequestParams 构造销户请求参数
func ConstructCancelAccountRequestParams(userNo string) (req *Req, err error) {
	// 1.生成订单号
	orderNo := utils.GenerateOrderNo()
	// 2.构造数据
	var data = &CancelAccountData{
		Mid:             MID,
		CustomerOrderNo: orderNo,
		BizUserNo:       userNo,
		BizType:         "CLOSE", // CLOSE
		NotifyUrl:       "https://api.dev.zijinwenchuang.com/ca",
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 3.生成AES Key
	aesKey, err := utils.RandomBytes(16)
	if err != nil {
		return nil, err
	}

	// 4.AES Key加密数据
	encryptedData, err := crypt.AESEncryptECB(rawData, aesKey)
	if err != nil {
		return nil, err
	}
	// 解密测试
	//encryptedDataBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	//if err != nil {
	//	return nil, err
	//}
	//decryptedData, err := crypt.AESDecryptECB(encryptedDataBytes, aesKey)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(string(rawData))
	//fmt.Println(string(decryptedData))
	// 5.RSA算法 加密 ACE Key
	encryptedKey, err := crypt.RSAEncryptECB(aesKey, "./cert/sand_public.cer")
	if err != nil {
		return nil, err
	}

	// 6.对数据进行 签名
	sign, err := Sign([]byte(encryptedData))
	if err != nil {
		return nil, err
	}

	//fmt.Println("data => ", rawData)
	//fmt.Println("key => ", string(aesKey))
	//fmt.Println("encrypted data => ", encryptedData)
	//fmt.Println("encrypted key => ", encryptedKey)
	//fmt.Println("sign => ", sign)

	return &Req{
		Mid:             MID,
		Sign:            sign,
		Timestamp:       time.Now().Format(Layout),
		Version:         Version,
		CustomerOrderNo: orderNo,
		SignType:        SignType,
		EncryptType:     EncryptType,
		EncryptKey:      encryptedKey,
		Data:            encryptedData,
	}, nil
}

func ConstructConfirmRequestParams(oriCustomerOrderNo, userNo, captcha string) (req *Req, err error) {
	// 1.生成订单号
	orderNo := utils.GenerateOrderNo()
	// 2.构造数据
	var data = &ConfirmData{
		Mid:                MID,
		CustomerOrderNo:    orderNo,
		BizUserNo:          userNo,
		SmsCode:            captcha,
		OriCustomerOrderNo: oriCustomerOrderNo,
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(rawData))

	// 3.生成AES Key
	aesKey, err := utils.RandomBytes(16)
	if err != nil {
		return nil, err
	}

	// 4.AES Key加密数据
	encryptedData, err := crypt.AESEncryptECB(rawData, aesKey)
	if err != nil {
		return nil, err
	}
	// 5.RSA算法 加密 ACE Key
	encryptedKey, err := crypt.RSAEncryptECB(aesKey, "./cert/sand_public.cer")
	if err != nil {
		return nil, err
	}

	// 6.对数据进行 签名
	sign, err := Sign([]byte(encryptedData))
	if err != nil {
		return nil, err
	}

	//fmt.Println("data => ", rawData)
	//fmt.Println("key => ", string(aesKey))
	//fmt.Println("encrypted data => ", encryptedData)
	//fmt.Println("encrypted key => ", encryptedKey)
	//fmt.Println("sign => ", sign)

	return &Req{
		Mid:             MID,
		Sign:            sign,
		Timestamp:       time.Now().Format(Layout),
		Version:         Version,
		CustomerOrderNo: orderNo,
		SignType:        SignType,
		EncryptType:     EncryptType,
		EncryptKey:      encryptedKey,
		Data:            encryptedData,
	}, nil
}

func ConstructQueryRequestParams(userNo string) (req *Req, err error) {
	// 1.生成订单号
	orderNo := utils.GenerateOrderNo()
	// 2.构造数据
	var data = &CancelAccountData{
		Mid:             MID,
		CustomerOrderNo: orderNo,
		BizUserNo:       userNo,
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 3.生成AES Key
	aesKey, err := utils.RandomBytes(16)
	if err != nil {
		return nil, err
	}

	// 4.AES Key加密数据
	encryptedData, err := crypt.AESEncryptECB(rawData, aesKey)
	if err != nil {
		return nil, err
	}

	// 5.RSA算法 加密 ACE Key
	encryptedKey, err := crypt.RSAEncryptECB(aesKey, "./cert/sand_public.cer")
	if err != nil {
		return nil, err
	}

	// 6.对数据进行 签名
	sign, err := Sign([]byte(encryptedData))
	if err != nil {
		return nil, err
	}

	return &Req{
		Mid:             MID,
		Sign:            sign,
		Timestamp:       time.Now().Format(Layout),
		Version:         Version,
		CustomerOrderNo: orderNo,
		SignType:        SignType,
		EncryptType:     EncryptType,
		EncryptKey:      encryptedKey,
		Data:            encryptedData,
	}, nil
}
