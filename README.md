<div align="center">

# coze-discord-proxy

_代理`Discord`对话`Coze-Bot`，实现API形式请求GPT4对话模型/微调模型_

_觉得有点意思的话 别忘了点个🌟_

🐞<a href="https://t.me/+LGKwlC_xa-E5ZDk9" style="font-size: 15px;">COZE-DISCORD-PROXY交流群</a>(群内有详细教程)

</div>

## 功能

- [x] 完美适配`NextChat`,`one-api`,`LobeChat`等对话面板。
- [x] 完美支持对话隔离。
- [x] 对话接口支持流式返回。
- [x] 支持创建 `discord`分类/频道/线程。
- [x] 支持和`openai`对齐的对话接口(`v1/chat/completions`)(也支持`dall-e-3`文生图)
- [x] 支持和`openai`对齐的图/文件生文接口(`v1/chat/completions`)(按照`GPT4V`
  图/文件生文接口的请求格式 [ 支持`url`或`base64` ])。
- [x] 支持和`openai`对齐的`dall-e-3`文生图接口(`v1/images/generations`)。
- [x] 支持每日`24`点定时任务自动活跃机器人。
- [x] 支持配置多机器人 (通过`PROXY_SECRET`/`model`指定) 详细请看[进阶配置](#进阶配置)。

### 接口文档:

`http://<ip>:<port>/swagger/index.html`

<span><img src="docs/img.png" width="800"/></span>

### 示例:

<span><img src="docs/img2.png" width="800"/></span>

## 如何使用

1. 打开 [discord开发者平台](https://discord.com/developers/applications) 。
2. 创建bot-A,并记录bot专属的`token`和`id(COZE_BOT_ID)`,此bot为被coze托管的bot。
3. 创建bot-B,并记录bot专属的`token(BOT_TOKEN)`,此bot为我们用来监听discord消息的bot。
4. 两个bot开通对应权限(`Administrator`)并邀请进服务器,记录服务器ID(`GUILD_ID`) (
   过程不在此赘述)。
5. 打开F12发送一次消息,在`https://discord.com/api/v9/channels/1206*******703/messages`
   接口header中获取`Authorization(USER_AUTHORIZATION)`。
6. 在discord中打开开发者模式，右键自己的用户名获取`用户Id(USER_ID)`。
7. 打开 [coze官网](https://www.coze.com) 创建自己bot。
8. 创建好后推送(`Auto-Suggestion`为`default`),配置discord-bot的`token`,即bot-A的`token`
   ,点击完成后在discord的服务器中可看到bot-A在线并可以@使用。
9. 配置环境变量,并启动本项目。
10. 访问接口地址即可开始调试。

## 如何集成NextChat

填 接口地址(ip:端口/域名) 及 API-Key(`PROXY_SECRET`),其它的随便填随便选。

> 如果自己没有搭建NextChat面板,这里有个已经搭建好的可以使用 [NextChat](https://ci.goeast.io/)

<span><img src="docs/img5.png" width="800"/></span>

## 如何集成one-api

填 `BaseURL`(ip:端口/域名) 及 密钥(`PROXY_SECRET`),其它的随便填随便选。

<span><img src="docs/img3.png" width="800"/></span>

## 部署

### 基于 Docker-Compose(All In One) 进行部署

```shell
docker-compose pull && docker-compose up -d
```

#### docker-compose.yml

```docker
version: '3.4'

services:
  coze-discord-proxy:
    image: deanxv/coze-discord-proxy:latest
    container_name: coze-discord-proxy
    restart: always
    ports:
      - "7077:7077"
    volumes:
      - ./data:/app/coze-discord-proxy/data
    environment:
      - USER_ID=1099*********055  # 必须修改为我们discord用户的ID
      - USER_AUTHORIZATION=MTA5OTg5N************uIfytxUgJfmaXUBHVI  # 必须修改为我们discord用户的授权密钥
      - BOT_TOKEN=MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w  # 必须修改为监听消息的Bot-Token
      - GUILD_ID=119xxxxxxxx796  # 必须修改为两个机器人所在的服务器ID
      - COZE_BOT_ID=119xxxxxxxx7  # 必须修改为由coze托管的机器人ID
      - CHANNEL_ID=119xxxxxx24  # 默认频道-(目前版本下该参数仅用来活跃机器人)
      - PROXY_SECRET=123456  # [可选]接口密钥-修改此行为请求头校验的值(多个请以,分隔)
      - TZ=Asia/Shanghai
```

### 基于 Docker 进行部署

```shell
docker run --name coze-discord-proxy -d --restart always \
-p 7077:7077 \
-v $(pwd)/data:/app/coze-discord-proxy/data \
-e USER_ID="1099*********055"  \
-e USER_AUTHORIZATION="MTA5OTg5N************uIfytxUgJfmaXUBHVI" \
-e BOT_TOKEN="MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w" \
-e GUILD_ID="119xxxxxxxx796" \
-e COZE_BOT_ID="119xxxxxxxx7" \
-e PROXY_SECRET="123456" \
-e CHANNEL_ID="119xxxxxx24" \
-e TZ=Asia/Shanghai \
deanxv/coze-discord-proxy
```

其中,`USER_ID`,`USER_AUTHORIZATION`,`BOT_TOKEN`,`GUILD_ID`,`COZE_BOT_ID`,`PROXY_SECRET`,`CHANNEL_ID`修改为自己的。

如果上面的镜像无法拉取，可以尝试使用 GitHub 的 Docker 镜像，将上面的 `deanxv/coze-discord-proxy`
替换为 `ghcr.io/deanxv/coze-discord-proxy` 即可。

### 部署到第三方平台

<details>
<summary><strong>部署到 Zeabur</strong></summary>
<div>

> Zeabur 的服务器在国外,自动解决了网络的问题,同时免费的额度也足够个人使用

点击一键部署:

[![Deploy on Zeabur](https://zeabur.com/button.svg)](https://zeabur.com/templates/GMU8C8?referralCode=deanxv)

**一键部署后 `USER_ID`,`USER_AUTHORIZATION`,`BOT_TOKEN`,`GUILD_ID`,`COZE_BOT_ID`,`PROXY_SECRET`,`CHANNEL_ID`变量也需要替换！**

或手动部署:

1. 首先 **fork** 一份代码。
2. 进入 [Zeabur](https://zeabur.com?referralCode=deanxv),使用github登录,进入控制台。
3. 在 Service -> Add Service,选择 Git（第一次使用需要先授权）,选择你 fork 的仓库。
4. Deploy 会自动开始,先取消。
5. 添加环境变量

   `USER_ID:1099*********055`  主动发送消息的discord用户的ID

   `USER_AUTHORIZATION:MTA5OTg5N************uIfytxUgJfmaXUBHVI`  主动发送消息的discord用户的授权密钥

   `BOT_TOKEN:MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w`  监听消息的Bot-Token

   `GUILD_ID:119xxxxxxxx796`  两个机器人所在的服务器ID

   `COZE_BOT_ID:119xxxxxxxx7` 由coze托管的机器人ID

   `CHANNEL_ID:119xxxxxx24`  # 默认频道-(目前版本下该参数仅用来活跃机器人)

   `PROXY_SECRET:123456` [可选]接口密钥-修改此行为请求头校验的值(多个请以,分隔)(与openai-API-KEY用法一致)

保存。

6. 选择 Redeploy。

</div>


</details>

<details>
<summary><strong>部署到 Render</strong></summary>
<div>

> Render 提供免费额度,绑卡后可以进一步提升额度

Render 可以直接部署 docker 镜像,不需要 fork 仓库：[Render](https://dashboard.render.com)

</div>
</details>

## 配置

### 环境变量
1. `USER_ID:1099*********055`  主动发送消息的discord用户的ID
2. `USER_AUTHORIZATION:MTA5OTg5N************uIfytxUgJfmaXUBHVI`  主动发送消息的discord用户的授权密钥
3. `BOT_TOKEN:MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w`  监听消息的Bot-Token
4. `GUILD_ID:119xxxxxxxx796`  两个机器人所在的服务器ID
5. `COZE_BOT_ID:119xxxxxxxx7`  由coze托管的机器人ID
6. `CHANNEL_ID:119xxxxxx24`  默认频道-(目前版本下该参数仅用来活跃机器人)
7. `CHANNEL_AUTO_DEL_TIME:60`  [可选]频道自动删除时间(秒) 此参数为每次对话完成后自动删除频道的时间(默认为5s)
   ,为0时则不删除,推荐不使用此环境变量
8. `COZE_BOT_STAY_ACTIVE_ENABLE:1`  [可选]是否开启每日`24`点活跃coze-bot的定时任务,默认开启,为0时则不开启,推荐不使用此环境变量
9. `PORT:7077`  [可选]端口,默认为7077
10. `PROXY_SECRET:123456`  [可选]接口密钥-修改此行为请求头校验的值(多个请以,分隔)(与openai-API-KEY用法一致)
11. `REQUEST_OUT_TIME:60`  [可选]对话接口非流响应下的请求超时时间,推荐不使用此环境变量
12. `STREAM_REQUEST_OUT_TIME:60`  [可选]对话接口流响应下的每次流返回超时时间,推荐不使用此环境变量
13. `USER_AGENT:Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36`  [可选]discord用户端Agent,使用自己的可能有效防止被ban，不设置时默认使用作者的 推荐使用此环境变量
14. `PROXY_URL:http://127.0.0.1:10801`  [可选]代理

## 进阶配置

### 配置多机器人

1. 部署前在`docker`/`docker-compose`部署同级目录下创建`data/config/bot_config.json`文件
2. 编写该`json`文件,`bot_config.json`格式如下

```shell
[
  {
    "proxySecret": "123", // 接口请求密钥(PROXY_SECRET)
    "cozeBotId": "12***************31", // coze托管的机器人ID
    "model": "GPT-3.5-16k", // coze托管的机器人模型名称(与请求参数中的model对应,如请求中的model在该json中未匹配到则会抛出异常)
    "channelId": "12***************56"  // [可选]discord频道ID(机器人必须在此频道所在的服务器)(目前版本下该参数仅用来活跃机器人)
  },
  {
    "proxySecret": "456",
    "cozeBotId": "12***************64",
    "model": "GPT-4-8k", 
    "channelId": "12***************78"
  },
  {
    "proxySecret": "789",
    "cozeBotId": "12***************12",
    "model": "GPT-4-Turbo-128k",
    "channelId": "12***************24"
  }
]
```

3. 重启服务

> 当有此配置时,会通过请求头携带的[请求密钥]+请求体中的[`model`]匹配此配置中的`cozeBotId`
> ,若匹配到多个则随机选择一个,所以当存在多用户使用时可对每个用户分发独立的请求密钥。配置很灵活,可以根据自己的需求进行配置。

第三方平台(如: `zeabur`)部署的服务需要[配置多机器人]
请参考[issue#30](https://github.com/deanxv/coze-discord-proxy/issues/30)

## Q&A

Q: 我们如何使用该服务托管多个Bot去请求多个由coze托管的Bot？

A: 首先用不同的端口部署多个`coze-discord-proxy`服务,对每个服务都[配置多机器人](#配置多机器人)
,并对每个服务设置不同的`BOT_TOKEN`,再部署[one-api](https://github.com/songquanpeng/one-api)
后[添加多个渠道](#如何集成one-api),利用[one-api](https://github.com/songquanpeng/one-api)
的轮询去请求我们的`coze-discord-proxy`服务。

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=deanxv/coze-discord-proxy&type=Date)](https://star-history.com/#deanxv/coze-discord-proxy&Date)

## 其他版本

**开源不易,若你参考此项目或基于此项目二开可否麻烦在你的项目文档中标识此项目呢？谢谢你！♥♥♥**

Java: https://github.com/oddfar/coze-discord

## 其他引用

Coze 官网 : https://www.coze.com

Discord 开发地址 : https://discord.com/developers/applications



