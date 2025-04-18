package config

type Application struct {
	Mode         string
	Host         string
	Name         string
	Port         string
	ReadTimeout  int
	WriteTimeout int
	EnableDp     bool
}

var ApplicationConfig = new(Application)
