package server

import (
	"chatroom"
	"chatroom/logic"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func homeHandleFunc(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFS(chatroom.Tmpl, "template/home.html")
	if err != nil {
		fmt.Fprint(writer, "模板解析错误")
		return
	}
	err = t.Execute(writer, nil)
	if err != nil {
		fmt.Fprint(writer, "模板解析错误")
		return
	}
	/*writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusBadRequest)
	writer.Write(chatroom.Template)*/
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