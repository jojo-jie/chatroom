package main

import (
	_ "chatroom"
	"chatroom/global"
	"chatroom/server"
	"flag"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strconv"
)

var (
	addr   int
	banner = ` ____              _____
			   |    |    |   /\     |
			   |    |____|  /  \    | 
			   |    |    | /----\   |
			   |____|    |/      \  |
Go语言编程之旅 —— 一起用Go做项目：ChatRoom，start on：%s`
)

func init() {
	flag.IntVar(&addr, "p", 2021, "port")
	flag.Parse()
	global.Init()
}

func main() {
	r := server.RegisterHandle()
	g := errgroup.Group{}
	s := &http.Server{
		Addr:    ":" + strconv.Itoa(addr),
		Handler: r,
	}
	g.Go(func() error {
		return s.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
