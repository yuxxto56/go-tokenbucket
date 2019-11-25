/**
 * 限流器业务逻辑
 * @author liyang<654516092@qq.com>
 */
package token

import (
	"fmt"
	redisgo "go-oversell/go"
	"time"
)

//定义接口
type Bucket interface {
	SetTicker()
	GetBucket()  bool
	pushBucket(timer *time.Ticker)
	Close()
}
//定义令牌桶结构体
type TokenBucket struct {
	key      string
	tLimiter int
	limiter  int
	second   time.Duration
	redis    *redisgo.Redis
	isTure   bool
}

//定义redis参数结构体
type RedisPro struct {
	Host string
	Port string
	Protocol string
	Exp uint64
}

//定义结构体初始化方法
func NewTokenBucket(key string,totalLimit int,limiter int,second time.Duration,redis *RedisPro) Bucket{
	red := &redisgo.Redis{
		Host:    redis.Host,
		Port:    redis.Port,
		Protocol: redis.Protocol,
		Exp:redis.Exp,
	}
	if err := red.GetConn();err != nil{
		panic(fmt.Sprintf("redis connect error:%s",err.Error()))
		return nil
	}
	result := red.SetNx(key,totalLimit)
	return &TokenBucket{
		tLimiter:totalLimit,
		limiter:limiter,
		second:second,
		redis:red,
		isTure:result,
		key:key,
	}
}
//关闭redis连接
func (t *TokenBucket) Close(){
	t.redis.Close()
}

//返回定时器
func (t *TokenBucket) SetTicker(){
	var timer *time.Ticker
	if t.isTure{
		//设置过期时间
		if t.redis.Exp > uint64(0) {
			t.redis.Expire(t.key, t.redis.Exp)
		}
		timer = time.NewTicker(t.second * time.Second)
		t.pushBucket(timer)
	}
}

//每隔多长时间往令牌桶里放令牌
func (t *TokenBucket) pushBucket(timer *time.Ticker){
	for {
		select {
			case <-timer.C:
				//检测key是否有过期
				ttl := t.redis.Ttl(t.key)
				if ttl == int64(-2) {
					timer.Stop()
					break
				}
				num := t.redis.IncrBy(t.key,t.limiter)
				if num > t.limiter {
					diffNum := num - t.tLimiter
					t.redis.DecrBy(t.key, diffNum)
				}
		}
	}
}

//获取令牌
func (t *TokenBucket) GetBucket() bool{
	num := t.redis.DecrBy(t.key,1)
	if num>=0{
		return true
	}else{
		t.redis.IncrBy(t.key,1)
		return false
	}
}


