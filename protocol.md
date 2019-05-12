# 协议
这里讲讲解一些内部的运行流程相关的东西
## 鉴权
如果后台鉴权成功，请往网关的`/login`服务发送鉴权相关信息，该API将会发布证书。请自行返回给用户

鉴权相关信息结构如下：
```json
{
  "secret":"this is secret", 
  "U": {
      "ID": 1,
      "username": "username"
  }
  
}
```
`secret`值务必和`conf.yaml`下的配置一样
## 路由
进入路由后，会检查配置。查看是否再白名单（待实现），如果不在。

将会解析`X-Headers-Session`，中的字段，如果不存在或者解析失败，则会报错

如果解析成功，所有请求将会封装成特殊的结构已`POST`请求发送到`conf.d/xxx.yaml`的中绑定的`backend`字段所指示的地址中。

其中发送地址的规则是 `backend` + `route`，如果`route`带有特殊规则，则会替换成匹配规则后的。

特殊结构如下：
```json
{
  "userinfo": {
    "ID": 1,
    "username": "username"
  },
  "method": "GET",
  "requestheader": {},
  "requestbody": "base63.Encode",
  "RequestForm": "aaa=12345&bbb=123456"
}
``` 
`method`是原始请求的请求方式，`requestheader`是原始请求Header字段的json版，但是内容是动态的。