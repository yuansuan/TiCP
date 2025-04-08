// Copyright (C) 2018 LambdaCal Inc.
// This file defines room operations in signalserver

package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp-signal-server/config"
)

// RoomManager room collection
type RoomManager struct {
	mu          *sync.RWMutex
	rooms       map[string]*Room
	turnServers []*config.TurnServer
}

type StatusReportInfo struct {
	ActionCount int64 `json:"action_count"`
	TStartMs    int64 `json:"t_start_ms"`
	TDurationMs int64 `json:"t_duration_ms"`
}

type StatusReportCtx struct {
	StatusReport StatusReportInfo `json:"status_report"`
}

// Room chat hub
type Room struct {
	uid         string
	hostid      string
	sessionStat VisSessionStatus
	peers       map[string]*Peer
	broadcast   chan *Message
	m           *sync.RWMutex
	quitchan    chan struct{}
}

type MType string

const (
	RWType    MType = "read/write"
	CloseType MType = "close"
)

// Message room chat message body
type Message struct {
	from  *Peer
	msg   []byte
	mType MType
}

// 警告！！在建立连接过程中不应该使用向broadcast channel发送消息，如果和正常的wsConn.WriteMessage混用，会导致消息顺序不可控制
func (room *Room) open(r *RoomManager) {
	defer room.close(r)

	logging.Default().Infof("[Room.open] room is open. room: %s", room.uid)
	for {
		select {
		case m := <-room.broadcast:
			for uid, p := range room.peers {
				if uid == m.from.uid {
					continue
				}
				p.WriteMessage(m.msg)
				logging.Default().Debugf("[Room.Broadcase] to room: %s, peer: %s, msg: %s", room.uid, p.uid, string(m.msg))
			}
		case <-room.quitchan:
			logging.Default().Infof("[Room.open] room is closing, room: %s", room.uid)
			room.sessionStat.ServerNum = 0
			room.sessionStat.ClientNum = 0
			room.sessionStat.RdpStatus = "closed"
			return
		}
	}
}

func (room *Room) close(r *RoomManager) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, p := range room.peers {
		if err := p.conn.WriteControl(websocket.CloseMessage, []byte("room is closed"), time.Now().Add(constWriteTimeoutS)); err != nil {
			logging.Default().Infof("[Room.close] write close msg to peer failed, err: %v", err)
		}
	}

	delete(r.rooms, room.uid)

	logging.Default().Infof("[Room.close] room is closed. room: %s", room.uid)
}

// AddPeer peer join the room
func (room *Room) AddPeer(p *Peer, turnInfos []*config.TurnServer) error {
	room.m.Lock()
	defer room.m.Unlock()
	if _, ok := room.peers[p.uid]; ok {
		// send ROOM_FULL to new joined peer
		logging.Default().Info("[Room.AddPeer] room already has peer")
		return errors.New("room_full")
	}

	logging.Default().Infof("[Room.AddPeer] broadcast join msg. room: %s, peer: %s", room.uid, p.uid)
	turnInfoStr := "TURN_CRED_INFO"
	for _, element := range turnInfos {
		//https://github.com/coturn/coturn, REST_API generate cred
		expired := time.Now().Unix() + element.Expire
		user_combo := fmt.Sprintf("%d:%s", expired, "yskj")
		mac := hmac.New(sha1.New, []byte(element.Secret))
		mac.Write([]byte(user_combo))
		userPassword := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		turn_info_each := fmt.Sprintf(" %s %s %s", element.Uri, user_combo, userPassword)
		turnInfoStr += turn_info_each
	}

	logging.Default().Infof("trying to broadcase turn info")
	room.peers[p.uid] = p
	// 由client这边来决定什么时候刷新两端的turn info，此处不能使用broadcast channel来进行消息推送，会导致消息顺序不可控
	// 此处需要保证 TURN_CRED_INFO 消息先于 ROOM_PEER_JOINED 消息前发送至双端
	if p.peerType == constTypeClient {
		for _, p := range room.peers {
			p.WriteMessage([]byte(turnInfoStr))
		}
	}

	for _, c := range room.peers {
		if c.uid == p.uid {
			continue
		}
		c.WriteMessage([]byte(fmt.Sprintf("ROOM_PEER_JOINED %s", p.uid)))
	}

	go p.read()

	logging.Default().Infof("[Room.AddPeer] addPeer, room: %s, peer: %s", room.uid, p.uid)
	return nil
}

// RemovePeer peer left the room
func (room *Room) RemovePeer(p *Peer) {
	room.m.Lock()
	defer room.m.Unlock()

	if _, ok := room.peers[p.uid]; !ok {
		// already left the room
		return
	}

	logging.Default().Infof("[Room.RemovePeer] broadcast leave msg. room: %s, peer: %s", room.uid, p.uid)

	for _, c := range room.peers {
		// 饱和式攻击，一端退出，通知所有端都退出，防止一端退出，另一端进程仍存在，但类似websocket请求被断开的异常情况
		logging.Default().Infof("[send.ROOM_PEER_LEFT] to room: %s, peer: %s", room.uid, p.uid)
		c.WriteMessage([]byte(fmt.Sprintf("ROOM_PEER_LEFT %s", p.uid)))
	}

	delete(room.peers, p.uid)

	if p.peerType == constTypeServer {
		room.sessionStat.ServerNum = 0
		room.sessionStat.RdpStatus = "closed"
	} else {
		room.sessionStat.ClientNum = len(room.peers) - 1
		if room.sessionStat.ClientNum == 0 {
			if room.sessionStat.ServerNum == 1 {
				room.sessionStat.RdpStatus = "ready"
			} else {
				room.sessionStat.RdpStatus = "closed"
			}
		}
	}

	logging.Default().Infof("[Room.RemovePeer] remove peer, room: %s, peer: %s", room.uid, p.uid)

	// close room if server is disconnected or room is empty
	if p.peerType == constTypeServer {
		logging.Default().Infof("[Room.RemovePeer] server is offline, closing rooom, room: %s", room.uid)
		p.room.quitchan <- struct{}{}
	} else if len(room.peers) == 0 {
		logging.Default().Infof("[Room.RemovePeer] room has no peer now, closing. room: %s", room.uid)
		room.quitchan <- struct{}{}
	}
}

func (r *RoomManager) IsOpened(roomId string) (ok bool) {
	r.mu.RLock()
	_, ok = r.rooms[roomId]
	r.mu.RUnlock()

	return
}

// AssignRoom get an exist room or create one for peer
func (r *RoomManager) AssignRoom(p *Peer, roomid string) (*Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	room, ok := r.rooms[roomid]
	if !ok {
		if p.peerType == constTypeClient {
			logging.Default().Infof("[RoomManager.AssignRoom] room has not server yet, room: %s, peer: %s, peertype: %d", roomid, p.uid, p.peerType)
			return nil, errors.New("room has no server yet")
		}

		room = &Room{
			uid:         roomid,
			hostid:      p.uid,
			sessionStat: VisSessionStatus{RoomId: roomid, LastActTime: time.Now().Format("2006-01-02 15:04:05"), RdpStatus: "ready", ServerNum: 1, ClientNum: 0},
			peers:       make(map[string]*Peer),
			broadcast:   make(chan *Message, 128),
			m:           &sync.RWMutex{},
			quitchan:    make(chan struct{}),
		}
		r.rooms[roomid] = room
		// open the room first
		go room.open(r)
	} else {
		for _, peer := range room.peers {
			if peer.peerType == constTypeServer && peer.peerType == p.peerType {
				logging.Default().Infof("[RoomManager.AssignRoom] room already has type %d, room: %s, peer: %s, peertype: %d", p.peerType, roomid, p.uid, peer.peerType)
				return nil, fmt.Errorf("room already has type %d", p.peerType)
			}
		}
	}
	// now p join the room, if success, will set p.room = room
	if err := room.AddPeer(p, r.turnServers); err != nil {
		return room, err
	}
	room.sessionStat.RdpStatus = "using"
	room.sessionStat.ClientNum = len(room.peers) - 1
	p.room = room

	logging.Default().Infof("[RoomManager.AssignRoom] peer: %s, peertype: %d, room: %s", p.uid, p.peerType, roomid)
	return room, nil
}
