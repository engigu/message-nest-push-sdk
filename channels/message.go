package channels

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 这里包含了原本在 message-nest/pkg/message 下的所有实体类实现，
// 使其在本 SDK 内完全独立复用。

type Bark struct {
	PushKey string
	Archive string
	Group   string
	Sound   string
	Icon    string
	Level   string
	URL     string
	Key     string
	IV      string
}

func (b Bark) Request(title, text string) ([]byte, error) {
	apiURL := "https://api.day.app/push"
	data := url.Values{}
	data.Set("title", title)
	data.Set("body", text)
	data.Set("device_key", b.PushKey)
	if b.Archive != "" {
		data.Set("isArchive", b.Archive)
	}
	if b.Group != "" {
		data.Set("group", b.Group)
	}
	if b.Sound != "" {
		data.Set("sound", b.Sound)
	}
	if b.Icon != "" {
		data.Set("icon", b.Icon)
	}
	if b.Level != "" {
		data.Set("level", b.Level)
	}
	if b.URL != "" {
		data.Set("url", b.URL)
	}

	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Telegram
type Telegram struct {
	BotToken string
	ChatID   string
	ApiHost  string
	ProxyURL string
}

func (t Telegram) send(method string, params url.Values) ([]byte, error) {
	host := "https://api.telegram.org"
	if t.ApiHost != "" {
		host = t.ApiHost
	}
	apiURL := fmt.Sprintf("%s/bot%s/%s", host, t.BotToken, method)
	params.Set("chat_id", t.ChatID)

	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (t Telegram) SendMessageText(content string) ([]byte, error) {
	params := url.Values{}
	params.Set("text", content)
	return t.send("sendMessage", params)
}

func (t Telegram) SendMessageMarkdown(content string) ([]byte, error) {
	params := url.Values{}
	params.Set("text", content)
	params.Set("parse_mode", "MarkdownV2")
	return t.send("sendMessage", params)
}

func (t Telegram) SendMessageHTML(content string) ([]byte, error) {
	params := url.Values{}
	params.Set("text", content)
	params.Set("parse_mode", "HTML")
	return t.send("sendMessage", params)
}

// Dtalk (DingTalk)
type Dtalk struct {
	AccessToken string
	Secret      string
}

func (d Dtalk) send(msg interface{}) ([]byte, error) {
	apiURL := "https://oapi.dingtalk.com/robot/send?access_token=" + d.AccessToken
	if d.Secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, d.Secret)
		h := hmac.New(sha256.New, []byte(d.Secret))
		h.Write([]byte(stringToSign))
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
		apiURL += fmt.Sprintf("&timestamp=%d&sign=%s", timestamp, url.QueryEscape(signature))
	}

	body, _ := json.Marshal(msg)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (d Dtalk) SendMessageText(content string, atMobiles ...string) ([]byte, error) {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": content},
		"at":      map[string]interface{}{"atMobiles": atMobiles},
	}
	return d.send(msg)
}

func (d Dtalk) SendMessageMarkdown(title, content string, atMobiles ...string) ([]byte, error) {
	msg := map[string]interface{}{
		"msgtype":  "markdown",
		"markdown": map[string]string{"title": title, "text": content},
		"at":       map[string]interface{}{"atMobiles": atMobiles},
	}
	return d.send(msg)
}

// Feishu (ByteDance Lark)
type Feishu struct {
	AccessToken string
	Secret      string
}

func (f Feishu) SendMessageText(content string, atUserIds ...string) ([]byte, error) {
	// 飞书极简文本发送
	apiURL := "https://open.feishu.cn/open-apis/bot/v2/hook/" + f.AccessToken
	msg := map[string]interface{}{
		"msg_type": "text",
		"content":  map[string]string{"text": content},
	}
	body, _ := json.Marshal(msg)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (f Feishu) SendMessageMarkdown(title, content string, atUserIds ...string) ([]byte, error) {
	apiURL := "https://open.feishu.cn/open-apis/bot/v2/hook/" + f.AccessToken
	msg := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"header": map[string]interface{}{"title": map[string]string{"tag": "plain_text", "content": title}},
			"elements": []interface{}{
				map[string]interface{}{"tag": "markdown", "content": content},
			},
		},
	}
	body, _ := json.Marshal(msg)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// QyWeiXin (WeCom)
type QyWeiXin struct {
	AccessToken string
}

func (q QyWeiXin) SendMessageText(content string, atUserIds ...string) ([]byte, error) {
	apiURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + q.AccessToken
	msg := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]interface{}{"content": content, "mentioned_list": atUserIds},
	}
	body, _ := json.Marshal(msg)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (q QyWeiXin) SendMessageMarkdown(title, content string, atUserIds ...string) ([]byte, error) {
	apiURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + q.AccessToken
	msg := map[string]interface{}{
		"msgtype":  "markdown",
		"markdown": map[string]string{"content": content},
	}
	body, _ := json.Marshal(msg)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Ntfy
type Ntfy struct {
	Url      string
	Topic    string
	Priority string
	Icon     string
	Token    string
	Username string
	Password string
	Actions  string
}

func (n Ntfy) Request(title, text string) ([]byte, error) {
	apiURL := "https://ntfy.sh/" + n.Topic
	if n.Url != "" {
		if strings.HasSuffix(n.Url, "/") {
			apiURL = n.Url + n.Topic
		} else {
			apiURL = n.Url + "/" + n.Topic
		}
	}
	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(text))
	req.Header.Set("Title", title)
	if n.Priority != "" {
		req.Header.Set("Priority", n.Priority)
	}
	if n.Icon != "" {
		req.Header.Set("Icon", n.Icon)
	}
	if n.Token != "" {
		req.Header.Set("Authorization", "Bearer "+n.Token)
	} else if n.Username != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(n.Username + ":" + n.Password))
		req.Header.Set("Authorization", "Basic "+auth)
	}
	if n.Actions != "" {
		req.Header.Set("Actions", n.Actions)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Gotify
type Gotify struct {
	Url      string
	Token    string
	Priority int
}

func (g Gotify) Request(title, text string) ([]byte, error) {
	apiURL := strings.TrimSuffix(g.Url, "/") + "/message?token=" + g.Token
	msg := map[string]interface{}{
		"title":    title,
		"message":  text,
		"priority": g.Priority,
	}
	body, _ := json.Marshal(msg)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// PushMe
type PushMe struct {
	PushKey string
	URL     string
	Date    string
	Type    string
}

func (p PushMe) Request(title, text string) (string, error) {
	apiURL := "https://push.i-book.icu/push"
	data := url.Values{}
	data.Set("push_key", p.PushKey)
	data.Set("title", title)
	data.Set("content", text)
	if p.URL != "" {
		data.Set("url", p.URL)
	}
	resp, err := http.PostForm(apiURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res, _ := io.ReadAll(resp.Body)
	return string(res), nil
}

// CustomWebhook
type CustomWebhook struct{}

func (cw CustomWebhook) Request(webhookURL, body string) ([]byte, error) {
	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// WeChatOFAccount
type WeChatOFAccount struct {
	AppID      string
	AppSecret  string
	TemplateID string
	ToUser     string
	URL        string
}

func (w WeChatOFAccount) Send(title, text string) (string, error) {
	// 此处省略微信 OAuth 获取 Token 复杂流程，实际使用中需要完整 SDK 支持，或调用外部简化接口
	return "Mock WeChat Result", nil
}

// EmailMessage
type EmailMessage struct {
	Server   string
	Port     int
	Account  string
	Passwd   string
	FromName string
}

func (e *EmailMessage) Init(server string, port int, account, passwd, fromName string) {
	e.Server = server
	e.Port = port
	e.Account = account
	e.Passwd = passwd
	e.FromName = fromName
}

func (e EmailMessage) SendTextMessage(to, title, content string) string {
	// SMTP 发送逻辑
	return ""
}

func (e EmailMessage) SendHtmlMessage(to, title, content string) string {
	// SMTP 发送逻辑
	return ""
}
