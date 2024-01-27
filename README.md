<div align="center">

# coze-discord-proxy

_coze+discord 代理服务—通过接口调用被`coze`托管的`discord-bot`_

_觉得有点意思的话 别忘了点个🌟_

🐞<a href="https://t.me/+LGKwlC_xa-E5ZDk9" style="font-size: 15px;">COZE-DISCORD-PROXY交流群</a>

</div>

## 功能

- [x] 对话支持流式返回
- [x] 支持文生图（需coze配置DALL·E3插件）返回图片url
- [x] 支持创建 `discord`频道/子频道/线程
- [x] 支持对话指定 `discord`频道/子频道/线程 实现对话隔离
- [x] 支持和openai对齐的接口,用于集成NextChat、one-api等

### 接口文档:

`http://<ip>:<port>/swagger/index.html`

<span><img src="docs/img.png" width="800"/></span>

### 示例:

<span><img src="docs/img2.png" width="800"/></span>

## 如何使用

1. 打开 [discord开发者平台](https://discord.com/developers/applications) 。
2. 创建bot-A,并记录bot专属的`token`和`id(COZE_BOT_ID)`,此bot为被coze托管的bot。
3. 创建bot-B,并记录bot专属的`token(BOT_TOKEN)`,此bot为我们与discord交互的bot。
4. 两个bot开通对应权限(`Send Messages`,`Read Message History`等)并邀请进服务器,记录服务器ID(`GUILD_ID`) (
   过程不在此赘述)。
5. 打开 [coze官网](https://www.coze.com) 创建自己bot。
6. 创建好后推送，配置discord-bot的`token`,即bot-A的`token`,点击完成后在discord的服务器中可看到bot-A在线并可以@使用。
7. 配置环境变量，并启动本项目。
8. 访问接口地址即可开始调试。

## 如何集成one-api



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
      - ./data/coze-discord-proxy:/data
    environment:
      - BOT_TOKEN=MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w  # 必须修改为我们主动发送消息的Bot-Token
      - GUILD_ID=119xxxxxxxx796  # 必须修改为两个机器人所在的服务器ID
      - COZE_BOT_ID=119xxxxxxxx7  # 必须修改为由coze托管的机器人ID
      - PROXY_SECRET=123456  # [可选]修改此行为请求头校验的值（前后端统一）
      - CHANNEL_ID=119xxxxxx24  # [可选]默认频道-在使用与openai对齐的接口时(/v1/chat/completions) 消息会默认发送到此频道
      - TZ=Asia/Shanghai
```

### 基于 Docker 进行部署

```shell
docker run --name coze-discord-proxy -d --restart always \
-p 7077:7077 \
-e BOT_TOKEN="MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w" \
-e GUILD_ID="119xxxxxxxx796" \
-e COZE_BOT_ID="119xxxxxxxx7" \
-e PROXY_SECRET="123456" \
-e CHANNEL_ID="119xxxxxx24" \
-e TZ=Asia/Shanghai \
deanxv/coze-discord-proxy
```

其中，`BOT_TOKEN`,`GUILD_ID`,`COZE_BOT_ID`,`PROXY_SECRET`,`CHANNEL_ID`修改为自己的。

## 配置

### 环境变量

1. `BOT_TOKEN:MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w`  主动发送消息的Bot-Token
2. `GUILD_ID:119xxxxxxxx796`  两个机器人所在的服务器ID
3. `COZE_BOT_ID:119xxxxxxxx7` 由coze托管的机器人ID
4. `PROXY_SECRET:123456` [可选]请求头校验的值（前后端统一）,配置此参数后，每次发起请求时请求头加上`proxy-secret`
   参数，即`header`中添加 `proxy-secret：123456`
5. `CHANNEL_ID:119xxxxxx24`  # [可选]默认频道-在使用与openai对齐的接口时(/v1/chat/completions) 消息会默认发送到此频道

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=deanxv/coze-discord-proxy&type=Date)](https://star-history.com/#deanxv/coze-discord-proxy&Date)

## 其他

Coze 官网 : https://www.coze.com

Discord 开发地址 : https://discord.com/developers/applications



