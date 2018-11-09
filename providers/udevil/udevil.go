package udevil

import (
	"github.com/pauldotknopf/automounter/providers"
)

func init() {
	providers.AddProvider(&udevil{})
}

type udevil struct {
}

func (p *udevil) Name() string {
	return "udevil"
}

func (p *udevil) Start() error {
	return nil
}

func (p *udevil) Stop() error {
	return nil
}

func (p *udevil) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	result = append(result, udevilMedia{})
	return result
}
