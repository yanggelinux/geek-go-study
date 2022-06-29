package test

import (
	"fmt"
	"geek-go-study/week08/pkg"
	"log"
	"testing"
)

func TestRedisOpt(t *testing.T) {
	rds := pkg.NewRedisHandler()

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("k%d", i)
		err := rds.Set(key, pkg.K10)
		if err != nil {
			log.Println("set error", err)
		}
	}
}
