# mydocker

《自己动手写 docker》笔记和源码


建议先了解一下 Docker 的核心原理大致分析，可以看这几篇文章：
* **核心原理**：[深入理解 Docker 核心原理：Namespace、Cgroups 和 Rootfs](https://mp.weixin.qq.com/s?__biz=Mzk0NzE5OTQyOQ==&mid=2247483699&idx=1&sn=177ce68bfe5b66676374450cca8a270c&chksm=c37bcd99f40c448fdd65a057160f8941c97d2a76f8607948fb7381a3d5089e61df8ff1e32ef7#rd)。
* **基于 namespace 的视图隔离**：[探索 Linux Namespace：Docker 隔离的神奇背后](https://mp.weixin.qq.com/s?__biz=Mzk0NzE5OTQyOQ==&mid=2247483717&idx=1&sn=e30fe959dfc9d7cd0dae0585004ec3e4&chksm=c37bcdeff40c44f94dbb08316f73feaba74f6aec354ba9d5afcb61f7ef821adf891c52e2941b#rd)
* **基于 cgroups 的资源限制**
    * [初探 Linux Cgroups：资源控制的奇妙世界](https://mp.weixin.qq.com/s?__biz=Mzk0NzE5OTQyOQ==&mid=2247483984&idx=1&sn=17e410280d893861656cffabe04aaf51&chksm=c37bcefaf40c47ec2fcebd11e72671a38bd668be686d107237dbe7e44a5cb0e4c001c910433b#rd)
    * [深入剖析 Linux Cgroups 子系统：资源精细管理](https://mp.weixin.qq.com/s?__biz=Mzk0NzE5OTQyOQ==&mid=2247484038&idx=1&sn=3e5c2917f67c4d42c2a5d3f8ca6ec371&chksm=c37bce2cf40c473a4987b805e623dd6c4bc219ab51549752fc80abaa9e1418a4562fd0df0f0b#rd)
    * [Docker 与 Linux Cgroups：资源隔离的魔法之旅](https://mp.weixin.qq.com/s?__biz=Mzk0NzE5OTQyOQ==&mid=2247484043&idx=1&sn=d2668d10623d329be62c5ef1e299c084&chksm=c37bce21f40c473786db38b655ebd28ca9897f7ce2ff073eb9f3a6d179f03c6a7948665b0e2c#rd)
* **基于 overlayfs 的文件系统**：[Docker 魔法解密：探索 UnionFS 与 OverlayFS](https://mp.weixin.qq.com/s?__biz=Mzk0NzE5OTQyOQ==&mid=2247484175&idx=1&sn=4c7c0105cdac469842774b0bb1495e2c&chksm=c37bcfa5f40c46b3a705412f832af86e09823a7bb6083b8c63b3e734a6ec9c8bce3f68d644c4#rd)
* **基于 veth pair、bridge、iptables 等等技术的 Docker 网络**：[揭秘 Docker 网络：手动实现 Docker 桥接网络](https://mp.weixin.qq.com/s?__biz=Mzk0NzE5OTQyOQ==&mid=2247484280&idx=1&sn=c33ce213e561486a1b61b9bbb42ed54d&chksm=c37bcfd2f40c46c45d36a570ea4d7540f7b24ec85dc6547664d5ca7607f3669f92815359a3f6#rd)

通过上述文章，大家对 Docker 的实现原理已经有了初步的认知，接下来我们就用 Golang 手动实现一下自己的 docker(mydocker)。

## 微信公众号：探索云原生

> 鸽了很久之后，终于开通了，欢迎关注。

一个云原生打工人的探索之路，专注云原生，Go，坚持分享最佳实践、经验干货。

扫描下面的二维码关注我的微信公众帐号，一起`探索云原生`吧~

![](https://img.lixueduan.com/about/wechat/qrcode_search.png)


## 实现 mydocker run 命令
搭配 [从零开始写 Docker：实现 run 命令](https://mp.weixin.qq.com/s?__biz=Mzk0NzE5OTQyOQ==&mid=2247484581&idx=1&sn=6474b3a088c9d0e4be6717b668c2b2cc&chksm=c37bc80ff40c4119becc95163201d2646b36eefa6a1010d0b078ab2df258cd56e479bcaedf29#rd) 食用更加~。

---

开发环境如下：
```bash
root@mydocker:~# lsb_release -a
No LSB modules are available.
Distributor ID:	Ubuntu
Description:	Ubuntu 20.04.2 LTS
Release:	20.04
Codename:	focal
root@mydocker:~# uname -r
5.4.0-74-generic
```
---


测试脚本如下：
```bash 
# 克隆代码
git clone -b feat-run https://github.com/lixd/mydocker.git
cd mydocker
# 拉取依赖并编译
go mod tidy
go build .
# 测试
./mydocker run -it /bin/ls # 需要 root 权限
```

正常结果
```bash
root@mydocker:~/mydocker# ./mydocker run -it /bin/ls
{"level":"info","msg":"init come on","time":"2024-01-08T09:32:52+08:00"}
{"level":"info","msg":"command: /bin/ls","time":"2024-01-08T09:32:52+08:00"}
{"level":"info","msg":"command:/bin/ls","time":"2024-01-08T09:32:52+08:00"}
LICENSE  Makefile  README.md  container  example  go.mod  go.sum  main.go  main_command.go  mydocker  run.go
root@mydocker:~/mydocker# ./mydocker run -it /bin/sh
{"level":"info","msg":"init come on","time":"2024-01-08T09:32:54+08:00"}
{"level":"info","msg":"command: /bin/sh","time":"2024-01-08T09:32:54+08:00"}
{"level":"info","msg":"command:/bin/sh","time":"2024-01-08T09:32:54+08:00"}
# ps -e
    PID TTY          TIME CMD
      1 pts/1    00:00:00 sh
      5 pts/1    00:00:00 ps
```

## 代码分析

mydocker 的代码分为以下几个部分：

```sh
.
├── container
│   ├── container_process.go # 构建容器进程运行参数
│   └── init.go # 初始化容器进程，并执行容器进程
├── example
│   └── main.go # 单独的文件，可编译成独立的可执行文件， 一个 Go 中调用 namespace 和 Cgroups 的例子，不牵涉其他 go 文件
├── go.mod
├── go.sum
├── LICENSE
├── main_command.go # 命令行解析，包含两个部分 run 和 init
├── main.go # main 函数入口
├── Makefile
├── README.md
└── run.go # 启动子进程
```

下面介绍基本的执行流程：

当执行`./mydocker run -it /bin/ls`时，会先执行到 `main_command.go::cli.Command::Action`，在里面会提取出`/bin/ls`命令（被保存在`cmd`变量中），tty变量则是用于确定是否需要打开新终端。

之后会调用`run.go::Run`函数，该函数会调用`container/container_process.go::NewParentProcess`函数，构建子进程运行参数。在构建参数时，会指示创建新的namespaces.
之后`run.go::Run`函数启动子进程。在`container/container_process.go::NewParentProcess`中构建的子进程参数如下：`/proc/self/exe init /bin/ls`，表示要创建的子进程就是自身，只不过要执行的命令是`init`，参数是`/bin/ls`。

在新创建的子进程中，mydocker会执行到`main_command.go::cli.Command::Action`，这里会调用`container/init.go::RunContainerInitProcess`函数。在该运行函数中，会挂在相应的目录，最后会执行`/bin/ls`命令。