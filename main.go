package main

import (
	"os"
	"fmt"
	"bufio"
	"github.com/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/json"
	"io"
)

type Label struct {
	ClusterName		string `json:"cluster_name"`
	ContainerName	string `json:"container_name"`
	InstanceId 		string `json:"instance_id"`
	NamespaceId 	string `json:"namespace_id"`
	PodId			string `json:"pod_id"`
	ProjectId		string `json:"project_id"`
	Zone 			string `json:"zone"`
}

type Resource struct{
	Label	Label	`json:"labels"`
	Type	string `json:"type"`
}

type Log struct {
	InsertId 		string 		`json:"insertId"`
	Labels			string		`json:"labels"`
	LogName 		string 		`json:"logName"`
	ReceiveTimeStamp	string	`json:"receiveTimeStamp"`
	Resource 		Resource 	`json:"resources"`
	Severity 		string 		`json:"severity"`
	TextPayLoad 	string 		`json:"textPayload"`
	TimeStamp 		string 		`json:"timestamp"`
}

type Record struct {
	countN, minCostTime, maxCostTime int
	costTime []int
}

var m				map[string]Record
var redirectCount	int

func main() {
	jsonFile, err := os.Open("D:\\GoProject\\src\\github.com\\HelloChenHZ\\logAnalysis\\S0.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	read := bufio.NewReader(jsonFile)
	countN := 0
	var log Log
	for {
		line, isPrefix, err := read.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		json.Unmarshal(line, &log)
		fmt.Printf("%v th log  is %v and isPrefix is %v\n", countN, log.TextPayLoad, isPrefix)
		countN ++
	}
}