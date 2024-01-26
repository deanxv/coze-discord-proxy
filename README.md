<div align="center">

# coze-discord-proxy

_coze+discord 代理服务—通过接口调用被`coze`托管的`discord-bot`_

_觉得有点意思的话 别忘了点个🌟_
</div>

## 功能

接口文档: `http://<ip>:<port>/swagger/index.html`

<span><img src="docs/img.png" width="500"/></span>

## 部署

### 基于 Docker-Compose(All In One) 进行部署

```shell
docker-compose pull && docker-compose up -d
```

#### docker-compose.yml

```docker
version: '3.4'

services:
  code-discord-proxy:
    image: deanxv/code-discord-proxy:latest
    container_name: code-discord-proxy
    restart: always
    ports:
      - "7077:7077"
    volumes:
      - ./data/code-discord-proxy:/data
    environment:
      - BOT_TOKEN=MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w  # 必须修改为我们主动发送消息的Bot-Token
      - GUILD_ID=119xxxxxxxx796  # 必须修改为两个机器人所在的服务器ID
      - COZE_BOT_ID=119xxxxxxxx7  # 必须修改为由coze托管的机器人ID
      - PROXY_SECRET=123456  # 修改此行为请求头校验的值（前后端统一）
      - TZ=Asia/Shanghai
```

### 基于 Docker 进行部署

```shell
docker run --name code-discord-proxy -d --restart always \
-p 7078:7077 \
-e BOT_TOKEN="MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w" \
-e GUILD_ID="119xxxxxxxx796" \
-e COZE_BOT_ID="119xxxxxxxx7" \
-e PROXY_SECRET="123445" \
deanxv/code-discord-proxy
```

其中，`BOT_TOKEN`,`GUILD_ID`,`COZE_BOT_ID`,`PROXY_SECRET`修改为自己的。

## 配置

### 环境变量

1. `COZE_BOT_ID：MTE5OTk2xxxxxxxxxxxxxxrwUrUWNbG63w`  主动发送消息的Bot-Token
2. `GUILD_ID：119xxxxxxxx796`  两个机器人所在的服务器ID
3. `COZE_BOT_ID：119xxxxxxxx7` 由coze托管的机器人ID
4. `PROXY_SECRET`:`123456` [可选]请求头校验的值（前后端统一）

## 其他

Coze 官网 : https://www.coze.com

Discord 开发地址 : https://discord.com/developers/applications

