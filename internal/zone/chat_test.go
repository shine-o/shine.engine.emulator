package zone

import "testing"

type chatType int

const (
	whisper chatType = iota
	guild
	party
	academy
)

type chatLine interface {
	text() string
	chatType() chatType
}

type whisperLine struct{}
type guildLine struct{}
type partyLine struct{}
type academyLine struct{}

func Test_Chat_Process_Whisper_Command(t *testing.T) {
	t.Fail()
}

func Test_Chat_Process_Guild_Command(t *testing.T) {
	t.Fail()
}

func Test_Chat_Process_Academy_Command(t *testing.T) {
	t.Fail()
}

func Test_Chat_Process_Party_Command(t *testing.T) {
	t.Fail()
}

// general
// INFO : 2021/04/24 23:55:01.426155 handlers.go:267: 2021-04-24 23:55:01.424044 +0200 CEST 39878->9120 outbound NC_ACT_CHAT_REQ {"packetType":"small","length":5,"department":8,"command":"1","opCode":8193,"data":"000177","rawData":"050120000177","friendlyName":""}
// INFO : 2021/04/24 2	3:55:01.549378 handlers.go:267: 2021-04-24 23:55:01.534468 +0200 CEST 9120->39878 inbound NC_ACT_SOMEONECHAT_CMD {"packetType":"small","length":10,"department":8,"command":"2","opCode":8194,"data":"002c240102000077","rawData":"0a0220002c240102000077","friendlyName":""}

// academy (not in any academy)
// INFO : 2021/04/24 23:56:04.057276 handlers.go:267: 2021-04-24 23:56:04.044831 +0200 CEST 39878->9120 outbound OperationCode(39016) {"packetType":"small","length":9,"department":38,"command":"68","opCode":39016,"data":"00056164736661","rawData":"09689800056164736661","friendlyName":""}

// guild (not in any guild )
// INFO : 2021/04/24 23:57:15.615055 handlers.go:267: 2021-04-24 23:57:15.608967 +0200 CEST 39878->9120 outbound NC_GUILD_ACADEMY_REWARD_STORAGE_WITHDRAW_CMD {"packetType":"small","length":8,"department":29,"command":"73","opCode":29811,"data":"000461616461","rawData":"087374000461616461","friendlyName":""}

// party (not in any party)
// INFO : 2021/04/24 23:58:09.306052 handlers.go:267: 2021-04-24 23:58:09.30302 +0200 CEST 39878->9120 outbound NC_ACT_PARTYCHAT_REQ {"packetType":"small","length":8,"department":8,"command":"14","opCode":8212,"data":"000461647361","rawData":"081420000461647361","friendlyName":""}
