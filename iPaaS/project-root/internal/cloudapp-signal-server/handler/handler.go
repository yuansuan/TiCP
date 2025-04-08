package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp-signal-server/config"
)

type Handler struct {
	cfg         *config.Config
	upgrader    *websocket.Upgrader
	roomManager *RoomManager
}

type VisSessionStatus struct {
	LastActTime string `json:"last_action_time"`
	RdpStatus   string `json:"rdp_status"`
	ServerNum   int    `json:"server_num"`
	ClientNum   int    `json:"client_num"`
	RoomId      string `json:"room_id"`
}

type Report2Web struct {
	Sessions []VisSessionStatus `json:"sessions"`
}

var (
	roomReady   = []byte(`{"ready":true}`)
	roomUnready = []byte(`{"ready":false}`)
)

func (h *Handler) checkReady(w http.ResponseWriter, r *http.Request) {
	roomId := r.URL.Query().Get("room_id")
	if roomId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if h.roomManager.IsOpened(roomId) {
		_, _ = w.Write(roomReady)
	} else {
		_, _ = w.Write(roomUnready)
	}
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	if !websocket.IsWebSocketUpgrade(r) {
		return
	}
	c, err := h.upgrader.Upgrade(w, r, nil)

	if err != nil {
		logging.Default().Errorf("[handler.handle] upgrade failed, err: %v", err)
		return
	}

	p, room, err := h.negotiate(c)
	if err != nil {
		if p != nil {
			p.close()
		} else {
			c.Close()
		}
		logging.Default().Warnf("[handler.handle] negotiate failed, conn: %#v, err: %v", c.RemoteAddr(), err)
		return
	}

	logging.Default().Infof("[handler.handle] peer is connected. peer: %s, room: %s", p.uid, room.uid)
}

func (h *Handler) negotiate(c *websocket.Conn) (*Peer, *Room, error) {
	var err error

	var p *Peer
	if p, err = h.hello(c); err != nil {
		logging.Default().Infof("[handler.negotiate] hello failed, err: %v", err)
		return nil, nil, err
	}

	var room *Room
	if room, err = h.waitingRoom(p); err != nil {
		logging.Default().Infof("[handler.negotiate] waitingRoom failed, peer: %s, err: %v", p.uid, err)
		return p, nil, err
	}

	return p, room, nil
}

// HELLO client-001 type["client", "server"]
func (h *Handler) hello(c *websocket.Conn) (*Peer, error) {
	_, msg, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	cmd := strings.SplitN(string(msg), " ", 3)
	if len(cmd) < 3 || cmd[0] != "HELLO" || (cmd[2] != "client" && cmd[2] != "server") {
		logging.Default().Infof("[handler.hello] invalid hello, msg: %s", string(msg))
		return nil, errors.New("invalid hello")
	}

	// HELLO resp
	if err := c.WriteMessage(websocket.TextMessage, []byte("HELLO")); err != nil {
		return nil, err
	}

	peerType := constTypeClient
	if cmd[2] == "server" {
		peerType = constTypeServer
	}

	p := &Peer{
		uid:       cmd[1],
		room:      nil,
		conn:      c,
		datachan:  make(chan *Message, 128),
		peerType:  peerType,
		writeLock: &sync.Mutex{},
	}

	logging.Default().Infof("[handler.hello] hello, peer: %s", p.uid)
	return p, nil
}

// ROOM room-001
func (h *Handler) waitingRoom(p *Peer) (*Room, error) {
	for {
		go p.write()

		p.conn.SetReadDeadline(time.Now().Add(constPongWait))
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			logging.Default().Infof("[handler.waitingRoom] read msg failed, peer: %s, err: %v", p.uid, err)
			return nil, err
		}

		logging.Default().Infof("[handler.waitingRoom] recv room cmd, %s, peer: %s", string(msg), p.uid)

		cmd := strings.SplitN(string(msg), " ", 2)
		if len(cmd) < 2 || cmd[0] != "ROOM" {
			p.WriteMessage([]byte("ERROR invalid command"))
			logging.Default().Infof("[handler.waitingRoom] invalid command. cmd: %s", string(msg))
			return nil, errors.New("ERROR invalid command")
		}

		roomid := cmd[1]

		room, err := h.roomManager.AssignRoom(p, roomid)
		if err != nil {
			p.WriteMessage([]byte(fmt.Sprintf("ERROR %v", err)))
			logging.Default().Infof("[handler.waitingRoom] assignRoom failed. err: %v", err)
			if err.Error() == "room has no server yet" {
				continue
			} else {
				return nil, err
			}
		}

		p.WriteMessage([]byte(fmt.Sprintf("ROOM_OK %s", room.hostid)))

		p.conn.SetCloseHandler(func(code int, text string) error {
			logging.Default().Infof("[close handler] receive close msg, shuting down, peer: %s", p.uid)
			p.close()
			return nil
		})

		return room, nil
	}
}

func (h *Handler) ExportHttp() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handle)
	mux.HandleFunc("/ready", h.checkReady)
	mux.HandleFunc("/healthy", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}

func New(cfg *config.Config) (*Handler, error) {
	return &Handler{
		cfg: cfg,
		roomManager: &RoomManager{
			mu:          &sync.RWMutex{},
			rooms:       make(map[string]*Room),
			turnServers: cfg.TurnServer,
		},
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
	}, nil
}
