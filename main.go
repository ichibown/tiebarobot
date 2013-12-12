package main

import (
	"tiebarobot/robot"
	"tiebarobot/tieba"

//
// "fmt"
// "tiebarobot/simsimi"
)

func main() {
	user := tieba.User{"942028639", "f1b254e1d55c53941386816804", "ZiWTVjaVQyT1JveHB1MHh6c1hoWUh-RVpVWHFsZjhkLWNNfn5yUGRITWtzdEJTQVFBQUFBJCQAAAAAAAAAAAEAAABfNyY4zeO2ubv6xvfIywAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACQlqVIkJalSZW"}
	robot.Start(&user, "豌豆机器人")
}
