package networking

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)


// PCList protocol command list
// friendly names for each Department and Commands within a Department
type PCList struct {
	Departments map[uint8]Department
	mu          sync.Mutex
}

// RawPCList struct used to unmarshal the protocol commands file
type RawPCList struct {
	Departments []Department `yaml:"departments,flow"`
}

// Department type used to unmarshal data from the protocol commands file
type Department struct {
	HexID             string `yaml:"hexId"`
	Name              string `yaml:"name"`
	RawCommands       string `yaml:"commands"`
	ProcessedCommands map[string]string
}

// Command type used to unmarshal data from the protocol commands file
type Command struct {
	Base CommandBase // common data in every command, like operation code and length
	//NcStruct interface{} // any kind of structure that is the representation in bytes of the network packet
	NcStruct structs.NC // any kind of structure that is the representation in bytes of the network packet
	sync.RWMutex
}

// CommandBase type used to store decoded data from a packet
type CommandBase struct {
	Length           int
	Department       uint16
	Command          uint16
	OperationCode    uint16
	ClientStructName string
	Data             []byte
}

// RawData of a packet that contains the length, operation code and packet data
func (pcb *CommandBase) RawData() []byte {
	var header []byte
	var data []byte

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, pcb.OperationCode); err != nil {
		log.Fatalf("failed writing operation code to buffer %v", err)
	}

	data = append(data, buf.Bytes()...)
	data = append(data, pcb.Data...)

	if len(data) > 255 { // means big packet
		header = append(header, byte(0))
		lenBuf := new(bytes.Buffer)
		if err := binary.Write(lenBuf, binary.LittleEndian, uint16(len(data))); err != nil {
			log.Fatalf("failed writing length for big packet to buffer %v", err)
		}
		header = append(header, lenBuf.Bytes()...)
	} else {
		header = append(header, byte(len(data)))
	}

	return append(header, data...)
}

func (pcb *CommandBase) EncryptedRawData() []byte {
	var header []byte
	var data []byte

	data = append(data, pcb.Data...)

	if len(data) > 255 { // means big packet
		header = append(header, byte(0))
		lenBuf := new(bytes.Buffer)
		if err := binary.Write(lenBuf, binary.LittleEndian, uint16(len(data))); err != nil {
			log.Fatalf("failed writing length for big packet to buffer %v", err)
		}
		header = append(header, lenBuf.Bytes()...)
	} else {
		header = append(header, byte(len(data)))
	}

	return append(header, data...)
}

// PacketLength of a packet, which includes de operation code bytes
func (pcb *CommandBase) PacketLength() int {
	return len(pcb.Data) + 2
}

func (pcb *CommandBase) String() string {
	type exportedPcb struct {
		PacketType    string `json:"packetType"`
		Length        int    `json:"length"`
		Department    uint16 `json:"department"`
		Command       string `json:"command"`
		OperationCode uint16 `json:"opCode"`
		Data          string `json:"data"`
		RawData       string `json:"rawData"`
		FriendlyName  string `json:"friendlyName"`
	}
	ePcb := exportedPcb{
		Length:        pcb.PacketLength(),
		Department:    pcb.OperationCode >> 10,
		Command:       fmt.Sprintf("%X", pcb.OperationCode&1023),
		OperationCode: pcb.OperationCode,
		Data:          hex.EncodeToString(pcb.Data),
		RawData:       hex.EncodeToString(pcb.RawData()),
		FriendlyName:  pcb.ClientStructName,
	}
	if pcb.PacketLength() > 255 {
		ePcb.PacketType = "big"
	} else {
		ePcb.PacketType = "small"
	}

	rawJSON, err := json.Marshal(&ePcb)

	if err != nil {
		log.Error(err)
		return ""
	}
	return string(rawJSON)
}

// ExportedPcb is a utility struct for logging network packets
type ExportedPcb struct {
	PacketType    string `json:"packetType"`
	Length        int    `json:"length"`
	Department    uint16 `json:"department"`
	Command       string `json:"command"`
	OperationCode uint16 `json:"opCode"`
	Data          string `json:"data"`
	RawData       string `json:"rawData"`
	FriendlyName  string `json:"friendlyName"`
}

// JSON representation of a processed network command
func (pcb *CommandBase) JSON() ExportedPcb {
	//department = opCode >> 10
	//command = opCode & 1023
	ePcb := ExportedPcb{
		Length:        pcb.PacketLength(),
		Department:    pcb.OperationCode >> 10,
		Command:       fmt.Sprintf("%X", pcb.OperationCode&1023),
		OperationCode: pcb.OperationCode,
		Data:          hex.EncodeToString(pcb.Data),
		RawData:       hex.EncodeToString(pcb.RawData()),
		FriendlyName:  pcb.ClientStructName,
	}
	return ePcb
}



// InitCommandList from protocol commands file
func InitCommandList(filePath string)  error {
	pcl := PCList{
		Departments: make(map[uint8]Department),
	}

	d, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Error(err)
		return  err
	}

	rPcl := &RawPCList{}

	if err = yaml.Unmarshal(d, rPcl); err != nil {
		log.Error(err)
		return  err
	}

	for _, d := range rPcl.Departments {

		dptHexVal := strings.ReplaceAll(d.HexID, "0x", "")

		dptIntVal, _ := strconv.ParseUint(dptHexVal, 16, 32)

		department := Department{
			HexID:             d.HexID,
			Name:              d.Name,
			ProcessedCommands: make(map[string]string),
		}
		cmdsRaw := d.RawCommands
		cmdsRaw = strings.ReplaceAll(cmdsRaw, "\n", "")
		cmdsRaw = strings.ReplaceAll(cmdsRaw, " ", "")
		cmdsRaw = strings.ReplaceAll(cmdsRaw, "0x", "")
		cmdsRaw = strings.ReplaceAll(cmdsRaw, "\t", "")

		cmds := strings.Split(cmdsRaw, ",")

		for _, c := range cmds {
			if c == "" {
				continue
			}
			cs := strings.Split(c, "=")
			department.ProcessedCommands[cs[1]] = cs[0]
		}
		pcl.Departments[uint8(dptIntVal)] = department
	}

	commandList = &pcl
	return nil
}
