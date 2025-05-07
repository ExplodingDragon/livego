package configure

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert/yaml"
	"os"
	"time"
)

type ServerCfg struct {
	Level string `yaml:"level,omitempty"`

	RTMPAddr string `yaml:"rtmp,omitempty"`
	WebAddr  string `yaml:"web,omitempty"`

	HLSKeepAfterEnd bool          `yaml:"hls_keep_after_end,omitempty"`
	HLSKeepTsCache  time.Duration `yaml:"hls_keep_ts_cache,omitempty"`

	WriteTimeout    int  `yaml:"write_timeout,omitempty"`
	EnableTLSVerify bool `yaml:"enable_tls_verify,omitempty"`
	GopNum          int  `yaml:"gop_num,omitempty"`

	HFlvInfo bool `yaml:"flv_info,omitempty"`

	Pusher        map[string]string `yaml:"pusher,omitempty"`
	MaxTsCacheNum int               `yaml:"hls_history,omitempty"`
}

var Cfg = &ServerCfg{
	RTMPAddr:        ":1935",
	WebAddr:         ":7001",
	HLSKeepAfterEnd: false,
	HLSKeepTsCache:  1 * time.Minute,
	WriteTimeout:    10,
	EnableTLSVerify: true,
	HFlvInfo:        false,
	GopNum:          1,
	MaxTsCacheNum:   60,
	Pusher:          make(map[string]string)}

var conf = flag.String("conf", "livego.yaml", "config path")

func init() {
	flag.Parse()
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
