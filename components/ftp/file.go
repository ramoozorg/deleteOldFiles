package ftp

import (
	"github.com/jlaffaye/ftp"
	"time"
)

type Entry struct {
	Name       string
	Path       string
	Size       uint64
	Type       ftp.EntryType
	ModifyTime time.Time
}
