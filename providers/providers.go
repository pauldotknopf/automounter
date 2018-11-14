package providers

import (
	"context"
	"errors"
)

var (
	// ErrIDNotFound An error indicating the given id wasn't found
	ErrIDNotFound = errors.New("Item not found")
)

// MediaProvider The type that will detect and mount media
type MediaProvider interface {
	Name() string
	Start(context.Context) error
	GetMedia() []Media
	Mount(id string) (MountSession, error)
	Unmount(id string) error
	MediaAddded() (<-chan Media, func())
	MediaRemoved() (<-chan string, func())
}

// MountSession represents a mount session for a media type
type MountSession interface {
	Release() error
	// Were the media is mounted
	Location() string
}

// Media A media type that can be mounted/used.
type Media interface {
	ID() string
	DisplayName() string
}
