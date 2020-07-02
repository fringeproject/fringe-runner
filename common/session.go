package common

import (
	"fmt"
)

type ModuleList []ModuleInterface

type Session struct {
	Modules ModuleList
}

func NewSession() (*Session, error) {
	sess := &Session{
		Modules: make([]ModuleInterface, 0),
	}

	return sess, nil
}

func (s *Session) Close() {
}

func (s *Session) RegisterModule(mod ModuleInterface) {
	s.Modules = append(s.Modules, mod)
}

func (s *Session) Module(name string) (mod ModuleInterface, err error) {
	for _, m := range s.Modules {
		if m.Slug() == name {
			return m, nil
		}
	}

	return nil, fmt.Errorf("Module %s not found.", name)
}

func (s *Session) GetModules() ([]Module, error) {
	modules := make([]Module, 0)

	for _, mod := range s.Modules {
		m := Module{
			Name:        mod.Name(),
			Slug:        mod.Slug(),
			Description: mod.Description(),
		}
		modules = append(modules, m)
	}

	return modules, nil
}
