package cqtotg

import (
	"net/http"
	"log"
	"encoding/json"
	"io"
)

type PublicHeartbeatRequest struct {
	PostType    string `json:"post_type"`
}

type UnheartbeatRequest struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	Time        int    `json:"time"`
	SelfID      int64  `json:"self_id"`
	SubType     string `json:"sub_type"`
	TargetID    int64  `json:"target_id"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
	Font        int    `json:"font"`
	Sender      struct {
		Age      int    `json:"age"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserID   int    `json:"user_id"`
	} `json:"sender"`
	MessageID int `json:"message_id"`
	UserID    int `json:"user_id"`
}

type HeartbeatRequest struct {
	PostType      string `json:"post_type"`
	MetaEventType string `json:"meta_event_type"`
	Time          int    `json:"time"`
	SelfID        int64  `json:"self_id"`
	Interval      int    `json:"interval"`
	Status        struct {
		AppEnabled     bool        `json:"app_enabled"`
		AppGood        bool        `json:"app_good"`
		AppInitialized bool        `json:"app_initialized"`
		Good           bool        `json:"good"`
		Online         bool        `json:"online"`
		PluginsGood    interface{} `json:"plugins_good"`
		Stat           struct {
			PacketReceived  int `json:"packet_received"`
			PacketSent      int `json:"packet_sent"`
			PacketLost      int `json:"packet_lost"`
			MessageReceived int `json:"message_received"`
			MessageSent     int `json:"message_sent"`
			LastMessageTime int `json:"last_message_time"`
			DisconnectTimes int `json:"disconnect_times"`
			LostTimes       int `json:"lost_times"`
		} `json:"stat"`
	} `json:"status"`
}

func Post(writer http.ResponseWriter, request *http.Request) {

	x, _ := io.ReadAll(request.Body)
	
	var jsonRet PublicHeartbeatRequest
	json.Unmarshal(x, &jsonRet)
	// log.Printf("%T: %v\n",jsonRet ,jsonRet)
	// log.Println(jsonRet.PostType)

	PostType := jsonRet.PostType

	switch PostType {
		// case "meta_event": {}

		case "message": {
			log.Printf("%v", string(x))
			var message UnheartbeatRequest
			json.Unmarshal(x, &message)
			
			log.Printf("[%v][%v][%v]said: ", message.UserID, message.MessageType, message.Sender.Nickname)
			log.Println(message.Message)


		}
	}
}
