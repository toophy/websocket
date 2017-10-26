package game

import "github.com/gorilla/websocket"

type AccountData struct {
	ID      int64
	Account string
	Pwd     string
	Mt      int
	C       *websocket.Conn
	PosX    float64
	PosY    float64
}
