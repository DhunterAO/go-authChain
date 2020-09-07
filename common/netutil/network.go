package netutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(ip, url string) string {
	fmt.Println("get", ip, url)

	resp, err := http.Get("http://" + ip + url)
	if err != nil {
		fmt.Println(err)
		return "err"
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//fmt.Println(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	//
	//fmt.Println(string(body))
	//fmt.Println(reflect.TypeOf(body))
	return string(body)
}

func BroadCastGet(ipList []string, url string) {
	fmt.Println(ipList)
	for _, ip := range ipList {
		go Get(ip, url)
	}
}

func GoPost(ip string, url string, json []byte) {
	//fmt.Println("post", url, string(json))
	go Post(ip, url, json)
}

func Post(ip string, url string, json []byte) bool {
	addr := "http://" + ip + url
	fmt.Println(addr)
	resp, err := http.Post(addr, "application/json", bytes.NewReader(json))
	fmt.Println(resp, err)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//fmt.Println(string(body))
	return true
}

func BroadCastPost(ipList []string, url string, json []byte) {
	for _, ip := range ipList {
		GoPost(ip, url, json)
	}
}
