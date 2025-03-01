package main

import (
	"fmt"

	"github.com/andrejacobs/go-aj/buildinfo"
)

func main() {
	fmt.Println(buildinfo.UsageNameAndVersion())
}
