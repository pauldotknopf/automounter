package muxer

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/pauldotknopf/automounter/providers"
)

type muxer struct {
	p []providers.MediaProvider
}

// Create a muxer from multiple providers
func Create(p ...providers.MediaProvider) providers.MediaProvider {
	return &muxer{
		p,
	}
}

func (s *muxer) Name() string {
	return "muxer"
}

func (s *muxer) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var eg errgroup.Group
	for _, provider := range s.p {
		p := provider
		eg.Go(func() error {
			fmt.Println("starting " + p.Name())
			err := p.Start(ctx)
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

func (s *muxer) Mount(id string) (providers.MountSession, error) {
	for _, provider := range s.p {
		session, err := provider.Mount(id)
		if err == providers.ErrIDNotFound {
			continue
		}
		if err == nil {
			return session, nil
		}
		return nil, err
	}
	return nil, providers.ErrIDNotFound
}

func (s *muxer) Unmount(id string) error {
	for _, provider := range s.p {
		err := provider.Unmount(id)
		if err == providers.ErrIDNotFound {
			continue
		}
		return err
	}
	return providers.ErrIDNotFound
}

func (s *muxer) MediaAddded() (<-chan providers.Media, func()) {
	out := make(chan providers.Media)

	type pairStruct struct {
		in   <-chan providers.Media
		canc func()
	}
	pairs := make([]pairStruct, 0)

	var wg sync.WaitGroup

	cancel := func() {
		for _, p := range pairs {
			p.canc()
		}
		wg.Wait()
		close(out)
	}

	for _, p := range s.p {
		mediaAddedChannel, mediaAddedCancel := p.MediaAddded()
		pairs = append(pairs, pairStruct{mediaAddedChannel, mediaAddedCancel})
	}

	// Start goroutines for all the inbound channels
	// to send them to the single outbound
	for _, p := range pairs {
		wg.Add(1)
		go func(pair pairStruct) {
			defer wg.Done()
			for m := range pair.in {
				out <- m
			}
		}(p)
	}

	return out, cancel
}

func (s *muxer) MediaRemoved() (<-chan string, func()) {
	out := make(chan string)

	type pairStruct struct {
		in   <-chan string
		canc func()
	}
	pairs := make([]pairStruct, 0)

	var wg sync.WaitGroup

	cancel := func() {
		for _, p := range pairs {
			p.canc()
		}
		wg.Wait()
		close(out)
	}

	for _, p := range s.p {
		mediaRemovedChannel, mediaRemovedCancel := p.MediaRemoved()
		pairs = append(pairs, pairStruct{mediaRemovedChannel, mediaRemovedCancel})
	}

	// Start goroutines for all the inbound channels
	// to send them to the single outbound
	for _, p := range pairs {
		wg.Add(1)
		go func(pair pairStruct) {
			defer wg.Done()
			for m := range pair.in {
				out <- m
			}
		}(p)
	}

	return out, cancel
}
