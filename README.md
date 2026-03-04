# Message Nest Push SDK

统一消息发送 SDK，可独立于 Message-Nest 业务层使用。

## 特点

- **零业务依赖**：不依赖数据库、路由、认证等业务代码
- **开箱即用**：直接传配置 + 消息即可发送
- **支持 12 种渠道**：Email, 钉钉, 企业微信, 飞书, Telegram, Bark, Ntfy, Gotify, PushMe, 自定义Webhook, 微信公众号, 阿里云短信
- **可扩展**：通过 `RegisterChannel` 注册自定义渠道

## 快速使用

```go
import "github.com/engigu/message-nest-push-sdk"

// 方式1: 直接发送
result, err := messenger.Send("Telegram", messenger.ChannelConfig{
    "bot_token": "your-bot-token",
    "chat_id":   "123456",
}, &messenger.Message{
    Title:    "告警通知",
    Text:     "服务器 CPU 超过 90%",
    Markdown: "**服务器 CPU** 超过 `90%`",
})

if err != nil {
    log.Fatal(err)
}
if !result.Success {
    log.Printf("发送失败: %s", result.Error)
}
```

### 使用 Client（预设默认配置）

```go
client := messenger.NewClient()

// 预设 Telegram 配置
client.SetDefaultConfig("Telegram", messenger.ChannelConfig{
    "bot_token": "your-bot-token",
    "chat_id":   "default-chat",
})

// 后续发送可以覆盖部分配置，也可以传 nil 使用全部默认值
result, err := client.Send("Telegram", nil, &messenger.Message{
    Title: "Hello",
    Text:  "World",
})
```

## 各渠道配置参数

### Email
| 参数 | 必填 | 说明 |
|------|------|------|
| `server` | ✅ | SMTP 服务地址 |
| `port` | ✅ | SMTP 端口 |
| `account` | ✅ | 邮箱账号 |
| `passwd` | ✅ | 邮箱密码 |
| `from_name` | ❌ | 发信人名称 |
| `to_account` | ✅ | 收件邮箱 |

### Dtalk（钉钉）
| 参数 | 必填 | 说明 |
|------|------|------|
| `access_token` | ✅ | 钉钉 access_token |
| `secret` | ❌ | 加签秘钥 |

### QyWeiXin（企业微信）
| 参数 | 必填 | 说明 |
|------|------|------|
| `access_token` | ✅ | 企业微信 access_token |

### Feishu（飞书）
| 参数 | 必填 | 说明 |
|------|------|------|
| `access_token` | ✅ | 飞书 access_token |
| `secret` | ❌ | 加签秘钥 |

### Telegram
| 参数 | 必填 | 说明 |
|------|------|------|
| `bot_token` | ✅ | Bot Token |
| `chat_id` | ✅ | Chat ID |
| `api_host` | ❌ | 自定义 API 地址 |
| `proxy_url` | ❌ | 代理地址 (http/https/socks5) |

### Bark
| 参数 | 必填 | 说明 |
|------|------|------|
| `push_key` | ✅ | Bark Push Key |
| `archive` | ❌ | 是否存档 |
| `group` | ❌ | 推送分组 |
| `sound` | ❌ | 推送声音 |
| `icon` | ❌ | 推送图标 |
| `level` | ❌ | 时效性 |
| `url` | ❌ | 跳转URL |
| `key` | ❌ | 加密Key |
| `iv` | ❌ | 加密IV |

### Ntfy
| 参数 | 必填 | 说明 |
|------|------|------|
| `topic` | ✅ | Topic |
| `url` | ❌ | 自定义 API 地址 |
| `priority` | ❌ | 优先级 |
| `icon` | ❌ | 图标 URL |
| `token` | ❌ | Token |
| `username` | ❌ | 用户名 |
| `password` | ❌ | 密码 |
| `actions` | ❌ | Actions |

### Gotify
| 参数 | 必填 | 说明 |
|------|------|------|
| `url` | ✅ | Gotify 服务地址 |
| `token` | ✅ | Token |
| `priority` | ❌ | 优先级 |

### PushMe
| 参数 | 必填 | 说明 |
|------|------|------|
| `push_key` | ✅ | PushMe Push Key |
| `url` | ❌ | 自定义 API 地址 |
| `date` | ❌ | 日期 |
| `type` | ❌ | 类型 |

### Custom（自定义Webhook）
| 参数 | 必填 | 说明 |
|------|------|------|
| `webhook` | ✅ | Webhook URL |
| `body` | ❌ | 请求体模板（`TEXT` 占位符会被替换） |

### WeChatOFAccount（微信公众号）
| 参数 | 必填 | 说明 |
|------|------|------|
| `appID` | ✅ | 公众号 AppID |
| `appsecret` | ✅ | 公众号 AppSecret |
| `tempid` | ❌ | 模板消息 ID |
| `to_account` | ✅ | 接收者 OpenID |

### AliyunSMS（阿里云短信）
| 参数 | 必填 | 说明 |
|------|------|------|
| `access_key_id` | ✅ | AccessKeyId |
| `access_key_secret` | ✅ | AccessKeySecret |
| `sign_name` | ✅ | 短信签名 |
| `region_id` | ❌ | 区域ID（默认 cn-hangzhou） |
| `phone_number` | ✅ | 手机号码 |
| `template_code` | ✅ | 短信模板 CODE |

## 自定义扩展

```go
// 注册自定义渠道
messenger.RegisterChannel("MyChannel", func() messenger.Channel {
    return &myCustomChannel{}
})
```
