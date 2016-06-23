package main

import (
	"runtime"
	"time"

	Aerospike "github.com/aerospike/aerospike-client-go"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func createStatement(value int) *Aerospike.Statement {
	stmt := Aerospike.NewStatement("DeliveryCriteria", "Delivery", "SpotId", "AdvId", "CPC", "Score", "SecondScore", "Removed")
	stmt.Addfilter(Aerospike.NewEqualFilter("SpotId", value))
	stmt.IndexName = "DeliveryCriteriaIX"
	return stmt
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// client, err := Aerospike.NewClientWithPolicyAndHost(Aerospike.NewClientPolicy(), Aerospike.NewHost("192.168.5.104", 3000), Aerospike.NewHost("192.168.5.105", 3000))
	client, err := Aerospike.NewClientWithPolicyAndHost(Aerospike.NewClientPolicy(), Aerospike.NewHost("192.168.16.174", 3000), Aerospike.NewHost("192.168.16.175", 3000), Aerospike.NewHost("192.168.16.176", 3000))
	panicOnError(err)

	querypolicy := Aerospike.NewQueryPolicy()
	querypolicy.RecordQueueSize = 10000

	for spot := 1; spot < 5000; spot++ {
		recordset, err := client.Query(querypolicy, createStatement(spot))
		panicOnError(err)

		<-recordset.Records
	}

	time.Sleep(10 * time.Minute)
}
