package server

import (
	"chatroom/global"
	"chatroom/logic"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func homeHandleFunc(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		fmt.Fprint(writer, "模板解析错误1")
		return
	}
	err = t.Execute(writer, nil)
	if err != nil {
		fmt.Fprint(writer, "模板解析错误")
		return
	}
}

func userListHandleFunc(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	userList :=logic.Broadcaster.GetUserList()
	b,err:=json.Marshal(userList)
	if err != nil {
		writer.Write([]byte(`[]`))
	} else {
		writer.Write(b)
	}
}