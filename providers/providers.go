package providers

import (
	"fmt"
)

var providers = struct {
	p []MediaProvider
}{}

// MediaProvider The type that will detect and mount media
type MediaProvider interface {
	Name() string
}

// AddProvider Add a provider
func AddProvider(provider MediaProvider) {
	providers.p = append(providers.p, provider)
}

func Test() {
	for _, provider := range providers.p {
		fmt.Println(provider.Name())
	}
}
