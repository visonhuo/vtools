package main

import (
	"math/rand"
	"time"

	"github.com/vtools/app/vecho/cmd"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	cmd.Execute()
}
