package config

import (
	"github.com/spf13/viper"
	"os"
)

type ProviderConf struct {
	Name     string   `mapstructure:"name"`
	Type     string   `mapstructure:"type"`
	EndPoint string   `mapstructure:"end_point"`
	APIKeys  []string `mapstructure:"api_keys"`
	Weight   int      `mapstructure:"weight"`
	Models   []string `mapstructure:"models"`
}

type Config struct {
	Env string `mapstructure:"env"`
	App struct {
		Name    string `mapstructure:"name"`
		Version string `mapstructure:"version"`
	} `mapstructure:"app"`
	HTTP struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"http"`
	Database struct {
		Driver string `mapstructure:"driver"`
		Dsn    string `mapstructure:"dsn"`
	} `mapstructure:"database"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	Security struct {
		ApiSign struct {
			AppKey      string `mapstructure:"app_key"`
			AppSecurity string `mapstructure:"app_security"`
		} `mapstructure:"api_sign"`
		Jwt struct {
			Key string `mapstructure:"key"`
		} `mapstructure:"jwt"`
	} `mapstructure:"security"`
	Log struct {
		Level         string `mapstructure:"log_level"`
		Encoding      string `mapstructure:"encoding"`
		LogPath       string `mapstructure:"log_path"`
		ErrorFileName string `mapstructure:"error_file_name"`
		FileName      string `mapstructure:"log_file_name"`
		MaxBackups    int    `mapstructure:"max_backups"`
		MaxAge        int    `mapstructure:"max_age"`
		MaxSize       int    `mapstructure:"max_size"`
		Compress      bool   `mapstructure:"compress"`
	} `mapstructure:"log"`
	ModelMapping        map[string][]string `mapstructure:"model_mapping"`
	ChatCOmpletionCheck []string            `mapstructure:"chat_completion_check"`
	Providers           []ProviderConf      `mapstructure:"providers"`
}

func LoadConfig(path string) *Config {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = path
	}
	conf := viper.New()
	conf.SetConfigFile(envConf)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var cfg Config
	err = conf.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
