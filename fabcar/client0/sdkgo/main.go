/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type TaskInformation struct {
	NumNodes     string
	NumIteration string
	BatchSize    string
	OutDimension string
}

type LocalEmbedding struct {
	Owener  string
	Predict string
}

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

	// config
	numNodes := 2
	numIter := 200
	batch := 2000
	outDim := 10
	// Init Task Information
	result, err := contract.SubmitTransaction("SubmitTaskInfo", "task_0", strconv.Itoa(numNodes), strconv.Itoa(numIter), strconv.Itoa(batch), strconv.Itoa(outDim))
	if err != nil {
		fmt.Printf("Failed to SubmitTaskInfo : %s\n", err)
		os.Exit(1)
	}
	fmt.Println(">> SubmitTaskInfo success :" + string(result))
	// QueryTaskInfo
	result, err = contract.EvaluateTransaction("QueryTaskInfo", "task_0")
	if err != nil {
		fmt.Printf("Failed to evaluate transaction QueryTaskInfo: %s\n", err)
		os.Exit(1)
	}
	taskInfo := new(TaskInformation)
	err = json.Unmarshal(result, taskInfo)
	if err != nil {
		fmt.Println("error unmarshal")
	}
	fmt.Println(">> QueryTaskInfo :", taskInfo)

	//开启面向训练脚本的服务器，等待其连接
	Port := "127.0.0.1:50000"
	connection, err := StartFabricServer(Port)
	if err != nil {
		fmt.Println("startFabricServer error.")
	}
	// 将训练任务发送给python
	_, err = connection.Write([]byte(fmt.Sprintf("TASK_INFO#%s#%s#%s#%s", taskInfo.NumIteration, taskInfo.BatchSize, taskInfo.OutDimension, taskInfo.NumNodes)))
	if err != nil {
		fmt.Println("TASK_INFO send fail.")
	}
	time.Sleep(time.Duration(50) * time.Millisecond)
	// FL
	numIter, err = strconv.Atoi(taskInfo.NumIteration)
	if err != nil {
		fmt.Println("numIter, err = strconv.Atoi(taskInfo.NumIteration) fail.")
	}
	batch, err = strconv.Atoi(taskInfo.BatchSize)
	if err != nil {
		fmt.Println("batch, err = strconv.Atoi(taskInfo.BatchSize) fail.")
	}
	indices := 60000.0
	indices = math.Ceil(indices / float64(batch))

	for iter := 0; iter < numIter; iter++ {
		fmt.Println("------------------------------------------------")
		fmt.Println(">> Iter : ", iter)
		sampler := iter % int(indices)
		fmt.Println("Send GO_TO_PYTHON_SAMPLER_INFO")
		_, err = connection.Write([]byte(fmt.Sprintf("GO_TO_PYTHON_SAMPLER_INFO#%s#%s", strconv.Itoa(sampler), strconv.Itoa(iter))))
		if err != nil {
			fmt.Println("SAMPLER_INFO send fail.")
		}

		fmt.Println(">> 从python脚本接收训练完毕的信号")
		buf := make([]byte, 512)
		lens, err := connection.Read(buf)
		if err != nil {
			log.Fatal(fmt.Sprintf("connection.Read error: %s", err.Error()))
		}
		signal := buf[:lens]
		//解析信号
		signalSlice := strings.Split(string(signal), "#")
		if signalSlice[0] != "PYTHON_TO_GO_PRED" {
			log.Fatal(errors.New("can't get PYTHON_TO_GO_PRED"))
		}
		fmt.Println("收到模型训练完成的信号")

		localEmbeddingFileName := "localPred_C" + signalSlice[1] + "_S" + signalSlice[2] + ".json"
		localEmbeddingAsbytes, err := ioutil.ReadFile("../results/" + localEmbeddingFileName)
		if err != nil {
			fmt.Println(localEmbeddingFileName + "err readfile")
		}
		fmt.Println("success readfile" + localEmbeddingFileName)
		result, err = contract.SubmitTransaction("SubmitLocalEmbedding", localEmbeddingFileName, "client"+signalSlice[1], string(localEmbeddingAsbytes))
		if err != nil {
			fmt.Printf("Failed to submit transaction--SubmitLocalEmbedding: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(localEmbeddingFileName + " SubmitLocalEmbedding success." + string(result))

		fmt.Println(">> QueryLocalEmbedding...")
		for i := 0; i < numNodes; i++ {
			ledgerPredFileName := "localPred_C" + strconv.Itoa(i) + "_S" + signalSlice[2] + ".json"
			result, err = contract.EvaluateTransaction("QueryLocalEmbedding", ledgerPredFileName)
			if err != nil {
				fmt.Printf(">> %s, Failed to evaluate transaction: %s\n", strconv.Itoa(i), err)
				continue
			}
			//fmt.Println("QueryLocalEmbedding client0_iter1:" + string(result))

			embedding := new(LocalEmbedding)
			err = json.Unmarshal(result, embedding)
			if err != nil {
				fmt.Println("error unmarshal")
			}
			// fmt.Println("success Unmarshal")
			// fmt.Println(reflect.TypeOf(embedding.Predict))
			err = ioutil.WriteFile("../ledger/"+ledgerPredFileName, []byte(embedding.Predict), 0644)
			if err != nil {
				fmt.Println("error writeFile")
			}
			fmt.Println(ledgerPredFileName + "writeFIle success.")
		}

		// 发送梯度帮助更新本地模型
		fmt.Println(">> 发送梯度帮助更新本地模型...")
		_, err = connection.Write([]byte(fmt.Sprintf("GO_TO_PYTHON_LEDGER#%s#%s", strconv.Itoa(sampler), strconv.Itoa(iter))))
		if err != nil {
			fmt.Println("SAMPLER_INFO send fail.")
		}
		time.Sleep(time.Duration(20) * time.Millisecond)
	}

}
