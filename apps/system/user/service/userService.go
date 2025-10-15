package service

import (
	"github.com/GoFurry/gofurry-nav-backend/apps/system/user/models"
	"github.com/GoFurry/gofurry-nav-backend/common"
)

type userService struct{}

var userSingleton = new(userService)

func GetUserService() *userService { return userSingleton }

// 用户登录
func (svc *userService) Login(req models.UserLoginRequest) (tokenStr string, err common.GFError) {
	return
}
