package service

import (
	"encoding/hex"
	"fmt"
	"github.com/mitchellh/go-homedir"
	protocol "github.com/shine-o/shine.engine.protocol"
	"github.com/spf13/viper"
	"os"
	"testing"
)

var cfgFile string

func TestMain(m *testing.M) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".shine.engine.login" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".shine.engine.login")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
		log.Infof("Using config file: %v", viper.ConfigFileUsed())
	}

	viper.SetDefault("serve.port", 9010)
	viper.SetDefault("crypt.xorKey", "0759694a941194858c8805cba09ecd583a365b1a6a16febddf9402f82196c8e99ef7bfbdcfcdb27a009f4022fc11f90c2e12fba7740a7d78401e2ca02d06cba8b97eefde49ea4e13161680f43dc29ad486d7942417f4d665bd3fdbe4e10f50f6ec7a9a0c273d2466d322689c9a520be0f9a50b25da80490dfd3e77d156a8b7f40f9be80f5247f56f832022db0f0bb14385c1cba40b0219dff08becdb6c6d66ad45be89147e2f8910b89360d860def6fe6e9bca06c1759533cfc0b2e0cca5ce12f6e5b5b426c5b2184f2a5d261b654df545c98414dc7c124b189cc724e73c64ffd63a2cee8c8149396cb7dcbd94e232f7dd0afc020164ec4c940ab156f5c9a934de0f3827bc81300f7b3825fee83e29ba5543bf6b9f1f8a4952187f8af888245c4fe1a830878e501f2fd10cb4fd0abcdc1285e252ee4a5838abffc63db960640ab450d54089179ad585cfec0d7e817fe3c3040122ec27ccfa3e21a654c8de00b6df279ff625340785bfa7a5a5e0830c3d5d2040af60a36456f305c41c7d3798c3e85a6e5885a49a6b6af4a37b619b09401e604b32d951a4fef95d4e4afb4ad47c330233d59dce5baa5a7cd8f805fa1f2b8c725750ae6c1989ca01fcfc299b61126863654626c45b50aa2bbeef9a790223752c2013fdd95a7623f10bb5b859f99f7ae606e9a53ab450bf165898b39a6e36ee8deb")
	viper.SetDefault("crypt.xorLimit", 350)
	viper.SetDefault("protocol.nc-data", "defaults/protocol-commands.yml")

	setup()

	os.Exit(m.Run())
}

// for an outbound stream, check that each detected packet has an available handler
func TestOutboundStreamHandlers(t *testing.T) {
	offset := 0
	stream := "040708b50004670c010003050c01270a0c0200494e4954494f000000000000000000000601504147454c000000000000000000000002550c0c063139322e3136382e312e3234380000009623011c4a3301705c193a097f160336ac27842f36596251431edf17fc48c65f3b51d57ec84597406664b97231776445294a2d27dd6186418d24e6213d71003db97b"

	data, _ := hex.DecodeString(stream)

	for offset != len(data) {
		var (
			skipBytes int
			pLen      int
			pType     string
			pd        []byte
		)

		pLen, pType = protocol.PacketBoundary(offset, data)

		if pType == "small" {
			skipBytes = 1
		} else {
			skipBytes = 3
		}

		nextOffset := offset + skipBytes + pLen
		if nextOffset > len(data) {
			break
		}

		pd = append(pd, data[offset+skipBytes:nextOffset]...)
		offset += skipBytes + pLen
		if pc, err := protocol.DecodePacket(pType, pLen, pd); err != nil {
			t.Fatal(err)
		} else {
			if _, ok := hw.handlers[pc.Base.OperationCode()]; !ok {
				t.Fatalf("handler not found for operationCode %v", pc.Base.OperationCode())
			}
		}
	}
}

// for an inbound stream, check that each detected packet has an available handler
func TestInboundStreamHandlers(t *testing.T) {
	var xorOffset uint16
	xorOffset = 181

	offset := 0
	stream := "42109901fff187d1fd94f823c2d4838611c5c1faeac4d1d8e49ab299b9d0840474336747119cc7741b2564fed63a2cb6709849b990aedcbd94e232f7dd0afc1a14cbe6003e0116986bd53b9ca7a934de0f3827bc81300f7b3825fee83e29ba5543bf6b9f1f8a4952187f8af888245c4fe1a830878e501f2fd10cb4fd0abcdc1285e252ee4a5838abffc63db960640ab450d54089179ad585cfec0d7e817fe3c3040122ec27ccfa3e21a654c8de0759694a941194858c8805cba09ecd583a365b1a6a16febddf9402f82196c8e99ef7bfbdcfcdb27a009f4022fc11f90c2e12fba7740a7d78401e2ca02d06cba8b97eefde49ea4e13161680f43dc29ad486d7942417f4d665bd3fdbe4e10f50f6ec7a9a0c273d2466d322689c9a520be0f9a50b25da80490dfd3e77d156a8b7f40f9be80f5247f56f832022db0f0bb14385c1cba40b0219dff08becdb6c6d669f748cba26181db027d9a657b955bfc1ca5da3f332a045f007aef882d1aac6fd12f6e5b5fb54acd571214b31261b654df545c98414dc7c12204f1481f417a50950cc940a6fafbac47e7a58869988d0d376c7eb3fcd313153ec0347980a"

	data, _ := hex.DecodeString(stream)

	for offset != len(data) {
		var (
			skipBytes int
			pLen      int
			pType     string
			pd        []byte
		)

		pLen, pType = protocol.PacketBoundary(offset, data)

		if pType == "small" {
			skipBytes = 1
		} else {
			skipBytes = 3
		}

		nextOffset := offset + skipBytes + pLen
		if nextOffset > len(data) {
			break
		}

		pd = append(pd, data[offset+skipBytes:nextOffset]...)
		protocol.XorCipher(pd, &xorOffset)

		offset += skipBytes + pLen
		if pc, err := protocol.DecodePacket(pType, pLen, pd); err != nil {
			t.Fatal(err)
		} else {
			if _, ok := hw.handlers[pc.Base.OperationCode()]; !ok {
				t.Fatalf("handler not found for operationCode %v", pc.Base.OperationCode())
			}
		}
	}
}
