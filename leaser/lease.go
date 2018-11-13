package leaser

import (
	"github.com/pauldotknopf/automounter/providers"
)

type mediaLease struct {
	mediaID string
	providers.MountSession
	leases []*mediaLeaseItem
}

type mediaLeaseItem struct {
	leaseID string
	media   *mediaLease
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
