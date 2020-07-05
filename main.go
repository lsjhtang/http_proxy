package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type MyHttp1 struct {
}

func (this *MyHttp1) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Print(err)
		}
	}()

	if req.URL.Path == "/favicon.ico" {
		return
	}

	target, _ := url.Parse(LB.RandRobin3().Host)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, req)

}

type MyHttp2 struct {
}

func (*MyHttp2) GetIp(req *http.Request) string {
	ips := req.Header.Get("x-forwarded-for")
	if ips != "" {
		return ips
	}
	return req.RemoteAddr
}

func (this *MyHttp2) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	auth := request.Header.Get("Authorization")
	if auth == "" {
		writer.Header().Set("WWW-Authenticate", `Basic realm="必须验证账号密码"`)
		writer.WriteHeader(401)
		return
	}
	authList := strings.Split(auth, " ")
	if len(authList) == 2 && authList[0] == "Basic" {
		result, err := base64.StdEncoding.DecodeString(authList[1])
		if err == nil && string(result) == "abc:123" {
			writer.Write([]byte(fmt.Sprintf("<h1>MyHttp2,ip=%s</h1>", this.GetIp(request))))
			return
		}
	}

	writer.Write([]byte("账号密码错误"))

}

type MyHttp3 struct {
}

func (this *MyHttp3) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	writer.Write([]byte("<h1>MyHttp3</h1>"))

}

type MyHttp4 struct {
}

func (this *MyHttp4) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	writer.Write([]byte("<h1>MyHttp4</h1>"))

}

func main() {
	c := make(chan os.Signal)
	go (func() {
		header := new(MyHttp1)
		err := http.ListenAndServe(":8080", header)
		if err != nil {
			fmt.Print(err.Error())
		}
	})()

	go (func() {
		header := new(MyHttp2)
		err := http.ListenAndServe(":9091", header)
		if err != nil {
			fmt.Print(err.Error())
		}
	})()

	go (func() {
		header := new(MyHttp3)
		err := http.ListenAndServe(":9092", header)
		if err != nil {
			fmt.Print(err.Error())
		}
	})()

	go (func() {
		header := new(MyHttp4)
		err := http.ListenAndServe(":9093", header)
		if err != nil {
			fmt.Print(err.Error())
		}
	})()

	s := <-c
	log.Print(s)
}
