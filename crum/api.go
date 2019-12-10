package crum

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)


func processGeneralApi(g ServiceGeneral, action string, w http.ResponseWriter, r *http.Request) {
	if action == "add" {
		g.Add(w, r)
	}
	if action == "remove" {
		g.Remove(w, r)
	}
	if action == "list" {
		g.List(w, r)
	}
}

func WsGeneralApiService(w http.ResponseWriter, r *http.Request) {
	_, dataString := GetPostedDataMapAndString(r)
	c := ServiceGeneral{}
	json.Unmarshal([]byte(dataString), &c)
	vars := mux.Vars(r)
	action, _ := vars["action"]
	processGeneralApi(c, action, w, r)
}
