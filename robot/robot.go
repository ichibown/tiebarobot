package robot

import (
	"log"
	"os"
	"strings"
	"tiebarobot/simsimi"
	"tiebarobot/tieba"
	"time"
)

const (
	GET_AT_TICK    = 5
	REPLY_AT_TICK  = 5
	TASK_MAX_RETRY = 5
)

var (
	fidMapping   map[string]string
	replyChannel chan tieba.AtNode
	repliedId    []string
	user         *tieba.User
	username     string
)

func init() {
	fidMapping = make(map[string]string)
	replyChannel = make(chan tieba.AtNode, 100)
	repliedId = make([]string, 20)
	logFile, _ := os.OpenFile("robot.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(logFile)
}

func getFid(fname string) string {
	fid, exist := fidMapping[fname]
	if !exist {
		fid, err := tieba.GetFid(fname)
		if err != nil {
			return ""
		}
		fidMapping[fname] = fid
	}
	return fid
}

func isReplied(task tieba.AtNode) bool {
	for _, fid := range repliedId {
		if task.Post_id == fid {
			return true
		}
	}
	return false
}

func getAtRoutine() {
	timer := time.Tick(time.Second * GET_AT_TICK)

	for _ = range timer {
		result, err := user.GetAtMe()
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
		} else {
			for index, node := range result {
				if !isReplied(node) {
					content := node.Content
					content = strings.Replace(content, username, "", -1)
					node.Reply = simsimi.Talk(content)
					log.Printf("NewTask: User: %s  From: %s\n", node.Quote_user.Name, node.Fname)
					replyChannel <- node
				}
				repliedId[index] = node.Post_id
			}
		}

	}
}

func replyAtRoutine() {
	timer := time.Tick(time.Second * REPLY_AT_TICK)

	for _ = range timer {
		task := <-replyChannel
		err := user.ReplyFloor(task.Post_id, task.Thread_id, simsimi.Talk2(task.Content), getFid(task.Fname), task.Fname)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			if task.Retry < TASK_MAX_RETRY {
				task.Retry += 1
				replyChannel <- task
				log.Printf("RetryTask: User: %s  From: %s  Content: %s\n", task.Quote_user.Name, task.Fname, task.Content)
			} else {
				log.Printf("FailTask: User: %s  From: %s  Content: %s  Reply: %s\n", task.Quote_user.Name, task.Fname, task.Content, task.Reply)
			}
		} else {
			log.Printf("DoneTask: User: %s  From: %s  Content: %s  Reply: %s\n", task.Quote_user.Name, task.Fname, task.Content, task.Reply)
		}
	}

}

func Start(u *tieba.User, name string) {
	user = u
	username = "@" + name
	exit := make(chan int)

	go getAtRoutine()
	go replyAtRoutine()

	<-exit
}
