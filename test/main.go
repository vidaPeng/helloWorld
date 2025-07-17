package main

import (
	"fmt"
	"os"
)

var pidPath = "data/pid"

func main() {
	pid := os.Getpid()
	fmt.Println(pid)
	err := os.WriteFile(pidPath, []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		fmt.Println("Error writing to .pid:", err)
		return
	}
}
