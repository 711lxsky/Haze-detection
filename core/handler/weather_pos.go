package handler

import (
	"core/constant"
	"core/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SearchPositionWithName(c *gin.Context) {
	//  拿到模糊查询的地名
	var posName *request.QueryPositionRequest
	if err := c.ShouldBind(&posName); err != nil {
		ResponseFail(c, http.StatusBadRequest, constant.DataParseError, err.Error())
		return
	}
	if posName.Position == "" {
		ResponseFail(c, http.StatusBadRequest, constant.DataParseError, "请输入查询地名")
		return
	}
	// 调用api
	requestHandler := NewAPIHandler()
	posList, err := requestHandler.SearchPosition(posName)
	if err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryPositionInfo, err.Error())
		return
	}
	data := map[string]interface{}{
		"pos_list": posList,
	}
	ResponseSuccessWithData(c, data)
}
