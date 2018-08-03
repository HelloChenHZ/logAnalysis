package main

import (
	"os"
	"fmt"
	"bufio"
	"io"
	"regexp"
			"strconv"
	"sort"
	"encoding/json"
	"time"
		"path/filepath"
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
	costTimes []int
	path string
}

var mPath			map[string]int
var records			[] Record
var redirectCount	int
var pathCounnt 		int

func main() {
	//init
	mPath = make(map[string]int)
	pathCounnt = 0
	redirectCount = 0

	//download api log files
	now := time.Now()
	fmt.Println("year is ", now.Year(), " month is ", int(now.Month()), " hour is ", now.Hour())
	year := strconv.Itoa(now.Year())
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(now.Day())
	if now.Month() < 10 {
		month = "0"+month
	}

	if now.Day() < 10 {
		day = "0" + day
	}

	//traversing files
	//exec.Command("gsutil ","cp -r gs://crash_log_unix_io/api-v2-master/", year, "/", month, "/", day, " ./")
	fmt.Println("command is sutil cp -r gs://crash_log_unix_io/api-v2-master/", year, "/", month, "/", day, " ./")
	dir := "/home/bitmart/"+day
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error{
		if err != nil {
			return err
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		fmt.Println("visit file /home/bitmart/",day,info.Name())
		jsonFile, err := os.Open("/home/bitmart/"+day+"/"+info.Name())
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()

		read := bufio.NewReader(jsonFile)
		countN := 0
		var log Log
		for {
			//line, isPrefix, err := read.ReadLine()
			line, _, err := read.ReadLine()
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Println(err)
				continue
			}

			json.Unmarshal(line, &log)
			countN ++

			if countN > 1000 {
				break
			}

			re := regexp.MustCompile(`redirect`)
			exist := re.Match([]byte(log.TextPayLoad))
			if exist {
				redirectCount ++
			} else {
				reg := regexp.MustCompile(`Path.*?,`)
				path := reg.FindAllStringSubmatch(log.TextPayLoad, 1)
				reg = regexp.MustCompile(`End.*?（`)
				costTime := reg.FindAllStringSubmatch(log.TextPayLoad, 1)
				//fmt.Println(path)
				//fmt.Println(costTime)
				//fmt.Println(reflect.TypeOf(costTime[0][0][5: len(costTime[0][0])-3]))
				if len(path) < 1 || len(costTime) < 1 || len(path[0]) < 1 || len(costTime[0]) < 1 || len(path[0][0]) < 7 || len(costTime[0][0]) < 8 {
					continue
				}

				costTimeInt, err := strconv.Atoi(costTime[0][0][5 : len(costTime[0][0])-3])
				if err != nil {
					continue
				}

				index, ok := mPath[path[0][0][6:len(path[0][0])-1]]
				if ok {
					records[index].costTimes = append(records[index].costTimes, costTimeInt)
					records[index].countN ++
					if costTimeInt > records[index].maxCostTime {
						records[index].maxCostTime = costTimeInt
					}

					if costTimeInt < records[index].minCostTime {
						records[index].minCostTime = costTimeInt
					}

				} else {

					record := Record{1, costTimeInt, costTimeInt, [] int{1}, path[0][0][6 : len(path[0][0])-1]}
					records = append(records, record)
					mPath[path[0][0][6:len(path[0][0])-1]] = pathCounnt
					pathCounnt++
				}
			}
		}

		return nil
	})

	fmt.Println("接口地址																					次数		最大时间		最小时间		平均时间		最小时间")
	for _, record := range records {
		sort.Ints(record.costTimes)
		totalTime := 0
		for _, recordCostTime := range record.costTimes {
			totalTime = totalTime + recordCostTime
		}
		fmt.Printf("%80v			%4v			%4v			%4v			%4v			%4v\n", record.path, record.countN, record.maxCostTime, record.minCostTime, totalTime/record.countN, record.costTimes[len(record.costTimes)/2])
	}
	fmt.Println("redirect count is ", redirectCount)
}