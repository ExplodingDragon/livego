package configure

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert/yaml"
	"os"
)

type Application struct {
	Appname string `yaml:"appname"`
	Live    bool   `yaml:"live"`
	Hls     bool   `yaml:"hls"`
	Flv     bool   `yaml:"flv"`
	Secret  string `yaml:"secret"`
}

type Applications []Application

type ServerCfg struct {
	Level           string `yaml:"level,omitempty"`
	FLVArchive      bool   `yaml:"flv_archive,omitempty"`
	FLVDir          string `yaml:"flv_dir,omitempty"`
	RTMPAddr        string `yaml:"rtmp_addr,omitempty"`
	HTTPFLVAddr     string `yaml:"httpflv_addr,omitempty"`
	HLSAddr         string `yaml:"hls_addr,omitempty"`
	HLSKeepAfterEnd bool   `yaml:"hls_keep_after_end,omitempty"`
	APIAddr         string `yaml:"api_addr,omitempty"`
	ReadTimeout     int    `yaml:"read_timeout,omitempty"`
	WriteTimeout    int    `yaml:"write_timeout,omitempty"`
	EnableTLSVerify bool   `yaml:"enable_tls_verify,omitempty"`
	GopNum          int    `yaml:"gop_num,omitempty"`

	ServerCfg Applications `yaml:"server"`
}

var Cfg = &ServerCfg{
	FLVArchive:      false,
	RTMPAddr:        ":1935",
	HTTPFLVAddr:     ":7001",
	HLSAddr:         ":7002",
	HLSKeepAfterEnd: false,
	APIAddr:         ":8090",
	WriteTimeout:    10,
	ReadTimeout:     10,
	EnableTLSVerify: true,
	GopNum:          1,
	ServerCfg:       make(Applications, 0)}

func ParseConfig(path string) (*ServerCfg, error) {
	result := &ServerCfg{
		FLVArchive:      false,
		RTMPAddr:        ":1935",
		HTTPFLVAddr:     ":7001",
		HLSAddr:         ":7002",
		HLSKeepAfterEnd: false,
		APIAddr:         ":8090",
		WriteTimeout:    10,
		ReadTimeout:     10,
		EnableTLSVerify: true,
		GopNum:          1,
		ServerCfg:       make(Applications, 0)}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, result); err != nil {
		return nil, err
	}
	if l, err := log.ParseLevel(result.Level); err == nil {
		log.SetLevel(l)
		log.SetReportCaller(l == log.DebugLevel)
	}
	return result, nil

}
