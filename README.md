# go-tokenbucket
go语言、redis、定时器结合实现限流器

### main.go中编写
```cassandraql
package main

import (
	"fmt"
	"go-tokenbucket/token"
	"log"
	"net/http"
)

func HanderOverSell(w http.ResponseWriter, r *http.Request){
	tokens := token.NewTokenBucket(
		"redis_limiter", //redis key名称
		10, //令牌桶中令牌总数
		5, //速率，每10s向令牌桶中新增令牌5个
		10, //时间，10s 
		&token.RedisPro{ //redis配置
		"127.0.0.1",//ip
		"6379", //端口号
		"tcp",//协议
		30,//当前key 过期时间,不过期则填写0即可，单位:秒
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
```
### 测试
* 编译
```cassandraql
1、go build main.go
2、./main &
```
* 模拟访问请求
```cassandraql
 ab -c 20 -n 20  http://127.0.0.1:8000/
 //模拟20次请求，20个用户并发
```
* 结果
```cassandraql
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token成功
2019/11/25 13:23:24 拿到token失败
2019/11/25 13:23:24 拿到token失败
2019/11/25 13:23:24 拿到token失败
2019/11/25 13:23:24 拿到token失败
2019/11/25 13:23:24 拿到token失败
2019/11/25 13:23:24 拿到token失败
2019/11/25 13:23:24 拿到token失败
2019/11/25 13:23:24 拿到token失败
2019/11/25 13:23:24 拿到token失败

```

