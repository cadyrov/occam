package config

type Config struct {
	Project Project `mapstructure:"project"`
}

type Project struct {
	Name            string `mapstructure:"name"`
	Log             Log    `mapstructure:"log"`
	Shift           int    `mapstructure:"shift"`
	Output          string `mapstructure:"output"`
	PrecisionSecond int64  `mapstructure:"output"`
}

type Log struct {
	Level  string `mapstructure:"level"`
	Output string `mapstructure:"output"`
}
