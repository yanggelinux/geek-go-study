package service

import (
	"context"
	"encoding/json"
	"geek/internal/cloudeye/api"
	pb "geek/proto"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

//grpcurl -plaintext -d '{"userName":"admin"}' localhost:8001 proto.UserService.GetUserList

func (t *UserService) GetUserList(ctx context.Context, r *pb.GetUserListRequest) (*pb.GetUserListReply, error) {
	u := api.NewUser()
	body, err := u.GetUser(ctx, r.GetUserName())
	if err != nil {
		return nil, err
	}
	userList := pb.GetUserListReply{}
	err = json.Unmarshal(body, &userList)
	if err != nil {
		return nil, err
	}
	return &userList, nil
}
