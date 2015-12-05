Beego Blog
=====================


###Ubuntu下配置和运行本博客指南
(Windows上基本类似)

######1 安装配置Golang 
[参考来源](https://golang.org:https://golang.org/doc/install)

- 1.1 移除已经存在的老版本(安装的过程反过来做)

- 1.2 下载golang压缩包[fileName] 
[下载地址](https://golang.org/dl/)

- 1.3 解压缩安装
<pre>
sudo tar -C /usr/local -xzf  [fileName]        //[fileName]是官网下载的golang压缩包
</pre>

- 1.4 配置环境变量
<pre>
sudo gedit /etc/profile
</pre>
<pre>
添加配置到profile文件末尾并保存：
export GOPATH=[yourOwnGopath]         //[yourOwnGopath]是自己选择的golang项目目录
export GOBIN=/usr/local/go/bin:${GOPATH//://bin:}/bin 　//"/usr/local/go/bin"是golang安装目录下的bin目录
export PATH=$PATH:$GOBIN
</pre>


######2 安装配置Beego Blog

- 2.1 Git安装
<pre>
sudo apt-get install git
</pre>

- 2.2 克隆本博客
<pre>
cd $GOPATH/src
git clone https://github.com/hopehook/beegoblog.git
</pre>

- 2.3 获取外部依赖包
<pre>
go get github.com/astaxie/beego
go get github.com/Unknwon/com
go get github.com/mattn/go-sqlite3
go get github.com/garyburd/redigo/redis   
</pre>

######3 运行本博客

- 3.1 编译和运行
<pre>
go build $GOPATH/src/beegoblog/main.go
./main
</pre>

- 3.2 浏览器查看
<pre>
firefox http://localhost:8080
</pre>


######4 启用Redis替换SQLite

- 4.1 安装Redis
<pre>
sudo apt-get update
sudo apt-get install redis-server
</pre>

- 4.2 运行Redis
<pre>
sudo redis-server
sudo redis-cli
</pre>

- 4.3 配置Redis
<pre>
sudo CONFIG SET requirepass 123
</pre>

- 4.4 替换SQLite
<pre>
sudo gedit $GOPATH/src/beegoblog/conf/app.conf
</pre>
<pre>
添加配置到app.conf文件末尾并保存：
database = redis
requirepass = 123   //Redis验证密码，与4.3保持一致即可
</pre>

</br>
</br>

### 本博客技术架构介绍

######0 总体
<pre>
html/JS/CSS + golang + DB
</pre>

######1 UI
<pre>
AngularJS;jQuery
</pre>

######2 后台
<pre>
Beego
</pre>

######3 数据库
<pre>
SQLite Or Redis
</pre>




