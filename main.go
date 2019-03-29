package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/BurntSushi/toml"
)

// Server is a binding structure that binds everything related to the broker
// runtime.
type Server struct {
	Config struct {
		Debug bool
		Log   struct {
			Path string
			UTC  bool
		}
		MQTTAddress   string
		MQTTSNAddress string
	}
}

var serv Server

// LoadConfig reads config from specified file and decodes it into a generic
// structure
func (s *Server) LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := toml.DecodeReader(file, &s.Config); err != nil {
		return err
	}

	return nil
}

// basically debugging hexa way lol
func debug(s ...interface{}) {
	if serv.Config.Debug {
		log.Println(s...)
	}
}

func main() {
	fmt.Println("GoMQTT Broker")
	fmt.Println("Copyright © 2019 Vladyslav Yamkovyi (Hexawolf)")
	fmt.Println()

	err := serv.LoadConfig("broker.cfg")
	if err != nil {
		fmt.Println("open broker.cfg: The system cannot find the file specified.")
		os.Exit(1)
	}
	if serv.Config.Log.Path != "" {
		logFile, err := os.OpenFile(serv.Config.Log.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			panic(errors.New("failed to open log file: " + err.Error()))
		}
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}
	if serv.Config.Log.UTC {
		log.SetFlags(log.LstdFlags | log.LUTC)
	}
	if serv.Config.Debug {
		log.SetFlags(log.Flags() | log.Lshortfile)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	if serv.Config.MQTTSNAddress != "" {
		go ListenUDP(serv.Config.MQTTSNAddress)
	}

	if serv.Config.MQTTAddress != "" {
		log.Println("TCP listener is not implemented yet!")
	}

	// Shutdown gracefully on signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	<-c
}
