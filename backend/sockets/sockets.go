package sockets

// This module contains functionality to create websocket connections that listen to redis connection

import (
	"errors"
	"log/slog"
	"scribl-clone/eventListener"

	"github.com/gorilla/websocket"
)

var EXPECTED_CLOSE_ERRORS = []int{
	websocket.CloseNormalClosure,
	websocket.CloseGoingAway,
	websocket.CloseProtocolError,
	websocket.CloseUnsupportedData,
	websocket.CloseNoStatusReceived,
	websocket.CloseAbnormalClosure,
	websocket.CloseInvalidFramePayloadData,
	websocket.ClosePolicyViolation,
	websocket.CloseMessageTooBig,
	websocket.CloseMandatoryExtension,
	websocket.CloseInternalServerErr,
	websocket.CloseServiceRestart,
	websocket.CloseTryAgainLater,
	websocket.CloseTLSHandshake,
}

type Connection struct {
	id       string
	channel  string
	ws       *websocket.Conn
	callback func(string)
}

func (c *Connection) SetCallback(callback func(string)) {
	c.callback = callback
}

var connections = make(map[string]*Connection)

func CreateConnection(channel string, id string, ws *websocket.Conn) (*Connection, error) {
	_, ok := connections[id]
	if ok {
		return nil, errors.New("connection with this id already exists")
	}

	conn := &Connection{
		channel:  channel,
		id:       id,
		ws:       ws,
		callback: func(s string) {},
	}
	connections[id] = conn

	eventListener.Subscribe(channel, id, func(data string) {
		ws.WriteMessage(websocket.TextMessage, []byte(data))
	})

	go readMessages(conn)

	return conn, nil
}

func readMessages(conn *Connection) {
	for {
		messageType, message, err := conn.ws.ReadMessage()
		if closeErr, ok := err.(*websocket.CloseError); ok {
			if websocket.IsUnexpectedCloseError(closeErr, EXPECTED_CLOSE_ERRORS...) {
				slog.Error(closeErr.Error())
			}
			CloseConnection(conn.id)
			return
		}
		if messageType == websocket.TextMessage {
			conn.callback(string(message))
		}
	}
}

func CloseConnection(id string) {
	conn, ok := connections[id]
	if !ok {
		return
	}
	conn.ws.Close()
	eventListener.Unsubscribe(conn.channel, conn.id)
	delete(connections, id)
	slog.Debug("closing connection", "connectionId", id)
}
