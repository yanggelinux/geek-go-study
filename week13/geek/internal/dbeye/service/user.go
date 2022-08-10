package service

import (
	"encoding/json"
	"geek/internal/dbeye/dao"
	"geek/internal/dbeye/pkg/ce"
	"geek/internal/dbeye/stroe/redis"
	"geek/pkg/log"
	"go.uber.org/zap"
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

type SyncUserReq struct {
	UserName *string `form:"userName"`
}

func (us *UserService) GetUserInfo(req GetUserReq) (error, int, map[string]interface{}) {
	userDao := dao.NewUserDao()
	err, code, data := userDao.GetUserList(req.UserName, req.Page, req.PageSize)
	return err, code, data
}

func (us *UserService) SyncUserInfo(req SyncUserReq) (error, int) {
	userDao := dao.NewUserDao()
	err, code, data := userDao.GetUserList(req.UserName, nil, nil)
	redisHandler := redis.NewRedisHandler()

	defer func() {
		err := redisHandler.CloseClient()
		if err != nil {
			log.Logger.Error("close redis error:", zap.Error(err))
		}
	}()
	sByte, err := json.Marshal(data)
	if err != nil {
		return err, ce.ERROR
	}
	key := "user:info"
	err = redisHandler.Set(key, sByte)
	if err != nil {
		return err, ce.ERROR
	}
	return err, code

}
