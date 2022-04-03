/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"os"
)

func main() {
	// The gateway SDK is a tool for applications to interact with blockchain networks.
	// It provides some simple APIs to submit transactions or query to the ledger.
	gw, err := getGateway("org1")
	if err != nil {
		fmt.Printf("getGateway error: %s", err.Error())
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("fabcar")

	// FL
	result, err := contract.EvaluateTransaction("queryAllCars")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = contract.SubmitTransaction("createCar", "CAR10", "VW", "Polo", "Grey", "Mary")
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("SubmitTransaction_CAR10_Mary" + string(result))

	result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("EvaluateTransaction_queryCar_CAR10" + string(result))

	_, err = contract.SubmitTransaction("changeCarOwner", "CAR10", "Archie")
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		os.Exit(1)
	}

	result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(result))

	result, err = contract.SubmitTransaction("SubmitLocalEmbedding", "client0_iter0", "client0", "[[0], [0]]")
	if err != nil {
		fmt.Printf("Failed to submit transaction--SubmitLocalEmbedding: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("SubmitLocalEmbedding: " + string(result))

	result, err = contract.EvaluateTransaction("QueryLocalEmbedding", "client0_iter0")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("QueryLocalEmbedding: " + string(result))
}
