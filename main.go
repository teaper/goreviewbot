package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
	"goreviewbot/code"
	"goreviewbot/msgc"
	"goreviewbot/rss"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

//4：抽取出常用的两个对象
type TeleBot struct {
	botAPI  *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

//入群验证内联键盘
var joinedInlineKeyboardMarkup = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("0", "0"),
		tgbotapi.NewInlineKeyboardButtonData("1", "1"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
		tgbotapi.NewInlineKeyboardButtonData("7", "7"),
		tgbotapi.NewInlineKeyboardButtonData("8", "8"),
		tgbotapi.NewInlineKeyboardButtonData("9", "9"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("人工通过", "人工通过"),
		tgbotapi.NewInlineKeyboardButtonData("人工拒绝", "人工拒绝"),
	),
)

//判断用户是否是管理员
func (t *TeleBot) IsAdministrator(chatID int64,userName string) (bool bool,status string){
	administrators, _ := t.botAPI.GetChatAdministrators(tgbotapi.ChatConfig{
		ChatID: chatID,
	})
	for _, user := range administrators {
		log.Printf("管理员名字：%s 职责：%s \n",user.User.UserName,user.Status)
		//creator 创建者 administrator 管理员
		if userName == user.User.UserName {
			return 	true,user.Status
		}
	}
	return false,""
}

// 发送警告 CallbackQuery
func (t *TeleBot) EmptyAnswer(CallbackQueryID string,text string) {
	configAlert := tgbotapi.NewCallback(CallbackQueryID, text)
	go t.botAPI.AnswerCallbackQuery(configAlert)
}

/*
ctk: 选择限制权限(restrict)还是拉黑踢出(kick)解除权限(unban)
chatID：群 id
userID：群里被处理的人的 id
untilDate：处理时间(永久封禁条件：> 366 day || < 30s)
 */
func (t *TeleBot) RestrictOrKickChatMember(ctk string,chatID int64,userID int,untilDate int64)  {
	switch ctk {
	case "restrict":
		go t.botAPI.RestrictChatMember(
			tgbotapi.RestrictChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: chatID,
					UserID: userID,
				},
				UntilDate: untilDate,
			},
		)
	case "kick":
		go t.botAPI.KickChatMember(
			tgbotapi.KickChatMemberConfig{
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: chatID,
					UserID: userID,
				},
				UntilDate: untilDate,
			},
		)
	default:
		log.Println("ctk 请选择 restrict || kick")
	}
}


//发送 RSS 新闻
func (t *TeleBot) SendRssNews()  {
	var news = rss.GetRssPage(cfg.Rss.ClientURL,&cfg.Rss.Pubdate)
	if news != "" {
		log.Println("拿到的go news信息：",news)
		//向 @golangzh 群发送消息
		msg := tgbotapi.NewMessageToChannel("@"+cfg.Channels.ChatUserName,news)
		msg.ParseMode = tgbotapi.ModeMarkdown
		send, _ := t.botAPI.Send(msg) //发送消息
		//消息置顶
		pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
			ChatID:              send.Chat.ID,
			MessageID:           send.MessageID,
			DisableNotification: true, //是否通知所有成员
		}
		go t.botAPI.PinChatMessage(pinChatMessageConfig)
	}
}



// 全局变量
var (
	cfg         Config //config.yaml 文件
	callNum     = 0    //回调匹配四次 codes 数组中的元素
	codeMsgsMap = make(map[int] msgc.CodeMessage)
)

//4：主体逻辑
func (t *TeleBot) sendAnswerCallbackQuery() {
	//5：获取 update 对象（消息的更新）
	for update := range t.updates {
		//6：如果 update 对象中没有更新的消息或者回调消息，就跳过当前 update，否则继续执行
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		//7：如果有新消息
		if update.Message != nil {
			log.Println("当前群的 ChatID (用于 conf.yaml 中) ==> ",update.Message.Chat.ID)
			//8：如果有新入群消息
			if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
				if update.Message.NewChatMembers != nil {
					//读取所有新加群成员
					for _, user := range *update.Message.NewChatMembers {
						log.Printf("添加入群验证用户：%s ==> %d \n", "@"+user.UserName, user.ID)
						//11:判断用户名长度和是否包含两位数字（一些广告账户）
						//正则匹配用户名中带两位连续数字的帐号
						reg, _ := regexp.Compile(`\D\d\d`)
						photos, _ := t.botAPI.GetUserProfilePhotos(tgbotapi.NewUserProfilePhotos(user.ID))
						//用户名中有两位数字，用户名为空，用户名长度超过 15 个字符，用户是机器人，用户头像图片数量为 0 ,一律拉黑踢出
						if (len(user.UserName) >= 15) || (reg.FindString(user.UserName) != "") || (user.UserName == "") || (user.IsBot) || (photos.TotalCount == 0) {
							go t.RestrictOrKickChatMember("kick",update.Message.Chat.ID,user.ID,time.Now().Unix()+1800) //踢出去 5 分钟
							continue
						}
						//11：验证时先添加限制
						go t.RestrictOrKickChatMember("restrict",update.Message.Chat.ID,user.ID,time.Now().Unix()) //默认永久禁言
						codes := [4] int32{}
						for {
							//9：生成验证码（随机4位数）
							ycode := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000)
							//用数组表示四张图的编号,但是我不希望这四个数字有重复，因为内联按钮就操作延迟（短时间内不能重复按同一个数字）
							codes = [4]int32{ycode / 1000 % 10, ycode / 100 % 10, ycode / 10 % 10, ycode / 1 % 10}
							if (codes[0] != codes[1]) && (codes[0] != codes[2]) && (codes[0] != codes[3]) && (codes[1] !=
								codes[2]) && (codes[1] != codes[3]) && (codes[2] != codes[3]) {
								break
							}
							log.Println("验证码 ==>>> ", codes)
						}
						code.CreateVerificationCode(codes)
						msg := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "code.gif")
						msg.Caption = "⚠️ " + "@"+user.UserName + " 35 秒时间输入图中验证码，超过时间或输入错误将被立即拉黑踢出 (仅一次机会)"
						//10：绑定内联按钮
						msg.ReplyMarkup = joinedInlineKeyboardMarkup
						log.Printf("发送图片的 msgID:%d , meg ==> %v",msg.ReplyToMessageID, msg)
						send, _ := t.botAPI.Send(msg)
						//将验证信息存储到 map
						codeMsgsMap[send.MessageID] = msgc.CodeMessage{
							MessageID: send.MessageID,
							Codes: codes,
							AuthUser: msgc.AuthUser{
								UserID: user.ID,
								UserName: user.UserName,
							},
							Enabled: true,
						}
						log.Printf("存储到 map 的验证信息：%v, 用户：%v \n",codeMsgsMap[send.MessageID],codeMsgsMap[send.MessageID].AuthUser)
						//长时间不操作
					}
				}
			}

			//检查话题 OT
			otif,ot := msgc.OtMessage(update.Message.Text)
			if otif == true {
				//提示话题#OT
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "<em>#OT</em> <strong>知音</strong> ⇉ "+ot)
				msg.ParseMode = tgbotapi.ModeHTML
				msg.ReplyToMessageID = update.Message.MessageID
				go t.botAPI.Send(msg)
			}

			//8：其他消息做 switch 匹配消息
			switch update.Message.Text {
			//9：如果是消息为 “/demo” ，初始化一个可操作消息
			case "/start":
				log.Println("/start")
			case "/demo":
				log.Println("/demo")
			case "/newbot":
				log.Println("创建一个机器人")
			}

		}

		//7：如果有回调消息
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID //群 id
			chatName := update.CallbackQuery.Message.Chat.UserName //群用户名
			CallMessageID := update.CallbackQuery.Message.MessageID
			CallformName := update.CallbackQuery.From.UserName
			log.Println("codeMsgsMap 中的所有值 ==>> ",codeMsgsMap[CallMessageID].AuthUser)
			//抽取点击者的信息
			codeMsgsID := codeMsgsMap[CallMessageID].MessageID
			newUser := codeMsgsMap[CallMessageID].AuthUser.UserName
			newUserID := codeMsgsMap[CallMessageID].AuthUser.UserID
			codes := codeMsgsMap[CallMessageID].Codes
			log.Printf("\n 回调消息 ChatID ==> %d, ChatName ==> %s, 点击验证码的人 ==> %s, 需要验证的人 ==> %s \n",chatID,chatName,CallformName,newUser)
			log.Printf("\n 被点击的消息ID ==> %d, 被点消息属于验证用户 ==> %s, 被点消息验证用户的ID ==> %d, 被点消息的验证码 ==> %v \n",codeMsgsID,newUser,newUserID,codes)

			//10：匹配回调
			switch update.CallbackQuery.Data {
			case "0","1","2","3","4","5","6","7","8","9":
				//判断点击的人是否是要验证的人
				if CallformName != newUser{
					//发送警告
					go t.EmptyAnswer(update.CallbackQuery.ID, "你戳疼人家了 (*/ω＼*)")
					continue
				}
				//12：正式匹配验证码
				log.Printf("点击数字键盘获取数字：%s  ==> 原来验证码中的值：%d",update.CallbackQuery.Data,codes)
				if callNum <= 3 {
					if update.CallbackQuery.Data == strconv.FormatInt(int64(codes[callNum]), 10) {
						callNum++
						log.Printf("点击第 %d 次通过 \n",callNum)
						if callNum == 4 {
							//验证通过
							callNum = 0
							log.Printf("验证结束 callNum 重置为 ==> %d \n",callNum)
							//发送提示
							go t.EmptyAnswer(update.CallbackQuery.ID, "@"+CallformName+" 「验证成功 欢迎入群」 🎉🎉🎉")
							//删除面板
							go t.botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: CallMessageID})
							//修改修改禁言时间为35秒(永久封禁：> 366 day || < 30s)，时间到了就能聊天了
							go t.RestrictOrKickChatMember("restrict",chatID,newUserID,time.Now().Unix() + 115)
							//删除 map
							delete(codeMsgsMap,CallMessageID)
						}
					} else {
						//点错了
						callNum = 0
						log.Printf("验证未通过 callNum 重置为 %d \n",callNum)
						//发送提示
						go t.EmptyAnswer(update.CallbackQuery.ID, "@"+CallformName+" 「验证失败 10分钟后再试」 💔💔💔")
						//删除面板
						go t.botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: CallMessageID})
						//踢出并拉黑成员
						log.Printf("踢出成员：UserID ==> %d \n",newUserID)
						go t.RestrictOrKickChatMember("kick",chatID,newUserID,time.Now().Unix()+1800) //10分钟后再试
						//删除 map
						delete(codeMsgsMap,CallMessageID)
					}
				}
			case "人工通过":
				//判断操作的人是否是管理员
				bl, status := t.IsAdministrator(chatID,CallformName)
				if bl == false {
					//发送警告
					go t.EmptyAnswer(update.CallbackQuery.ID, "您不是 "+status+" 无法操作")
					continue
				}
				//人工通过
				callNum = 0
				log.Printf("人工通过 callNum 重置为 ==> %d \n",callNum)
				//删除面板
				log.Printf("codeMsgsID => %d ,messageId => %d ,codeMsgsMap[messageId] => %v \n",codeMsgsID,CallMessageID,codeMsgsMap[CallMessageID])
				go t.botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: CallMessageID})
				//修改权限
				//go t.RestrictOrKickChatMember("restrict",chatID,newUserID,time.Now().Unix()+115)
				go t.RestrictOrKickChatMember("restrict",chatID,newUserID,time.Now().Unix()+115)
				//删除 map
				delete(codeMsgsMap,CallMessageID)
			case "人工拒绝":
				bl, status := t.IsAdministrator(chatID,CallformName)
				if bl == false {
					//发送警告
					go t.EmptyAnswer(update.CallbackQuery.ID, "您不是 "+status+" 无法操作")
					continue
				}
				//人工拒绝
				callNum = 0
				log.Printf("人工拒绝 callNum 重置为 %d \n",callNum)
				//删除面板
				log.Printf("codeMsgsID => %d ,messageId => %d ,codeMsgsMap[messageId] => %v \n",codeMsgsID,CallMessageID,codeMsgsMap[CallMessageID])
				go t.botAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: chatID, MessageID: CallMessageID})
				//踢出并拉黑（1800 >= 5分钟，官方时间不准）成员
				log.Printf("踢出成员：UserID ==> %d \n",newUserID)
				go t.RestrictOrKickChatMember("kick",chatID,newUserID,time.Now().Unix()+1800)
				//删除map
				delete(codeMsgsMap,CallMessageID)
			}
		}

		//检查 RSS 是否有最新消息(计时器)
		c := cron.New()
		err := c.AddFunc("@every 10m", func() {
			log.Println("启动 RSS 消息推送（10 分钟检查一次）")
			go t.SendRssNews()
		})
		if err != nil {
			log.Println(err)
		}
		c.Start()
		time.Sleep(time.Second * 5)
	}
}
/** ---------------------------------------------------------------------------
	说明：goreviewbot 是一个 telegram 群管理审查机器人，主要功能如下：
	1：入群验证
	2：删除敏感消息
	3：话题 #OT 提醒
	4：订阅 Go 语言 RSS 消息并推送
	5：菜单助手
 */
func main() {
	//1：加载配置文件 config.yaml
	log.Println("加载 yaml 文件中的 token：",cfg.GetConf().Bot.Token)
	//2：传入 token 并抛出 err
	bot, err := tgbotapi.NewBotAPI(cfg.GetConf().Bot.Token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	teleBot := TeleBot{
		botAPI: bot,
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	teleBot.updates, err = bot.GetUpdatesChan(u)
	//3：正式主体逻辑（匹配消息，送出菜单，匹配菜单回调，处理结果）
	teleBot.sendAnswerCallbackQuery()
}
