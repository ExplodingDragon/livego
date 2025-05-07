package httpflv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gwuhaolin/livego/av"
	"github.com/gwuhaolin/livego/protocol/rtmp"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	handler av.Handler
}

type stream struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type streams struct {
	Publishers []stream `json:"publishers"`
	Players    []stream `json:"players"`
}

func NewServer(h av.Handler) *Server {
	return &Server{
		handler: h,
	}
}

// 获取发布和播放器的信息
func (server *Server) getStreams() *streams {
	rtmpStream := server.handler.(*rtmp.RtmpStream)
	if rtmpStream == nil {
		return nil
	}
	msgs := new(streams)

	rtmpStream.GetStreams().Range(func(key, val interface{}) bool {
		if s, ok := val.(*rtmp.Stream); ok {
			if s.GetReader() != nil {
				msg := stream{key.(string), s.GetReader().Info().UID}
				msgs.Publishers = append(msgs.Publishers, msg)
			}
		}
		return true
	})

	rtmpStream.GetStreams().Range(func(key, val interface{}) bool {
		ws := val.(*rtmp.Stream).GetWs()

		ws.Range(func(k, v interface{}) bool {
			if pw, ok := v.(*rtmp.PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					msg := stream{key.(string), pw.GetWriter().Info().UID}
					msgs.Players = append(msgs.Players, msg)
				}
			}
			return true
		})
		return true
	})

	return msgs
}

func (server *Server) Rooms() map[string]string {
	result := make(map[string]string)
	for _, publisher := range server.getStreams().Publishers {
		result[strings.TrimSuffix(publisher.Key, "/live")] = publisher.Id
	}
	return result
}

func (server *Server) GetStream(w http.ResponseWriter, _ *http.Request) {
	msgs := server.getStreams()
	if msgs == nil {
		return
	}
	resp, _ := json.Marshal(msgs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (server *Server) HandleConn(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("http flv HandleConn panic: ", r)
		}
	}()

	url := r.URL.String()
	u := r.URL.Path
	errMsg := fmt.Sprintf("invalid path: %s", u)
	if pos := strings.LastIndex(u, "."); pos < 0 || u[pos:] != ".flv" {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	path := strings.TrimSuffix(strings.TrimPrefix(u, "/flv/"), ".flv")
	paths := strings.SplitN(path, "/", 2)
	log.Info("url:", u, "path:", path, "paths:", paths)

	if len(paths) != 2 {
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// 判断视屏流是否发布,如果没有发布,直接返回404
	msgs := server.getStreams()
	if msgs == nil || len(msgs.Publishers) == 0 {
		http.Error(w, errMsg, http.StatusNotFound)
		return
	} else {
		include := false
		for _, item := range msgs.Publishers {
			if item.Key == path {
				include = true
				break
			}
		}
		if include == false {
			http.Error(w, errMsg, http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	writer := NewFLVWriter(paths[0], paths[1], url, w)

	server.handler.HandleWriter(writer)
	writer.Wait()
}
