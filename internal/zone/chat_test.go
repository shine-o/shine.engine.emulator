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

type (
	whisperLine struct{}
	guildLine   struct{}
	partyLine   struct{}
	academyLine struct{}
)

func TestChatProcessWhisperCommand(t *testing.T) {
	t.Fail()
}

func TestChatProcessGuildCommand(t *testing.T) {
	t.Fail()
}

func TestChatProcessAcademyCommand(t *testing.T) {
	t.Fail()
}

func TestChatProcessPartyCommand(t *testing.T) {
	t.Fail()
}
