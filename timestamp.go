package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().UTC()

	timestamp := now.Format(time.RFC3339)

	fmt.Println(timestamp)
}
