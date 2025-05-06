package main

import (
	"flag"
	"fmt"
	"net"
	"path"
	"runtime"
	"time"

	"github.com/gwuhaolin/livego/configure"
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	"github.com/gwuhaolin/livego/protocol/rtmp"

	log "github.com/sirupsen/logrus"
)

var ConfigPath = "./livego.yaml"

func init() {
	flag.StringVar(&ConfigPath, "conf", ConfigPath, "config path")
}

type Server struct {
	Config *configure.ServerCfg
}

func NewServer(cfg string) (*Server, error) {
	config, err := configure.ParseConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Server{config}, nil

}

func (s *Server) startHls() *hls.Server {
	hlsListen, err := net.Listen("tcp", s.Config.HLSAddr)
	if err != nil {
		log.Fatal(err)
	}
	hlsServer := hls.NewServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HLS server panic: ", r)
			}
		}()
		log.Info("HLS listen On ", s.Config.HLSAddr)
		hlsServer.Serve(hlsListen)
	}()
	return hlsServer
}

func (s *Server) startRtmp(stream *rtmp.RtmpStream, hlsServer *hls.Server) {
	rtmpListen, err := net.Listen("tcp", s.Config.RTMPAddr)
	if err != nil {
		log.Fatal(err)
	}

	var rtmpServer *rtmp.Server

	if hlsServer == nil {
		rtmpServer = rtmp.NewRtmpServer(stream, nil)
		log.Info("HLS server disable....")
	} else {
		rtmpServer = rtmp.NewRtmpServer(stream, hlsServer)
		log.Info("HLS server enable....")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("RTMP server panic: ", r)
		}
	}()
	log.Info("RTMP Listen On ", s.Config.RTMPAddr)
	rtmpServer.Serve(rtmpListen)
}

func (s *Server) startHTTPFlv(stream *rtmp.RtmpStream) {

	flvListen, err := net.Listen("tcp", s.Config.HTTPFLVAddr)
	if err != nil {
		log.Fatal(err)
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HTTP-FLV server panic: ", r)
			}
		}()
		log.Info("HTTP-FLV listen On ", s.Config.HTTPFLVAddr)
		hdlServer.Serve(flvListen)
	}()
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})
}

func main() {

	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			log.Error("livego panic: ", r)
			time.Sleep(1 * time.Second)
		}
	}()
	server, err := NewServer(ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	configure.Cfg = server.Config

	for _, app := range server.Config.ServerCfg {
		stream := rtmp.NewRtmpStream()
		var hlsServer *hls.Server
		if app.Hls {
			hlsServer = server.startHls()
		}
		if app.Flv {
			server.startHTTPFlv(stream)
		}
		if app.Secret != "" {
		}
		server.startRtmp(stream, hlsServer)
	}
}
