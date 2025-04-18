package config

import (
	"context"
	"go-project/core/config/source"

	"github.com/bytedance/sonic/loader"
	"github.com/go-admin-team/go-admin-core/config/reader"
)

type Entity interface {
	OnChange()
}

type Options struct {
	Loader  loader.Loader
	Reader  reader.Reader
	Source  []source.Source
	Context context.Context
	Entity  Entity
}

type Option func(*Options)

type Watcher interface {
	Next() (reader.Value, error)
	Stop() error
}

type Config interface {
	reader.Values
	Init(opts ...Option) error
	Options() Options
	Close() error
	Load(source ...source.Source) error
	Sync() error
	Watch(path ...string) (Watcher, error)
}

func NewConfig(options ...Option) (Config, error) {
	return newConfig(options...)
}
