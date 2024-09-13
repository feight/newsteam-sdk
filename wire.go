package newsteam

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"buf.build/gen/go/dgroux/newsteam/connectrpc/go/v1/v1connect"
	"buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/admin"
	v1 "buf.build/gen/go/dgroux/newsteam/protocolbuffers/go/v1"
	buf "connectrpc.com/connect"
	"github.com/fatih/color"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	state map[string]Feed
	mux   *http.ServeMux
	mu    sync.Mutex
)

type Feed interface {
	Id() string
	GetEnv() (string, error)
	GetLogfiles(state *admin.Cursor) ([][]byte, error)
	ProcessLogfile(*admin.Feed, []byte) []*admin.Article
}

func register(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mu.Lock()
	defer mu.Unlock()

	mux.HandleFunc(pattern, handler)
}

func InitializeFeeds(feeds []Feed) {
	state = map[string]Feed{}

	for _, feed := range feeds {
		state[feed.Id()] = feed
	}

	mux = http.NewServeMux()
	mux.Handle(v1connect.NewWireServiceHandler(&WireService{}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	server := &http.Server{
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
		Addr:         ":" + port,
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
	}

	fmt.Println("⚡ Running [http://localhost:3333]. Waiting for incoming connection from Newsteam...")

	server.ListenAndServe()
}

type WireService struct{}

func (s *WireService) GetEnv(ctx context.Context, r *buf.Request[v1.GetEnvRequest]) (*buf.Response[v1.GetEnvResponse], error) {
	if feed, ok := state[r.Msg.FeedId]; ok {

		_, err := feed.GetEnv()
		if err != nil {
			return nil, err
		}

		return buf.NewResponse(&v1.GetEnvResponse{}), nil
	}

	return nil, errors.New("Feed does not exist")
}

func (s *WireService) GetLogfiles(ctx context.Context, r *buf.Request[v1.GetLogfilesRequest]) (*buf.Response[v1.GetLogfilesResponse], error) {
	if feed, ok := state[r.Msg.FeedId]; ok {

		logfiles, err := feed.GetLogfiles(r.Msg.Cursor)
		if err != nil {
			return nil, err
		}

		return buf.NewResponse(&v1.GetLogfilesResponse{
			Data:   logfiles,
			Cursor: r.Msg.Cursor,
		}), nil
	}

	return nil, errors.New("Feed does not exist")
}

func (s *WireService) ProcessLogfile(ctx context.Context, r *buf.Request[v1.ProcessLogfileRequest]) (*buf.Response[v1.ProcessLogfileResponse], error) {
	if feed, ok := state[r.Msg.Feed.Id]; ok {
		return buf.NewResponse(&v1.ProcessLogfileResponse{
			Articles: feed.ProcessLogfile(r.Msg.Feed, r.Msg.Data),
		}), nil
	}

	return nil, errors.New("Feed does not exist")
}

func (s *WireService) OnEvent(ctx context.Context, r *buf.Request[v1.WireEvent]) (*buf.Response[v1.Success], error) {

	message := r.Msg.Message

	switch r.Msg.Type {
	case v1.WireEventType_WIRE_EVENT_TYPE_ARTICLE_CREATED:
		color.HiYellow(message)
	case v1.WireEventType_WIRE_EVENT_TYPE_ARTICLE_UPDATED:
		color.HiCyan(message)
	case v1.WireEventType_WIRE_EVENT_TYPE_ARTICLE_CREATE_ERROR:
		color.HiRed(message)
	case v1.WireEventType_WIRE_EVENT_TYPE_ARTICLE_UPDATE_ERROR:
		color.HiRed(message)
	case v1.WireEventType_WIRE_EVENT_TYPE_ARTICLE_PUBLISH_ERROR:
		color.HiRed(message)
	case v1.WireEventType_WIRE_EVENT_TYPE_LOGFILE_CREATED:
		color.HiYellow(message)
	case v1.WireEventType_WIRE_EVENT_TYPE_LOGFILE_UNCHANGED:
		color.HiGreen(message)
	case v1.WireEventType_WIRE_EVENT_TYPE_LOGFILE_PROCESSED:
		color.HiBlue(message)
	case v1.WireEventType_WIRE_EVENT_TYPE_UPLOADING_IMAGE:
		fmt.Println(message)
	}

	return buf.NewResponse(&v1.Success{}), nil
}
