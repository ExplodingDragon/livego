package main

import (
	_ "embed"
	"errors"
	"maps"
	"net"
	"net/http"
	"slices"
	"text/template"

	"github.com/gwuhaolin/livego/configure"
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	"github.com/gwuhaolin/livego/protocol/rtmp"
	log "github.com/sirupsen/logrus"
)

var (
	//go:embed hls.min.js
	hlsJS []byte
	//go:embed index.html.tmpl
	index []byte
)
var parse, err = template.New("index").Parse(string(index))

func init() {
	if err != nil {
		panic(err)
	}
}

func startRtmp(stream *rtmp.RtmpStream, hlsServer *hls.Server) {
	rtmpListen, err := net.Listen("tcp", configure.Cfg.RTMPAddr)
	if err != nil {
		log.Fatal(err)
	}

	var rtmpServer *rtmp.Server

	rtmpServer = rtmp.NewRtmpServer(stream, hlsServer)
	defer func() {
		if r := recover(); r != nil {
			log.Error("RTMP server panic: ", r)
		}
	}()
	log.Info("RTMP Listen On ", configure.Cfg.RTMPAddr)
	rtmpServer.Serve(rtmpListen)
}

func startWeb(hFlv *httpflv.Server, hls *hls.Server) error {
	web, err := net.Listen("tcp", configure.Cfg.WebAddr)
	log.Info("Web Listen On ", configure.Cfg.WebAddr)

	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/flv/", hFlv.HandleConn)
	if configure.Cfg.HFlvInfo {
		mux.HandleFunc("/flv/streams", hFlv.GetStream)
	}
	mux.HandleFunc("/hls/", hls.Handle)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/js/hls.min.js" {
			w.Header().Set("Content-Type", "application/javascript")
			_, _ = w.Write(hlsJS)
			return
		}

		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			parse.Execute(w, map[string]interface{}{
				"rooms": configure.Cfg.Pusher,
			})
		}
	})

	go func() {
		err = http.Serve(web, mux)
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Error("Web server panic: ", err)
		}
	}()
	return nil
}

func main() {
	stream := rtmp.NewRtmpStream()

	log.Infof("Find Pusher: %s", slices.Sorted(maps.Keys(configure.Cfg.Pusher)))
	hdlServer := httpflv.NewServer(stream)
	hlsServer := hls.NewServer()
	if err := startWeb(hdlServer, hlsServer); err != nil {
		log.Fatal(err)
	}
	startRtmp(stream, hlsServer)
}
