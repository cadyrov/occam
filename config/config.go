package config

type Config struct {
	Project Project `mapstructure:"project"`
}

type Project struct {
	Name  string `mapstructure:"name"`
	Log   Log    `mapstructure:"log"`
	Shift int    `mapstructure:"shift"`
}

type Log struct {
	Level string `mapstructure:"level"`
}
