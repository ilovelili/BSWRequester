package main

import (
	"fmt"
	"runtime"

	Aerospike "github.com/aerospike/aerospike-client-go"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func createStatement(value int) *Aerospike.Statement {
	stmt := Aerospike.NewStatement("DeliveryCriteria", "Delivery", "SpotId", "AdvId")
	stmt.Addfilter(Aerospike.NewEqualFilter("SpotId", value))
	stmt.IndexName = "DeliveryCriteriaIX"
	return stmt
}

func benchmark(client *Aerospike.Client, policy *Aerospike.QueryPolicy) chan *Aerospike.Record {
	ade := make(chan *Aerospike.Record)
	go func() {
		for spot := 1; ; spot++ {
			recordset, _ := client.Query(policy, createStatement(spot))
			fmt.Println(<-recordset.Records)
		}
	}()
	return ade
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := Aerospike.NewClientWithPolicyAndHost(Aerospike.NewClientPolicy(), Aerospike.NewHost("192.168.5.104", 3000), Aerospike.NewHost("192.168.5.105", 3000))
	// client, err := Aerospike.NewClientWithPolicyAndHost(Aerospike.NewClientPolicy(), Aerospike.NewHost("192.168.16.174", 3000), Aerospike.NewHost("192.168.16.175", 3000), Aerospike.NewHost("192.168.16.176", 3000))
	panicOnError(err)

	querypolicy := Aerospike.NewQueryPolicy()
	querypolicy.RecordQueueSize = 10000

	ade := make(chan *Aerospike.Record)
	ade = benchmark(client, querypolicy)

	for index := 0; ; index++ {
		fmt.Println(<-ade)
	}
}
