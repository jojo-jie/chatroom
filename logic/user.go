package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"io"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

var globalUID uint32 = 0

type User struct {
	UID            int           `json:"uid"`
	Nickname       string        `json:"nickname"`
	EnterAt        time.Time     `json:"enter_at"`
	Addr           string        `json:"addr"`
	MessageChannel chan *Message `json:"-"`
	Token          string        `json:"token"`
	conn           *websocket.Conn
	isNew          bool
}

// System 系统用户，代表是系统主动发送的消息
var System = &User{}

func NewUser(conn *websocket.Conn, token, nickname, addr string) *User {
	user := &User{
		Nickname:       nickname,
		Addr:           addr,
		EnterAt:        time.Now(),
		MessageChannel: make(chan *Message, 32),
		Token:          token,
		conn:           conn,
	}
	if user.Token != "" {
		uid, err := parTokenAndValidate(token, nickname)
		if err == nil {
			user.UID = uid
		}
	}
	if user.UID == 0 {
		user.UID = int(atomic.AddUint32(&globalUID, 1))
		user.Token = getToken(user.UID, nickname)
		user.isNew = true
	}
	return user
}

func (u *User) SendMessageChannel(msg *Message) {
	u.MessageChannel <- msg
}

func (u *User) SendMessage(ctx context.Context) {
	for msg := range u.MessageChannel {
		wsjson.Write(ctx, u.conn, msg)
	}
}

// CloseMessageChannel 避免goroutine 泄漏
func (u *User) CloseMessageChannel() {
	close(u.MessageChannel)
}

func (u *User) ReceiveMessage(ctx context.Context) error {
	var (
		receiveMsg map[string]string
		err        error
	)
	for {
		err = wsjson.Read(ctx, u.conn, &receiveMsg)
		if err != nil {
			//判断连接是否关闭了，正常关闭，不认为是错误
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		// 内容发送到聊天室
		sendMsg := NewMessage(u, receiveMsg["content"], receiveMsg["send_time"])
		sendMsg.Content = FilterSensitive(sendMsg.Content)

		//解析 content , 看看"@"了谁
		reg, err := regexp.Compile(`@[^\s@]{2,20}`)
		if err != nil {
			return err
		}
		log.Println(sendMsg.Content)
		sendMsg.Ats = reg.FindAllString(sendMsg.Content, -1)
		log.Println(sendMsg.Ats)
		Broadcaster.Broadcast(sendMsg)
	}
}

func parTokenAndValidate(token, nickname string) (int, error) {
	pos := strings.LastIndex(token, "uid")
	messageMAC, err := base64.StdEncoding.DecodeString(token[:pos])
	if err != nil {
		return 0, err
	}
	uid := cast.ToInt(token[pos+3:])
	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)
	ok := validateMAC([]byte(message), []byte(messageMAC), []byte(secret))
	if ok {
		return uid, nil
	}
	return 0, errors.New("token is illegal")
}

func getToken(uid int, nickname string) string {
	secret := viper.GetString("token-secret")
	message := fmt.Sprintf("%s%s%d", nickname, secret, uid)
	messageMAC := macSha256([]byte(message), []byte(secret))
	return fmt.Sprintf("%suid%d", base64.StdEncoding.EncodeToString(messageMAC), uid)
}

func macSha256(message, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func validateMAC(message, messageMAC, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
