package api

import (
	"context"
	"fmt"
)

type User struct {
	api *API
}

func NewUser() *User {
	api := NewAPI("http://127.0.0.1:8973")
	return &User{api: api}
}

func (u *User) GetUser(ctx context.Context, userName string) ([]byte, error) {

	body, err := u.api.httpGet(ctx, fmt.Sprintf("%s?userName=%s", "api/v1/user", userName))
	if err != nil {
		return nil, err
	}

	return body, nil
}
