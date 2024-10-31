
<p align="right">
   <strong>English</strong> | <a href="./README.zh.md">中文</a>
</p>

<div align="center">

# coze-discord-proxy

<a href="https://trendshift.io/repositories/7350" target="_blank"><img src="https://trendshift.io/api/badge/repositories/7350" alt="deanxv%2Fcoze-discord-proxy | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>

_Proxies `Discord` conversations to `Coze-Bot`, enabling API requests to the GPT4 model, offering functionalities such as conversation, text-to-image, image-to-text, and knowledge base retrieval._

_If you find this interesting, don't forget to give it a 🌟_

📄<a href="https://cdp-docs.pages.dev" style="font-size: 15px;">CDP Project Documentation</a> (Essential Tutorial)

🐞<a href="https://t.me/+LGKwlC_xa-E5ZDk9" style="font-size: 15px;">CDP Project - Discussion Group</a> (Discussion)

📢<a href="https://t.me/+0fYkYY_zUZYzNzRl" style="font-size: 15px;">CDP Project - Notification Channel</a> (Notifications)

</div>

## Features (The project features are now stable, iterations will not be frequent, feel free to raise an issue if you find bugs!)

- [x] Perfectly compatible with `NextChat`, `one-api`, `LobeChat`, and other chat panels.
- [x] Fully supports conversation isolation.
- [x] Conversation API supports streaming responses.
- [x] Supports creating `discord` categories/channels/threads.
- [x] Supports conversation interfaces aligned with `openai` (`v1/chat/completions`) (also supports `dall-e-3` text-to-image) (supports specifying `discord-channel`).
- [x] Supports image-to-text/image-to-image/file-to-text interfaces aligned with `openai` (`v1/chat/completions`) (using the `GPT4V` request format [supports `url` or `base64`]) (supports specifying `discord-channel`).
- [x] Supports `dall-e-3` text-to-image interface aligned with `openai` (`v1/images/generations`).
- [x] Supports daily tasks at `9` AM to keep the bot active.
- [x] Supports configuring multiple Discord user `Authorization` (environment variable `USER_AUTHORIZATION`) for request load balancing (**currently, each Discord user's use of coze-bot is limited to 24 hours, configure multiple users to stack request limits and load balance**).
- [x] Supports configuring multiple coze bots for response load balancing (specified through `PROXY_SECRET`/`model`) see [Advanced Configuration](#advanced-configuration) for details.

### API Documentation:

`http://<ip>:<port>/swagger/index.html`

<span><img src="docs/img.png" width="800"/></span>

### Example:

<span><img src="docs/img2.png" width="800"/></span>

## How to Use

1. Open [Discord's official website](https://discord.com/app), log in, click settings -> advanced settings -> developer mode -> turn on.
2. Create a Discord server, right-click this server to choose `Copy Server ID (GUILD_ID)` and record it, create a default channel in this server, right-click this channel to choose `Copy Channel ID (CHANNEL_ID)` and record it.
3. Open [Discord Developer Portal](https://discord.com/developers/applications) and log in.
4. Create a new application -> Bot named `COZE-BOT`, and record its exclusive `token` and `id (COZE_BOT_ID)`, this bot will be managed by coze.
5. Create a new application -> Bot named `CDP-BOT`, and record its exclusive `token (BOT_TOKEN)`, this bot will listen for Discord messages.
6. Grant appropriate permissions (`Administrator`) to both bots and invite them to the created Discord server (this process is not described here).
7. Open [Discord's official website](https://discord.com/app), enter the server, press F12 to open developer tools, send a message in any channel, find the request `https://discord.com/api/v9/channels/1206*******703/messages` in the developer tools-`Network`, get `Authorization (USER_AUTHORIZATION)` from the header of this interface and record it.
8. Open [coze's official website](https://www.coze.com) to create a bot and configure it (note `Auto-Suggestion` should be `Default/on` (usually no need to change)).
9. After configuration, choose to publish to Discord, fill in the `COZE-BOT`'s `token`, after publishing, you can see `COZE-BOT` online in the Discord server and can be used with @.
10. Use the recorded parameters to start configuring [environment variables](#environment-variables) and [deploy](#deployment) this project.
11. Access the API documentation address, and you can start debugging or integrating other projects.

## How to Integrate with NextChat

Fill in the interface address (ip:port/domain) and API-Key (`PROXY_SECRET`), other fields are optional.

> If you haven't set up a NextChat panel yourself, here's one that's already set up and available for use: [NextChat](https://ci.goeast.io/)

<span><img src="docs/img5.png" width="800"/></span>

## How to Integrate with one-api

Fill in `BaseURL` (ip:port/domain) and the key (`PROXY_SECRET`), other fields are optional.

<span><img src="docs/img3.png" width="800"/></span>

## Deployment

### Deploying with Docker-Compose (All In One)

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
      - USER_AUTHORIZATION=MTA5OTg5N************aXUBHVI  # Must be modified to your Discord user's authorization key (multiple keys separated by commas)
      - BOT_TOKEN=MTE5OT************UrUWNbG63w  # Must be modified to the listening bot's token
      - GUILD_ID=11************96  # Must be modified to the server ID where both bots are located
      - COZE_BOT_ID=11************97  # Must be modified to the bot ID managed by coze
      - CHANNEL_ID=11************94  # [Optional] Default channel - (in the current version, this parameter is only used to keep the bot active)
      - PROXY_SECRET=123456  # [Optional] API key - modify this line to the value checked in the request header (multiple keys separated by commas)
      - TZ=Asia/Shanghai
```

### Deploying with Docker

```docker
docker run --name coze-discord-proxy -d --restart always \
-p 7077:7077 \
-v $(pwd)/data:/app/coze-discord-proxy/data \
-e USER_AUTHORIZATION="MTA5OTg5N************uIfytxUgJfmaXUBHVI" \
-e BOT_TOKEN="MTE5OTk2************rUWNbG63w" \
-e GUILD_ID="11************96" \
-e COZE_BOT_ID="11************97" \
-e PROXY_SECRET="123456" \
-e CHANNEL_ID="11************24" \
-e TZ=Asia/Shanghai \
deanxv/coze-discord-proxy
```

Modify `USER_AUTHORIZATION`, `BOT_TOKEN`, `GUILD_ID`, `COZE_BOT_ID`, `PROXY_SECRET`, `CHANNEL_ID` to your own values.

If the above image cannot be pulled, you may try using the GitHub Docker image by replacing `deanxv/coze-discord-proxy` with `ghcr.io/deanxv/coze-discord-proxy`.

### Deployment on Third-Party Platforms

<details>
<summary><strong>Deploying on Zeabur</strong></summary>
<div>

> Zeabur's servers are located abroad, automatically solving network issues, and the free tier is sufficient for personal use.

Click to deploy:

[![Deploy on Zeabur](https://zeabur.com/button.svg)](https://zeabur.com/templates/GMU8C8?referralCode=deanxv)

**After one-click deployment, variables like `USER_AUTHORIZATION`, `BOT_TOKEN`, `GUILD_ID`, `COZE_BOT_ID`, `PROXY_SECRET`, `CHANNEL_ID` need to be replaced!**

Or manually deploy:

1. First **fork** a copy of the code.
2. Enter [Zeabur](https://zeabur.com?referralCode=deanxv), log in with GitHub, and go to the console.
3. In Service -> Add Service, select Git (authorization is required for the first use), choose the repository you forked.
4. Deployment will automatically start, cancel it first.
5. Add environment variables:

   `USER_AUTHORIZATION:MTA5OTg5N************uIfytxUgJfmaXUBHVI`  Authorization key for the Discord user initiating messages (separate multiple keys with commas)

   `BOT_TOKEN:MTE5OTk************WNbG63w`  Token for the bot that listens for messages

   `GUILD_ID:11************96`  Server ID where both bots are located

   `COZE_BOT_ID:11************97`  Bot ID managed by coze

   `CHANNEL_ID:11************24`  [Optional] Default channel - (in the current version, this parameter is only used to keep the bot active)

   `PROXY_SECRET:123456` [Optional] API key - modify this line to the value checked in the request header (separate multiple keys with commas) (used in the same way as the openai API-KEY)

Save.

6. Choose Redeploy.

</div>


</details>

<details>
<summary><strong>Deploying on Render</strong></summary>
<div>

> Render offers a free tier, and linking a card can further increase the limit.

Render can directly deploy Docker images without needing to fork the repository: [Render](https://dashboard.render.com)

</div>
</details>

## Configuration

### Environment Variables

1. `USER_AUTHORIZATION=MTA5OTg5N************uIfytxUgJfmaXUBHVI`  Authorization key for the Discord user initiating messages (separate multiple keys with commas)
2. `BOT_TOKEN=MTE5OTk2************rUWNbG63w`  Token for the bot that listens for messages
3. `GUILD_ID=11************96`  Server ID where all bots are located
4. `COZE_BOT_ID=11************97`  Bot ID managed by coze
5. `PORT=7077`  [Optional] Port, default is 7077
6. `SWAGGER_ENABLE=1`  [Optional] Whether to enable Swagger API documentation [0: No; 1: Yes] (default is 1)
7. `ONLY_OPENAI_API=0`  [Optional] Whether to expose only the interfaces aligned with OpenAI [0: No; 1: Yes] (default is 0)
8. `CHANNEL_ID=11************24`  [Optional] Default channel - (in the current version, this parameter is only used to keep the bot active)
9. `PROXY_SECRET=123456`  [Optional] API key - modify this line to the value checked in the request header (separate multiple keys with commas) (used in the same way as the openai API-KEY), **it is recommended to use this environment variable**
10. `DEFAULT_CHANNEL_ENABLE=0`  [Optional] Whether to enable the default channel [0: No; 1: Yes] (default is 0) If enabled, every conversation will take place in the default channel, **session isolation will be ineffective**, **it is recommended not to use this environment variable**
11. `ALL_DIALOG_RECORD_ENABLE=1`  [Optional] Whether to enable full context recording [0: No; 1: Yes] (default is 1) If disabled, each conversation will only send the last `content` of `role` as `user` in `messages`, **it is recommended not to use this environment variable**
12. `CHANNEL_AUTO_DEL_TIME=5`  [Optional] Auto-delete time for channels (seconds) This parameter sets the time to automatically delete the channel after each conversation (default is 5s), 0 means no deletion, **it is recommended not to use this environment variable**
13. `COZE_BOT_STAY_ACTIVE_ENABLE=1`  [Optional] Whether to enable a daily task at `9` AM to keep the coze-bot active [0: No; 1: Yes] (default is 1), **it is recommended not to use this environment variable**
14. `REQUEST_OUT_TIME=60`  [Optional] Request timeout for non-stream response under conversation interface, **it is recommended not to use this environment variable**
15. `STREAM_REQUEST_OUT_TIME=60`  [Optional] Timeout for each stream return under stream response of conversation interface, **it is recommended not to use this environment variable**
16. `REQUEST_RATE_LIMIT=60`  [Optional] Request rate limit per minute per IP, default: 60 requests/min
17. `USER_AGENT=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36`  [Optional] User agent for Discord users, using your own might effectively prevent being banned, if not set, the default is the author's, **it is recommended to use this environment variable**
18. `NOTIFY_TELEGRAM_BOT_TOKEN=6232***********Niz9c`  [Optional] Token for a Telegram Bot used for notifications (notification events: 1. No available `user_authorization`; 2. `BOT_TOKEN` associated BOT triggers risk control)
19. `NOTIFY_TELEGRAM_USER_ID=10******35`  [Optional] The `Telegram-Bot` associated with `NOTIFY_TELEGRAM_BOT_TOKEN` will push to the `Telegram-User` associated with this variable (**`NOTIFY_TELEGRAM_BOT_TOKEN` must not be empty if this variable is set**)
20. `PROXY_URL=http://127.0.0.1:10801`  [Optional] Proxy (supports http only)

## Advanced Configuration

### Configuring Multiple Bots

1. Before deploying, create a `data/config/bot_config.json` file in the same directory as your `docker`/`docker-compose` deployment.
2. Write the following `json` file, `bot_config.json` format as shown below:

```shell
[
  {
    "proxySecret": "123", // API request key (PROXY_SECRET) (Note: this key must exist in the environment variable PROXY_SECRET for this bot to be matched!)
    "cozeBotId": "12***************31", // Bot ID managed by coze
    "model": ["gpt-3.5","gpt-3.5-16k"], // Model names (array format) (matches the `model` in the request body, if the model in the request is not matched in this json, an exception will be thrown)
    "channelId": "12***************56"  // [Optional] Discord channel ID (the bot must be in the server where this channel is located) (in the current version, this parameter is only used to keep the bot active)
  },
  {
    "proxySecret": "456",
    "cozeBotId": "12***************64",
    "model": ["gpt-4","gpt-4-16k"],
    "channelId": "12***************78"
  },
  {
    "proxySecret": "789",
    "cozeBotId": "12***************12",
    "model": ["dall-e-3"],
    "channelId": "12***************24"
  }
]
```

3. Restart the service.

> When this json configuration is present, it will match the `cozeBotId` through the request header's [request key] + request body's [`model`].
> If multiple matches are found, one will be selected at random. The configuration is very flexible and can be adjusted according to your needs.

For services deployed on third-party platforms (such as `zeabur`) that need [configuring multiple bots], please refer to [issue#30](https://github.com/deanxv/coze-discord-proxy/issues/30).

## Limitations

Current details of coze's free and paid subscriptions: https://www.coze.com/docs/guides/subscription?_lang=en

You can configure multiple Discord user `Authorization` (refer to [Environment Variables](#environment-variables) `USER_AUTHORIZATION`) or [configure multiple bots](#configuring-multiple-bots) to stack request limits and achieve load balancing.

## Q&A

Q: How should I configure the service for high concurrency?

A: First, [configure multiple bots](#configuring-multiple-bots) to be used as response bots for load balancing. Next, prepare multiple Discord accounts for request load balancing and invite them to the same server. Obtain `Authorization` for each account, separate them with commas, and configure them in the environment variable `USER_AUTHORIZATION`. Each request will then pick one Discord account to initiate the conversation, effectively achieving load balancing.

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=deanxv/coze-discord-proxy&type=Date)](https://star-history.com/#deanxv/coze-discord-proxy&Date)

## Related

[GPT-Content-Audit](https://github.com/deanxv/gpt-content-audit): Aggregates platforms like OpenAI, Alibaba Cloud, Baidu Intelligent Cloud, Qiniu Cloud, etc., providing content moderation services aligned with OpenAI's request format.

## Other

**Open source is challenging; if you refer to this project or base your project on it, could you please credit this project in your project documentation? Thank you!**

Java: https://github.com/oddfar/coze-discord (currently unavailable
