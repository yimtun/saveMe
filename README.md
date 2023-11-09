2019 code, make a record





# saveMe



saveMe a client for save docker iamge   with the    remote  docker-engine  from   image repo







### data stream



![image-20220511112241613](README.assets/image-20220511112241613.png)





###  build



#### for linux

```
GOOS=linux  GOARCH=amd64 go build  saveMe.go
```



#### for mac



```
GOOS=darwin  GOARCH=amd64 go build  saveMe.go
```



#### for windows



```
GOOS=windows  GOARCH=amd64 go build  saveMe.go
```







### usage



config a remote docker engine server



```
edit file

vim /lib/systemd/system/docker.service

ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock



ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock -H tcp://0.0.0.0:2375




systemctl  restart docker

systemctl daemon-reload

docker -H 127.0.0.1:2375  info
```



ex:

```
sudo ./saveme -i hub.xxx.cn/nginx:1.16-alpine  -h 172.16.100.7:2375 -u image-repo-usename -p image-repo-passwd
```



windows 平台推荐使用 git bash 终端



### test image tar



on docker engine server 

```
docker load  -i ./nginx-1.16-alpine.tar 
```

