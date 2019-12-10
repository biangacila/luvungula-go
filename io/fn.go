package io

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"strings"
	"time"
)

func QueryCassandra(qry string) string {

	var err error

	cluster := gocql.NewCluster(CusterHost1, CusterHost2, CusterHost3)
	cluster.Keyspace = "system"
	cluster.Timeout = 0
	cluster.Consistency = gocql.LocalQuorum
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println("Cassandra cluster.CreateSession err : ", err)
		return "[]"
	}

	defer func() {
		RecoverMe("QueryCassandra")
		session.Close()
	}()

	iter := session.Query(qry).Iter()
	myrow, err := iter.SliceMap()
	if err != nil {
		fmt.Println("Cassandra session.Query 2  Error --->>> ", err, " > ", qry)
		return "[]"
	}
	str, _ := json.Marshal(myrow)
	return string(str)
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
func GetDateAndTimeString() (string, string) {
	mydate := time.Now()
	arr := strings.Split(fmt.Sprintln(mydate.Format("2006-01-02 15:04:05")), " ")
	date := arr[0]
	time := arr[1]
	return strings.TrimSpace(date), strings.TrimSpace(time)
}



