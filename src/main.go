package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func HelloWord(resp http.ResponseWriter, req *http.Request) {
	var version []byte
	var err error
	log.Println("host: ", req.Host)
	for s := range req.Header {
		log.Println(s + ":" + req.Header.Get(s))
		resp.Header().Set(s, req.Header.Get(s))
	}
	command := exec.Command("/bin/bash", "-c", "cat /proc/version")
	if version, err = command.Output(); err == nil {
		resp.Header().Set("VERSION", string(version))
	} else {
		fmt.Println(err)
	}
	resp.WriteHeader(200)
}

func start() {
	http.HandleFunc("/healthz", HelloWord)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err.Error())
	}
}

func main() {
	start()
}
