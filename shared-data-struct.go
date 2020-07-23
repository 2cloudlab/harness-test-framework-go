package main

type EventParams struct {
	LambdaFunctionName    string `json:"functionName"`
	DynamoDBName          string `json:"dynamoDBName"`
	Iteration             int    `json:"iteration"`
	CountInSingleInstance int    `json:"countInSingleInstance"`
}
