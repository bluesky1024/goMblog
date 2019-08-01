package chatService

import (
	"fmt"
	"testing"
)

var (
	chatSrv ChatServicer
)

func init() {
	chatSrv, err := NewChatServicer()
	if err != nil {
		fmt.Println(err)
		panic("abc")
	}
}

func TestChatService_AddRoom(t *testing.T) {

}

func TestChatService_GetRoomConfigByRoomId(t *testing.T) {

}
