package service

import (
	"context"
	"encoding/json"
	"fmt"
	"geek/internal/cloudeye/api"
	pb "geek/proto"
)

type UserRes struct {
	ID          uint32 `json:"id"`
	UserName    string `json:"userName"`
	FullName    string `json:"fullName"`
	UpdatedTime string `json:"updatedTime"`
	CreatedTime string `json:"createdTime"`
}

type UserData struct {
	ResList []UserRes `json:"resList"`
}

type UserResp struct {
	Data   UserData `json:"data"`
	Msg    string   `json:"msg"`
	Status int      `json:"status"`
}

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

//grpcurl -plaintext -d '{"userName":"admin"}' localhost:8001 proto.UserService.GetUserList
func (t *UserService) GetUserList(ctx context.Context, r *pb.GetUserListRequest) (*pb.GetUserListReply, error) {
	u := api.NewUser()
	body, err := u.GetUser(ctx, r.GetUserName())
	fmt.Println("bbbbbb", string(body), err)
	if err != nil {
		return nil, err
	}

	userList := pb.GetUserListReply{}
	err = json.Unmarshal(body, &userList)
	fmt.Println("UserResp", userList)
	if err != nil {
		return nil, err
	}

	return &userList, nil
}
