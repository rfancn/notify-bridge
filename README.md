# 消息通知桥接器

接收推送过来的webhook消息，解析消息并发送到目标应用中。

当前支持：
- 接收禅道的webhook，查询Bug获取当前Bug指派人
- 通过机器人在Mattermost中与指派人建立私聊并发送通知

## 功能

- 接收禅道 Bug webhook（`POST /<route.name>`）
- 强制共享密钥鉴权（`Authorization: Bearer` 或 `X-Webhook-Token`）
- 从 webhook 文本或 `objectID` 提取 Bug ID
- 查询禅道 Bug 详情，使用 `assignedTo.account` 作为 Mattermost 用户名
- 与指派人建立 Mattermost 私聊（Direct Message）
- 按 `route + 用户` 缓存私聊频道
- 进程内 SHA-256 去重（带 TTL 和容量上限）
- 禅道 token 缓存，401 时自动刷新一次
- TOML 配置，支持 `HD_` 环境变量覆盖
- 优雅退出

## 构建

要求 Go 1.24+。

```bash
make build
```

产物为 `./notify-bridge`。

也可以直接：

```bash
go build -o notify-bridge .
```

## 配置

复制示例配置：

```bash
cp config.toml.example config.toml
```

然后修改 `app.token` 以及各路由的 `source` / `target` 配置。

启动：

```bash
./notify-bridge run -c config.toml
```

### 配置示例

```toml
[sdk.log]
    level = "info"
    filename = "notify-bridge.log"
    dir = "logs"
    [sdk.log.rotate]
        max_age = 168
        max_backup = 5
        max_size = 100

[app]
    listen = ":8080"
    token = "CHANGE_ME"

    [app.routes.zentao-mattermost-notify]
        [app.routes.zentao-mattermost-notify.source]
            type = "zentao"
            url = "http://zentao"
            account = "YOUR_ZENTAO_ACCOUNT"
            password = "YOUR_ZENTAO_PASSWORD"
        
        [app.routes.zentao-mattermost-notify.target]
            type = "mattermost"
            url = "http://mattermost"
            # 机器人访问令牌
            bot_token = "YOUR_MATTERMOST_BOT_TOKEN"
            # 机器人用户名
            bot_user = "YOUR_MATTERMOST_BOT_USER"
```

`app.routes` 下的每一个 key 都会注册为一个 webhook 路径：`POST /<key>`，例如 `[app.routes.zentao-mattermost-notify]` 对应 `POST /zentao-mattermost-notify`。`source.type` 决定消息解析方式，`target.type` 决定通知发送方式。当前支持 `zentao` → `mattermost`。

路由 key 会去除首尾空白，且只允许字母、数字、`-` 和 `_`。

`[sdk.log]` 由 hd SDK 使用。`dir` 必须是运行用户可写的目录；省略时 Linux 默认尝试 `/var/log`。

## 运行参数

```text
notify-bridge run [--config <file>] [--app-address <addr>]
notify-bridge version
```

- `--config, -c`：配置文件路径
- `--app-address, -d`：监听地址，覆盖 `app.listen`
- `version`：打印版本、commit 和构建时间

## webhook 接口

支持以下任意一种方式携带 token：

```text
POST /zentao-mattermost-notify
Authorization: Bearer <app.token>
Content-Type: application/json
```

```text
POST /zentao-mattermost-notify
X-Webhook-Token: <app.token>
Content-Type: application/json
```

```text
POST /zentao-mattermost-notify?token=<app.token>
Content-Type: application/json
```

```text
POST /zentao-mattermost-notify/<app.token>
Content-Type: application/json
```

请求体示例：

```json
{
  "objectType": "bug",
  "objectID": 3,
  "action": "edited",
  "actor": "admin",
  "date": "2026-10-08 17:43:58",
  "text": "admin编辑了Bug [#3::test_bug1](http://127.0.0.1/index.php?m=bug&f=view&id=1)"
}
```

`objectID` 支持字符串或数字。

### 响应语义

| 情况 | 状态码 |
|---|---|
| 鉴权失败 | 401 |
| 非 POST | 405 |
| Content-Type 非 `application/json` | 415 |
| JSON 解析失败、缺少 Bug ID | 400 |
| 请求体超过 固定 1MB | 413 |
| 禅道或 Mattermost 调用失败 | 500，禅道会按自身策略重试 |
| 发送成功 | 200 `{"ok":true}` |

## 禅道配置

在禅道的 webhook 设置中：

- URL 设置为如下其中一种方式：
  - `http://<服务地址>:8080/zentao-mattermost-notify?token=<app.token>`
  - `http://<服务地址>:8080/zentao-mattermost-notify/<app.token>`

解析出的消息会直接通知到对应 target，不再按动作过滤。

## Mattermost 权限

Mattermost token 对应的机器人账号需要：

1. 按用户名查询用户
2. 创建/获取与指定用户的私聊频道
3. 在私聊频道中发帖

## 去重

去重始终启用，无需配置。去重键是完整 webhook 请求体的 SHA-256，保存在进程内，默认 TTL 10 分钟、最多 2048 条。

- 处理成功：记录保留到 TTL 过期
- 处理失败：立即移除记录，允许禅道重试
- 进程重启后记录丢失；多实例之间不共享

## 禅道 token

服务使用 `account` + `password` 调用 `/api.php/v1/tokens` 获取 token，并在进程内缓存。token 失效（接口返回 401）时会自动清除缓存并重新获取一次。

```toml
[app.routes.zentao-mattermost-notify.source]
    type = "zentao"
    url = "http://127.0.0.1"
    account = "admin"
    password = "YOUR_PASSWORD"
```

## Docker

```bash
docker build -t notify-bridge .
docker run --rm -p 8080:8080 \
  -v "$PWD/config.toml:/app/config.toml:ro" \
  notify-bridge
```

也可以用 `-e HD_APP_TOKEN=<token>` 覆盖配置中的 `app.token`。

## 测试

```bash
go test -race ./...
go test -cover ./...
golangci-lint run
gosec ./...
```
