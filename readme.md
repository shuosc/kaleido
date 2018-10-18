# kaleido

> A high-speed http redirector for open-source mirror site.

本工具将会根据您的IP所处在的地理位置，将您的请求自动定向到较近的镜像站。

![](https://camo.githubusercontent.com/320706ea50cf1c2ebff0ea31c82fd9dcad4a3f4c/68747470733a2f2f692e763265782e636f2f513174526f30516e2e706e67)

## How to use

将对应的官方源的URL中的域名改成mirrors.rocks即可。

以Ubuntu为例，更改`/etc/apt/sources.list`文件中Ubuntu 默认的源地址 <http://archive.ubuntu.com/> 为 http://mirrors.rocks 即可。

⚠️：暂时没有HTTPS的支持，对于有些镜像（如pip）请强制允许使用http协议。

