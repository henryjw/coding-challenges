package main

import (
	"fmt"
	"memcached-server/cache"
)

func main() {
	c := cache.New()
	fmt.Println("Hello, World!")
}
