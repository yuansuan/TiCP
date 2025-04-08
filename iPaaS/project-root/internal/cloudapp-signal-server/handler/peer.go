// Copyright (C) 2018 LambdaCal Inc.
// This file defines peer operations in signalserver

package handler

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

const (
	constPongWait      = time.Second * 60
	constPingPeriod    = (constPongWait * 8) / 10
	constWriteTimeoutS = time.Second * 60

	constTypeClient = 0
	constTypeServer = 1
)

// Peer room member datastructure
type Peer struct {
	conn      *websocket.Conn
	uid       string
	room      *Room
	datachan  chan *Message
	peerType  int
	writeLock *sync.Mutex
}

func (p *Peer) read() {
	defer p.close()

	p.conn.SetReadDeadline(time.Now().Add(constPongWait))
	p.conn.SetPongHandler(func(string) error {
		p.conn.SetReadDeadline(time.Now().Add(constPongWait))
		return nil
	})
	for {
		mt, msg, err := p.conn.ReadMessage()

		if err != nil {
			logging.Default().Infof("[Peer.read] read msg failed. err: %v, Peer: %s", err, p.uid)
			return
		}

		if mt != websocket.TextMessage {
			logging.Default().Infof("[Peer.read] get binary msg. Peer: %s", p.uid)
			continue
		}

		if p.peerType == constTypeServer {
			var jsonObj map[string]interface{}
			if err := json.Unmarshal([]byte(msg), &jsonObj); err == nil {
				if jsonObj["status_report"] != nil {
					// status report from server
					logging.Default().Infof("recv from native: %s", msg)
					var statusReportCtx StatusReportCtx
					if err := json.Unmarshal([]byte(msg), &statusReportCtx); err == nil {
						if statusReportCtx.StatusReport.ActionCount != 0 {
							p.room.sessionStat.LastActTime = time.Now().Format("2006-01-02 15:04:05")
						}
					} else {
						logging.Default().Warnf("[Peer.read] parse status_report error %s. Peer: %s", msg, p.uid)
					}
					continue
				}
			}
		}
		p.room.broadcast <- &Message{msg: msg, from: p}
	}
}

func (p *Peer) write() {
	ticker := time.NewTicker(constPingPeriod)
	defer func() {
		ticker.Stop()
		p.close()
	}()

	for {
		select {
		case msg := <-p.datachan:
			logging.Default().Debugf("msgType: %s, msg: %s", msg.mType, string(msg.msg))
			if err := p.dealMsgFromDataChan(msg); err != nil {
				logging.Default().Infof("deal msg from data channel failed, %v", err)
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(constWriteTimeoutS))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logging.Default().Warnf("[Peer.write] write ping msg failed, err: %v, Peer: %s", err, p.uid)
				return
			}
		}
	}
}

func (p *Peer) dealMsgFromDataChan(msg *Message) error {
	switch msg.mType {
	case RWType:
		if err := p.conn.WriteMessage(websocket.TextMessage, msg.msg); err != nil {
			logging.Default().Infof("[Peer.write] write normal msg failed, err: %v, Peer: %s", err, p.uid)
			return err
		}
	case CloseType:
		p.conn.Close()
		err := fmt.Errorf("connection close, Peer: %s", p.uid)
		return err
	}

	return nil
}

func (p *Peer) close() {
	if p.room != nil {
		p.room.RemovePeer(p)
	}
	p.sendCloseSignal()
	logging.Default().Infof("[Peer.close] Peer is closed. Peer: %s", p.uid)
}

// WriteMessage write message to peer
// just put it into datachan
func (p *Peer) WriteMessage(msg []byte) {
	p.datachan <- &Message{
		msg:   msg,
		mType: RWType,
	}
}

func (p *Peer) sendCloseSignal() {
	p.datachan <- &Message{
		mType: CloseType,
	}
}
