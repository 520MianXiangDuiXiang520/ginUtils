package middleware

import (
	"fmt"
	"github.com/520MianXiangDuiXiang520/GoTools/crypto"
	"github.com/520MianXiangDuiXiang520/agingMap"
	"github.com/520MianXiangDuiXiang520/ginUtils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ThrottledFunc func(ctx *gin.Context) (interface{}, bool)

// Throttled 节流中间件，f 第二个返回值为 false 时，拦截请求，响应状态码
// 为 429（StatusTooManyRequests），f 的第一个返回值是拦截请求后给客户端返回的内容。
func Throttled(f ThrottledFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		if resp, ok := f(context); !ok {
			context.Abort()
			context.JSON(http.StatusTooManyRequests, resp)
		}
	}
}

// 节流规则
type ThrottledRule int

const (
	// 只根据 IP 节流
	ThrottledRuleByIP ThrottledRule = iota + 1

	// 只根据 UA 节流
	ThrottledRuleByUserAgent

	// 根据 UA 和 IP 节流
	ThrottledRuleByUserAgentAndIP

	// 根据自定义字段节流, 若要使用自定义字段
	// 请使用 ThrottledCustom() 方法
	throttledRuleByCustomField
)

func getKey(rule ThrottledRule, ctx *gin.Context, customFields []string) string {
	switch rule {
	case ThrottledRuleByIP:
		return crypto.MD5([]string{ctx.ClientIP()})
	case ThrottledRuleByUserAgent:
		return crypto.MD5([]string{ctx.Request.UserAgent()})
	case ThrottledRuleByUserAgentAndIP:
		return crypto.MD5([]string{ctx.ClientIP(), ctx.Request.UserAgent()})
	case throttledRuleByCustomField:
		fields := make([]string, len(customFields))
		for i, v := range customFields {
			fields[i] = ctx.Request.Header.Get(v)
		}
		return crypto.MD5(fields)
	}
	return ""
}

func parseRate(rate string) (int, time.Duration, error) {
	r := strings.Split(rate, "/")
	if len(r) != 2 {
		msg := fmt.Sprintf("The rate string does not comply with the rules," +
			" please use a style similar to 1s")
		return -1, 0, fmt.Errorf(msg)
	}
	n, d := r[0], r[1]
	frequency, err := strconv.Atoi(n)
	if err != nil {
		return -1, 0, fmt.Errorf("fail to parse from string(%s) to int; error: %w", n, err)
	}
	ts := d[:len(d)-1]
	t, err := strconv.Atoi(ts)
	if err != nil {
		return -1, 0, fmt.Errorf("fail to parse from string(%s) to int; error: %w", n, err)
	}
	duration := map[string]time.Duration{"s": time.Second, "m": time.Minute,
		"h": time.Hour, "d": time.Hour * 24}[strings.ToLower(d)[len(d)-1:]]
	return frequency, duration * time.Duration(t), nil
}

func initCache() *agingMap.AgingMap {
	return agingMap.NewBaseAgingMap(time.Second*5, 0.5)
}

var cache = initCache()

// SimpleThrottle 会根据 rule 为每一次请求生成一个 key, key 相同
// 就会被认为是同一个用户，最简单的情况下，我们会以 IP 和 UA 来判断请求
// 是不是来自同一个人，这时，rule 可以使用下面三个选项：
//   - ThrottledRuleByUserAgent: 只要 UA 相同即视为同一用户
//   - ThrottledRuleByIP: 根据客户端 IP 判断是否是统一用户
//   - ThrottledRuleByUserAgentAndIP: UA 和 IP 都匹配才视为同一个用户
// 这样最简单但页不安全，因为用户可以随便更换自己使用的 UA 和 IP.
// rate 用来控制同一个用户的访问频率，如使用 16/3m 表示三分钟内最多可以访问 16 次
// “/” 后面的 Duration 支持 s(秒)，m(分)，h(小时)，d(天)，大小写不敏感
func SimpleThrottle(rule ThrottledRule, rate string) ThrottledFunc {
	return BaseThrottled(rate, getKeyFunc(rule, nil), SimpleTooManyReqResp)
}

type KeyFunc func(ctx *gin.Context) string

// TooManyRequestResponseFunc 用来返回访问频率过快时的返回值
// waitingSeconds 表示用户还应该等待多长时间才能进行下一次正常访问
type TooManyRequestResponseFunc func(waitingSeconds float64) interface{}

func getKeyFunc(rule ThrottledRule, fields []string) KeyFunc {
	return func(ctx *gin.Context) string {
		return getKey(rule, ctx, fields)
	}
}

// 访问频率太快时的默认返回值：
//   {
//       "code": 429,
//       "msg": "您的请求太快了，休息一下吧 ^_^ (5s)"
//   }
func SimpleTooManyReqResp(waitingSeconds float64) interface{} {
	return ginUtils.BaseRespHeader{Code: http.StatusTooManyRequests,
		Msg: fmt.Sprintf("您的请求太快了，休息一下吧 ^_^ (%ds)", int64(waitingSeconds))}
}

// SimpleThrottledWithFields 允许你传入一组请求头字符串，这些请求头的值
// 将会被合在一起作为 key (参考 SimpleThrottle)
func SimpleThrottledWithFields(rate string, fields []string) ThrottledFunc {
	return BaseThrottled(rate, getKeyFunc(throttledRuleByCustomField, fields), SimpleTooManyReqResp)
}

// SimpleThrottledWithKeyFunc 允许传入一个函数 keyFunc 该函数应该返回一个 string 类型的 key
// 相同的 key 表示同一个用户
func SimpleThrottledWithKeyFunc(rate string, keyFunc KeyFunc) ThrottledFunc {
	return BaseThrottled(rate, keyFunc, SimpleTooManyReqResp)
}

// BaseThrottled 是一个更基础打函数，相比于 SimpleThrottledWithKeyFunc，
// BaseThrottled 允许自定义用户访问频率太快时的返回值
func BaseThrottled(rate string, keyFunc KeyFunc, respFunc TooManyRequestResponseFunc) ThrottledFunc {
	return func(ctx *gin.Context) (interface{}, bool) {
		frequency, duration, err := parseRate(rate)
		if err != nil {
			panic("Unable to parse rate")
		}
		key := keyFunc(ctx)

		if key == "" {
			return nil, true
		}
		history, deadline, ok := cache.TermLoadOrStore(key, 1, duration, func(val interface{}, ok bool) bool {
			// 如果 key 不存在（第一次访问），存储 k, value 置为 1
			if !ok {
				return true
			}
			his := val.(int)
			// 未达到限流的次数
			if his < frequency {
				cache.Store(key, his+1, duration)
			}
			return false
		})
		// 第一次访问
		if ok {
			return nil, true
		}

		his := history.(int)
		if his >= frequency {
			// 直接拦截，不计入 cache
			resp := respFunc(deadline)
			return resp, false
		}
		return nil, true
	}
}
