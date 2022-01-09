package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type clientinfo struct {
	ipaddr string
	//未找到获取方式
	httpCode string
}

func remoteAddrSplitIp(s string) string{
	idx:=strings.LastIndex(s,":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

/*
func getRequestIPCode(w http.ResponseWriter,r *http.Request) (string,string){
	xrealip:=r.Header.Get("X-Real-IP")
	xforwardfor:=r.Header.Get("X-Forwarded-For")
	if xrealip == "" && xforwardfor==""{
		xrealip=remoteAddrSplitIp(r.RemoteAddr)
	}
	if xforwardfor !="" {
		parts:=strings.Split(xforwardfor,",")
		for i,p:=range parts{
			parts[i]=strings.TrimSpace(p)
		}
		xrealip=parts[0]
	}

}
*/

func Handler(w http.ResponseWriter,r * http.Request){

	//获取 request 的 header 写入 response 的 header
	//接收客户端 request，并将 request 中带的 header 写入 response header
	if len(r.Header)>0{
		for k,v :=range r.Header{
			w.Header().Add(k,v[0])
		}
	}

	//读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	//获取客户端主机名
	name,err:=os.Hostname()

	ver:=os.Getenv("Version")
	if len(ver)==0 {
        ver="not found version"
	}
	w.Header().Add("Version:",ver)

	if err==nil{
		w.Header().Add("hostname:",name)
	}
	w.Header().Add("OperatorSystemInfo:",runtime.GOOS+ "--" +runtime.GOARCH  )




	//给客户端回复数据
	w.Write([]byte("hello one"))
	//
	io.WriteString(w," ok\n")
	io.WriteString(w,remoteAddrSplitIp(r.RemoteAddr))


}

func Log(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)

		handler.ServeHTTP(w, r)

	})

}

func secHandler(w http.ResponseWriter,r * http.Request){
	w.WriteHeader(http.StatusOK)
	//io.WriteString(w,string(http.StatusOK))
	io.WriteString(w,string("the value is 200"))
}

func main() {
	//注册处理函数，用户连接，自动调用指定的处理函数
	http.HandleFunc("/",Handler)
	http.HandleFunc("/healthz",secHandler)

	//监听绑定
	//err:=http.ListenAndServe(":29010",nil)
	err:=http.ListenAndServe(":29010",Log(http.DefaultServeMux))
	if err != nil{
		log.Fatal(err)
	}
}
