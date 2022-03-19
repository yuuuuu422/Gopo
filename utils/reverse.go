package utils

import (
	"Gopo/utils/proto"
	"bytes"
	"context"
	"fmt"
	"github.com/jweny/xhttp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ceyeApi    string
	ceyeDomain string
)

type ceye struct {
	Domain string `yaml:"domain"`
	Api    string `yaml:"api"`
}

func InitCeyeApi() bool {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		Error(err)
		return false
	}
	c := &ceye{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		Error(err)
		return false
	}
	if c.Api == "" || c.Domain == "" {
		return false
	}
	ceyeApi = c.Api
	ceyeDomain = c.Domain
	return true
}

func newReverse() *proto.Reverse {
	flag := RandLetters(8)
	if ceyeDomain == "" {
		return &proto.Reverse{}
	}
	urlStr := fmt.Sprintf("http://%s.%s", flag, ceyeDomain)
	u, _ := url.Parse(urlStr)
	return &proto.Reverse{
		Flag:               flag,
		Url:                ParseUrl(u),
		Domain:             u.Hostname(),
		Ip:                 "",
		IsDomainNameServer: false,
	}
}

func reverseCheck(r *proto.Reverse, timeout int64) bool {
	if ceyeApi == "" || r.Domain == "" {
		return false
	}
	time.Sleep(time.Second * time.Duration(timeout))
	//请求转化为小写
	sub := strings.ToLower(r.Flag)
	urlStr := fmt.Sprintf("http://api.ceye.io/v1/records?token=%s&type=dns&filter=%s", ceyeApi, sub)
	hr, _ := http.NewRequest("GET", urlStr, nil)
	req := &xhttp.Request{RawRequest: hr}
	ctx := context.Background()
	oResp, err := Client.Do(ctx, req)
	if err != nil {
		Error(err)
		return false
	}
	if !bytes.Contains(oResp.Body, []byte(`"data": []`)) { // api返回结果不为空
		return true
	}
	return false
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"
const letterNumberBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lowletterNumberBytes = "0123456789abcdefghijklmnopqrstuvwxyz"

func RandFromChoices(n int, choices string) string {

	randSource := rand.New(rand.NewSource(time.Now().Unix()))
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	randBytes := make([]byte, n)
	for i, cache, remain := n-1, randSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSource.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			randBytes[i] = choices[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(randBytes)

}

// RandLetters 随机小写字母
func RandLetters(n int) string {
	return RandFromChoices(n, letterBytes)
}

// RandLetterNumbers 随机大小写字母和数字
func RandLetterNumbers(n int) string {
	return RandFromChoices(n, letterNumberBytes)
}

// RandLowLetterNumber 随机小写字母和数字
func RandLowLetterNumber(n int) string {
	return RandFromChoices(n, lowletterNumberBytes)
}
