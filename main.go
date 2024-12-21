package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type sout string

func (s sout) Write(p []byte) (n int, err error) {
	fmt.Println("hook")
	return 0, nil
}

func main() {
	s := sout("hello")
	logger := logrus.New()
	logger.Out = s
	logger.Info("hello")
	logger.Info("hello2")
	logger.Info("hello3")
}
