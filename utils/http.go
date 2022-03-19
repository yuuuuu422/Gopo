package utils

import (
	"Gopo/utils/proto"
	"github.com/xiecat/xhttp"
	"net/url"
	"strings"
)

var Client *xhttp.Client

func InitHttp(cookie string, proxy string, debug bool) error {
	var err error
	options := xhttp.DefaultClientOptions()
	if proxy != "" {
		if strings.HasPrefix(proxy, "http://") || strings.HasPrefix(proxy, "https://") {
		} else {
			proxy = "http://" + proxy
		}
		options.Proxy = proxy
	}
	if cookie != "" {
		options.Headers["Cookie"] = cookie
	}
	options.Debug = debug
	Client, err = xhttp.NewClient(options, nil)
	return err
}

func ParseUrl(u *url.URL) *proto.UrlType {
	nu := &proto.UrlType{}
	nu.Scheme = u.Scheme
	nu.Domain = u.Hostname()
	nu.Host = u.Host
	nu.Port = u.Port()
	nu.Path = u.EscapedPath()
	nu.Query = u.RawQuery
	nu.Fragment = u.Fragment
	return nu
}

//func ParseRequest(task Task) *proto.Request {
//	req:=&proto.Request{}
//	u,_:= url.Parse(task.Target)
//	req.Url=ParseUrl(task.Target)
//	fmt.Println(req.Url)
//
//	return req
//}

func ParseResponse(xResp *xhttp.Response) *proto.Response {
	var resp proto.Response
	oResp := xResp.RawResponse
	header := make(map[string]string)
	resp.Status = int32(oResp.StatusCode)
	resp.Url = ParseUrl(oResp.Request.URL)
	for k := range oResp.Header {
		header[k] = oResp.Header.Get(k)
	}
	resp.Headers = header
	//resp.ContentType = oResp.Header.Get("Content-Type")
	resp.Body = xResp.Body
	return &resp
}

func UrlTypeToString(u *proto.UrlType) string {
	var buf strings.Builder
	if u.Scheme != "" {
		buf.WriteString(u.Scheme)
		buf.WriteByte(':')
	}
	if u.Scheme != "" || u.Host != "" {
		if u.Host != "" || u.Path != "" {
			buf.WriteString("//")
		}
		if h := u.Host; h != "" {
			buf.WriteString(u.Host)
		}
	}
	path := u.Path
	if path != "" && path[0] != '/' && u.Host != "" {
		buf.WriteByte('/')
	}
	if buf.Len() == 0 {
		if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i], '/') == -1 {
			buf.WriteString("./")
		}
	}
	buf.WriteString(path)

	if u.Query != "" {
		buf.WriteByte('?')
		buf.WriteString(u.Query)
	}
	if u.Fragment != "" {
		buf.WriteByte('#')
		buf.WriteString(u.Fragment)
	}
	return buf.String()
}
