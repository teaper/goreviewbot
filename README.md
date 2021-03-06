### Telegram 机器人
[Golang 中文研习社](https://t.me/golangzh) 机器人 「 [阿茶](https://t.me/GoReviewBot) 」  
#### 颜值
![](https://i.loli.net/2021/09/02/f8gjIGSHLdW1ebC.png)
#### 功能  
- [x] 入群验证  
    - [x] 数字验证码  
    - [x] 人工通过  
    - [x] 人工拒绝  
    - [ ] 超时踢出  
- [x] RSS 解析  
    - [x] Golang 开发进度推送  
    - [x] 消息置顶  
- [ ] 操作命令  
- [x] 消息审查  
    - [x] 网络敏感词  
    - [x] 英文消息自动翻译
    - [x] 感叹号复读机
#### 编译
```bash
make all-amd64 #编译适用于各个平台的版本在 bin 目录中
make clean  #清楚 bin 目录中编译后的文件
```
#### 部署
```bash
curl -LO https://github.com/teaper/goreviewbot/releases/download/v1.0.0/tgbot-linux-amd64 #下载程序
chmod +x tgbot-linux-amd64
./tgbot-linux-amd64 #初次运行会自动生成一个 conf.yaml 模板
vim conf.yaml #配置模板中的 token 
nohup ./tgbot-linux-amd64 > tgbot.log 2>&1 & #借助 nohub 后台运行
ps -aux | grep "tgbot-linux-amd64" #查看运行状态
cat tgbot.log #查看运行日志
```
#### 说明  
* 关于配置 `conf.yaml`，至少应该先配置 `token` 启动程序后，随意发送一条消息触发 bot ，会在控制台打印出 群组的 `Chat ID` 信息
* 机器人需要添加到群里，并且给管理员权限： `修改群组信息` `删除消息` `封禁用户` `生成邀请链接` `置顶消息` `Manage Voice Chats` `保持匿名`
* 建议：如果不想其他人使用你的 bot，可以去 [@BotFather](https://t.me/botfather) 输入 `/mybot` 找到 Bot Settings 将机器人设为私有

#### YAML 模板
```yaml
enabled: true
bot:
  #bot token
  token: 1889811505:AAE_Z2tOlROqkAeXC6Vdf5pnTnZ-Z4vJKiE
channels:
  #测试群 chat_id: -1001102843992 (配置 token 启动 bot 后控制台会打印出来)
  chat_id: -1001102843652
  #测试群 chat_user_name: https://t.me/golangzh
  chat_user_name: golangzh
  #群类型（supergroup）(channel)
  chat_type: supergroup
  #群创建者
  creator: teaper
rss:
  #RSS 订阅地址
  client_url: https://nitter.net/golang/rss
  #pubdate: RSS 最新消息的发布时间
  pubdate: Thu, 05 Aug 2021 18:29:17 GMT
```

#### 参考  
[Golang Telegram API 交流群](https://t.me/go_telegram_bot_api)  
[Telegram bot API](https://core.telegram.org/bots/api)  
[Golang Telegram bot API](https://github.com/go-telegram-bot-api/telegram-bot-api)  
[Golang Telegram bot API wiki](https://github.com/go-telegram-bot-api/telegram-bot-api/wiki)  
[验证码图片](https://count.getloli.com/)  
[内联菜单](https://zwindr.blogspot.com/2018/09/go-telegram-bot_22.html)  
[yaml 解析](https://gopkg.in/yaml.v3)  
[RSS 解析](https://www.youtube.com/watch?v=YynNUr1t6io)  
[Embedding Files](https://pkg.go.dev/embed@master)  





