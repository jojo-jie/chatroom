package global

import (
	"bytes"
	"chatroom"
	"github.com/spf13/viper"
)

var (
	SensitiveWords []string
	MessageQueueLen  = 1024
)

func initConfig()  {

	/*viper.SetConfigName("chatroom")
	viper.AddConfigPath(RootDir+"/config")
	err := viper.ReadInConfig()*/

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader(chatroom.C))
	if err != nil {
		panic(err)
	}
	SensitiveWords = viper.GetStringSlice("sensitive")
	MessageQueueLen = viper.GetInt("message-queue")
	/*viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		viper.ReadInConfig()
		SensitiveWords = viper.GetStringSlice("sensitive")
	})*/
}
