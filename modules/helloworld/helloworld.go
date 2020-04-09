package helloworld

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/common/assets"
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
	asset, err := ctx.GetAssetAsRawString()
	if err != nil {
		logrus.Debug("Cannot get asset as a raw string")
		logrus.Debug(err)
		return err
	}

	err = ctx.CreateNewAsset(fmt.Sprintf("Hello, %s!", asset), assets.AssetTypes["raw"])

	return err
}
