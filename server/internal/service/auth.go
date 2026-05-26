package service

import (
	"errors"
	"time"

	"go-pure-admin/server/internal/global"
	"go-pure-admin/server/internal/model"
	"go-pure-admin/server/internal/pkg/jwtutil"
	"go-pure-admin/server/internal/util"

	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token string `json:"token"`
}

type AuthInfo struct {
	User        model.SysUser   `json:"user"`
	Roles       []model.SysRole `json:"roles"`
	Menus       []model.SysMenu `json:"menus"`
	Permissions []string        `json:"permissions"`
}

func (s AuthService) Login(req LoginReq, ip, userAgent string) (*LoginResp, error) {
	logSvc := LogService{}
	record := func(status int16, msg string, user *model.SysUser) {
		entry := &model.SysLoginLog{
			Username:  req.Username,
			Status:    status,
			Message:   msg,
			IP:        ip,
			UserAgent: userAgent,
		}
		if user != nil {
			entry.UserID = &user.ID
			entry.Nickname = user.Nickname
		}
		logSvc.RecordLogin(entry)
	}

	var user model.SysUser
	if err := global.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			record(model.LogStatusFail, "用户名或密码错误", nil)
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}
	if !user.Enabled {
		record(model.LogStatusFail, "账号已禁用", &user)
		return nil, errors.New("账号已禁用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		record(model.LogStatusFail, "用户名或密码错误", &user)
		return nil, errors.New("用户名或密码错误")
	}
	now := time.Now()
	_ = global.DB.Model(&user).Update("last_login_at", now)

	token, err := jwtutil.Sign(user.ID, user.SuperAdmin, global.Cfg.JWT.Secret, global.Cfg.JWT.ExpireHours)
	if err != nil {
		return nil, err
	}
	record(model.LogStatusSuccess, "登录成功", &user)
	return &LoginResp{Token: token}, nil
}

type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func (AuthService) loadUserProfile(userID int64) (*model.SysUser, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles").Preload("Dept").First(&user, userID).Error; err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}

func (s AuthService) GetProfile(userID int64) (*model.SysUser, error) {
	return s.loadUserProfile(userID)
}

func (s AuthService) UpdateProfile(userID int64, req UpdateProfileReq) (*model.SysUser, error) {
	fields := map[string]any{
		"nickname": req.Nickname,
		"phone":    req.Phone,
		"email":    req.Email,
		"avatar":   req.Avatar,
	}
	if err := global.DB.Model(&model.SysUser{}).Where("id = ?", userID).Updates(fields).Error; err != nil {
		return nil, err
	}
	return s.loadUserProfile(userID)
}

func (AuthService) ChangePassword(userID int64, req ChangePasswordReq) error {
	var user model.SysUser
	if err := global.DB.Select("id", "password").First(&user, userID).Error; err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("原密码错误")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return global.DB.Model(&user).Update("password", string(hash)).Error
}

func (AuthService) Info(userID int64, superAdmin bool) (*AuthInfo, error) {
	var user model.SysUser
	if err := global.DB.Preload("Roles").Preload("Dept").First(&user, userID).Error; err != nil {
		return nil, err
	}
	user.Password = ""

	menus, perms, err := loadUserMenus(userID, superAdmin)
	if err != nil {
		return nil, err
	}
	routeMenus := util.BuildMenuTree(util.FilterMenusForRoute(menus), nil)

	return &AuthInfo{
		User:        user,
		Roles:       user.Roles,
		Menus:       routeMenus,
		Permissions: perms,
	}, nil
}

func loadUserMenus(userID int64, superAdmin bool) ([]model.SysMenu, []string, error) {
	var menus []model.SysMenu
	q := global.DB.Where("enabled = ?", true).Order("sort asc, id asc")
	if superAdmin {
		if err := q.Find(&menus).Error; err != nil {
			return nil, nil, err
		}
		return menus, util.CollectPermCodes(menus), nil
	}

	var roleIDs []int64
	if err := global.DB.Model(&model.SysUserRole{}).Where("user_id = ?", userID).Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, nil, err
	}
	if len(roleIDs) == 0 {
		return nil, []string{}, nil
	}

	var menuIDs []int64
	if err := global.DB.Model(&model.SysRoleMenu{}).Where("role_id IN ?", roleIDs).Distinct().Pluck("menu_id", &menuIDs).Error; err != nil {
		return nil, nil, err
	}
	if len(menuIDs) == 0 {
		return nil, []string{}, nil
	}

	// include ancestor menus for complete tree
	allIDs := append([]int64{}, menuIDs...)
	var grantedMenus []model.SysMenu
	if err := global.DB.Where("id IN ?", menuIDs).Find(&grantedMenus).Error; err != nil {
		return nil, nil, err
	}
	for _, m := range grantedMenus {
		pid := m.ParentID
		for pid != nil {
			allIDs = append(allIDs, *pid)
			var parent model.SysMenu
			if err := global.DB.Select("id", "parent_id").First(&parent, *pid).Error; err != nil {
				break
			}
			pid = parent.ParentID
		}
	}
	allIDs = lo.Uniq(allIDs)

	if err := global.DB.Where("id IN ? AND enabled = ?", allIDs, true).Order("sort asc, id asc").Find(&menus).Error; err != nil {
		return nil, nil, err
	}
	// 仅返回 role_menus 中显式授权的按钮权限码
	grantedSet := make(map[int64]struct{}, len(menuIDs))
	for _, id := range menuIDs {
		grantedSet[id] = struct{}{}
	}
	perms := util.CollectPermCodes(lo.Filter(menus, func(m model.SysMenu, _ int) bool {
		if m.MenuType != 2 {
			return false
		}
		_, ok := grantedSet[m.ID]
		return ok
	}))

	return menus, perms, nil
}
