package models

import (
	"log"
	"time"

	"react-go-admin-backend/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// User 用户模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Realname  string    `gorm:"size:50" json:"realname"`
	Email     string    `gorm:"size:100" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Status    int       `gorm:"default:1" json:"status"` // 1:正常 0:禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联关系
	Role   *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	RoleID *uint `json:"role_id"`
}

// Role 角色模型
type Role struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}

// Permission 权限模型
type Permission struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	ParentCode  string    `gorm:"size:50" json:"parent_code"` // 父权限代码，用于构建权限树
	Path        string    `gorm:"size:100" json:"path"`       // 路由路径
	Type        int       `gorm:"default:1" json:"type"`      // 1:菜单 2:功能 3:按钮
	Sort        int       `gorm:"default:0" json:"sort"`      // 排序
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Roles []Role `gorm:"many2many:role_permissions" json:"roles,omitempty"`
}

// InitDB 初始化数据库
func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.GetDBPath()), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移
	if err := DB.AutoMigrate(&User{}, &Role{}, &Permission{}); err != nil {
		return err
	}

	// 初始化默认数据
	initDefaultData()

	log.Println("数据库初始化成功")
	return nil
}

// initDefaultData 初始化默认数据
func initDefaultData() {
	// 检查是否已有管理员用户
	var count int64
	DB.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	// 创建默认管理员
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	admin := User{
		Username: "admin",
		Password: string(hashedPassword),
		Realname: "管理员",
		Email:    "admin@example.com",
		Status:   1,
	}
	DB.Create(&admin)

	// 创建默认角色
	roles := []Role{
		{Name: "超级管理员", Code: "admin", Description: "系统超级管理员"},
		{Name: "普通用户", Code: "user", Description: "普通用户"},
	}
	DB.Create(&roles)

	// 创建默认权限
	permissions := []Permission{
		// 系统管理
		{Name: "系统管理", Code: "system", ParentCode: "", Path: "/system", Type: 1, Sort: 0, Description: "系统管理模块"},
		{Name: "用户管理", Code: "system:user", ParentCode: "system", Path: "/system/user", Type: 1, Sort: 1, Description: "用户管理"},
		{Name: "角色管理", Code: "system:role", ParentCode: "system", Path: "/system/role", Type: 1, Sort: 2, Description: "角色管理"},
		{Name: "权限管理", Code: "system:permission", ParentCode: "system", Path: "/system/permission", Type: 1, Sort: 3, Description: "权限管理"},

		// 用户管理功能权限
		{Name: "用户查看", Code: "system:user:view", ParentCode: "system:user", Path: "", Type: 2, Sort: 1, Description: "查看用户列表"},
		{Name: "用户新增", Code: "system:user:add", ParentCode: "system:user", Path: "", Type: 2, Sort: 2, Description: "新增用户"},
		{Name: "用户编辑", Code: "system:user:edit", ParentCode: "system:user", Path: "", Type: 2, Sort: 3, Description: "编辑用户"},
		{Name: "用户删除", Code: "system:user:delete", ParentCode: "system:user", Path: "", Type: 2, Sort: 4, Description: "删除用户"},

		// 角色管理功能权限
		{Name: "角色查看", Code: "system:role:view", ParentCode: "system:role", Path: "", Type: 2, Sort: 1, Description: "查看角色列表"},
		{Name: "角色新增", Code: "system:role:add", ParentCode: "system:role", Path: "", Type: 2, Sort: 2, Description: "新增角色"},
		{Name: "角色编辑", Code: "system:role:edit", ParentCode: "system:role", Path: "", Type: 2, Sort: 3, Description: "编辑角色"},
		{Name: "角色删除", Code: "system:role:delete", ParentCode: "system:role", Path: "", Type: 2, Sort: 4, Description: "删除角色"},

		// 权限管理功能权限
		{Name: "权限查看", Code: "system:permission:view", ParentCode: "system:permission", Path: "", Type: 2, Sort: 1, Description: "查看权限列表"},
		{Name: "权限分配", Code: "system:permission:assign", ParentCode: "system:permission", Path: "", Type: 2, Sort: 2, Description: "分配权限"},

		// 仪表盘
		{Name: "仪表盘", Code: "dashboard", ParentCode: "", Path: "/dashboard", Type: 1, Sort: 0, Description: "仪表盘模块"},
	}
	DB.Create(&permissions)

	// 为超级管理员分配所有权限
	var adminRole Role
	if err := DB.Where("code = ?", "admin").First(&adminRole).Error; err == nil {
		DB.Model(&adminRole).Association("Permissions").Append(permissions)
	}

	log.Println("默认数据初始化成功")
}
