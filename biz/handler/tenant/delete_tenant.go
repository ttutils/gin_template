package tenant

import (
	"gin_template/biz/dal"
	"gin_template/biz/response"
	"gin_template/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteReq struct {
	ID string `uri:"id" binding:"required"`
}

// DeleteTenant 删除命名空间
// @Tags 命名空间
// @Summary 删除命名空间
// @Description 删除命名空间
// @Accept application/json
// @Produce application/json
// @Param id path string true "命名空间ID"
// @Success 200 {object} response.CommonResp
// @Security ApiKeyAuth
// @router /api/tenant/delete/:id [DELETE]
func DeleteTenant(c *gin.Context) {
	req := new(DeleteReq)
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp := new(response.CommonResp)

	// 检查是否为管理员
	err := utils.IsAdmin(c)
	if err != nil {
		c.JSON(http.StatusOK, &response.CommonResp{
			Code: response.Code_Unauthorized,
			Msg:  err.Error(),
		})
		return
	}

	// 检查租户下是否还有配置
	id, _ := strconv.Atoi(req.ID)
	tenantInfo, err := dal.GetTenantById(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, &response.CommonResp{Code: response.Code_DBErr, Msg: "查询命名空间失败: " + err.Error()})
		return
	}
	if tenantInfo == nil {
		c.JSON(http.StatusOK, &response.CommonResp{Code: response.Code_Err, Msg: "命名空间不存在"})
		return
	}

	if err = dal.DeleteTenant(uint(id)); err != nil {
		c.JSON(http.StatusOK, &response.CommonResp{Code: response.Code_DBErr, Msg: "删除命名空间失败: " + err.Error()})
		return
	}

	resp.Code = response.Code_Success
	resp.Msg = "删除命名空间成功"

	c.JSON(http.StatusOK, resp)
}
