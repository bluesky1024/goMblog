package chatService

import (
	"fmt"
	"testing"
)

var (
	chatSrv ChatServicer
	err     error
)

func init() {
	chatSrv, err = NewChatServicer()
	if err != nil {
		fmt.Println(err)
		panic("abc")
	}
}

func TestChatService_AddRoom(t *testing.T) {
	err := chatSrv.AddRoom("fang's room", 666, 2317487850917888, 2)
	fmt.Println(err)
}

func TestChatService_GetRoomConfigByRoomId(t *testing.T) {
	info, err := chatSrv.GetRoomConfigByRoomId(666)
	fmt.Println(info, err)
}
