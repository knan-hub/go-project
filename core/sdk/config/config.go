package config

import (
	"go-project/core/config/source"
	"log"
)

type Config struct {
	Application *Application `yaml:"application"`
	Logger      *Logger      `yaml:"logger"`
	Jwt         *Jwt         `yaml:"jwt"`
	Database    *Database    `yaml:"database"`
}

type Settings struct {
	Settings  Config `yaml:"settings"`
	callbacks []func()
}

func (s *Settings) runCallback() {
	for _, cb := range s.callbacks {
		cb()
	}
}

func (s *Settings) init() {
	s.Settings.Logger.Setup()
	s.runCallback()
}

func (s *Settings) OnChange() {
	s.init()
	log.Println("config change and reload")
}

func (s *Settings) Init() {
	s.init()
	log.Println("config init")
}

var _cfg *Settings

func setup(s source.Source, fs ...func()) {
	_cfg = &Settings{
		Settings: Config{
			Application: ApplicationConfig,
			Logger:      LoggerConfig,
			Jwt:         JwtConfig,
			Database:    DatabaseConfig,
		},
		callbacks: fs,
	}

	var err error

}
