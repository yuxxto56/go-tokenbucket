package main

import (
	"fmt"
	"go-tokenbucket/token"
	"log"
	"net/http"
)

func HanderOverSell(w http.ResponseWriter, r *http.Request){
	tokens := token.NewTokenBucket(
		"redis_limiter",
		10,
		5,
		10,
		&token.RedisPro{
		"127.0.0.1",
		"6379",
		"tcp",
		30,
	})
	defer tokens.Close()
	tokens.SetTicker()
	bools := tokens.GetBucket()
	if bools{
		log.Println("拿到token成功")
		//处理业务逻辑
	}else{
		log.Println("拿到token失败")
	   //服务器流量太大，请稍后再试
	}
}

func main() {
	//ticker := time.NewTicker(5*time.Second)
	http.HandleFunc("/", HanderOverSell)
	err := http.ListenAndServe("0.0.0.0:8000", nil)
	if(err != nil){
		fmt.Println("start Http Error,err is ",err)
	}
	fmt.Println("start Http,Success.0.0.0:8000")
}
