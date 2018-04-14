# simplecontainer

## 概要

一个简单的容器实现。<br>本着学习容器技术的心态，开发的一个使用的感觉上更贴近于虚拟机的容器。
本容器实现，只是为了满足我自己的学习需求和工作上的某些需求，建议只使用在测试或者开发中。
<br>本容器的部分代码参考了 https://github.com/xianlubird/mydocker 这里面的代码。
<br>也就是《自己动手写Docker》这本书中的代码。
<br>购买地址：https://www.amazon.cn/dp/B072ZDHK9S/

## 容器的实现

* 命名空间的隔离<br>
* 通过cgroups支持cpu,mem的资源限制<br>
* 每个容器拥有自己的rootfs<br>
* 通过macvlan实现对网络的支持，可以选择是静态指定还是通过DHCP获取IP<br>
* 对镜像以及容器进行管理，支持导入docker export出来的镜像 <br>

## 环境要求

现在只在 debian9 和 centos7 中做过测试<br>
一般来说只需要您的linux满足以下几点即可支持本容器<br>
* 支持AUFS或OVERLAY<br>
* 支持macvlan<br>

## 安装

##### 取得代码

``` bash
go get github.com/staroffish/simplecontainer

```

##### 用root执行代码目录底下初始化脚本

``` bash
##### 参数为一个文件夹，用来存放容器的数据
perl $GOPATH/src/github.com/staroffish/simplecontainer/init.pl /root/container

```

##### 执行上述脚本后，会在/etc下面创建一个sc.json的配置文件，<br>
##### 以及再指定的路径(这里是/root/container)下创建一系列目录

## 相关命令演示
下面用包中的busybox进行一些命令的演示
```
NAME:
   simplecontainer - A very simple container runtime implatemention.

USAGE:
   simplecontainer [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     run      Create a new container
              simplecontainer run [-m memroy_limit(mega)] [-cpu core_num] [-name container_name] [-net dhcp|static -ip ip -parent parent_dev -gateway gateway_ip] imagename
     init     Init container process run user's process in container. Do not call it outside
     exec     Execute the command in the container
              simplecontainer exec container_name command
     start    Start container
              simplecontainer start container_name
     stop     Stop container
              simplecontainer stop container_name
     rm       Remove container
              simplecontainer rm container_name
     ps       List up containers
              simplecontainer ps [-a]
     image    Image operation
               simplecontainer image [command]
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

```

### 导入容器
```
##### 导入的容器必须为tar.gz或者tgz格式,可以将docker export出来的镜像进行压缩然后直接导入进去
simplecontainer image import $GOPATH/src/github.com/staroffish/simplecontainer/example_image/busybox.tar.gz
```

### 创建一个配置了dhcp网络的容器
```
simplecontainer run -net dhcp -parent ens33 -name 容器名 busybox
```

### 查看容器状况
```
simplecontainer ps -a
```

### 进入容器
```
simplecontainer exec 容器名 命令
```

### 启动停止容器
```
simplecontainer start/stop 容器名
```

### 删除容器
```
simplecontainer rm 容器名
```

## 注意点
为了感觉上更像虚拟机，
所以这个容器和docker很大的一个不同在于，<br>
这个容器启动容器进程以后不会exec一个命令，<br>
而是只会在那里一直等待结束的signal的到来。<br>
这样做有违容器的常规用法，但是对不熟悉容器的人来说或许更好接受点
