package netutil

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	//defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
	//	fmt.Println("c")
	//	if err:=recover();err!=nil{
	//		fmt.Println(err) // 这里的err其实就是panic传入的内容，55
	//	}
	//	fmt.Println("d")
	//}()
	ip := "127.0.0.1"
	port := "8880"
	url := "/blockchain"
	fmt.Println(Get(ip, port, url))
}

func TestPostAuthorization(t *testing.T) {
	ip := "127.0.0.1"
	port := "8880"
	url := "/api/authorization"
	auth := types.NewAuthorization()
	fmt.Println(Post(ip, port, url, jsonData))
}