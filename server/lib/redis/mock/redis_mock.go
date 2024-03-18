package redismock

import "github.com/alicebob/miniredis/v2"

func MockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return s
}
