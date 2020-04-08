package modules

import (
	"github.com/fringeproject/fringe-runner/modules/helloworld"

	"github.com/fringeproject/fringe-runner/session"
)

func LoadModules(sess *session.Session) {
	sess.RegisterModule(helloworld.NewHelloWorld())
}
