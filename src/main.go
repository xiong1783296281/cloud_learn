package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func HelloWord(resp http.ResponseWriter, req *http.Request) {
	//var version []byte
	//var err error
	var ip = ClientIP(req)
	log.Println("host: ", ip)
	for s := range req.Header {
		log.Println(s + ":" + req.Header.Get(s))
		resp.Header().Set(s, req.Header.Get(s))
	}
	//command := exec.Command("/bin/bash", "-c", "cat /proc/version")
	//if version, err = command.Output(); err == nil {
	//	resp.Header().Set("VERSION", string(version))
	//} else {
	//	fmt.Println(err)
	//}
	getenv := os.Getenv("VERSION")
	resp.Header().Set("VERSION", getenv)

}

func Healthz(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(200)
}

func start() {
	http.HandleFunc("/", Healthz)
	http.HandleFunc("/healthz", HelloWord)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err.Error())
	}
}

func main() {
	start()
}
