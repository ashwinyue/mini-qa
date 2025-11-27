package api

import (
	"net/http"
	"strconv"

	"react-go-admin-backend/models"
	"react-go-admin-backend/services"
	"react-go-admin-backend/utils"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{
		userService: &services.UserService{},
	}
}

// GetList 获取用户列表
func (ctrl *UserController) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	users, total, err := ctrl.userService.GetUserList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, utils.Error("获取用户列表失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(utils.NewPageData(users, total, page, pageSize)))
}

// GetDetail 获取用户详情
func (ctrl *UserController) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(user))
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Realname string `json:"realname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
}

// Create 创建用户
func (ctrl *UserController) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("参数错误"))
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Realname: req.Realname,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Status:   1,
	}

	if err := ctrl.userService.CreateUser(user); err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(user))
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username"`
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   *int   `json:"status"`
}

// Update 更新用户
func (ctrl *UserController) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Error("参数错误"))
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Realname != "" {
		updates["realname"] = req.Realname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := ctrl.userService.UpdateUser(uint(id), updates); err != nil {
		c.JSON(http.StatusOK, utils.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Success(nil))
}

// Delete 删除用户
func (ctrl *UserController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := ctrl.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusOK, utils.Error("删除用户失败"))
		return
	}

	c.JSON(http.StatusOK, utils.Success(nil))
}
