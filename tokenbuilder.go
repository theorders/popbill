package popbill

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/imroc/req"
	"github.com/theorders/aefire"
	"io/ioutil"
	"net/http"
)

type TokenBuilder struct {
	context.Context
	ServiceURL      string
	linkId          string
	RecentServiceID string
	accessId        string
	recentScope     []string
	secret string
}

var instance *TokenBuilder

func (c *Client) Instance() *TokenBuilder {
	if instance == nil {
		instance = &TokenBuilder{}
	}

	instance.Context = c.Context
	instance.linkId = c.LinkId
	instance.accessId = c.CorpNum
	instance.ServiceURL = defaultServiceURL
	instance.secret = c.Secret
	if c.Test {
		instance.RecentServiceID = serviceIdTest
	} else {
		instance.RecentServiceID = serviceIDReal
	}
	instance.AddScopes("member")

	return instance
}

func (builder *TokenBuilder) AddScopes(newScopes ...string) *TokenBuilder {
	for _, newScope := range newScopes {
		exist := false
		for _, scope := range builder.recentScope {
			if scope == newScope {
				exist = true
				break
			}
		}

		if !exist {
			builder.recentScope = append(builder.recentScope, newScope)
		}
	}

	return builder
}

func (builder *TokenBuilder) Build(forwardedIP string) (token *SessionToken, err error) {
	fmt.Printf("build")
	if builder.RecentServiceID == "" {
		return nil, errors.New("서비스아이디가 입력되지 않았습니다")
	}

	uri := "/" + builder.RecentServiceID + "/Token"

	postJson := aefire.ToJson(aefire.MapOf(
		"access_id", builder.accessId,
		"scope", builder.recentScope,
	))
	postData := []byte(postJson)

	invokeTime, err := builder.ServerTime()
	if err != nil {
		return nil, errors.New("팝빌 서버 시간 조회 실패: " + err.Error())
	}

	signTarget := "POST" + newLine
	postDataHash, err := aefire.MD5Base64(postData)
	if err != nil {
		return nil, err
	}
	signTarget += postDataHash + newLine
	signTarget += invokeTime + newLine

	if forwardedIP != "" {
		signTarget += forwardedIP + newLine
	}

	signTarget += apiVersion + newLine
	signTarget += uri

	b64Enc := base64.StdEncoding
	secretKeyData, err := b64Enc.DecodeString(builder.secret)
	if err != nil {
		return nil, err
	}

	signature, err := aefire.HMacSha1(secretKeyData, []byte(signTarget))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(builder.Context,
		http.MethodPost,
		builder.ServiceURL+uri,
		bytes.NewReader(postData),
	)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("x-lh-date", invokeTime)
	//req.Header.Set("x-lh-version", apiVersion)
	req.Header["x-lh-date"] = []string{invokeTime}
	req.Header["x-lh-version"] = []string{apiVersion}

	if forwardedIP != "" {
		//req.Header.Set("x-lh-forwarded", forwardedIP)
		req.Header["x-lh-forwarded"] = []string{forwardedIP}
	}

	req.Header.Set("Authorization", "LINKHUB "+builder.linkId+" "+signature)
	req.Header.Set("Content-Type", "application/json; charset=utf8")

	//println(odsvalue.ToString(req.Header))

	//lowerCaseHeader := make(http.Header)
	//for key, value := range req.Header {
	//	lowerCaseHeader[strings.ToLower(key)] = value
	//}
	//
	//req.Header = lowerCaseHeader

	//fmt.Printf("%v", req.Header)
	//println(odsvalue.ToString(lowerCaseHeader))
	//println(odsvalue.ToString(postData,true))

	httpcli := &http.Client{}

	res, err := httpcli.Do(req)

	//fmt.Printf(odsvalue.ToString(resData))
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	resData, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Printf(aefire.ToJson(resData))
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, errors.New(aefire.ToJson(resData))
	}

	//fmt.Printf(odsvalue.ToString(resData))

	token = &SessionToken{}
	err = json.Unmarshal(resData, token)
	if err != nil {
		fmt.Printf("linkhub token unmarshal failed")
		fmt.Printf(aefire.ToJson(resData))

		return nil, err
	}

	if token.Code < 0 {
		return nil, errors.New(token.Message)
	}

	return token, nil
}

func (builder *TokenBuilder) ServerTime() (string, error) {
	res, err := req.Get(builder.ServiceURL+"/Time", builder.Context)
	if err != nil {
		return "", errors.New("링크허브 서버 접속 실패: " + err.Error())
	}

	return res.String(), nil
}
