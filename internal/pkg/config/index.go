package config

import (
	"adorable-star/internal/pkg/util"
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config = &Configuration{}

func Init() {
	// Init viper
	v := viper.New()
	v.SetConfigFile(util.GetCwd() + "/config/config.yaml")
	v.SetConfigType("yaml")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	// Watch changes for config file and enable hot reload
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		// Reload config
		if err := v.Unmarshal(&Config); err != nil {
			log.Printf("ERR failed to hot reload config file: %s\n", err.Error())
		}
	})

	// Initialize global var for config
	if err := v.Unmarshal(&Config); err != nil {
		log.Printf("ERR failed to load config file: %s\n", err.Error())
	}
}
