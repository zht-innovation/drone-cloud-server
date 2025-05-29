package gateway

import (
	"fmt"
	"sync"
	"time"

	"github.com/aler9/gortsplib/v2"
	"github.com/aler9/gortsplib/v2/pkg/format"
	"github.com/aler9/gortsplib/v2/pkg/url"
	"github.com/gorilla/websocket"
	"github.com/pion/rtp"
)

type StreamRelay struct {
	rtspURL    string
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	rtspClient *gortsplib.Client
	mutex      sync.Mutex
}

func NewStreamRelay(rtspURL string) *StreamRelay {
	return &StreamRelay{
		rtspURL:   rtspURL,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
		mutex:     sync.Mutex{},
	}
}

func (s *StreamRelay) RegisterClient(conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clients[conn] = true
}

func (s *StreamRelay) UnregisterClient(conn *websocket.Conn) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.clients[conn]; ok {
		delete(s.clients, conn)
		conn.Close()
	}
}

// StartRTSPClient 开启RTSP客户端
func (s *StreamRelay) StartRTSPClient() error {
	u, err := url.Parse(s.rtspURL)
	if err != nil {
		return err
	}

	s.rtspClient = &gortsplib.Client{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = s.rtspClient.Start(u.Scheme, u.Host)
	if err != nil {
		return err
	}

	medias, _, _, err := s.rtspClient.Describe(u)
	if err != nil {
		s.rtspClient.Close()
		return err
	}

	var format *format.H264
	media := medias.FindFormat(&format)
	if media == nil {
		s.rtspClient.Close()
		return fmt.Errorf("no H264 format found")
	}

	s.rtspClient.OnPacketRTP(media, format, func(pkt *rtp.Packet) {
		// TODO: Handle RTP packets
	})

	return nil
}
