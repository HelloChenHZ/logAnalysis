package main

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
	LogName 		string 		`json:"logName"`
	Resource 		Resource 	`json:"resources"`
	Severity 		string 		`json:"severity"`
	TextPayLoad 	string 		`json:"textPayload"`
	TimeStamp 		string 		`json:"timestamp"`
}

func main() {
	
}