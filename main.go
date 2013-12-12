package main

import (
	"tiebarobot/robot"
	"tiebarobot/tieba"
)

func main() {
	user := tieba.User{"你的贴吧ID", "贴吧TBS，没有也无所谓", "贴吧BDUSS"}
	robot.Start(&user, "你的贴吧昵称")
}
