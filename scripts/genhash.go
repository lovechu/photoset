package main

import (
	"fmt"
	"photoset/internal/pkg/password"
)

func main() {
	hash, err := password.Hash("admin123")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(hash)
}
