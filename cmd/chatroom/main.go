package main

import (
	"chatroom/global"
	"chatroom/server"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

var (
	addr   = ":2021"
	banner = ` ____              _____
			   |    |    |   /\     |
			   |    |____|  /  \    | 
			   |    |    | /----\   |
			   |____|    |/      \  |
Go语言编程之旅 —— 一起用Go做项目：ChatRoom，start on：%s`
)

func init()  {
	global.Init()
}

func main() {
	r := server.RegisterHandle()
	g := errgroup.Group{}
	s:=&http.Server{
		Addr: addr,
		Handler: r,
	}
	g.Go(func() error {
		return s.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
