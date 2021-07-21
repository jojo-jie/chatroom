package chatroom

import (
	"embed"
)

//go:embed template
var Tmpl embed.FS

//go:embed config
var Config embed.FS

var C []byte

func Init() {
	var err error
	C, err = Config.ReadFile("config/chatroom.yaml")
	if err != nil {
		panic(err)
		return
	}
}

func init()  {
	Init()
}
