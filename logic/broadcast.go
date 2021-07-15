package logic

import (
	"chatroom/global"
	"log"
)

// 广播器
type broadcaster struct {
	// 所有聊天室用户
	users map[string]*User
	// 所有channel 同一管理避免外部乱用
	enteringChannel chan *User
	leavingChannel  chan *User
	messageChannel  chan *Message

	// 判断该昵称用户是否进入聊天室（重复与否）true能 false不能
	checkUserChannel      chan string
	checkUserCanInChannel chan bool

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel:  make(chan *User),
	messageChannel:  make(chan *Message, global.MessageQueueLen),

	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

// Start 启动一个广播器
// 需要在一个新 goroutine 中运行，因为他不会返回
func (b *broadcaster) Start() {
	log.Println("广播器启动中...")
	for {
		select {
		case user := <-b.enteringChannel:
			// 新用户进入
			b.users[user.Nickname] = user
			log.Println("users", b.users)
			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			// 用户离开
			delete(b.users, user.Nickname)
			// 避免 goroutine 泄漏
			user.CloseMessageChannel()
		case msg := <-b.messageChannel:
			// 给所有在线用户发消息
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.SendMessageChannel(msg)
			}
			OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- true
			} else {
				b.checkUserCanInChannel <- false
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}
			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <- u
}

func (b broadcaster) Broadcast(msg *Message) {
	log.Println("Broadcast.msg", msg)
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("broadcast queue 满了")
	}
	b.messageChannel <- msg
}

func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname
	return <-b.checkUserCanInChannel
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}
