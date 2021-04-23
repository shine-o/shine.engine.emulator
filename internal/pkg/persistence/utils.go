package persistence

import (
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NewDummyCharacter to be used only for testing purposes where a character is needed
func NewDummyCharacter(class string, initialItems bool) *Character {
	var (
		bitField byte
		name     string
	)

	switch class {
	case "mage":
		bitField = byte(1 | 16<<2 | 1<<7)
		name = fmt.Sprintf("mage%v", 1)
		break
	case "fighter":
		bitField = byte(1 | 1<<2 | 1<<7)
		name = fmt.Sprintf("fighter%v", 1)
		break
	case "archer":
		bitField = byte(1 | 11<<2 | 1<<7)
		name = fmt.Sprintf("archer%v", 1)
		break
	case "cleric":
		bitField = byte(1 | 6<<2 | 1<<7)
		name = fmt.Sprintf("cleric%v", 1)
		break
	}

	c := structs.NcAvatarCreateReq{
		SlotNum: byte(0),
		Name: structs.Name5{
			Name: name,
		},
		Shape: structs.ProtoAvatarShapeInfo{
			BF:        bitField,
			HairType:  6,
			HairColor: 0,
			FaceShape: 0,
		},
	}

	char, err := NewCharacter(1, &c, initialItems)
	if err != nil {
		log.Fatal(err)
	}

	return char
}
