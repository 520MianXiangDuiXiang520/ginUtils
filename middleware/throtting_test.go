package middleware

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func doSomething(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

var client = http.Client{}

func clientGet(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"+
		" AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.190 Safari/537.36")
	req.Header.Add("X-Forwarded-For", "127.0.0.1")
	req.Header.Add("X-Real-Ip", "127.0.0.1")
	return client.Do(req)
}

func TestSimpleThrottle(t *testing.T) {
	go func() {
		r := gin.Default()
		r.GET(
			"/ping",
			Throttled(SimpleThrottle(ThrottledRuleByUserAgentAndIP, "35/1s")),
			doSomething,
		)
		_ = r.Run(":8888")
	}()
	// 等待服务端启动
	time.Sleep(time.Second * 3)
	wg := sync.WaitGroup{}
	wg.Add(35)
	for i := 0; i < 35; i++ {
		go func() {
			defer wg.Done()
			resp, err := clientGet("http://127.0.0.1:8888/ping")
			if err != nil {
				t.Errorf("请求失败！%v", err)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("响应状态码为 %d", resp.StatusCode)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}
			if string(body) != "{\"message\":\"pong\"}" {
				t.Errorf("响应结果错误：%s", string(body))
			}
		}()
	}
	wg.Wait()
	resp, err := clientGet("http://127.0.0.1:8888/ping")
	if err != nil {
		t.Errorf("请求失败！%v", err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("响应状态码为 %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body)[:10] != "{\"code\":429,\"msg\":\"您的请求太快了，休息一下吧 ^_^ (3s)\"}"[:10] {
		t.Errorf("响应结果错误：%s", string(body))
	}

}
