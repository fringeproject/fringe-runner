package session

import (
	"fmt"

	"github.com/fringeproject/fringe-runner/common"
)

type ModuleList []common.ModuleInterface
type MiddlewareList []common.MiddlewareInterface

type Session struct {
	Modules     ModuleList
	Middlewares MiddlewareList
}

func NewSession() (*Session, error) {
	sess := &Session{
		Modules: make([]common.ModuleInterface, 0),
	}

	return sess, nil
}

func (s *Session) Close() {
}

func (s *Session) RegisterModule(mod common.ModuleInterface) {
	s.Modules = append(s.Modules, mod)
}

func (s *Session) Module(name string) (mod common.ModuleInterface, err error) {
	for _, m := range s.Modules {
		if m.Slug() == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("Module %s not found.", name)
}

func (s *Session) GetModules() ([]common.Module, error) {
	modules := make([]common.Module, 0)

	for _, mod := range s.Modules {
		m := common.Module{
			Name:        mod.Name(),
			Slug:        mod.Slug(),
			Description: mod.Description(),
		}
		modules = append(modules, m)
	}

	return modules, nil
}

func (s *Session) RegisterMiddleware(mod common.MiddlewareInterface) {
	s.Middlewares = append(s.Middlewares, mod)
}
