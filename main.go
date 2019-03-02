package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// preInit function reads configuration and initializes logging.
func preInit() {
	viper.SetConfigName("config")
	viper.SetConfigType("TOML")
	if runtime.GOOS == "linux" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0])) // XXX: Unreliable
		if err != nil {                                    // XXX: unexpected
			panic(err)
		}
		if strings.Contains(dir, "bin") {
			viper.SetDefault("Log.file", "/var/log/mqtt/broker.log")
			viper.AddConfigPath("/etc/gomqtt/")
		} else {
			viper.AddConfigPath(".")
		}
		// TODO: viper.SetDefault("Log.system", false) // system logging
	} else {
		viper.AddConfigPath(".")
		viper.SetDefault("Log.file", "./broker.log")
	}
	viper.SetDefault("Log.UTC", true)
	err := viper.ReadInConfig()
	if err != nil {
		panic(errors.New("failed to read config file: " + err.Error()))
	}

	logFilePath := viper.GetString("Log.file")
	finalFlag := log.Ldate | log.Ltime
	if viper.GetBool("Log.UTC") {
		finalFlag |= log.LUTC
	}
	if viper.GetBool("debug") {
		finalFlag |= log.Lshortfile
	}
	log.SetFlags(finalFlag)

	if strings.ContainsAny(logFilePath, "/\\") {
		logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			panic(errors.New("failed to open log file: " + err.Error()))
		}
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	} // else path is probably invalid. Log to stdout only.
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Initializing logging and configuration manager.
	preInit()

	log.Println("GoMQTT Broker")
	log.Println("Copyright © 2019 Vladyslav Yamkovyi (Hexawolf)")
}
