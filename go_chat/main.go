package main

import (
	"go_chat/chat"
	"go_chat/sqldb"

)

//
// skelton is from https://www.cetus-media.info/article/2021/line-chat/
// melody sample https://github.com/olahol/melody/tree/master/examples/multichat
//
// the gin & melody are used to create this application
//

func main() {
	// data base create and make table to store chat messages
	sql_db.DbCreate()
	chat.Run()
}
