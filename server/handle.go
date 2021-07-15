package server

import (
	"chatroom/logic"
	"github.com/gorilla/mux"
)

func RegisterHandle() *mux.Router {
	// 广播消息处理
	go logic.Broadcaster.Start()

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandleFunc)
	r.HandleFunc("/user_list", userListHandleFunc)
	r.HandleFunc("/ws", websocketHandleFunc)
	return r
}
