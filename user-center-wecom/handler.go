package wecom

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/incubator-answer-plugins/user-center-wecom/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/gin-gonic/gin"
)

// RespBody response body.
type RespBody struct {
	// http code
	Code int `json:"code"`
	// reason key
	Reason string `json:"reason"`
	// response message
	Message string `json:"msg"`
	// response data
	Data interface{} `json:"data"`
}

// NewRespBodyData new response body with data
func NewRespBodyData(code int, reason string, data interface{}) *RespBody {
	return &RespBody{
		Code:   code,
		Reason: reason,
		Data:   data,
	}
}

func (uc *UserCenter) GetRedirectURL(ctx *gin.Context) {
	authorizeUrl := fmt.Sprintf("%s/answer/api/v1/user-center/login/callback", plugin.SiteURL())
	redirectURL := uc.Company.GetRedirectURL(authorizeUrl)
	state := genNonce()
	redirectURL = strings.ReplaceAll(redirectURL, "STATE", state)
	ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "success", map[string]string{
		"redirect_url": redirectURL,
		"key":          state,
	}))
}

func (uc *UserCenter) Sync(ctx *gin.Context) {
	uc.syncCompany()
	if uc.syncSuccess {
		ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "success", map[string]any{
			"message": plugin.Translate(ctx, i18n.ConfigSyncNowSuccessResponse),
		}))
		return
	}
	errRespBodyData := NewRespBodyData(http.StatusBadRequest, "error", map[string]any{
		"err_type": "toast",
	})
	errRespBodyData.Message = plugin.Translate(ctx, i18n.ConfigSyncNowFailedResponse)
	ctx.JSON(http.StatusBadRequest, errRespBodyData)
}

func (uc *UserCenter) Data(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "success", map[string]any{
		"dep":  uc.Company.DepartmentMapping,
		"user": uc.Company.EmployeeMapping,
	}))
}

func (uc *UserCenter) CheckUserLogin(ctx *gin.Context) {
	key := ctx.Query("key")
	val, exist := uc.Cache.Get(key)
	if !exist {
		ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "success", map[string]any{
			"is_login": false,
			"token":    "",
		}))
		return
	}
	token := ""
	externalID, _ := val.(string)
	tokenStr, exist := uc.Cache.Get(externalID)
	if exist {
		token, _ = tokenStr.(string)
	}
	ctx.JSON(http.StatusOK, NewRespBodyData(http.StatusOK, "success", map[string]any{
		"is_login": len(token) > 0,
		"token":    token,
	}))
}

// 随机生成 nonce
func genNonce() string {
	bytes := make([]byte, 10)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
