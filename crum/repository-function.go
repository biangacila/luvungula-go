package crum

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
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
}
type repositoryCassandraRequest struct {
	DbName      string
	Tables      []string
	TableSelect string
	Ref         string
	In          interface{}
	Conditions  []ParamsFilter
	Fields      []string

	CusterHost1 string
	CusterHost2 string
	CusterHost3 string

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
			//fmt.Println(":> ",table," > ",key," > ",mp," > ",val)
			if mp == "string" {
				tmpVal = fmt.Sprintf("'%s'", val)
			} else if mp == "float" {
				tmpVal = fmt.Sprintf("%0.2f", val)
			} else if mp == "bool" {
				tmpVal = fmt.Sprintf("%v", val)
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

		obj.queryCassandra(strQry)
	}

	return nil
}

func (obj repositoryCassandraRequest) queryCassandra(qry string) string {
	defer RecoverMe("QueryCassandra")
	var err error

	cluster := gocql.NewCluster(obj.CusterHost1, obj.CusterHost2, obj.CusterHost3)
	cluster.Keyspace = "system"
	cluster.Timeout = 0
	cluster.Consistency = gocql.LocalQuorum
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println("Cassandra cluster.CreateSession err : ", err)
		return "[]"
	}

	iter := session.Query(qry).Iter()
	myrow, err := iter.SliceMap()
	if err != nil {
		fmt.Println("Cassandra session.Query 2  Error --->>> ", err, " > ", qry)
		return "[]"
	}
	str, _ := json.Marshal(myrow)
	return string(str)
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

	obj.queryCassandra(query)

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

	res := obj.queryCassandra(query)

	var ls []interface{}

	json.Unmarshal([]byte(res), &ls)

	return ls
}

func repositoryCassandraProcessQuery(g repositoryCassandra, action string) interface{} {
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
