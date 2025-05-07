package configure

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
)

type ServerCfg struct {
	Level string `yaml:"level,omitempty"`

	WebAddr string `yaml:"web_addr,omitempty"`

	HLSKeepAfterEnd bool          `yaml:"hls_keep_after_end,omitempty"`
	HLSKeepTsCache  time.Duration `yaml:"hls_keep_ts_cache,omitempty"`
	WriteTimeout    int           `yaml:"rtmp_write_timeout,omitempty"`
	MaxTsCacheNum   int           `yaml:"hls_history_count,omitempty"`

	RTMPAddr        string `yaml:"rtmp_addr,omitempty"`
	EnableTLSVerify bool   `yaml:"rtmp_enable_tls_verify,omitempty"`
	GopNum          int    `yaml:"rtmp_gop_num,omitempty"`

	HFlvInfo bool `yaml:"flv_api_info,omitempty"`

	Pusher map[string]string `yaml:"pusher,omitempty"`
}

var Cfg = &ServerCfg{
	RTMPAddr:        "0.0.0.0:1935",
	WebAddr:         "127.0.0.1:7001",
	HLSKeepAfterEnd: false,
	HLSKeepTsCache:  1 * time.Minute,
	WriteTimeout:    10,
	EnableTLSVerify: true,
	HFlvInfo:        false,
	GopNum:          1,
	MaxTsCacheNum:   3,
	Pusher:          make(map[string]string),
}

var (
	conf     = flag.String("conf", "livego.yaml", "config path")
	generate = flag.Bool("generate", false, "generate config file")
)

func init() {
	flag.Parse()
	if *generate {
		Cfg.Pusher["example"] = "Ciallo～(∠・ω< )⌒☆ "
		out, _ := yaml.Marshal(Cfg)
		if conf != nil && *conf != "livego.yaml" {
			if err := os.WriteFile(*conf, out, 0o600); err != nil {
				panic(err)
			}
		} else {
			fmt.Println(string(out))
		}
		os.Exit(0)
	}
	if err := InitConfig(*conf); err != nil {
		log.Fatal(err)
	}
}

func InitConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, Cfg); err != nil {
		return err
	}
	if l, err := log.ParseLevel(Cfg.Level); err == nil {
		log.SetLevel(l)
		log.SetReportCaller(l == log.DebugLevel)
	}
	return nil
}
