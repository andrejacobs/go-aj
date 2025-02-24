package main

import (
	"fmt"

	"github.com/andrejacobs/go-micropkg/buildinfo"
)

func main() {
	fmt.Println(buildinfo.UsageNameAndVersion())
}
