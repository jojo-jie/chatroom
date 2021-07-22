package server

import (
	"chatroom/global"
	"chatroom/logic"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"unsafe"
)

func websocketHandleFunc(writer http.ResponseWriter, request *http.Request) {
	//Accept 从客户端接收websocket 握手，并连接升级到websocket
	//如果 Origin 域与主机不同，Accept 将拒绝握手，除非设置了
	//InsecureSkipVerify 第三个选项用于跨域访问设置
	conn, err := websocket.Accept(writer, request, &websocket.AcceptOptions{OriginPatterns: global.AllowOrigins})
	log.Printf("conn:%+v", unsafe.Pointer(conn))
	if err != nil {
		log.Println(err)
		return
	}
	// 1.新用户进来构建该用户实例
	token := request.FormValue("token")
	nickname := request.FormValue("nickname")

	if l := len(nickname); l < 2 || l > 20 {
		wsjson.Write(request.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：2-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}

	if logic.Broadcaster.CanEnterRoom(nickname) {
		wsjson.Write(request.Context(), conn, logic.NewErrorMessage("该昵称已经已存在！"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal")
		return
	}

	userHasToken := logic.NewUser(conn, token, nickname, request.RemoteAddr)

	// 2.开启给用户发消息的goroutine
	go userHasToken.SendMessage(request.Context())
	// 3.给当前用户发送欢迎消息
	userHasToken.SendMessageChannel(logic.NewWelcomeMessage(userHasToken))

	// 避免token 泄漏
	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""
	// 给所有用户告知新用户到来
	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)
	// 4.将该用户加入到广播器列表中
	logic.Broadcaster.UserEntering(user)
	// 接收用户消息
	err = user.ReceiveMessage(request.Context())
	// 用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	//根据读取时的错误执行不同的clsoe
	if err == nil {
		conn.Close(websocket.StatusNormalClosure, "")
	} else {
		conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}
