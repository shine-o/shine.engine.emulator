package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/google/logger"
	ps "github.com/shine-o/shine.engine.emulator/internal/packet-sniffer"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"math"
)

// used to run with debugger and check the rows are loaded correctly
func main() {
	//shn sample
	//shnTest()
	//packet data sample
	packetDataTest()
	//shbdTest()
	//itemAttributes()
}

func itemAttributes() {
	//ds := "2609230c0024b5b300000000000000000801248675e09304000802245975a0860100050324750e110504247f0e10050524840e07050624890e240507248e0e1c050824980e120509249d0e12050a24a20e12050b24790e02050c24487501050d24d8750b050e24847705050f24ca0903051024c90901171124f18b2c01ffecbb7600000000000000000000050000051224f81a14051324331b000c1424a5bc0000000000000000051524de1e1f48162440a009000000000000250000000000ffff00000000ffff000000000000417664614b656461767261000000000000000000008aeeff53eeffffffff0200000000000103030e001a1724dcd600000000000000000000000007030c000210000110001a1824ddd600000000000000000000000007010a00040600030c00171924ded600000000000000000000000005030400010a00301a2422a20000000000000000000000010000000000000000000000000000000000000000000000000005030900040900301b2425a20000000000000000000000010000000000000000000000000000000000000000000000000005020300030a002a1c24d5080000000000000000000000000000000000000000000000000000000000000000000000000000051d24591f14051e2453ee13061f24df0c010006202459810100112124e7050000000000000000000000000130222466a500000000000000000900000901030004070002010003020000000000000000000000000001050302000207004e232469e100000000000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000000000000000ffffffffffffffffff0200000000000007020c000113000007001a282413cf000000000000000000000000070003000409000210004e29244aa000000000000000ffff00000000ffff00000000ffff000000000000000000000000000000000000000000000000000000ffffffffffffffffff0100000000000007030d00040c00010400"
	ds := "0109d3480024fa000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003050500"

	d, _ := hex.DecodeString(ds)

	nc := structs.NcCharClientItemCmd{}

	err := structs.Unpack(d, &nc)

	if err != nil {
		logger.Error(err)
	}

	attr := structs.ShineItemAttrAmulet{
		DeleteTime:       0,
		IsBound:          0,
		Upgrade:          0,
		Strengthen:       0,
		UpgradeFailCount: 0,
		//UpgradeOption:            structs.ItemOptionStorage{},
		RandomOptionChangedCount: 0,
	}
	b, err := structs.Pack(&attr)

	if err != nil {
		logger.Error(err, b)
	}

	//tnc := structs.ShineItemAttrAmulet{}
	//err = structs.Unpack(nc.Items[34].ItemAttr, &tnc)

	tnc := structs.ShineItemAttrWeapon{}
	err = structs.Unpack(nc.Items[0].ItemAttr, &tnc)
	//err = structs.Unpack(nc.Items[35].ItemAttr, &tnc)

	//tnc := structs.ShineItemAttrArmor{}
	//err = structs.Unpack(nc.Items[33].ItemAttr, &tnc)

	if err != nil {
		logger.Error(err)
	}

	fmt.Println(b)
}

func packetFilter() {
	captured := make(chan ps.CapturedPacket, 1500)
	p := ps.Params{
		WatchCommands: make(map[uint16]interface{}),
		Send:          captured,
	}

	p.WatchCommands[9217] = true
	p.WatchCommands[9218] = true
	p.WatchCommands[9280] = true
	p.WatchCommands[9294] = true
	p.WatchCommands[9303] = true
	p.WatchCommands[9287] = true
	p.WatchCommands[9224] = true

	go ps.ExtendedCapture(&p)

	for {
		select {
		case cp := <-captured:
			var jData, structName string
			switch cp.Command.Base.OperationCode {
			case 9217:
				jData = packetData(cp.Command.Base.Data, &structs.NcBatTargetInfoReq{})
				structName = "NcBatTargetInfoReq"
			case 9218:
				jData = packetData(cp.Command.Base.Data, &structs.NcBatTargetInfoCmd{})
				structName = "NcBatTargetInfoCmd"

			case 9280:
				jData = packetData(cp.Command.Base.Data, &structs.NcBatSkillBashObjCastReq{})
				structName = "NcBatSkillBashObjCastReq"

			case 9294:
				jData = packetData(cp.Command.Base.Data, &structs.NcBatSkillBashHitObjStartCmd{})
				structName = "NcBatSkillBashHitObjStartCmd"

			case 9303:
				jData = packetData(cp.Command.Base.Data, &structs.NcBatSkillBashHitBlastCmd{})
				structName = "NcBatSkillBashHitBlastCmd"

			case 9287:
				jData = packetData(cp.Command.Base.Data, &structs.NcBatSwingDamageCmd{})
				structName = "NcBatSwingDamageCmd"
			case 9277:
				jData = packetData(cp.Command.Base.Data, &structs.NcBatCeaseFireCmd{})
				structName = "NcBatCeaseFireCmd"
			case 9224:
				jData = packetData(cp.Command.Base.Data, &structs.NcBatUnTargetReq{})
				structName = "NcBatUnTargetReq"
			}

			logger.Infof("%v %v %v %v", cp.Seen, cp.Direction, structName, jData)
		}
	}
}

func packetData(data []byte, nc interface{}) string {

	err := structs.Unpack(data, nc)

	if err != nil {
		logger.Error(err)
	}

	jData, err := json.Marshal(nc)

	if err != nil {
		logger.Error(err)
	}

	return string(jData)
}

//
func shnTest() {
	var mi data.ShineMobInfoServer
	err := data.Load("assets/MobInfoServer.shn", &mi)
	if err != nil {
		logger.Error(err)
	}
}

func packetDataTest() {
	var s structs.NcBriefInfoMobCmd
	//hexS := "2424252413cf00000000000000000000000007000300040900021000"
	hexS := "04040002e215470f0000b32400009600000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010002e4151d180000b911000055000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300026800e92100008d0f000022000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200028627ca140000941100005901526f75000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

	d, err := hex.DecodeString(hexS)

	if err != nil {
		logger.Error(err)
	}

	err = structs.Unpack(d, &s)

	if err != nil {
		logger.Error(err)
	}

	//var (
	//	s1 structs.CharBriefInfoNotCamp
	//	sdata = make([]byte, 45)
	//)
	//
	//for i, sd := range s.ShapeData.Data {
	//	sdata[i] = sd
	//}
	//
	//err = structs.Unpack(sdata, &s1)
	//
	//if err != nil {
	//	logger.Error(err)
	//}

}

func shbdTest() {
	m := "EldPri01"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	img, err := data.SHBDToImage(s)
	if err != nil {
		logger.Error(err)
	}

	err = data.SaveBmpFile(img, "./", m)

	if err != nil {
		logger.Error(err)
	}

	rs := data.ImageToSHBD(img)

	data.SaveSHBDFile(&rs, "./", m)

	xbm, ybm, err := walkingPositions(s)

	if err != nil {
		logger.Error(err)
	}
	testWalk(xbm, ybm)
}

func canWalk(x, y *roaring.Bitmap, rX, rY uint32) bool {
	if x.ContainsInt(int(rX)) && y.ContainsInt(int(rY)) {
		return true
	}
	return false
}

func testWalk(walkableX *roaring.Bitmap, walkableY *roaring.Bitmap) {
	igX := 5868
	igY := 10462

	rX := (igX * 8.0) / 50.0

	//fmt.Printf("%.6f", float64(rX))
	// 1589 = (x * 8) / 50
	rstX := (rX * 50.0) / 8.0

	fmt.Printf("%.6f", float64(rstX))

	rY := (igY * 8.0) / 50

	if canWalk(walkableX, walkableY, uint32(rX), uint32(rY)) {
		fmt.Printf("\nrX: %v, rY: %v", rX, rY)
		fmt.Printf("\nigX: %v, igY: %v", igX, igY)
	}
}

// WalkingPositions creates two X,Y roaring bitmaps with walkable coordinates
func walkingPositions(s *data.SHBD) (*roaring.Bitmap, *roaring.Bitmap, error) {
	walkableX := roaring.BitmapOf()
	walkableY := roaring.BitmapOf()

	r := bytes.NewReader(s.Data)

	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return walkableX, walkableY, err
			}
			for i := 0; i < 8; i++ {
				if b&byte(math.Pow(2, float64(i))) == 0 {
					rX := x*8 + i
					rY := y
					walkableX.Add(uint32(rX))
					walkableY.Add(uint32(rY))
				}
			}
		}
	}
	return walkableX, walkableY, nil
}
