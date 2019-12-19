package io

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"time"
)

var ConnCass *gocql.Session
var Cluster *gocql.ClusterConfig
var CusterHost1 string
var CusterHost2 string
var CusterHost3 string

func init() {
	if CusterHost1 == "" {
		CusterHost1 = "voip.easipath.com"
	}
	if CusterHost2 == "" {
		CusterHost2 = "voip2.easipath.com"
	}
	if CusterHost3 == "" {
		CusterHost3 = "joseph.easipath.com"
	}
	cluster := gocql.NewCluster(CusterHost1, CusterHost2, CusterHost3)
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.Compressor = &gocql.SnappyCompressor{}
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{NumRetries: 3}
	cluster.Consistency = gocql.One
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.ProtoVersion = 4
	cluster.DisableInitialHostLookup = true
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println("Cassandra cluster.CreateSession err : ", err)
	}
	ConnCass = session
}

func setup_cluster(arrayHost []string) *gocql.ClusterConfig {
	cluster := gocql.NewCluster(arrayHost...)
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.Compressor = &gocql.SnappyCompressor{}
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{NumRetries: 3}
	cluster.Consistency = gocql.One
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.ProtoVersion = 4
	cluster.DisableInitialHostLookup = true
	Cluster = cluster
	return cluster
}
func RunQueryCass(qry string, arrayHost []string) string {
	if Cluster == nil {
		if len(arrayHost) == 0 {
			arrayHost = []string{"voip.easipath.com", "voip2.easipath.com", "joseph.easipath.com"}
		}
		setup_cluster(arrayHost)
	}
	session, _ := Cluster.CreateSession()
	defer session.Close()
	qResult := "[]"
	var err error
	iter := session.Query(qry).Iter()
	myrow, err := iter.SliceMap()
	if err != nil {
		log.Println("RunQueryCass Cassandra  session.Query Error --->>> ", err, " > ", qry)
		qResult = "[]"
	}
	str, _ := json.Marshal(myrow)
	qResult = string(str)
	return qResult
}
