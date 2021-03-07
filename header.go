package ginUtils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	SuccessRespHeader        = BaseRespHeader{Code: http.StatusOK, Msg: "ok"}
	SystemErrorRespHeader    = BaseRespHeader{Code: http.StatusInternalServerError, Msg: "系统异常"}
	ParamErrorRespHeader     = BaseRespHeader{Code: http.StatusBadRequest, Msg: "参数错误"}
	ForbiddenErrorRespHeader = BaseRespHeader{Code: http.StatusForbidden, Msg: "拒绝服务"}
	UnauthorizedRespHeader   = BaseRespHeader{Code: http.StatusUnauthorized, Msg: "未登录"}
)

type BaseRespHeader struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type BaseReqInter interface {
	JSON(ctx *gin.Context) error
}

type BaseRespInter interface{}
