package persistence

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NewDummyCharacter to be used only for testing purposes where a character is needed
func NewDummyCharacter(class string, initialItems bool, name string) *Character {
	var bitField byte

	switch class {
	case "mage":
		bitField = byte(1 | 16<<2 | 1<<7)
		break
	case "fighter":
		bitField = byte(1 | 1<<2 | 1<<7)
		break
	case "archer":
		bitField = byte(1 | 11<<2 | 1<<7)
		break
	case "cleric":
		bitField = byte(1 | 6<<2 | 1<<7)
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
