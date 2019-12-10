package crum

import (
	"encoding/json"
	"net/http"
)

type ServiceGeneral struct {
	Organisation, Company, Ref string
	Conditions                 []ParamsFilter
	Fields                     []string
	DbName, TableName          string
	In                         interface{}
}

func (obj *ServiceGeneral) Add(w http.ResponseWriter, r *http.Request) (error, string) {
	_, dataString := GetPostedDataMapAndString(r)
	var in interface{}
	json.Unmarshal([]byte(dataString), &obj)
	in = obj.In
	tableName := obj.TableName
	dbname := obj.DbName
	tables := []string{tableName}
	LibCassInsert(dbname, tables, in)

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = nil
	my["RESULT"] = obj
	PublishToReact(w, r, my, 200)

	return nil, "ADDED"
}
func (obj *ServiceGeneral) Remove(w http.ResponseWriter, r *http.Request) (error, string) {
	_, dataString := GetPostedDataMapAndString(r)
	json.Unmarshal([]byte(dataString), &obj)

	LibCassDelete(obj.DbName, obj.TableName, obj.Conditions)

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = nil
	my["RESULT"] = obj
	PublishToReact(w, r, my, 200)

	return nil, "Removed"
}
func (obj *ServiceGeneral) List(w http.ResponseWriter, r *http.Request) interface{} {
	_, dataString := GetPostedDataMapAndString(r)
	json.Unmarshal([]byte(dataString), &obj)
	ls := LibCassSelect(obj.DbName, obj.TableName, obj.Fields, obj.Conditions)

	my := make(map[string]interface{})
	my["STATUS"] = "OK"
	my["REQUEST"] = nil
	my["RESULT"] = ls
	PublishToReact(w, r, my, 200)

	return ls
}
