package service

import (
	"fmt"
	"geek/internal/dbeye/dao"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

type GetUserReq struct {
	UserName *string `form:"userName"`
	Page     *int    `form:"page"`
	PageSize *int    `form:"pageSize"`
}

func (us *UserService) GetUserInfo(req GetUserReq) (error, int, map[string]interface{}) {
	userDao := dao.NewUserDao()
	err, code, data := userDao.GetUserList(req.UserName, req.Page, req.PageSize)
	fmt.Println(err, code, data)
	return err, code, data
}
