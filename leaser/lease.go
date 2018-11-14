package leaser

import (
	"time"

	"github.com/pauldotknopf/automounter/providers"
)

type mediaLease struct {
	mediaID string
	providers.MountSession
	leases []*mediaLeaseItem
	// When the last lease is closed,
	// we want to clean up this mount session
	// after a period of time
	lastClosedTime time.Time
}

type mediaLeaseItem struct {
	leaseID     string
	mediaItemID string
	media       *mediaLease
}

func (s *mediaLeaseItem) ID() string {
	return s.leaseID
}

func (s *mediaLeaseItem) MediaID() string {
	return s.media.mediaID
}

func (s *mediaLeaseItem) MountPath() string {
	return s.media.Location()
}

func (s *mediaLeaseItem) IsValid() bool {
	return s.media != nil
}
