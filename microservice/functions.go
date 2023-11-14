package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"github.com/segmentio/kafka-go"
	"sort"
	"strconv"
	"strings"
)

func StartKafka(kafkaHost, kafkaPort, topic string, feedback func([]byte)) {
	conf := kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%v:%v", kafkaHost, kafkaPort)},
		GroupID:  "1",
		Topic:    topic,
		MaxBytes: 1000,
	}
	reader := kafka.NewReader(conf)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Some error occured", err)
			continue
		}
		feedback(m.Value)
	}
}
func ConvertObjectToMapInterface(objIn interface{}) map[string]interface{} {
	var responseObject = make(map[string]interface{})
	b, _ := json.Marshal(objIn)
	_ = json.Unmarshal(b, &responseObject)
	return responseObject
}
func MakeLogFile(dataIn interface{}, ref string) {
	b, _ := json.Marshal(dataIn)
	dt, hr := global.GetDateAndTimeString()
	str := string(b)
	fileName := fmt.Sprintf("./logs-event/%v--%v--%v.json", ref, dt, hr)
	global.CreateFolderUploadIfNotExist("logs-event")
	global.WriteNewLineToLogFile(str, fileName)
}
func ConvertStringToByte(str string) []byte {
	return []byte(str)
}
func MicroServiceSaveInfoDb(dbHost, dbName, table string, record interface{}, MidwareCassandraQueryUrl string) string {
	var hub ServiceCassOrm
	hub.DbHost = dbHost
	hub.DbName = dbName
	hub.DbTable = table
	hub.Body = record

	urlCass := fmt.Sprintf("%v/insert", MidwareCassandraQueryUrl)
	post := ConvertObjectToMapInterface(hub)
	res := global.HttpContentPost(urlCass, post)
	return res
}
func MicroServiceGenerateNewDbTableColCode(dbName, table, col, prefix string, startWith float64, dbHost string, MidwareCassandraQueryUrl string) string {
	var hub ServiceCassOrm
	hub.DbHost = dbHost
	hub.DbName = dbName
	hub.DbTable = table
	hub.SelectFieldList = []string{col}

	urlCass := fmt.Sprintf("%v/select-all", MidwareCassandraQueryUrl)
	post := ConvertObjectToMapInterface(hub)
	res := global.HttpContentPost(urlCass, post)

	type ResponseFromCassType struct {
		REQUEST interface{}
		RESULT  []interface{}
	}
	var response ResponseFromCassType
	_ = json.Unmarshal(ConvertStringToByte(res), &response)

	fmt.Println("res ******> ", res)
	var code string
	if len(response.RESULT) == 0 {
		code = fmt.Sprintf("%v%v", prefix, startWith)
	} else {
		var maps = make(map[string]string)
		var infos []map[string]string
		bb, _ := json.Marshal(response.RESULT)
		_ = json.Unmarshal(bb, &infos)
		var prevs []int
		for _, row := range infos {
			b, _ := json.Marshal(row)
			_ = json.Unmarshal(b, &maps)
			value, _ := maps[col]
			value = strings.ReplaceAll(value, prefix, "")
			v, _ := strconv.Atoi(value)
			prevs = append(prevs, v)
		}
		sort.Ints(prevs)
		lastIndex := len(prevs) - 1
		if lastIndex < 0 {
			lastIndex = 0
		}
		global.DisplayObject("Prevs", prevs)
		lastValue := prevs[lastIndex]
		nextValue := lastValue + 1
		code = fmt.Sprintf("%v%v", prefix, nextValue)
	}

	fmt.Println(">>>>> ", code)
	return code
}
func MicroServiceSaveToEvents(topic, action string, record interface{}, dbHost, dbName, dbTable, dbAppName, messageBrokerUrl string) string {
	var hub ServiceSaveEventsCassandra
	hub.DbHost = dbHost
	hub.DbName = dbName
	hub.DbTable = dbTable
	hub.AppName = dbAppName
	hub.Topic = topic
	hub.Payload = map[string]interface{}{
		"Action":  action,
		"Payload": ConvertObjectToMapInterface(record),
	}
	urlCass := fmt.Sprintf(messageBrokerUrl)
	post := ConvertObjectToMapInterface(hub)
	res := global.HttpContentPost(urlCass, post)
	fmt.Println("::::->", res)
	return res
}
