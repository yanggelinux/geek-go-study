package test

import (
	"geek-go-study/week08/pkg"
	"log"
	"testing"
)

func BenchmarkRedisOpt10k(b *testing.B) {
	rds := pkg.NewRedisHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := rds.Set("k10", pkg.K10)
		if err != nil {
			log.Println("set error", err)
		}
		rds.Get("k10")
	}

}

func BenchmarkRedisOpt20k(b *testing.B) {
	rds := pkg.NewRedisHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := rds.Set("k20", pkg.K20)
		if err != nil {
			log.Println("set error", err)
		}
		rds.Get("k20")
	}

}

func BenchmarkRedisOpt50k(b *testing.B) {
	rds := pkg.NewRedisHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := rds.Set("k50", pkg.K50)
		if err != nil {
			log.Println("set error", err)
		}
		rds.Get("k50")
	}

}

func BenchmarkRedisOpt100k(b *testing.B) {
	rds := pkg.NewRedisHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := rds.Set("k100", pkg.K100)
		if err != nil {
			log.Println("set error", err)
		}
		rds.Get("k100")
	}

}

func BenchmarkRedisOpt200k(b *testing.B) {
	rds := pkg.NewRedisHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := rds.Set("k200", pkg.K200)
		if err != nil {
			log.Println("set error", err)
		}
		rds.Get("k200")
	}

}

func BenchmarkRedisOpt1000k(b *testing.B) {
	rds := pkg.NewRedisHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := rds.Set("k1000", pkg.K1000)
		if err != nil {
			log.Println("set error", err)
		}
		rds.Get("k1000")
	}
}

func BenchmarkRedisOpt5000k(b *testing.B) {
	rds := pkg.NewRedisHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := rds.Set("k5000", pkg.K5000)
		if err != nil {
			log.Println("set error", err)
		}
		rds.Get("k5000")
	}

}

//go test redis_size_test.go -bench=. -benchmem

//benchmark结果

//goos: darwin
//goarch: amd64
//cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
//BenchmarkRedisOpt10k-8     	      25	  50390763 ns/op	    3093 B/op	      74 allocs/op
//BenchmarkRedisOpt20k-8     	      30	  45369373 ns/op	    2696 B/op	      65 allocs/op
//BenchmarkRedisOpt50k-8     	      25	  49682093 ns/op	    3137 B/op	      74 allocs/op
//BenchmarkRedisOpt100k-8    	      28	  37599046 ns/op	    2573 B/op	      67 allocs/op
//BenchmarkRedisOpt200k-8    	      22	  51187296 ns/op	    3655 B/op	      81 allocs/op
//BenchmarkRedisOpt1000k-8   	      27	  38980763 ns/op	    3941 B/op	      70 allocs/op
//BenchmarkRedisOpt5000k-8   	      32	  52023704 ns/op	    7594 B/op	      61 allocs/op
//PASS
//ok  	command-line-arguments	17.583s
//可以看出val不同字节大小 性能影响不大
