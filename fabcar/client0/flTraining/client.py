from statistics import mean
import torch.nn as nn
import torch.nn.functional as F
import importlib
import torch
import pickle  # for pkl file reading
import os
import sys
import numpy as np
from torch.autograd import Variable
from torch.utils.data import DataLoader
from torch.utils.data import Dataset
from PIL import Image
import time
import scipy.io
import torchvision.transforms as transforms
from config import SERVER_ADDR, SERVER_PORT, ClientID
from utils import recv_msg, send_msg
import socket
import struct
from torchvision import transforms
import math
import json


# read data set from pkl files

def read_data(data_dir):
    """Parses data in given train and test data directories

    Assumes:
        1. the data in the input directories are .json files with keys 'users' and 'user_data'
        2. the set of train set users is the same as the set of test set users

    Return:
        clients: list of client ids
        groups: list of group ids; empty list if none found
        train_data: dictionary of train data (ndarray)
        test_data: dictionary of test data (ndarray)
    """

    # clients = []
    # groups = []
    data = {}
    print('>>> Read data from:', data_dir)

    # open training dataset pkl files
    with open(data_dir, 'rb') as inf:
        cdata = pickle.load(inf)
    # print('cdata :', cdata)
    data.update(cdata)

    data = MiniDataset(data['x'], data['y'])

    return data


class MiniDataset(Dataset):
    def __init__(self, data, labels):
        super(MiniDataset, self).__init__()
        self.data = np.array(data)
        self.labels = np.array(labels).astype("int64")

        self.data = self.data.astype("float32")
        self.transform = None

    def __len__(self):
        return len(self.labels)

    def __getitem__(self, index):
        data, target = self.data[index], self.labels[index]

        if self.data.ndim == 4 and self.data.shape[3] == 3:
            data = Image.fromarray(data)

        if self.transform is not None:
            data = self.transform(data)

        return data, target

def calculateLoss(preds, y):
    '''

    Args:
        preds: [pred1,pred2,...,predn]
        y: labels

    Returns: CrossEntropyLoss(preds, y)

    '''

    pred = sum(preds)
    criterion = torch.nn.CrossEntropyLoss()
    batch_loss = criterion(pred, y)

    # print("batch_loss:", batch_loss, batch_loss.grad_fn)
    return batch_loss


def test_inference(model, testloader):
    """ Returns the test accuracy and loss.
    """

    model.eval()
    loss, total, correct = 0.0, 0.0, 0.0
    device = 'cpu'
    criterion = torch.nn.CrossEntropyLoss()
    for batch_idx, (images, labels) in enumerate(testloader):
        images, labels = images.to(device), labels.to(device)

        # Inference
        outputs = model(images)
        batch_loss = criterion(outputs, labels)
        loss += batch_loss.item()

        # Prediction
        _, pred_labels = torch.max(outputs, 1)
        pred_labels = pred_labels.view(-1)
        correct += torch.sum(torch.eq(pred_labels, labels)).item()
        total += len(labels)

    accuracy = correct / total
    loss = loss / total
    return accuracy, loss


# Model for MQTT_IOT_IDS dataset
class Logistic(nn.Module):
    def __init__(self, in_dim, out_dim):
        super(Logistic, self).__init__()
        self.layer = nn.Linear(in_dim, out_dim)

    def forward(self, x):
        logit = self.layer(x)
        return logit


# socket
sock = socket.socket(socket.AF_INET,socket.SOCK_STREAM)
sock.connect((SERVER_ADDR, SERVER_PORT))
print('---------------------------------------------------------------------------')
try:
    while True:
        msg = recv_msg(sock, 'TASK_INFO')
        print('Received message from server:', msg)
        num_iter = int(msg[1])
        batch_size = int(msg[2])
        out_dim = int(msg[3])
        lr_rate = 0.3
        n_nodes = int(msg[4])
        # Import the data set
        file_name = './VFLMNIST/K' + str(n_nodes) + '_' + str(ClientID) + '.pkl'
        train_data = read_data(file_name)

        print('data read successfully')
        # Init model: make the data loader, divide dataset by batch_size; shuffle = False
        test_loader = torch.utils.data.DataLoader(dataset=train_data,
                                                  batch_size=batch_size,
                                                  shuffle=False)
        trainX, trainY = [], []
        for idx, (x, y) in enumerate(test_loader):
            trainX.append(x)
            trainY.append(y)
        test_x, test_y = next(iter(test_loader))
        # init model
        in_dim = len(test_x[0])
        model = Logistic(in_dim, out_dim)

        optimizer = torch.optim.SGD(model.parameters(), lr=lr_rate)
        # Learning rate decay , lr = lr * gamma
        gamma = 0.9
        step_size = 100
        scheduler = torch.optim.lr_scheduler.StepLR(optimizer, step_size, gamma, last_epoch=-1)
        criterion = torch.nn.CrossEntropyLoss()

        cv_acc, cv_loss = [], []
        cv_time = []
        while True:
            print('---------------------------------------------------------------------------')
            iter_start = time.time()
            msg = recv_msg(sock, 'GO_TO_PYTHON_SAMPLER_INFO')
            if msg == "MISSION_END":
                saveTitle = 'Blockclient' + str(ClientID) + '_T' + str(num_iter) \
                            + 'B' + str(batch_size)
                saveVariableName = saveTitle
                scipy.io.savemat(saveTitle + '_acc' + '.mat', mdict={saveVariableName + '_acc': cv_acc})
                scipy.io.savemat(saveTitle + '_loss' + '.mat', mdict={saveVariableName + '_loss': cv_loss})
                scipy.io.savemat(saveTitle + '_time' + '.mat', mdict={saveVariableName + '_time': cv_time})
                print("mission end")
                print("每轮训练平均用时:", mean(cv_time))
                sys.exit(0)
            sampler = int(msg[1])
            cur_iter = int(msg[2])
            print(">>  Round ", cur_iter)
            print(">>   lr  =", optimizer.param_groups[0]['lr'])
            # if cur_iter == num_iter:
            #     print("mission end")
            #     sys.exit(0)

            x, y = trainX[sampler], trainY[sampler]
            print("idx: ", sampler, " y:", y)
            print('Make dataloader successfully')

            # calculate predication
            model.train()
            optimizer.zero_grad()
            pred = model(x)
            # print(">> pred\n", pred, type(pred))
            pred_json = {"local_embedding": pred.tolist()}
            # print(">> pred_json\n", pred_json, type(pred_json))
            jsonfilename = '../results/localPred_C'+str(ClientID)+'_S'+str(sampler)+'.json'
            # 本地模型写入文件
            print(">> C{}_S{}本地模型写入文件.".format(ClientID, sampler))
            json_str = json.dumps(pred_json)
            # print(">> json_str\n", json_str, type(json_str))
            with open(jsonfilename, 'w') as json_file:
                json_file.write(json_str)
            print(">> 模型写入JSON成功.")
            # send pred to server
            msg = "PYTHON_TO_GO_PRED#"+str(ClientID)+"#"+str(sampler)
            send_msg(sock, msg)

            msg = recv_msg(sock, 'GO_TO_PYTHON_LEDGER')
            preds = []
            for i in range(n_nodes):
                jsonfilename = "../ledger/"+ "localPred_C" + str(i) +"_S"+str(sampler)+".json"
                if os.path.exists(jsonfilename):
                    f = open(jsonfilename, encoding='utf-8')
                    readstr = f.readline()
                    readData = json.loads(readstr)
                    f.close()
                    # print(">> readData\n", readData["local_embedding"], type(readData["local_embedding"]))
                    temp_pred = torch.Tensor(readData["local_embedding"])
                    preds.append(temp_pred)
            global_loss = calculateLoss(preds, y)

            # for backward
            # calculate a loss value casually
            loss = criterion(pred, y)
            msg = "PYTHON_TO_GO_CONTRI#" + str(ClientID)+ "#" +str(math.exp(-1*loss.item()))
            send_msg(sock, msg)
            # global_loss.data -> loss.data
            loss.data = torch.full_like(loss, global_loss.item())

            loss.backward()
            optimizer.step()
            # learning rate decay
            scheduler.step()

            testaccuracy, testloss = test_inference(model, test_loader)
            cv_acc.append(testaccuracy)
            cv_loss.append(testloss)
            print("------ acc =", testaccuracy, "------ loss =", testloss)
            iter_end = time.time()
            cv_time.append(iter_end - iter_start)

except (struct.error, socket.error):
    print('Server has stopped')
    pass

