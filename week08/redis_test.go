package test

import (
	"log"
	"testing"
)

func BenchmarkRedisOpt10k(b *testing.B) {
	rds := NewRedisHandler()
	err := rds.Set("k10", k10)
	if err != nil {
		log.Println("set error", err)
	}
	rds.Get("k10")
}

func BenchmarkRedisOpt20k(b *testing.B) {
	rds := NewRedisHandler()
	err := rds.Set("k20", k20)
	if err != nil {
		log.Println("set error", err)
	}
	rds.Get("k20")

}

func BenchmarkRedisOpt50k(b *testing.B) {
	rds := NewRedisHandler()
	err := rds.Set("k50", k50)
	if err != nil {
		log.Println("set error", err)
	}
	rds.Get("k50")

}

func BenchmarkRedisOpt100k(b *testing.B) {
	rds := NewRedisHandler()
	err := rds.Set("k100", k100)
	if err != nil {
		log.Println("set error", err)
	}
	rds.Get("k100")

}

func BenchmarkRedisOpt200k(b *testing.B) {
	rds := NewRedisHandler()
	err := rds.Set("k200", k200)
	if err != nil {
		log.Println("set error", err)
	}
	rds.Get("k200")
}

func BenchmarkRedisOpt1000k(b *testing.B) {
	rds := NewRedisHandler()
	err := rds.Set("k1000", k1000)
	if err != nil {
		log.Println("set error", err)
	}
	rds.Get("k1000")
}

func BenchmarkRedisOpt5000k(b *testing.B) {
	rds := NewRedisHandler()
	err := rds.Set("k5000", k5000)
	if err != nil {
		log.Println("set error", err)
	}
	rds.Get("k5000")
}

//benchmark结果
//goos: darwin
//goarch: amd64
//pkg: geek-go-study/week08
//cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
//BenchmarkRedisOpt10k-8     	1000000000	         0.05039 ns/op
//BenchmarkRedisOpt20k-8     	1000000000	         0.05855 ns/op
//BenchmarkRedisOpt50k-8     	1000000000	         0.05466 ns/op
//BenchmarkRedisOpt100k-8    	1000000000	         0.03884 ns/op
//BenchmarkRedisOpt200k-8    	1000000000	         0.05049 ns/op
//BenchmarkRedisOpt1000k-8   	1000000000	         0.05510 ns/op
//BenchmarkRedisOpt5000k-8   	1000000000	         0.04186 ns/op
//PASS
//ok  	geek-go-study/week08	3.309s

//可以看出val不同字节大小 性能影响不大
