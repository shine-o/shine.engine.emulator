package persistence

import (
	"github.com/go-pg/pg/v10"
	"time"
)

// todo: move to inventory package
type Item struct {
	tableName struct{} `pg:"world.character_items"`
	//ID         				 uint64
	//UUID          string `pg:",pk,notnull,type:uuid,default:gen_random_uuid"`
	CharacterID   uint64 `pg:",pk,notnull"`
	InventoryType int    `pg:",pk,notnull"`
	// box 2 = reward  inventory
	// box 3 = mini house furniture
	// box 8 = equipped items  / 1-29
	// box 9 = inventory, storage  // 9216 - 9377 (24 slots per page)
	// box 12 = mini houses // 12288 equipped minihouse, 12299-12322 available slots
	// box 13 = mini house accessories
	// box 14 = mini house tile all inventory
	// box 15 = premium actions inventory(dances)
	// box 16 = mini house mini game inventory
	Slot      uint16     `pg:",pk,notnull"`
	Character *Character `pg:"rel:belongs-to"`

	ShnID     uint16 `pg:",notnull"`
	Stackable bool   `pg:",notnull,use_zero"`
	Amount    uint32
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `pg:",soft_delete"`
}

type ItemData struct {
	Attributes   ItemAttributes
	Enchantments Enchantments
	Licences     Licences
}

type Enchantments struct {
	Max          byte
	Used         byte
	Licences     byte
	Count        byte
	Enchantments []EnchantmentStone `struct:"sizefrom=Count"`
}

type Licences struct {
	Max      byte
	Used     byte
	Count    byte
	Licences []Licence `struct:"sizefrom=Count"`
}

type Licence struct {
	Name   string
	Damage byte
	// mobs killed required
	NextLevelReq uint32
	// current count of killed mobs
	CurrentCount uint32
}

// These are added stats to the ones defined in SHN files
type ItemAttributes struct {
	Strength     Stat
	Dexterity    Stat
	Endurance    Stat
	Spirit       Stat
	Intelligence Stat
	//CriticalRate Stat
	AimBonus Stat
	MDefense Stat
	PDefense Stat
	Aim      Stat
	Evasion  Stat
	MAttack  Stat
	PAttack  Stat
	MaxHP    Stat
	MaxSP    Stat
}

type Stat struct {
	Base uint32
	// added value that goes into parentheses
	Extra uint32
}

type EnchantmentStone struct {
}

func getItemWhere(db *pg.DB, clauses map[string]interface{}, deleted bool) (*Item, error) {
	var item Item

	query := db.Model(&item)

	for k, v := range clauses {
		query.Where(k, v)
	}

	if !deleted {
		query.Where("deleted_at = null")
	}

	err := query.Select()

	return &item, err
}
