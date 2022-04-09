package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TaskInformation describes basic details of training task
type TaskInformation struct {
	NumNodes     string
	NumIteration string
	BatchSize    string
	OutDimension string
}

// LocalEmbedding describes basic details of what makes up a LocalEmbedding
type LocalEmbedding struct {
	Owener  string
	Predict string
}

// SubmitLocalEmbedding adds a new LocalEmbedding to the world state with given details
func (s *SmartContract) SubmitTaskInfo(ctx contractapi.TransactionContextInterface, taskInfoNumber string, numNodes string, numIter string, batchsize string, outDimension string) error {
	taskInfo := TaskInformation{
		NumNodes:     numNodes,
		NumIteration: numIter,
		BatchSize:    batchsize,
		OutDimension: outDimension,
	}
	taskInfoAsBytes, _ := json.Marshal(taskInfo)

	return ctx.GetStub().PutState(taskInfoNumber, taskInfoAsBytes)
}

// QueryTaskInfo returns the TaskInformation stored in the world state with given id
func (s *SmartContract) QueryTaskInfo(ctx contractapi.TransactionContextInterface, taskInfoNumber string) (*TaskInformation, error) {
	taskInfoAsBytes, err := ctx.GetStub().GetState(taskInfoNumber)

	if err != nil {
		return nil, fmt.Errorf("queryTaskInfo failed to read from world state. %s", err.Error())
	}

	if taskInfoAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", taskInfoNumber)
	}

	taskInfo := new(TaskInformation)
	_ = json.Unmarshal(taskInfoAsBytes, taskInfo)

	return taskInfo, nil
}

// SubmitLocalEmbedding adds a new LocalEmbedding to the world state with given details
func (s *SmartContract) SubmitLocalEmbedding(ctx contractapi.TransactionContextInterface, localEmbeddingNumber string, owner string, predict string) error {
	localEmbedding := LocalEmbedding{
		Owener:  owner,
		Predict: predict,
	}
	localEmbeddingAsBytes, _ := json.Marshal(localEmbedding)

	return ctx.GetStub().PutState(localEmbeddingNumber, localEmbeddingAsBytes)
}

// QueryLocalEmbedding returns the LocalEmbedding stored in the world state with given id
func (s *SmartContract) QueryLocalEmbedding(ctx contractapi.TransactionContextInterface, localEmbeddingNumber string) (*LocalEmbedding, error) {
	localEmbeddingAsBytes, err := ctx.GetStub().GetState(localEmbeddingNumber)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if localEmbeddingAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", localEmbeddingNumber)
	}

	localEmbedding := new(LocalEmbedding)
	_ = json.Unmarshal(localEmbeddingAsBytes, localEmbedding)

	return localEmbedding, nil
}
