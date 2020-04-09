package modules

import (
	"github.com/fringeproject/fringe-runner/modules/awss3"
	"github.com/fringeproject/fringe-runner/modules/crtsh"
	"github.com/fringeproject/fringe-runner/modules/helloworld"

	"github.com/fringeproject/fringe-runner/session"
)

func LoadModules(sess *session.Session) {
	sess.RegisterModule(awss3.NewAWSS3())
	sess.RegisterModule(crtsh.NewCrtsh())
	sess.RegisterModule(helloworld.NewHelloWorld())
}
