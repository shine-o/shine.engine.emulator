package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"testing"
)

func Test_Move_Entity_A_B(t *testing.T) {
	zm, err := loadMap(1)

	if err != nil {
		t.Fatal(err)
	}

	persistence.CleanDB()

	char := persistence.NewDummyCharacter("mage", false)

	p := &player{
		baseEntity: baseEntity{},
	}

	err = p.load(char.Name)

	zm.entities.players.add(p)

	x := 4089
	y := 3214

	err = p.move(zm, x, y)

	if err != nil {
		t.Fatal(err)
	}

	if p.baseEntity.current.x != x || p.baseEntity.current.y != y {
		t.Fatalf("mismatched coordinates %v %v", p.baseEntity.current.x, p.baseEntity.current.y)
	}

}

func Test_Move_Entity_Collision(t *testing.T) {

}

func Test_Move_Entity_A_B_SpeedHack(t *testing.T) {
	// generate raw path from A to B
	//
}

func Test_Bitmap_Coordinates_Conversion(t *testing.T) {
	ogx := 4089
	ogy := 3214

	bx, by := bitmapCoordinates(ogx, ogy)

	gx, gy := gameCoordinates(bx, by)

	log.Infof("%v %v", gx, gy)

	if gx != 4087 ||  gy != 3212  {
		t.Errorf("mismatched coordinates bx=%v by=%v ogx=%v ogy=%v gx=%v gy=%v",bx, by, ogx, ogy, gx, gy)
	}
}

func Test_Add_Entity_Within_Range(t *testing.T) {}

func Test_Remove_Entity_Outside_Range(t *testing.T) {

}