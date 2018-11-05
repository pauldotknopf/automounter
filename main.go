package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("sdf")
	exec.Command("udevil", "--monitor")
}
