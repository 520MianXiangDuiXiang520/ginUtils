package ginUtils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type TestReq struct {
	Test string `json:"test"`
}

type TestResp struct {
	BaseRespHeader
	Test string `json:"test"`
}

func (t *TestReq) JSON(ctx *gin.Context) error {
	return ctx.BindJSON(t)
}

func server(req *Request, resp *Response) error {
	r := req.Req.(*TestReq)
	resp.Resp = TestResp{
		BaseRespHeader: SuccessRespHeader,
		Test:           r.Test + "Resp",
	}
	fmt.Println(resp.Resp)
	return nil
}

func check(req *Request, resp *Response) error {
	r := req.Req.(*TestReq)
	if r.Test != "test" {
		resp.Resp = ParamErrorRespHeader
		return errors.New("")
	}
	resp = &Response{Resp: SuccessRespHeader}
	return nil
}

func easyCheck(ctx *gin.Context, req BaseReqInter) (BaseRespInter, error) {
	r := req.(*TestReq)
	if r.Test != "test" {
		return errors.New(""), nil
	}
	resp := SuccessRespHeader
	return resp, nil
}

func easyServer(ctx *gin.Context, req BaseReqInter) BaseRespInter {
	r := req.(*TestReq)
	resp := TestResp{
		BaseRespHeader: SuccessRespHeader,
		Test:           r.Test + "Resp",
	}
	return resp
}

func ExampleHandler() {
	engine := gin.Default()
	defer engine.Run(":8888")
	URLPatterns(engine, "/ping", func(g *gin.RouterGroup) {
		g.POST("/", Handler(server, check, TestReq{}))
	})
}

func ExampleEasyHandler() {
	engine := gin.Default()
	defer engine.Run(":8888")
	URLPatterns(engine, "/ping", func(g *gin.RouterGroup) {
		g.POST("/", EasyHandler(easyCheck, easyServer, TestResp{}))
	})
}

func TestHandler(t *testing.T) {
	go func() {
		engine := gin.Default()
		defer engine.Run(":8888")
		cf := func(req *Request, resp *Response) error {
			r := req.Req.(*TestReq)
			if r.Test != "test" {
				resp.Resp = ParamErrorRespHeader
				return errors.New("")
			}
			resp = &Response{Resp: SuccessRespHeader}
			return nil
		}

		lf := func(req *Request, resp *Response) error {
			r := req.Req.(*TestReq)
			fmt.Println(r)
			resp.Resp = TestResp{
				BaseRespHeader: SuccessRespHeader,
				Test:           r.Test + "Resp",
			}
			fmt.Println(resp.Resp)
			return nil
		}

		URLPatterns(engine, "/ping", func(g *gin.RouterGroup) {
			g.POST("/", Handler(cf, lf, TestReq{}))
		})
	}()
	time.Sleep(time.Second * 3)
	client := http.Client{}
	b, _ := json.Marshal(TestReq{Test: "test"})
	req, err := http.NewRequest("POST", "http://127.0.0.1:8888/ping", bytes.NewBuffer(b))

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	r := TestResp{}
	err = json.Unmarshal(response, &r)
	if err != nil {
		t.Error(err)
	}
	if r.Code != 200 || r.Test != "testResp" {
		t.Error("false")
	}

}
