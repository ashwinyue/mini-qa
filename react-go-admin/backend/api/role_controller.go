package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"react-go-admin-backend/models"
	"react-go-admin-backend/services"
	"react-go-admin-backend/utils"
)

// RoleController 角色控制器
type RoleController struct {
	roleService *services.RoleService
}

// NewRoleController 创建角色控制器
func NewRoleController() *RoleController {
	return &RoleController{
		roleService: &services.RoleService{},
	}
}

// GetList 获取角色列表
func (ctrl *RoleController) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	roles, total, err := ctrl.roleService.GetRoleList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, utils.Error("获取角色列表失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(utils.NewPageData(roles, total, page, pageSize)))
}

// GetDetail 获取角色详情
func (ctrl *RoleController) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	role, err := ctrl.roleService.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(role))
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
}

// Create 创建角色
func (ctrl *RoleController) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("参数错误"))
		return
	}

	role := &models.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	}

	if err := ctrl.roleService.CreateRole(role); err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(role))
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Update 更新角色
func (ctrl *RoleController) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("参数错误"))
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := ctrl.roleService.UpdateRole(uint(id), updates); err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(nil))
}

// Delete 删除角色
func (ctrl *RoleController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := ctrl.roleService.DeleteRole(uint(id)); err != nil {
		c.JSON(http.StatusOK, utils.Error("删除角色失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(nil))
}
