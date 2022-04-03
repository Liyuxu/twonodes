package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Car describes basic details of what makes up a car
type LocalEmbedding struct {
	Owener  string
	Predict string
}

// CreateCar adds a new car to the world state with given details
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
