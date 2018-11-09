package udevil

import (
	"github.com/pauldotknopf/automounter/providers"
)

func init() {
	providers.AddProvider(&udevilProvider{})
}

type udevilProvider struct {
}

func (p *udevilProvider) Name() string {
	return "udevil"
}
