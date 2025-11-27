package services

import (
	"errors"
	"react-go-admin-backend/models"
)

// RoleService 角色服务
type RoleService struct{}

// GetRoleList 获取角色列表
func (s *RoleService) GetRoleList(page, pageSize int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	offset := (page - 1) * pageSize

	if err := models.DB.Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := models.DB.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := models.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByCode 根据代码获取角色
func (s *RoleService) GetRoleByCode(code string) (*models.Role, error) {
	var role models.Role
	if err := models.DB.Preload("Permissions").Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(role *models.Role) error {
	// 检查角色代码是否已存在
	var count int64
	models.DB.Model(&models.Role{}).Where("code = ?", role.Code).Count(&count)
	if count > 0 {
		return errors.New("角色代码已存在")
	}

	return models.DB.Create(role).Error
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint, updates map[string]interface{}) error {
	return models.DB.Model(&models.Role{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	return models.DB.Delete(&models.Role{}, id).Error
}
