package services

import "torrentRenamer"

type Service interface {
	Name() string
	Search(torrentRenamer.Video) (torrentRenamer.Video, error)
	IsAvailable() bool
	GetNewName(torrentRenamer.Video) (string, error)
	IsDefault() bool
}
