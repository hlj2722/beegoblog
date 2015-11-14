Beego Blog学习实践
=====================


###Ubuntu下配置和运行本博客指南

######0,安装Golang，配置环境变量 
[参见](golang.org:https://golang.org/doc/install)
- 0.1 移除已经存在的老版本
- 0.2 下载golang压缩包[fileName] 
[下载地址](https://golang.org/dl/)
- 0.3 解压缩 
<pre>
终端命令：
tar -C /usr/local -xzf  [fileName]        //[fileName]是官网下载的golang压缩包
</pre>
- 0.4 配置环境变量
<pre>
终端命令：
sudo gedit /etc/profile
</pre>
<pre>
添加配置：
export PATH=$PATH:/usr/local/go/bin
export GOPATH=[yourOwnGopath]         //[yourOwnGopath]是自己选择的golang项目目录
export GOBIN=$GOPATH/bin
</pre>


######1,克隆本工程文件到本地的GOPATH目录

######2,获取外部依赖包
<pre>
终端命令：
go get github.com/Unknwon/com
go get github.com/astaxie/beego
</pre>

######3,运行本博客
<pre>
终端命令：
go build $GOPATH/src/beegoblog/main.go
./main
</pre>

######4,浏览器查看
<pre>
访问地址：
http://localhost:8080
</pre>



### 本博客技术架构介绍

######0，UI
<pre>
AngularJS
</pre>

######1,后台
<pre>
Golang
</pre>

######2,数据库
<pre>
SQLite Or Redis
</pre>

######3，总体框架
<pre>
AngularJS + beego + DB
</pre>


