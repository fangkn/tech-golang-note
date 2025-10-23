package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-demo-service/types"
)

func GetInfo(c *gin.Context) {

	status := http.StatusOK
	ret := types.HTTPCommonHead{}
	req := types.GetInfoRequest{}
	resp := types.GetInfoResponse{}

	defer func() {
		c.JSON(status, types.HTTPResponse{Head: ret, Body: resp})
	}()

	err := c.ShouldBindQuery(&req)
	if err != nil {
		fmt.Printf("param error: %v", err)
		ret.Code = types.E_PARAM
		ret.Msg = err.Error()
		return
	}

	return
}
