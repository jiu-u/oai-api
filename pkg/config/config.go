package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type ProviderConf struct {
	Name     string   `mapstructure:"name" yaml:"name"`
	Type     string   `mapstructure:"type" yaml:"type"`
	EndPoint string   `mapstructure:"end_point" yaml:"end_point"`
	APIKeys  []string `mapstructure:"api_keys" yaml:"api_keys"`
	Weight   int      `mapstructure:"weight" yaml:"weight"`
	Models   []string `mapstructure:"models" yaml:"models"`
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
	//ModelMapping        map[string][]string `mapstructure:"model_mapping"`
	//ChatCompletionCheck []string            `mapstructure:"chat_completion_check"`
	//Providers           []ProviderConf      `mapstructure:"providers"`
}

const prefix = "OAI"

func LoadConfig(path string) *Config {
	fmt.Println("LoadEnv", os.Getenv(prefix+"_OAUTH_LINUX_DO_CLIENT_ID"))
	fmt.Println("LoadEnv", os.Getenv(prefix+"_OAUTH_LINUX_DO_CLIENT_SECRET"))
	//os.Setenv(prefix+"_OAUTH_LINUX_DO_CLIENT_ID", "123131213131213")
	//os.Setenv(prefix+"_OAUTH_LINUX_DO_CLIENT_SECRET", "123131213131213")
	envConf := os.Getenv(prefix + "_APP_CONF")
	if envConf == "" {
		envConf = path
	}
	fmt.Printf("envConf: %s\n", envConf)
	conf := viper.New()
	// 设置环境变量前缀
	conf.SetEnvPrefix("OAI")
	// 设置环境变量的分隔符
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	conf.BindEnv("oauth.linux_do.client_id", "OAI_OAUTH_LINUX_DO_CLIENT_ID")
	conf.BindEnv("oauth.linux_do.client_secret", "OAI_OAUTH_LINUX_DO_CLIENT_SECRET")
	conf.SetConfigFile(envConf)
	conf.AutomaticEnv()
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = conf.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	//loadDataYaml(&cfg)
	return &cfg
}

//func loadDataYaml(cfg *Config) {
//	type List struct {
//		ModelMapping        map[string][]string `yaml:"model_mapping"`
//		ChatCompletionCheck []string            `yaml:"chat_completion_check"`
//		Providers           []ProviderConf      `yaml:"providers"`
//	}
//	_, err := os.Stat("./data/conf/data.yaml")
//	if err != nil {
//		return
//	}
//	file, err := os.OpenFile("./data/conf/data.yaml", os.O_RDONLY, 0644)
//	if err != nil {
//		return
//	}
//	defer file.Close()
//	var data List
//	err = yaml.NewDecoder(file).Decode(&data)
//	if err != nil {
//		return
//	}
//	cfg.ModelMapping = data.ModelMapping
//	cfg.ChatCompletionCheck = data.ChatCompletionCheck
//	cfg.Providers = data.Providers
//}
