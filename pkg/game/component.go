package game

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrComponentNotFound = errors.New("component not found in component state")
)

type Component interface {
	Name() string
	Update(*Game) error
	Draw(*ebiten.Image)
}

type ComponentState struct {
	order      []string
	components map[string]Component
}

func NewComponentState() *ComponentState {
	return &ComponentState{
		order:      []string{},
		components: map[string]Component{},
	}
}

func (s *ComponentState) RegisterComponent(c Component) {
	name := c.Name()
	s.components[name] = c
	s.order = append(s.order, name)
}

func (s *ComponentState) UpdateAll(g *Game) error {
	for _, name := range s.order {
		component, ok := s.components[name]
		if !ok {
			continue
		}
		err := component.Update(g)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ComponentState) DrawAll(screen *ebiten.Image) {
	for _, name := range s.order {
		component, ok := s.components[name]
		if !ok {
			continue
		}
		component.Draw(screen)
	}
}

func (s *ComponentState) GetComponent(name string) (Component, error) {
	c, ok := s.components[name]
	if !ok {
		return nil, fmt.Errorf("component %s not found", name)
	}
	return c, nil
}
