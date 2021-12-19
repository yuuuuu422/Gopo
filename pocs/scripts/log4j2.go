package scripts

import "Gopo/utils"

const pocName2 = "Log4j2-cve-2021-44228"

var header = `
Accept-Charset
Accept-Datetime
Accept-Encoding
Accept-Language
Ali-CDN-Real-IP
Authorization
Cache-Control
Cdn-Real-Ip
Cdn-Src-Ip
CF-Connecting-IP
Client-IP
Contact
Cookie
DNT
Fastly-Client-Ip
Forwarded-For-Ip
Forwarded-For
Forwarded
Forwarded-Proto
From
If-Modified-Since
Max-Forwards
Originating-Ip
Origin
Pragma
Proxy-Client-IP
Proxy
Referer
TE
True-Client-Ip
True-Client-IP
Upgrade
User-Agent
Via
Warning
WL-Proxy-Client-IP
X-Api-Version
X-Att-Deviceid
X-ATT-DeviceId
X-Client-IP
X-Client-Ip
X-Client-IP
X-Cluster-Client-IP
X-Correlation-ID
X-Csrf-Token
X-CSRFToken
X-Do-Not-Track
X-Foo-Bar
X-Foo
X-Forwarded-By
X-Forwarded-For-Original
X-Forwarded-For
X-Forwarded-Host
X-Forwarded
X-Forwarded-Port
X-Forwarded-Protocol
X-Forwarded-Proto
X-Forwarded-Scheme
X-Forwarded-Server
X-Forwarded-Ssl
X-Forwarder-For
X-Forward-For
X-Forward-Proto
X-Frame-Options
X-From
X-Geoip-Country
X-Host
X-Http-Destinationurl
X-Http-Host-Override
X-Http-Method-Override
X-HTTP-Method-Override
X-Http-Method
X-Http-Path-Override
X-Https
X-Htx-Agent
X-Hub-Signature
X-If-Unmodified-Since
X-Imbo-Test-Config
X-Insight
X-Ip
X-Ip-Trail
X-Leakix
X-Original-URL
X-Originating-IP
X-ProxyUser-Ip
X-Real-Ip
X-Remote-Addr
X-Remote-IP
X-Requested-With
X-Request-ID
X-True-IP
X-UIDH
X-Wap-Profile
X-WAP-Profile
X-XSRF-TOKEN`

func init(){
	scriptRegister(pocName2,log4j2)
}

func log4j2(target string){
	utils.Debug(pocName2,target)
}

/*
由于种种原因，先把log4j2下了
*/