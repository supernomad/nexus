package common

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

// Config object for base level options
type Config struct {
	ConfigFile string
	PidFile    string
	NumWorkers int
	ListenIP   string
	ListenPort int

	RollingRestart string
	ReuseFDS       bool

	ListenAddress syscall.Sockaddr

	notSet map[string]bool
}

func (cfg *Config) stringVar(name, def string) string {
	env := EnvironmentVariablePrefix + strings.ToUpper(strings.Replace(name, "-", "_", 10))
	out := os.Getenv(env)
	if out == "" {
		cfg.notSet[name] = true
		return def
	}
	return out
}

func (cfg *Config) intVar(name string, def int) int {
	defStr := strconv.Itoa(def)
	out, err := strconv.Atoi(cfg.stringVar(name, defStr))
	if err != nil {
		panic(err)
	}
	return out
}

func (cfg *Config) durationVar(name string, def time.Duration) time.Duration {
	defStr := def.String()
	out, err := time.ParseDuration(cfg.stringVar(name, defStr))
	if err != nil {
		panic(err)
	}
	return out
}

func (cfg *Config) boolVar(name string, def bool) bool {
	defStr := strconv.FormatBool(def)
	out, err := strconv.ParseBool(cfg.stringVar(name, defStr))
	if err != nil {
		panic(err)
	}
	return out
}

func (cfg *Config) parseCliArguments() {
	flag.StringVar(&cfg.ConfigFile, "conf-file", cfg.stringVar("conf-file", ""), "The json or yaml file to load base configuration data from.")
	flag.StringVar(&cfg.PidFile, "pid-file", cfg.stringVar("pid-file", "/var/run/nexusd.pid"), "The pid file to use for tracking rolling restarts in systemd and other supervision mechanisms.")
	flag.IntVar(&cfg.NumWorkers, "num-workers", cfg.intVar("num-workers", runtime.NumCPU()), "The number of workers to use for load balancing packets.")
	flag.StringVar(&cfg.ListenIP, "listen-ip", cfg.stringVar("listen-ip", "0.0.0.0"), "The ip addresses to listen on.")
	flag.IntVar(&cfg.ListenPort, "listen-port", cfg.intVar("listen-port", 53), "The port to listen on.")
	flag.Parse()
}

func (cfg *Config) handleFileData(data map[string]string) error {
	for k, v := range data {
		if _, ok := cfg.notSet[k]; ok {
			switch k {
			case "conf-file":
				cfg.ConfigFile = v
			case "pid-file":
				cfg.PidFile = v
			case "num-workers":
				i, err := strconv.Atoi(v)
				if err != nil {
					return err
				}
				cfg.NumWorkers = i
			case "listen-ip":
				cfg.ListenIP = v
			case "listen-port":
				i, err := strconv.Atoi(v)
				if err != nil {
					return err
				}
				cfg.ListenPort = i
			}
		}
	}

	return nil
}

func (cfg *Config) parseFileArguments() error {
	if cfg.ConfigFile != "" {
		buf, err := ioutil.ReadFile(cfg.ConfigFile)
		if err != nil {
			return err
		}

		data := make(map[string]string)
		ext := path.Ext(cfg.ConfigFile)
		switch {
		case ".json" == ext:
			err = json.Unmarshal(buf, &data)
		case ".yaml" == ext || ".yml" == ext:
			err = yaml.Unmarshal(buf, &data)
		default:
			return errors.New("The configuration file is not in a supported format.")
		}

		return cfg.handleFileData(data)
	}
	return nil
}

func (cfg *Config) parseComputedArguments() error {
	runtime.GOMAXPROCS(cfg.NumWorkers)

	cfg.RollingRestart = os.Getenv(RollingRestart)
	if cfg.RollingRestart != "" {
		cfg.ReuseFDS = true
	}

	ip := net.ParseIP(cfg.ListenIP)

	var sa syscall.Sockaddr
	if ipv4 := ip.To4(); ipv4 != nil {
		var addr [4]byte
		copy(addr[:], ipv4)
		sa = &syscall.SockaddrInet4{
			Port: cfg.ListenPort,
			Addr: addr,
		}
	} else {
		var addr [16]byte
		ipv6 := ip.To16()
		copy(addr[:], ipv6)
		sa = &syscall.SockaddrInet6{
			Port: cfg.ListenPort,
			Addr: addr,
		}
	}

	cfg.ListenAddress = sa

	pid := os.Getpid()
	ioutil.WriteFile(cfg.PidFile, []byte(strconv.Itoa(pid)), os.ModePerm)

	return nil
}

// NewConfig object
func NewConfig(log *Logger) (*Config, error) {
	cfg := &Config{notSet: make(map[string]bool)}

	cfg.parseCliArguments()
	if err := cfg.parseFileArguments(); err != nil {
		return nil, err
	}
	if err := cfg.parseComputedArguments(); err != nil {
		return nil, err
	}

	return cfg, nil
}
