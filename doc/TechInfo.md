# 技术信息

本文档是为对kaleido背后技术信息感兴趣，准备阅读甚至贡献代码的人准备的。

## 技术栈

kaleido主要使用Go语言开发。

使用Postgres数据库来进行数据持久化。

我们幸运地从[这个网站](http://ipcn.chacuo.net)上找到了IP段与地区和运营商的对应关系。

IP地址 与 镜像名称 到 镜像站url的对应关系表为protobuf格式，放在阿里云OSS上，这样做是为了避免Master宕掉导致整个服务中断。

Dashboard部分使用Vue开发，UI框架为Vuetify。

Dashboard和Master部分的数据交换采用GraphQL。

## Big picture

kaleido分为三个部分：

- Master，负责
  - 定期检测某个镜像站是否还在运行
  - 定期检测某个镜像站有哪些镜像
  - 如果上述信息有变化，会制作 IP地址 与 镜像名称 到 镜像站url 的对应关系表
  - 提供查询镜像/镜像站情况的web接口给Dashboard
- Node，负责
  - 如果发现OSS上的对应关系表比目前的新，则从线上拉取新的关系表
  - 接受用户请求，根据关系表将用户的请求重定向到较为靠近的，有用户所请求的镜像的镜像站
- Dashboard，负责
  - 展示镜像/镜像站的情况

## Master

Master有两个模块：

- 爬虫模块，负责检测镜像站情况

  针对每个数据库中的镜像站，每过10s就会访问其包含镜像信息的页面（主页或json数据源页面），访问失败则判断为这个镜像站挂掉了，将其Alive字段标记为false，否则将其Alive字段标记为true，分析页面上的信息来确定这个镜像站有哪些镜像可用，并写入数据库。

  如果镜像站的信息（是否Alive、对应镜像站列表）变更，会告知关系表制作模块。

- 对应关系表制作模块

  如果收到镜像站的信息变更信息，则会制作新的对应关系表，并上传OSS。

## 对应关系表

对应关系表使用protobuf格式，代码及解释如下：

```protobuf
message Address_AreaId {
	// mask过的IP地址 -> 地区Id
    map<uint32, uint32> Address_AreaId = 1;
}

message MirrorStationGroup {
	// 一组镜像站，常用来表示几个距离某个地区距离相同的镜像站
    repeated uint32 Stations = 1;
}

message Mirror {
	// 某个镜像的Fallback 镜像站 Id
    uint32 DefaultMirrorStationId = 1;
    // 某个地区的ID->将这个地区对这一Mirror的请求重定向到的镜像站组的 Id
    map<uint32, MirrorStationGroup> AreaId_MirrorStationGroup = 2;
}

message KaleidoMessage {
	// 所有镜像的列表
    map<string, Mirror> Mirrors = 1;
    // mask长度->mask过的IP地址->地区Id
    // protobuf不支持map的key为message类型，否则应该写作IP(mask,ip)->AreaId
    map<uint32, Address_AreaId> Mask_Address_AreaID = 2;
    // 镜像站Id->镜像站根URL
    map<uint32, string> MirrorStationId_Url = 3;
}
```

## Node

Node也有两个模块：

- 检测并更新重定向表模块

  每10s检测OSS上的重定向表是否有变（通过获取其Last-Modified属性），有变则下载更新整个重定向表

- 对外服务模块

  对用户发送过来的对某个镜像的请求，根据用户的IP地址，查询到其所在的地区，再重定向到有这个镜像的最近镜像站，如果没有查询到其所在地区（如境外用户、IPV6[^1]用户），则会定向到这个镜像的Fallback 镜像站

## 数据源及距离计算方法

- Area这个词是指的是“网络上的区域”，包含地理区域和网络运营商两方面的信息。

我们从[这里](http://ipcn.chacuo.net)爬到了某个运营商和省市对应的IP段。

我们直接使用[百度地图接口](http://api.map.baidu.com/geocoder?address=上海&output=json&key=37492c0ee6f924cb5e934fa08c6b1676&city=北京市)来获取省市的经纬度信息，并使用球面距离公式计算出两个省市间的距离。

```
两个Area之间的距离=两个地理省市间的距离+两个Area同一个运营商?0:10;
```

最终Area之间的距离信息在数据库中构成一张类似邻接表的表，供Master使用。

[^1]: IPV6的支持需要IPV6对应的地区数据，我们目前没有找到这样的数据源