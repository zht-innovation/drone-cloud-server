<div align="center" id="top"> 
  <!-- <img src="./.github/app.gif" alt="Cloud" /> -->

  &#xa0;

  <!-- <a href="https://cloud.netlify.app">Demo</a> -->
</div>

<h1 align="center">无人机云控系统服务端</h1>

<p align="center">
  <img alt="Github top language" src="https://img.shields.io/github/languages/top/WenRunning/zht_cloud_server?color=56BEB8"> &#xa0;
  <img alt="Github language count" src="https://img.shields.io/github/languages/count/WenRunning/zht_cloud_server?color=56BEB8"> &#xa0;
  <img alt="Repository size" src="https://img.shields.io/github/repo-size/WenRunning/zht_cloud_server?color=56BEB8"> &#xa0;
  <img alt="License" src="https://img.shields.io/github/license/WenRunning/zht_cloud_server?color=56BEB8"> &#xa0;

  <!-- <img alt="Github issues" src="https://img.shields.io/github/issues/{{YOUR_GITHUB_USERNAME}}/cloud?color=56BEB8" /> -->

  <!-- <img alt="Github forks" src="https://img.shields.io/github/forks/{{YOUR_GITHUB_USERNAME}}/cloud?color=56BEB8" /> -->

  <!-- <img alt="Github stars" src="https://img.shields.io/github/stars/{{YOUR_GITHUB_USERNAME}}/cloud?color=56BEB8" /> -->
</p>

<!-- Status -->

<h4 align="center"> 
	🚧  云端 🚀 正在开发中 ……  🚧
</h4> 

<p align="center">
  <a href="#简介">简介</a> &#xa0; | &#xa0; 
  <a href="#技术">技术</a> &#xa0; | &#xa0;
  <a href="#前置条件">前置条件</a> &#xa0; | &#xa0;
  <a href="#checkered_flag-starting">开始</a> &#xa0; | &#xa0;
  <a href="https://github.com/WenRunning" target="_blank">作者</a>
</p>

<br>

## 简介 

本项目作为无人机云控系统的后端服务器，能够将无人机设备中的飞控等数据转发至前台 Web 端与 App 端，除转发数据的功能外，还以 RESTful 风格开发所用接口给用户端。

## 技术 ##

本项目使用到的开发工具如下:

- [Golang](https://golang.google.cn)

## 前置条件 ##

在运行项目之前，需要配置[Git](https://git-scm.com)， [Golang](https://golang.google.cn)， [Docker](https://www.docker.com)。

```bash
# docker 安装
$ curl -fsSL https://get.docker.com -o get-docker.sh
$ sudo sh get-docker.sh

# 启动Docker服务
$ sudo systemctl start docker
$ sudo systemctl enable docker

# 配置 mqtt 服务器
# dockerhub 需要配置镜像才可以访问，可以在/etc/docker/daemon.json目录下配置，或docker pull时指定 --registry-mirror参数
$ docker pull emqx/emqx:latest
$ docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx:latest


```

## 开始 ##

```bash
# 克隆该项目
$ git clone https://github.com/WenRunning/zht_cloud_server
$ cd zht_cloud_server

# 运行
$ make start
```

<a href="#top">Back to top</a>
