package crum

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func GetDateAndTimeString() (string, string) {
	mydate := time.Now()
	arr := strings.Split(fmt.Sprintln(mydate.Format("2006-01-02 15:04:05")), " ")
	date := arr[0]
	time := arr[1]
	return strings.TrimSpace(date), strings.TrimSpace(time)
}

func fnWitchDataType(myInterface interface{}) string {
	switch myInterface.(type) {
	case int:
		return "int"
	case int64:
		return "int"
	case int32:
		return "int"
	case float64:
		return "float"
	case float32:
		return "float"
	case string:
		return "string"
	case bool:
		return "bool"
	case map[string]float64:
		return "map"
	case map[string]string:
		return "map"
	case map[string]interface{}:
		return "map"
	default:
		return "none"
	}
}

func RecoverMe(module string) {
	if r := recover(); r != nil {
		//fmt.Println("-->> RECOVERD FROM "+module+" > , ",r)
	}
}
func GetPostedDataMapAndString(r *http.Request) (map[string]interface{}, string) {
	mymap := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mymap)
	if err != nil {
		//log.Println("ERROR CREDIT NEW > ",err)
		emp := make(map[string]interface{})
		return emp, "{}"
	}
	defer r.Body.Close()

	strP, _ := json.Marshal(mymap)
	strJ := string(strP)
	//fmt.Println("GetPostedDataMapAndString ******> ",strJ)
	return mymap, strJ
}
func PublishToReact(w http.ResponseWriter, r *http.Request, obj interface{}, htmlcode int) {
	type ReactResponse struct {
		Data interface{}
	}
	myresp := ReactResponse{}
	myresp.Data = obj

	w.WriteHeader(htmlcode)
	myw, _ := json.Marshal(obj)
	w.Write(myw)
}