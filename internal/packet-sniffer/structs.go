package packet_sniffer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/go-restruct/restruct"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

var ocs *opCodeStructs

type opCodeStructs struct {
	structs map[uint16]string
	mu      sync.Mutex
}

type ncRepresentation struct {
	UnpackedData string `json:"unpacked_data"`
}

func generateOpCodeSwitch() {
	type processedStructs struct {
		List map[uint16]bool `json:"processedStructs"`
	}
	// load processed structs
	filePath, err := filepath.Abs("config/processed-structs.json")
	if err != nil {
		log.Fatal(err)
	}

	d, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var ps processedStructs
	err = json.Unmarshal(d, &ps)
	if err != nil {
		log.Fatal(err)
	}
	start := `
	package generated
	func ncStructRepresentation(opCode uint16, data []byte) {
	switch opCode {` + "\n"
	ocs.mu.Lock()
	fmt.Println(ocs.structs)
	for k, v := range ocs.structs {
		if _, processed := ps.List[k]; processed {
			continue
		}
		caseStmt := fmt.Sprintf("\t"+`case %v:`+"\n", k)
		caseStmt += fmt.Sprintf("\t"+"// %v\n", v)
		caseStmt += "\t" + "// return ncStructData(&nc, data)\n"
		caseStmt += "\t" + "break\n"
		start += caseStmt
	}
	ocs.mu.Unlock()
	end := "}}"

	pathName, err := filepath.Abs("output/opcodes-switch.go")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(pathName, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write([]byte(start + end))
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}

func ncStructData(nc interface{}, data []byte) (ncRepresentation, error) {
	err := structs.Unpack(data, nc)
	if err != nil {
		log.Error(err)
		n, err := restruct.SizeOf(nc)
		if err != nil {
			log.Error(err)
		}
		hexString := hex.EncodeToString(data)
		log.Error(hexString)
		log.Errorf("struct: %v, size: %v", reflect.TypeOf(nc).String(), n)
		return ncRepresentation{}, err
	}

	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	nr := ncRepresentation{
		UnpackedData: string(sd),
	}

	return nr, nil
}
