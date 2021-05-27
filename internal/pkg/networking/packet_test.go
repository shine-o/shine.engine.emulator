package networking

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/crypto"
)

func TestMain(m *testing.M) {
	xorKey, _ = hex.DecodeString("0759694a941194858c8805cba09ecd583a365b1a6a16febddf9402f82196c8e99ef7bfbdcfcdb27a009f4022fc11f90c2e12fba7740a7d78401e2ca02d06cba8b97eefde49ea4e13161680f43dc29ad486d7942417f4d665bd3fdbe4e10f50f6ec7a9a0c273d2466d322689c9a520be0f9a50b25da80490dfd3e77d156a8b7f40f9be80f5247f56f832022db0f0bb14385c1cba40b0219dff08becdb6c6d66ad45be89147e2f8910b89360d860def6fe6e9bca06c1759533cfc0b2e0cca5ce12f6e5b5b426c5b2184f2a5d261b654df545c98414dc7c124b189cc724e73c64ffd63a2cee8c8149396cb7dcbd94e232f7dd0afc020164ec4c940ab156f5c9a934de0f3827bc81300f7b3825fee83e29ba5543bf6b9f1f8a4952187f8af888245c4fe1a830878e501f2fd10cb4fd0abcdc1285e252ee4a5838abffc63db960640ab450d54089179ad585cfec0d7e817fe3c3040122ec27ccfa3e21a654c8de00b6df279ff625340785bfa7a5a5e0830c3d5d2040af60a36456f305c41c7d3798c3e85a6e5885a49a6b6af4a37b619b09401e604b32d951a4fef95d4e4afb4ad47c330233d59dce5baa5a7cd8f805fa1f2b8c725750ae6c1989ca01fcfc299b61126863654626c45b50aa2bbeef9a790223752c2013fdd95a7623f10bb5b859f99f7ae606e9a53ab450bf165898b39a6e36ee8deb")
	xorLimit = 350
	os.Exit(m.Run())
}

func TestStreamPacketBoundary(t *testing.T) {
	expectedValues := make([]map[string]int, 5)

	iteration1 := make(map[string]int)
	iteration1["pLen"] = 4
	iteration1["offset"] = 5

	iteration2 := make(map[string]int)
	iteration2["pLen"] = 4
	iteration2["offset"] = 10

	iteration3 := make(map[string]int)
	iteration3["pLen"] = 3
	iteration3["offset"] = 14

	iteration4 := make(map[string]int)
	iteration4["pLen"] = 39
	iteration4["offset"] = 54

	iteration5 := make(map[string]int)
	iteration5["pLen"] = 85
	iteration5["offset"] = 140

	expectedValues[0] = iteration1
	expectedValues[1] = iteration2
	expectedValues[2] = iteration3
	expectedValues[3] = iteration4
	expectedValues[4] = iteration5

	resultValues := make([]map[string]int, 0)

	offset := 0
	stream := "040708b50004670c010003050c01270a0c0200494e4954494f000000000000000000000601504147454c000000000000000000000002550c0c063139322e3136382e312e3234380000009623011c4a3301705c193a097f160336ac27842f36596251431edf17fc48c65f3b51d57ec84597406664b97231776445294a2d27dd6186418d24e6213d71003db97b"

	data, _ := hex.DecodeString(stream)

	// for each
	for offset != len(data) {
		var skipBytes int
		var pLen uint16
		var pd []byte

		pLen, skipBytes = PacketBoundary(offset, data)

		nextOffset := offset + skipBytes + int(pLen)

		if nextOffset > len(data) {
			break
		}

		pd = append(pd, data[offset+skipBytes:nextOffset]...)
		offset += skipBytes + int(pLen)

		// test code
		currentIteration := make(map[string]int)
		currentIteration["pLen"] = int(pLen)
		currentIteration["offset"] = offset
		resultValues = append(resultValues, currentIteration)
	}
	for i, ev := range expectedValues {
		if ev["pLen"] != resultValues[i]["pLen"] {
			t.Errorf("Failed to assert that expectedValues[%v][\"pLen\"] equals resultValues[%v][\"pLen\"], that is  %v != %v ", i, i, ev["pLen"], resultValues[i]["pLen"])
		}

		if ev["offset"] != resultValues[i]["offset"] {
			t.Errorf("Failed to assert that expectedValues[%v][\"offset\"] equals resultValues[%v][\"offset\"], that is  %v != %v ", i, i, ev["offset"], resultValues[i]["offset"])
		}
	}
}

func TestDecodeOutboundStream(t *testing.T) {
	results := make([]string, 0)

	expectedResults := []string{
		`{"packetType":"small","length":4,"department":2,"command":"7","opCode":2055,"data":"b500","rawData":"040708b500","friendlyName":""}`,
		`{"packetType":"small","length":4,"department":3,"command":"67","opCode":3175,"data":"0100","rawData":"04670c0100","friendlyName":""}`,
		`{"packetType":"small","length":3,"department":3,"command":"5","opCode":3077,"data":"01","rawData":"03050c01","friendlyName":""}`,
		`{"packetType":"small","length":39,"department":3,"command":"A","opCode":3082,"data":"0200494e4954494f000000000000000000000601504147454c000000000000000000000002","rawData":"270a0c0200494e4954494f000000000000000000000601504147454c000000000000000000000002","friendlyName":""}`,
		`{"packetType":"small","length":85,"department":3,"command":"C","opCode":3084,"data":"063139322e3136382e312e3234380000009623011c4a3301705c193a097f160336ac27842f36596251431edf17fc48c65f3b51d57ec84597406664b97231776445294a2d27dd6186418d24e6213d71003db97b","rawData":"550c0c063139322e3136382e312e3234380000009623011c4a3301705c193a097f160336ac27842f36596251431edf17fc48c65f3b51d57ec84597406664b97231776445294a2d27dd6186418d24e6213d71003db97b","friendlyName":""}`,
	}

	offset := 0
	stream := "040708b50004670c010003050c01270a0c0200494e4954494f000000000000000000000601504147454c000000000000000000000002550c0c063139322e3136382e312e3234380000009623011c4a3301705c193a097f160336ac27842f36596251431edf17fc48c65f3b51d57ec84597406664b97231776445294a2d27dd6186418d24e6213d71003db97b"

	data, _ := hex.DecodeString(stream)

	// for each
	for offset != len(data) {
		var skipBytes int
		var pLen uint16
		var pd []byte

		pLen, skipBytes = PacketBoundary(offset, data)

		nextOffset := offset + skipBytes + int(pLen)

		if nextOffset > len(data) {
			break
		}

		pd = append(pd, data[offset+skipBytes:nextOffset]...)
		offset += skipBytes + int(pLen)

		if pc, err := DecodePacket(pd); err != nil {
			t.Fatal(err)
		} else {
			results = append(results, pc.Base.String())
		}
	}

	for i, er := range expectedResults {
		if er != results[i] {
			t.Errorf("Could not assert that given result %v is equal to expected result %v", results[i], er)
		}
	}
}

func TestDecodeInboundStream(t *testing.T) {
	var xorOffset uint16
	xorOffset = 181

	results := make([]string, 0)
	expectedResults := []string{
		`{"packetType":"small","length":66,"department":3,"command":"65","opCode":3173,"data":"32303135313131363134313632370073e2a5ee8cfeffffff6cfc190010a84f750c09000050fc19000100000058fc190080fc190000000000000000001815af0a","rawData":"42650c32303135313131363134313632370073e2a5ee8cfeffffff6cfc190010a84f750c09000050fc19000100000058fc190080fc190000000000000000001815af0a","friendlyName":""}`,
		`{"packetType":"big","length":318,"department":3,"command":"5A","opCode":3162,"data":"61646d696e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003231323332663239376135376135613734333839346130653461383031666333000000004f726967696e616c000000000000000000000000","rawData":"003e015a0c61646d696e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003231323332663239376135376135613734333839346130653461383031666333000000004f726967696e616c000000000000000000000000","friendlyName":""}`,
		`{"packetType":"small","length":32,"department":3,"command":"4","opCode":3076,"data":"1d3333423534334230434136453743343145354431443036353133303700","rawData":"20040c1d3333423534334230434136453743343145354431443036353133303700","friendlyName":""}`,
		`{"packetType":"small","length":3,"department":3,"command":"B","opCode":3083,"data":"00","rawData":"030b0c00","friendlyName":""}`,
	}

	offset := 0
	stream := "42109901fff187d1fd94f823c2d4838611c5c1faeac4d1d8e49ab299b9d0840474336747119cc7741b2564fed63a2cb6709849b990aedcbd94e232f7dd0afc1a14cbe6003e0116986bd53b9ca7a934de0f3827bc81300f7b3825fee83e29ba5543bf6b9f1f8a4952187f8af888245c4fe1a830878e501f2fd10cb4fd0abcdc1285e252ee4a5838abffc63db960640ab450d54089179ad585cfec0d7e817fe3c3040122ec27ccfa3e21a654c8de0759694a941194858c8805cba09ecd583a365b1a6a16febddf9402f82196c8e99ef7bfbdcfcdb27a009f4022fc11f90c2e12fba7740a7d78401e2ca02d06cba8b97eefde49ea4e13161680f43dc29ad486d7942417f4d665bd3fdbe4e10f50f6ec7a9a0c273d2466d322689c9a520be0f9a50b25da80490dfd3e77d156a8b7f40f9be80f5247f56f832022db0f0bb14385c1cba40b0219dff08becdb6c6d669f748cba26181db027d9a657b955bfc1ca5da3f332a045f007aef882d1aac6fd12f6e5b5fb54acd571214b31261b654df545c98414dc7c12204f1481f417a50950cc940a6fafbac47e7a58869988d0d376c7eb3fcd313153ec0347980a"

	data, _ := hex.DecodeString(stream)

	// for each
	for offset != len(data) {
		var skipBytes int
		var pLen uint16
		var pd []byte

		pLen, skipBytes = PacketBoundary(offset, data)

		nextOffset := offset + skipBytes + int(pLen)

		if nextOffset > len(data) {
			break
		}

		pd = append(pd, data[offset+skipBytes:nextOffset]...)

		crypto.XorCipher(pd, xorPos, &xorOffsetgi)

		offset += skipBytes + int(pLen)

		if pc, err := DecodePacket(pd); err != nil {
			t.Fatal(err)
		} else {
			results = append(results, pc.Base.String())
		}
	}

	for i, er := range expectedResults {
		if er != results[i] {
			t.Errorf("Could not assert that given result %v is equal to expected result %v", results[i], er)
		}
	}
}
