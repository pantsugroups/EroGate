#EroGate ---- Ero的網關入口
使用了Echo框架

詳細文檔地址請訪問：https://echo.labstack.com/

# 协议
然而本Gateway使用了另外一套框架，详情请访问[协议](protocol.md)

# TODO
  1. HTTP错误信息还未完善
  2. 路由规则还未完善,白名单机制等

# HOT TO USE?
編譯：
```bash
./install.sh || ./install.bat
```
## 第一次运行
第一次运行，必须使用命令`gateway setup`进行安装。这个命令将会在目录下生成一个`conf.yaml`总配置文件。

该文件内容如下：
```yaml
base:
  secret: this is a secret
  port: 80
route:
  login: /login
  backend: http://127.0.0.1:5000/
```
其中，`port`是当前gateway监听的端口。`login`是后台UserAPI的鉴权地址，
`secret`是用于token生成，以及gateway判断后台UserAPI鉴权成功的字段。UserAPI鉴权成功后，必须要返回该值内容，gateway才会签发令牌

## 添加路由
使用`gateway add`命令，会在`/conf.d`下面動態添加yaml文件。這個作爲路由使用，詳細結構如下：
```yaml
route: /website
backend: https://localhost:8080/
```

可以使用`route`标签和`backend`标签手动指明内容，同时还需要一个额外的标签`name`作为文件名
例如：
```bash
gateway add --name config1 --route /site --route http://localhost/
```

之后将会生成一个`/conf.d/config1.yaml`，如果程序在运行则`/site`将会添加到处理事件

## 运行
使用`gateway run`即可运行，但是前提得保证先执行了`gateway setup`


## 测试
使用`gateway test`命令，无条件的签发令牌或者解析令牌。

签发令牌需要`username`字段和`uid`字段，例如
```bash
gateway test --username aaa --uid 1
```
将会返回令牌：
```bash
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NTczMzE5MjEsImlkIjoxLCJuYmYiOjE1NTczMzE5MjEsInVzZXJuYW1lIjoiYWFhIn0.ZphjRwwfcVrvKevwkA1FESMGpWjZbaECUkqInlEKZNc
```

如果需要解析令牌，则需要使用`token`字段
```bash
gateway test --token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NTczMzE5MjEsImlkIjoxLCJuYmYiOjE1NTczMzE5MjEsInVzZXJuYW1lIjoiYWFhIn0.ZphjRwwfcVrvKevwkA1FESMGpWjZbaECUkqInlEKZNc
```
将会返回
`name: aa uid: 1`