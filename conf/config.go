package conf

import "github.com/spf13/viper"

type Config struct {
	Gin    *Gin
	DB     *DB
	Logger *Logger
}

type DB struct {
	DSN     string
	MaxIdle int
	MaxOpen int
}

type Gin struct {
	Mode string
	Port string
}

type Logger struct {
	Level      string
	FileLog    bool
	Path       string
	ConsoleLog bool
}

var GlobalConfig *Config

func InitConfig(path string) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		panic(err)
	}
}
