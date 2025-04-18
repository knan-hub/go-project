package source

import "time"

type ChangeSet struct {
	Data      []byte
	CheckSum  string
	Format    string
	Source    string
	Timestamp time.Time
}

type Watcher interface {
	Next() (*ChangeSet, error)
	Stop() error
}

type Source interface {
	Read() (*ChangeSet, error)
	Write(*ChangeSet) error
	Watch() (Watcher, error)
	String() string
}
