package General

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"os"
	"strings"
)

type (
	MessageInfo struct {
		Date        string `json:"data"`        // 时间
		Status      bool   `json:"status"`      // 群聊属性
		AtMe        bool   `json:"atme"`        // 是否@我
		RoomName    string `json:"roomname"`    // 群聊名称
		RoomID      string `json:"roomid"`      // 群聊ID
		UserName    string `json:"username"`    // 用户名称
		UserID      string `json:"userid"`      // 用户ID
		Content     string `json:"content"`     // 聊天内容
		AutoInfo    string `json:"autoinfo"`    // 信息一览
		Reply       string `json:"reply"`       // 自动回复
		ReplyStatus bool   `json:"replystatus"` // 自动回复状态
	}
)

var Messages = MessageInfo{}

/*
	EncodeMessage()
	对消息内容进行编码
*/
func EncodeMessage(message *user.Message) {
	if message.Type() != schemas.MessageTypeText {
		Messages.Content = "[未知消息类型: " + message.Type().String() + "] " + message.MentionText()
	} else {
		Messages.Content = message.MentionText()
	}
	Messages.Date = message.Date().Format("2006-01-02 15:04:05")
	Messages.Status = false
	Messages.AtMe = false
	Messages.UserName = message.From().Name()
	Messages.UserID = message.From().ID()
	Messages.AutoInfo = fmt.Sprintf("用户ID: [%v] 用户名称: [%v] 说: [%v] 回复: ", Messages.UserID, Messages.UserName, strings.Replace(Messages.Content, "\u2005", " ", -1))
	Messages.RoomName = ""
	Messages.RoomID = ""
	Messages.Reply = ""
	Messages.ReplyStatus = false
	if message.MentionSelf() {
		Messages.AtMe = true
	}
	if message.Room() != nil {
		Messages.RoomName = message.Room().Topic()
		Messages.RoomID = message.Room().ID()
		Messages.Status = true
		Messages.AutoInfo = fmt.Sprintf("群聊ID: [%v] 群聊名称: [%v] %v", Messages.RoomID, Messages.RoomName, Messages.AutoInfo)
	}
}

/*
	ExportMessages()
	对消息内容进行存储
*/
func ExportMessages(message *user.Message) {
	var (
		fp       *os.File
		filename = viper.GetString("rootPath") + "/data.json"
		result   []byte
	)
	if result, err = json.Marshal(Messages); err != nil {
		log.Errorf("[ExportMessages] Json 解析失败! Error: [%v]", err)
		return
	}
	if _, err = os.Stat(filename); err != nil {
		log.Errorf("[ExportMessages] 聊天备份文件不存在,正在创建! Error: [%v]", err)
	}
	if fp, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err != nil {
		log.Errorf("[ExportMessages] 打开聊天备份文件失败! Error: [%v]", err)
	}
	defer func(fp *os.File) {
		if err = fp.Close(); err != nil {
			log.Errorf("[ExportMessages] 关闭聊天备份文件失败! Error: [%v]", err)
		}
	}(fp)
	if _, err = fp.Write(result); err != nil {
		log.Errorf("[ExportMessages] 写入聊天记录到聊天备份文件失败! Error: [%v]", err)
		return
	}
	if _, err = fp.WriteString("\n"); err != nil {
		log.Errorf("[ExportMessages] 写入换行符到聊天备份文件失败! Error: [%v]", err)
		return
	}
	log.Printf(Messages.AutoInfo)
}
