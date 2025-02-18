package main

import (
	"fmt"
	"grpc-sevice/internal/config"
)

func main() {
	config := config.MustLoad()
	fmt.Println(config)
}
