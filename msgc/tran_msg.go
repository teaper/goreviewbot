package msgc

import (
	"github.com/Conight/go-googletrans"
	"log"
	"regexp"
	"unicode"
)

//判断字符串是否包含中文字符
func IsChineseChar(msg string) bool {
	for _, r := range msg {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

//翻译英文消息成中文
func TranEn(msg string) string {
	t := translator.New()
	result, err := t.Translate(msg, "auto", "zh-CN")
	if err != nil {
		log.Println("翻译错误：",err)
	}
	return result.Text
}
