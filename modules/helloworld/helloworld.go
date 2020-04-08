package helloworld

import ()

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

func (m *HelloWorld) Run() error {
	return nil
}
