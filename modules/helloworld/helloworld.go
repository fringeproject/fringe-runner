package helloworld

import (
	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type HelloWorld struct {
}

func NewHelloWorld() *HelloWorld {
	mod := &HelloWorld{}

	return mod
}

func (m *HelloWorld) Name() string {
	return "Hello World"
}

func (m *HelloWorld) Slug() string {
	return "helloworld"
}

func (m *HelloWorld) Description() string {
	return "This module output Hello World."
}

func (m *HelloWorld) Run(ctx *common.ModuleContext) error {
	asset := ctx.Asset
	logrus.Infof("Hello, %s!", asset)
	return nil
}
