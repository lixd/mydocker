# mydocker

## 1. 概述

参考[《自己动手写 docker》](https://github.com/xianlubird/mydocker),自己动手从零开始实现一个简易的 docker 以及配套教程。
> 再次感谢几位作者大佬

具体差异如下：

* UnionFS替换：从AUFS 替换为 Overlayfs
* 依赖管理更新：从 go vendor 替换为 Go Module
* 一些写法上的优化调整


### 微信公众号：探索云原生

> 鸽了很久之后，终于开通了，欢迎关注。

一个云原生打工人的探索之路，专注云原生，Go，坚持分享最佳实践、经验干货。

`从零开始写 Docker` 系列持续更新中，扫描下面二维码，关注我即时获取更新~

![](https://img.lixueduan.com/about/wechat/qrcode_search.png)



### 个人博客：指月小筑(探索云原生)
在线阅读：[指月小筑(探索云原生)](https://www.lixueduan.com/categories/docker/)


## 2. 基础知识

推荐阅读以下文章对 Docker 核心原理有一个大致认识：
* **核心原理**：[深入理解 Docker 核心原理：Namespace、Cgroups 和 Rootfs](https://www.lixueduan.com/posts/docker/03-container-core/)
* **基于 namespace 的视图隔离**：[探索 Linux Namespace：Docker 隔离的神奇背后](https://www.lixueduan.com/posts/docker/05-namespace/)
* **基于 cgroups 的资源限制**
    * [初探 Linux Cgroups：资源控制的奇妙世界](https://www.lixueduan.com/posts/docker/06-cgroups-1/)
    * [深入剖析 Linux Cgroups 子系统：资源精细管理](https://www.lixueduan.com/posts/docker/07-cgroups-2/)
    * [Docker 与 Linux Cgroups：资源隔离的魔法之旅](https://www.lixueduan.com/posts/docker/08-cgroups-3/)
* **基于 overlayfs 的文件系统**：[Docker 魔法解密：探索 UnionFS 与 OverlayFS](https://www.lixueduan.com/posts/docker/09-ufs-overlayfs/)
* **基于 veth pair、bridge、iptables 等等技术的 Docker 网络**：[揭秘 Docker 网络：手动实现 Docker 桥接网络](https://www.lixueduan.com/posts/docker/10-bridge-network/)

通过上述文章，大家对 Docker 的实现原理已经有了初步的认知，接下来我们就用 Golang 手动实现一下自己的 docker(mydocker)。


## 3. 具体实现

### 构造容器

本章构造了一个简单的容器，具有基本的 Namespace 隔离，确定了基本的开发架构，后续在此基础上继续完善即可。

第一篇：
* [从零开始写 Docker：实现 run 命令](https://www.lixueduan.com/posts/docker/mydocker/01-mydocker-run/)
* 代码分支 [feat-run](https://github.com/lixd/mydocker/tree/feat-run)

第二篇：
* [从零开始写 Docker(二)---优化：使用匿名管道传参](https://www.lixueduan.com/posts/docker/mydocker/02-passing-param-by-pipe/)
* 代码分支 [opt-passing-param-by-pipe](https://github.com/lixd/mydocker/tree/opt-passing-param-by-pipe)

第三篇：
* [从零开始写 Docker(三)---基于 cgroups 实现资源限制](https://www.lixueduan.com/posts/docker/mydocker/03-resource-limit-by-cgroups/)
* 代码分支 [feat-cgroup](https://github.com/lixd/mydocker/tree/feat-cgroup)





### 构造镜像

本章首先使用 busybox 作为基础镜像创建了一个容器，理解了什么是 rootfs，以及如何使用 rootfs 来打造容器的基本运行环境。

然后，使用 OverlayFS 来构建了一个拥有二层模式的镜像，对于最上层可写层的修改不会影响到基础层。这里就基本解释了镜像分层存储的原理。

之后使用 -v 参数做了一个 volume 挂载的例子，介绍了如何将容器外部的文件系统挂载到容器中，并且让它可以访问。

最后实现了一个简单版本的容器镜像打包。

这一章主要针对镜像的存储及文件系统做了基本的原理性介绍，通过这几个例子，可以很好地理解镜像是如何构建的，第 5 章会基于这些基础做更多的扩展。

第四篇：

* [从零开始写 Docker(四)---使用 pivotRoot 切换 rootfs 实现文件系统隔离](https://www.lixueduan.com/posts/docker/mydocker/04-change-rootfs-by-pivot-root/)
* 代码分支 [feat-rootfs](https://github.com/lixd/mydocker/tree/feat-rootfs)

第五篇：

* [从零开始写 Docker(五)---基于 overlayfs 实现写操作隔离](https://www.lixueduan.com/posts/docker/mydocker/05-isolate-operate-by-overlayfs/)
* 代码分支 [feat-overlayfs](https://github.com/lixd/mydocker/tree/feat-overlayfs)

第六篇：

* [从零开始写 Docker(六)---实现 mydocker run -v 支持数据卷挂载](https://www.lixueduan.com/posts/docker/mydocker/06-volume-by-bind-mount/)
* 代码分支 [feat-volume](https://github.com/lixd/mydocker/tree/feat-volume)

第七篇：

* [从零开始写 Docker(七)---实现 mydocker commit 打包容器成镜像](https://www.lixueduan.com/posts/docker/mydocker/07-mydocker-commit/)
* 代码分支 [feat-commit](https://github.com/lixd/mydocker/tree/feat-commit)


### 构建容器进阶

本章实现了容器操作的基本功能。

* 首先实现了容器的后台运行，然后将容器的状态在文件系统上做了存储。
* 通过这些存储信息，又可以实现列出当前容器信息的功能。
* 并且， 基于后台运行的容器，我们可以去手动停止容器，并清除掉容器的存储信息。
* 最后修改了上一章镜像的存储结构，使得多个容器可以并存，且存储的内容互不干扰。

第八篇：

* [从零开始写 Docker(八)---实现 mydocker run -d 支持后台运行容器](https://www.lixueduan.com/posts/docker/mydocker/08-mydocker-run-d/)
* 代码分支 [feat-run-d](https://github.com/lixd/mydocker/tree/feat-run-d)

第九篇：

* [从零开始写 Docker(九)---实现 mydocker ps 查看运行中的容器](https://www.lixueduan.com/posts/docker/mydocker/09-mydocker-ps/)
* 代码分支 [feat-ps](https://github.com/lixd/mydocker/tree/feat-ps)


第十篇：

* [从零开始写 Docker(十)---实现 mydocker logs 查看容器日志](https://www.lixueduan.com/posts/docker/mydocker/10-mydocker-logs/)
* 代码分支 [feat-logs](https://github.com/lixd/mydocker/tree/feat-logs)



第十一篇：

* [从零开始写 Docker(十一)---实现 mydocker exec 进入容器内部](https://www.lixueduan.com/posts/docker/mydocker/11-mydocker-exec/)
* 代码分支 [feat-exec](https://github.com/lixd/mydocker/tree/feat-exec)



第十二篇：

* [从零开始写 Docker(十二)---实现 mydocker stop 停止容器](https://www.lixueduan.com/posts/docker/mydocker/12-mydocker-stop/)
* 代码分支 [feat-stop](https://github.com/lixd/mydocker/tree/feat-stop)



第十三篇：

* [从零开始写 Docker(十三)---实现 mydocker rm 删除容器](https://www.lixueduan.com/posts/docker/mydocker/13-mydocker-rm/)
* 代码分支 [feat-rm](https://github.com/lixd/mydocker/tree/feat-rm)



第十四篇：

* [从零开始写 Docker(十四)---重构：实现容器间 rootfs 隔离](https://www.lixueduan.com/posts/docker/mydocker/14-isolation-rootfs-between-containers/)
* 代码分支 [refactor-isolate-rootfs](https://github.com/lixd/mydocker/tree/refactor-isolate-rootfs)
> refactor: 文件系统重构,为不同容器提供独立的rootfs. feat: 更新rm命令，删除容器时移除对应文件系统. feat: 更新commit命令，实现对不同容器打包.



第十五篇：

* [从零开始写 Docker(十五)---实现 mydocker run -e 支持环境变量传递](https://www.lixueduan.com/posts/docker/mydocker/15-mydocker-run-e/)
* 代码分支 [feat-run-e](https://github.com/lixd/mydocker/tree/feat-run-e)


### 容器网络

在这一章中，首先手动给一个容器配置了网路，并通过这个过程了解了 Linux 虚拟网络设备和操作。然后构建了容器网络的概念模型和模块调用关系、IP 地址分配方案，以及网络模块的接口设计和实现，并且通过实现 Bridge
驱动给容器连上了“网线”。

前置：[揭秘 Docker 网络：手动实现 Docker 桥接网络](https://www.lixueduan.com/posts/docker/10-bridge-network/)

* 第十六篇：

* [从零开始写 Docker(十六)---容器网络实现(上)：为容器插上”网线“](https://www.lixueduan.com/posts/docker/mydocker/16-network-1/)
* 代码分支 [feat-network-1](https://github.com/lixd/mydocker/tree/feat-network1)

第十七篇：

* [从零开始写 Docker(十七)---容器网络实现(中)：为容器插上”网线“](https://www.lixueduan.com/posts/docker/mydocker/16-network-2/)
* 代码分支 [feat-network-2](https://github.com/lixd/mydocker/tree/feat-network2)


第十八篇：

* [从零开始写 Docker(十八)---容器网络实现(下)：为容器插上”网线“](https://www.lixueduan.com/posts/docker/mydocker/16-network-3/)
* 代码分支 [feat-network-3](https://github.com/lixd/mydocker/tree/feat-network3)


第十九篇：

* [从零开始写 Docker(十九)---增加 cgroup v2 支持](https://www.lixueduan.com/posts/docker/mydocker/19-cgroup-v2/)
* 代码分支 [feat-cgroup-v2](https://github.com/lixd/mydocker/tree/feat-cgroup-v2)


---
最后打个广告，扫描下面二维码，关注我即时获取更多文章~

![](https://img.lixueduan.com/about/wechat/qrcode_search.png)
