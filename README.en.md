<p align="right">
   <strong>English</strong> | <a href="./README.md">中文</a>
</p>

<div align="center">

# coze-discord-proxy

<a href="https://trendshift.io/repositories/7350" target="_blank"><img src="https://trendshift.io/api/badge/repositories/7350" alt="deanxv%2Fcoze-discord-proxy | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>

_Proxy `Discord` conversations for `Coze-Bot`, enabling API requests to the GPT4 model with features like conversation, text-to-image, image-to-text, and knowledge base retrieval._

_If you find this interesting, don't forget to give it a 🌟_

📄<a href="https://cdp-docs.pages.dev" style="font-size: 15px;">CDP Project Documentation Site</a> (Must-read tutorial)

🐞<a href="https://t.me/+LGKwlC_xa-E5ZDk9" style="font-size: 15px;">CDP Project - Discussion Group</a> (Discussion)

📢<a href="https://t.me/+0fYkYY_zUZYzNzRl" style="font-size: 15px;">CDP Project - Notification Channel</a> (Notifications)

</div>

## Features (The project's features are now stable, updates will not be frequent, feel free to raise an issue if you find bugs!)

- [x] Perfectly compatible with `NextChat`, `one-api`, `LobeChat` and other conversation panels.
- [x] Perfect support for conversation isolation.
- [x] Conversation interface supports streaming responses.
- [x] Supports creating `discord` categories/channels/threads.
- [x] Supports conversation interface aligned with `openai` (`v1/chat/completions`) (also supports `dall-e-3` text-to-image) (supports specifying `discord-channel`).
- [x] Supports image-to-text/image-to-image/file-to-text interfaces aligned with `openai` (`v1/chat/completions`) (following the `GPT4V` interface request format [ supports `url` or `base64` ])(supports specifying `discord-channel`).
- [x] Supports `dall-e-3` text-to-image interface aligned with `openai` (`v1/images/generations`).
- [x] Supports daily `9 AM` scheduled tasks to keep the bot active.
- [x] Supports configuring multiple discord user `Authorization` (environment variable `USER_AUTHORIZATION`) for request load balancing (**currently each discord user has a 24-hour limit on coze-bot calls, configure multiple users to stack request counts and balance load**).
- [x] Supports configuring multiple coze bots for response load balancing (specified through `PROXY_SECRET`/`model`), see [Advanced Configuration](#advanced-configuration) for details.

### API Documentation:

`http://<ip>:<port>/swagger/index.html`

<span><img src="docs/img.png" width="800"/></span>

### Example:

<span><img src="docs/img2.png" width="800"/></span>

## How to Use

1. Open [Discord's official website](https://discord.com/app), log in, click settings-advanced settings-developer mode-turn on.
2. Create a discord server, right-click this server to select `Copy Server ID (GUILD_ID)` and record it, create a default channel in this server, right-click this channel to select `Copy Channel ID (CHANNEL_ID)` and record it.
3. Open [Discord Developer Portal](https://discord.com/developers/applications) and log in.
4. Create a new application-Bot, i.e., `COZE-BOT`, and record its unique `token` and `id (COZE_BOT_ID)`, this bot will be managed by coze.
5. Create a new application-Bot, i.e., `CDP-BOT`, and record its unique `token (BOT_TOKEN)`, this bot will listen for discord messages.
6. Grant corresponding permissions (`Administrator`) to both bots and invite them to the created discord server (the process is not described here).
7. Open [Discord's official website](https://discord.com/app), enter the server, press F12 to open developer tools, send a message in any channel, find the request `https://discord.com/api/v9/channels/1206*******703/messages` in developer tools-`Network`, get `Authorization (USER_AUTHORIZATION)` from the header of this interface and record it.
8. Open [Coze's official website](https://www.coze.com), create and configure a bot (note `Auto-Suggestion` should be `Default/on` (usually no need to change)).
9. After configuration, choose to publish to discord, fill in the `token` of `COZE-BOT`, after publishing, you can see `COZE-BOT` online and can be used with @ in the discord server.
10. Start configuring [environment variables](#environment-variables) and [deploy](#deployment) this project using the recorded parameters.
11. Visit the API documentation address, and you can start debugging or integrating other projects.

## How to Integrate with NextChat

Fill in the interface address (ip:port/domain) and API-Key (`PROXY_SECRET`), other fields are optional.

> If you haven't set up a NextChat panel yourself, here's one already set up that you can use: [NextChat](https://ci.goeast.io/)

<span><img src="docs/img5.png" width="800"/></span>

## How to Integrate with one-api

Fill in `BaseURL` (ip:port/domain) and key (`PROXY_SECRET`), other fields are optional.

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
      - USER_AUTHORIZATION=MTA5OTg5N************aXUBHVI  # Must modify to your discord user's authorization key (multiple keys separated by commas)
      - BOT_TOKEN=MTE5OT************UrUWNbG63w  # Must modify to the listening bot's token
      - GUILD_ID=11************96  # Must modify to the server ID where both bots are located
      - COZE_BOT_ID=11************97  # Must modify to the bot ID managed by coze
      - CHANNEL_ID=11************94  # [Optional] Default channel - (currently this parameter is only used to keep the bot active)
      - PROXY_SECRET=123456  # [Optional] API key - modify this line to the value used for request header verification (multiple keys separated by commas)
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

If the above image cannot be pulled, try using the GitHub Docker image by replacing `deanxv/coze-discord-proxy` with `ghcr.io/deanxv/coze-discord-proxy`.

### Deploying to a Third-Party Platform

<details>
<summary><strong>Deploy to Zeabur</strong></summary>
<div>

> Zeabur's servers are located abroad, automatically solving network issues, and the free tier is sufficient for personal use.

Click to deploy:

[![Deploy on Zeabur](https://zeabur.com/button.svg)](https://zeabur.com/templates/GMU8C8?referralCode=deanxv)

**After one-click deployment, the variables `USER_AUTHORIZATION`, `BOT_TOKEN`, `GUILD_ID`, `COZE_BOT_ID`, `PROXY_SECRET`, `CHANNEL_ID` must also be replaced!**

Or manually deploy:

1. First **fork** a copy of the code.
2. Enter [Zeabur](https://zeabur.com?referralCode=deanxv), log in with GitHub, go to the console.
3. In Service -> Add Service, choose Git (authorize first if it's your first time), select the repository you forked.
4. Deployment will automatically start, cancel it first.
5. Add environment variables

   `USER_AUTHORIZATION:MTA5OTg5N************uIfytxUgJfmaXUBHVI`  Authorization key for discord users initiating messages (separated by commas)

   `BOT_TOKEN:MTE5OTk************WNbG63w`  Token for the bot listening to messages

   `GUILD_ID:11************96`  Server ID where both bots are located

   `COZE_BOT_ID:11************97`  Bot ID managed by coze

   `CHANNEL_ID:11************24`  # [Optional] Default channel - (currently this parameter is only used to keep the bot active)

   `PROXY_SECRET:123456` [Optional] API key - modify this line to the value used for request header verification (separated by commas) (similar to the openai-API-KEY)

Save.

6. Choose Redeploy.

</div>
</details>

<details>
<summary><strong>Deploy to Render</strong></summary>
<div>

> Render provides a free tier, and linking a card can further increase the limit.

Render can directly deploy Docker images without needing to fork the repository: [Render](https://dashboard.render.com)

</div>
</details>

## Configuration

### Environment Variables

1. `USER_AUTHORIZATION=MTA5OTg5N************uIfytxUgJfmaXUBHVI`  Authorization key for discord users initiating messages (separated by commas)
2. `BOT_TOKEN=MTE5OTk2************rUWNbG63w`  Token for the bot listening to messages
3. `GUILD_ID=11************96`  Server ID where all bots are located
4. `COZE_BOT_ID=11************97`  Bot ID managed by coze
5. `PORT=7077`  [Optional] Port, default is 7077
6. `SWAGGER_ENABLE=1`  [Optional] Enable Swagger API documentation [0: No; 1: Yes] (default is 1)
7. `ONLY_OPENAI_API=0`  [Optional] Expose only interfaces aligned with openai [0: No; 1: Yes] (default is 0)
8. `CHANNEL_ID=11************24`  [Optional] Default channel - (currently this parameter is only used to keep the bot active)
9. `PROXY_SECRET=123456`  [Optional] API key - modify this line to the value used for request header verification (separated by commas) (similar to the openai-API-KEY), **recommended to use this environment variable**
10. `DEFAULT_CHANNEL_ENABLE=0`  [Optional] Enable default channel [0: No; 1: Yes] (default is 0) If enabled, each conversation will occur in the default channel, **session isolation will be ineffective**, **not recommended to use this environment variable**
11. `ALL_DIALOG_RECORD_ENABLE=1`  [Optional] Enable full context [0: No; 1: Yes] (default is 1) If disabled, each conversation will only send the last `content` in `messages` where `role` is `user`, **not recommended to use this environment variable**
12. `CHANNEL_AUTO_DEL_TIME=5`  [Optional] Channel auto-delete time (seconds) This parameter is for automatically deleting the channel after each conversation (default is 5s), if set to 0 then it will not delete, **not recommended to use this environment variable**
13. `COZE_BOT_STAY_ACTIVE_ENABLE=1`  [Optional] Enable daily `9 AM` task to keep coze-bot active [0: No; 1: Yes] (default is 1), **not recommended to use this environment variable**
14. `REQUEST_OUT_TIME=60`  [Optional] Non-stream response timeout for conversation interface, **not recommended to use this environment variable**
15. `STREAM_REQUEST_OUT_TIME=60`  [Optional] Stream response timeout for each stream return in conversation interface, **not recommended to use this environment variable**
16. `REQUEST_RATE_LIMIT=60`  [Optional] Request rate limit per minute per IP, default: 60 requests/min
17. `USER_AGENT=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36`  [Optional] Discord user agent, using your own might help prevent being banned, if not set, defaults to the author's, **recommended to use this environment variable**
18. `NOTIFY_TELEGRAM_BOT_TOKEN=6232***********Niz9c`  [Optional] Token for a Telegram bot used for notifications (Notification events: 1. No available `user_authorization`; 2. `BOT_TOKEN` related BOT triggers risk control)
19. `NOTIFY_TELEGRAM_USER_ID=10******35`  [Optional] The `Telegram-Bot` associated with `NOTIFY_TELEGRAM_BOT_TOKEN` will push to the `Telegram-User` associated with this variable (**`NOTIFY_TELEGRAM_BOT_TOKEN` must not be empty if this variable is used**)
20. `PROXY_URL=http://127.0.0.1:10801`  [Optional] Proxy (supports http only)

## Advanced Configuration

### Configuring Multiple Bots

1. Before deployment, create a `data/config/bot_config.json` file in the same directory as the `docker`/`docker-compose` deployment
2. Write the `json` file, `bot_config.json` format as follows

```shell
[
  {
    "proxySecret": "123", // API request key (PROXY_SECRET) (Note: this key must exist in the environment variable PROXY_SECRET for this Bot to be matched!)
    "cozeBotId": "12***************31", // Bot ID managed by coze
    "model": ["gpt-3.5","gpt-3.5-16k"], // Model names (array format) (if the model in the request does not match any in this json, an exception will be thrown)
    "channelId": "12***************56"  // [Optional] Discord channel ID (the bot must be in the server where this channel is located) (currently this parameter is only used to keep the bot active)
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

3. Restart the service

> When this json configuration is present, the bot will be matched through the [request key] carried in the request header and the [`model`] in the request body.
> If multiple matches are found, one will be randomly selected. The configuration is very flexible and can be adjusted according to your needs.

For services deployed on third-party platforms (such as `zeabur`) that need [configuring multiple bots], please refer to [issue#30](https://github.com/deanxv/coze-discord-proxy/issues/30)

## Limitations

Current details of coze's free and paid subscriptions: https://www.coze.com/docs/guides/subscription?_lang=en

You can configure multiple discord users `Authorization` (refer to [Environment Variables](#environment-variables) `USER_AUTHORIZATION`) or [configure multiple bots](#configuring-multiple-bots) to stack request counts and balance load.

## Q&A

Q: How should I configure for high concurrency?

A: First, [configure multiple bots](#configuring-multiple-bots) to serve as response bots. Secondly, prepare multiple discord accounts to serve as request load and invite them into the same server, obtain the `Authorization` for each account, separate them with commas, and configure them in the environment variable `USER_AUTHORIZATION`. Each request will then pick one discord account to initiate the conversation, effectively achieving load balancing.

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=deanxv/coze-discord-proxy&type=Date)](https://star-history.com/#deanxv/coze-discord-proxy&Date)

## Related

[GPT-Content-Audit](https://github.com/deanxv/gpt-content-audit): An aggregation of Openai, Alibaba Cloud, Baidu Intelligent Cloud, Qiniu Cloud, and other open platforms, providing content audit services aligned with `openai` request formats.

## Others

**Open source is not easy, if you refer to this project or base your project on it, could you please mention this project in your project documentation? Thank you!**

Java: https://github.com/oddfar/coze-discord (Currently unavailable)

## References

Coze Official Website: https://www.coze.com

Discord Development Address: https://discord.com/developers/applications
