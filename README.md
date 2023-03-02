## twonodes

### 安装依赖

1、安装docker，docker-compose，golang
	docker版本号：20.10.1	`docker --version 查询版本`
	docker-compose版本号：1.24.1	`docker-compose --version 查询版本`
	golang版本号：1.15.5	`go version  查询版本`



2、执行以下脚本：

`curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s -- 2.2.0 1.4.7`

或者手动拉取docker镜像文件，分别是：

hyperledger/fabric-tools       2.2.0        
hyperledger/fabric-peer        2.2.0    
hyperledger/fabric-orderer     2.2.0        
hyperledger/fabric-ccenv       2.2.0        
hyperledger/fabric-baseos      2.2.0       
hyperledger/fabric-nodeenv     2.2.0        
hyperledger/fabric-javaenv     2.2.0        
hyperledger/fabric-ca    1.4.7        
hyperledger/fabric-zookeeper   0.4.10       
hyperledger/fabric-kafka       0.4.10       
hyperledger/fabric-couchdb  0.4.10

```shell
sudo docker pull hyperledger/fabric-tools:2.2

sudo docker pull hyperledger/fabric-peer:2.2

sudo docker pull hyperledger/fabric-orderer:2.2

sudo docker pull hyperledger/fabric-ccenv:2.2

sudo docker pull hyperledger/fabric-baseos:2.2

sudo docker pull hyperledger/fabric-nodeenv:2.2

sudo docker pull hyperledger/fabric-javaenv:2.2

sudo docker pull hyperledger/fabric-ca:1.4.7

sudo docker pull hyperledger/fabric-zookeeper:0.4.10

sudo docker pull hyperledger/fabric-kafka:0.4.10

sudo docker pull hyperledger/fabric-couchdb:0.4.10
```

需要的docker镜像：

```shell
===> List out hyperledger docker images
REPOSITORY                   TAG       IMAGE ID       CREATED         SIZE
hyperledger/fabric-tools     <none>    068c59b6dd2d   12 days ago     458MB
busybox                      latest    66ba00ad3de8   8 weeks ago     4.87MB
couchdb                      3.1.1     a81efb6c8280   17 months ago   191MB
hyperledger/fabric-ca        <none>    dbbc768aec79   2 years ago     158MB
hyperledger/fabric-tools     2.2       5eb2356665e7   2 years ago     519MB
hyperledger/fabric-tools     2.2.0     5eb2356665e7   2 years ago     519MB
hyperledger/fabric-tools     latest    5eb2356665e7   2 years ago     519MB
hyperledger/fabric-peer      2.2       760f304a3282   2 years ago     54.9MB
hyperledger/fabric-peer      2.2.0     760f304a3282   2 years ago     54.9MB
hyperledger/fabric-peer      latest    760f304a3282   2 years ago     54.9MB
hyperledger/fabric-orderer   2.2       5fb8e97da88d   2 years ago     38.4MB
hyperledger/fabric-orderer   2.2.0     5fb8e97da88d   2 years ago     38.4MB
hyperledger/fabric-orderer   latest    5fb8e97da88d   2 years ago     38.4MB
hyperledger/fabric-ccenv     2.2       aac435a5d3f1   2 years ago     586MB
hyperledger/fabric-ccenv     2.2.0     aac435a5d3f1   2 years ago     586MB
hyperledger/fabric-ccenv     latest    aac435a5d3f1   2 years ago     586MB
hyperledger/fabric-baseos    2.2       aa2bdf8013af   2 years ago     6.85MB
hyperledger/fabric-baseos    2.2.0     aa2bdf8013af   2 years ago     6.85MB
hyperledger/fabric-baseos    latest    aa2bdf8013af   2 years ago     6.85MB
hyperledger/fabric-ca        1.4       743a758fae29   2 years ago     154MB
hyperledger/fabric-ca        1.4.7     743a758fae29   2 years ago     154MB
hyperledger/fabric-ca        latest    743a758fae29   2 years ago     154MB
```



### 运行脚本

```shell
cd ./fabVFL
# 启动区块链网络
./startFabric.sh

# 开启节点1
cd client0
# 启动peer0
cd sdkgo
./runfabVFL.sh
# 启动python脚本
cd flTraining
python client.py

# 节点2开启方法同上
```



### 可能遇到的问题

1、默认GoPROXY配置是：GOPROXY=https://proxy.golang.org,direct

由于国内访问不到 https://proxy.golang.org 所以我们需要换一个PROXY，这里推荐使用`go env -w GOPROXY=https://goproxy.cn,direct`

2、`ERROR: manifest for hyperledger/fabric-orderer:latest not found: manifest unknown: manifest unknown`

进入https://hub.docker.com/r/hyperledger/fabric-orderer

查看Tags

运行`sudo docker pull hyperledger/fabric-orderer:2.2`

等待拉取完成

运行`sudo docker tag hyperledger/fabric-orderer:2.2 hyperledger/fabric-orderer:latest`，更改fabric-orderer:2.2 名称为 fabric-orderer:latest

如果fabric-peer、fabric-tools也报同样错误，就重复以上步骤
