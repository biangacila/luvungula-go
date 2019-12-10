package io

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type ParamsFilter struct {
	Key  string
	Val  interface{}
	Type string
}

type repositoryCassandra interface {
	Add() error
	Remove() error
	Find() interface{}
	List() []interface{}
	InsertQuery() string
}
type repositoryCassandraRequest struct {
	DbName      string
	Tables      []string
	TableSelect string
	Ref         string
	In          interface{}
	Conditions  []ParamsFilter
	Fields      []string
}

func (obj repositoryCassandraRequest) InsertQuery() string {
	strQry := ""
	for _, table := range obj.Tables {
		dt, hr := GetDateAndTimeString()
		OrgDateTime := fmt.Sprintf("%s %s", dt, hr)
		id := uuid.New()
		status := "active"

		o := make(map[string]interface{})
		str, _ := json.Marshal(obj.In)
		json.Unmarshal(str, &o)

		strQry = fmt.Sprintf("insert into %s.%s ", obj.DbName, table)

		strCol := "("
		strVal := " values("
		x := 0

		for key, val := range o {
			mp := fnWitchDataType(val)
			//todo let our default value
			if key == "Id" && val == "" {
				val = id
			}
			if key == "Date" && val == "" {
				val = dt
			}
			if key == "Time" && val == "" {
				val = hr
			}
			if key == "OrgDateTime" && val == "" {
				val = OrgDateTime
			}
			if key == "Status" && val == "" {
				val = status
			}

			tmpVal := val
			//fmt.Println(":> ",table," > ",key," > ",mp," > ",val)
			if mp == "string" {
				tmpVal = fmt.Sprintf("'%s'", val)
			} else if mp == "float" {
				tmpVal = fmt.Sprintf("%v", val)
			} else if mp == "float" {
				tmpVal = fmt.Sprintf("%v", val)
			} else if mp == "int" {
				tmpVal = fmt.Sprintf("%v", val)
			} else if mp == "none" {
				if val == nil {
					tmpVal = fmt.Sprintf("{}")
				}
			} else if mp == "map" {
				m := make(map[string]interface{})
				valInner, _ := json.Marshal(val)
				json.Unmarshal(valInner, &m)
				st, _ := json.Marshal(m)
				val1 := string(st)

				val1 = strings.Replace(fmt.Sprintf("%v", val1), `"`, `'`, 5000)
				val1 = strings.Replace(fmt.Sprintf("%v", val1), `\`, ``, 5000)
				val1 = strings.Replace(fmt.Sprintf("%v", val1), `<nil>`, ``, 5000)

				//fmt.Println(":> ", table, " > ", key, " > ", mp, " > ", val1)

				tmpVal = fmt.Sprintf("%v", val1)
			} else {
				tmpVal = fmt.Sprintf("'%v'", val)
			}

			val = tmpVal

			v1 := fmt.Sprintf("%v", val)
			if x == 0 {
				strCol = strCol + " " + key
				strVal = strVal + " " + v1 + " "
			} else {
				strCol = strCol + ", " + key
				strVal = strVal + ", " + v1 + " "
			}
			x++
		}
		strCol = strCol + ")"
		strVal = strVal + ") "
		strQry = strQry + strCol + strVal

		strQry = strings.Replace(strQry, "<nil>", "", 5000)

		//fmt.Println("QRY =>> ", strQry)

		//QueryCassandra(strQry)
	}

	return strQry
}
func (obj repositoryCassandraRequest) Add() error {

	for _, table := range obj.Tables {
		dt, hr := GetDateAndTimeString()
		OrgDateTime := fmt.Sprintf("%s %s", dt, hr)
		id := uuid.New()
		status := "active"

		o := make(map[string]interface{})
		str, _ := json.Marshal(obj.In)
		json.Unmarshal(str, &o)

		strQry := fmt.Sprintf("insert into %s.%s ", obj.DbName, table)

		strCol := "("
		strVal := " values("
		x := 0

		for key, val := range o {
			mp := fnWitchDataType(val)
			//todo let our default value
			if key == "Id" && val == "" {
				val = id
			}
			if key == "Date" && val == "" {
				val = dt
			}
			if key == "Time" && val == "" {
				val = hr
			}
			if key == "OrgDateTime" && val == "" {
				val = OrgDateTime
			}
			if key == "Status" && val == "" {
				val = status
			}

			tmpVal := val
			fmt.Println("xxxxxxxx:> ", table, " > ", key, " > ", mp, " > ", val)
			if mp == "string" {
				tmpVal = fmt.Sprintf("'%s'", val)
			} else if mp == "float" {
				tmpVal = fmt.Sprintf("%v", val)
			} else if mp == "bool" {
				tmpVal = fmt.Sprintf("%v", val)
			} else if mp == "none" {
				fmt.Println("none :> ", table, " > ", key, " > ", mp, " > ", val)
				if val == nil {
					tmpVal = fmt.Sprintf("{}")
				}
			} else if mp == "map" {
				m := make(map[string]interface{})
				valInner, _ := json.Marshal(val)
				json.Unmarshal(valInner, &m)
				st, _ := json.Marshal(m)
				val1 := string(st)

				val1 = strings.Replace(fmt.Sprintf("%v", val1), `"`, `'`, 5000)
				val1 = strings.Replace(fmt.Sprintf("%v", val1), `\`, ``, 5000)
				val1 = strings.Replace(fmt.Sprintf("%v", val1), `<nil>`, ``, 5000)

				fmt.Println(":> ", table, " > ", key, " > ", mp, " > ", val1)

				tmpVal = fmt.Sprintf("%v", val1)
			} else {
				tmpVal = fmt.Sprintf("'%v'", val)
			}

			val = tmpVal

			v1 := fmt.Sprintf("%v", val)
			if x == 0 {
				strCol = strCol + " " + key
				strVal = strVal + " " + v1 + " "
			} else {
				strCol = strCol + ", " + key
				strVal = strVal + ", " + v1 + " "
			}
			x++
		}
		strCol = strCol + ")"
		strVal = strVal + ") "
		strQry = strQry + strCol + strVal

		strQry = strings.Replace(strQry, "<nil>", "", 5000)

		fmt.Println("QRY =>> ", strQry)

		QueryCassandra(strQry)
	}

	return nil
}
func (obj repositoryCassandraRequest) Remove() error {
	query := fmt.Sprintf("delete ")
	where := ""

	//TODO define our condition
	if len(obj.Conditions) == 0 {
		where = "  "
	} else {
		x := 0
		for _, item := range obj.Conditions {
			if x == 0 {
				val := fmt.Sprintf("'%v'", item.Val)
				if item.Type != "string" {
					val = fmt.Sprintf("%v", item.Val)
				}
				where = fmt.Sprintf("where %v=%v", item.Key, val)
			} else {
				val := fmt.Sprintf("'%v'", item.Val)
				if item.Type != "string" {
					val = fmt.Sprintf("%v", item.Val)
				}
				where = fmt.Sprintf("%v and %v=%v", where, item.Key, val)
			}
			x++
		}
	}

	query = fmt.Sprintf("delete from  %v.%v %v", obj.DbName, obj.TableSelect, where)

	//fmt.Println("ListQuery ==> ", query)

	QueryCassandra(query)

	return nil
}
func (obj repositoryCassandraRequest) Find() interface{} {

	return nil
}
func (obj repositoryCassandraRequest) List() []interface{} {
	query := fmt.Sprintf("select ")
	col := ""
	where := ""

	//TODO define our fields
	if len(obj.Fields) == 0 {
		col = " * "
	} else {
		x := 0
		for _, item := range obj.Fields {
			if x == 0 {
				col = fmt.Sprintf("%v", item)
			} else {
				col = fmt.Sprintf("%v, %v", col, item)
			}
			x++
		}
	}

	//TODO define our condition
	if len(obj.Conditions) == 0 {
		where = "  "
	} else {
		x := 0
		for _, item := range obj.Conditions {
			if x == 0 {
				val := fmt.Sprintf("'%v'", item.Val)
				if item.Type != "string" {
					val = fmt.Sprintf("%v", item.Val)
				}
				where = fmt.Sprintf("where %v=%v", item.Key, val)
			} else {
				val := fmt.Sprintf("'%v'", item.Val)
				if item.Type != "string" {
					val = fmt.Sprintf("%v", item.Val)
				}
				where = fmt.Sprintf("%v and %v=%v", where, item.Key, val)
			}
			x++
		}
	}

	query = fmt.Sprintf("select %v from  %v.%v %v", col, obj.DbName, obj.TableSelect, where)

	//fmt.Println("ListQuery ==> ", query)

	res := QueryCassandra(query)

	var ls []interface{}

	json.Unmarshal([]byte(res), &ls)

	return ls
}

func repositoryCassandraProcessQuery(g repositoryCassandra, action string) interface{} {
	if action == "query" {
		return g.InsertQuery()
	}
	if action == "add" {
		return g.Add()
	}
	if action == "remove" {
		return g.Remove()
	}
	if action == "select" {
		return g.List()
	}
	if action == "find" {
		return g.Find()
	}
	return nil
}
func LibCassSelect(dbname string, tableName string, fields []string, conditions []ParamsFilter) interface{} {
	c := repositoryCassandraRequest{
		DbName:      dbname,
		TableSelect: tableName,
		Conditions:  conditions,
		Fields:      fields,
	}
	return repositoryCassandraProcessQuery(c, "select")
}
func LibCassDelete(dbname string, tableName string, conditions []ParamsFilter) interface{} {
	c := repositoryCassandraRequest{
		DbName:      dbname,
		TableSelect: tableName,
		Conditions:  conditions,
	}
	return repositoryCassandraProcessQuery(c, "remove")
}
func LibCassInsert(dbname string, tableName []string, in interface{}) interface{} {

	c := repositoryCassandraRequest{
		DbName: dbname,
		Tables: tableName,
		In:     in,
	}
	return repositoryCassandraProcessQuery(c, "add")
}
func LibCassInsertQuery(dbname string, tableName []string, in interface{}) interface{} {

	c := repositoryCassandraRequest{
		DbName: dbname,
		Tables: tableName,
		In:     in,
	}
	return repositoryCassandraProcessQuery(c, "query")
}

func GetTableStructure(dbName, tableName string) map[string]string {
	type StrObj struct {
		Column_name string
		Type        string
	}
	ls := []StrObj{}
	fields := []string{"column_name", "type"}
	conditions := []ParamsFilter{}
	condDB := ParamsFilter{Val: dbName, Key: "keyspace_name", Type: "string"}
	condTB := ParamsFilter{Val: tableName, Key: "table_name", Type: "string"}

	conditions = append(conditions, condDB)
	conditions = append(conditions, condTB)
	res := LibCassSelect("system_schema", "columns", fields, conditions)
	str, _ := json.Marshal(res)

	json.Unmarshal(str, &ls)

	my := make(map[string]string)
	for _, row := range ls {
		my[row.Column_name] = row.Type
	}
	return my
}
