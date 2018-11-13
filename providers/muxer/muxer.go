package muxer

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/pauldotknopf/automounter/providers"
)

type muxer struct {
	p []providers.MediaProvider
}

// Create a muxer from multiple providers
func Create(p []providers.MediaProvider) providers.MediaProvider {
	return &muxer{
		p,
	}
}

func (s *muxer) Initialize() error {
	for _, provider := range s.p {
		err := provider.Initialize()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *muxer) Name() string {
	return "muxer"
}

func (s *muxer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var eg errgroup.Group
	for _, provider := range s.p {
		eg.Go(func() error {
			err := provider.Start(ctx)
			if err != nil {
				cancel()
				return err
			}
			return nil
		})
	}

	return eg.Wait()
}

func (s *muxer) GetMedia() []providers.Media {
	result := make([]providers.Media, 0)
	for _, provider := range s.p {
		result = append(result, provider.GetMedia()...)
	}
	return result
}

func (s *muxer) Mount(media providers.Media) (providers.MountSession, error) {
	return nil, fmt.Errorf("not implemented")
}
