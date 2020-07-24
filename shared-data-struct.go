package main

type EventParams struct {
	RequestID             string `json:"requestID"`
	LambdaFunctionName    string `json:"functionName"`
	DynamoDBName          string `json:"dynamoDBName"`
	Iteration             int    `json:"iteration"`
	CountInSingleInstance int    `json:"countInSingleInstance"`
	RawJson               string `json:"rawJson"`
}
