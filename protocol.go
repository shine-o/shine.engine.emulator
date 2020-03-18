package networking

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

type Settings struct {
	XorKey           []byte
	XorLimit         uint16
	CommandsFilePath string
}

type PCList struct {
	Departments map[uint8]Department
	mu          sync.Mutex
}

type RawPCList struct {
	Departments []Department `yaml:"departments,flow"`
}

type Department struct {
	HexId             string `yaml:"hexId"`
	Name              string `yaml:"name"`
	RawCommands       string `yaml:"commands"`
	ProcessedCommands map[string]string
}

type Command struct {
	Base     CommandBase // common data in every command, like operation code and length
	NcStruct interface{} // any kind of structure that is the representation in bytes of the network packet
}

type CommandBase struct {
	PacketType    string
	Length        int
	Department    uint16
	Command       uint16
	OperationCode uint16
	Data          []byte
}

// reassemble packet raw data
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
		//if err := binary.Write(lenBuf, binary.LittleEndian, uint16(buf.Len())); err != nil {
		if err := binary.Write(lenBuf, binary.LittleEndian, uint16(len(data))); err != nil {
			log.Fatalf("failed writing length for big packet to buffer %v", err)
		}
		header = append(header, lenBuf.Bytes()...)
	} else {
		header = append(header, byte(len(data)))
	}

	return append(header, data...)
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
	//department = opCode >> 10
	//command = opCode & 1023
	ePcb := exportedPcb{
		PacketType:    pcb.PacketType,
		Length:        pcb.Length,
		Department:    pcb.OperationCode >> 10,
		Command:       fmt.Sprintf("%X", pcb.OperationCode&1023),
		OperationCode: pcb.OperationCode,
		Data:          hex.EncodeToString(pcb.Data),
		RawData:       hex.EncodeToString(pcb.RawData()),
	}

	if (&PCList{}) != commandList {
		commandList.mu.Lock()
		if dpt, ok := commandList.Departments[uint8(ePcb.Department)]; ok {
			ePcb.FriendlyName = dpt.ProcessedCommands[ePcb.Command]
		} else {
			log.Warningf("Missing friendly name for command with: operationCode %v,  department %v, commmand %v, ", ePcb.OperationCode, ePcb.Department, ePcb.Command)
		}
		commandList.mu.Unlock()
	}

	if rawJson, err := json.Marshal(&ePcb); err != nil {
		log.Error(err)
		return ""
	} else {
		return string(rawJson)
	}
}

func (s *Settings) Set() {
	if cl, err := InitCommandList(s.CommandsFilePath); err != nil {
		log.Error(err)
	} else {
		commandList = &cl
	}
	xorKey = s.XorKey
	xorLimit = s.XorLimit
}

// struct information about captured network packets
func InitCommandList(filePath string) (PCList, error) {
	pcl := PCList{
		Departments: make(map[uint8]Department),
	}
	if d, err := ioutil.ReadFile(filePath); err != nil {
		log.Error(err)
		return PCList{}, err
	} else {

		rPcl := &RawPCList{}

		err = yaml.Unmarshal(d, rPcl)

		for _, d := range rPcl.Departments {

			dptHexVal := strings.ReplaceAll(d.HexId, "0x", "")

			dptIntVal, _ := strconv.ParseUint(dptHexVal, 16, 32)

			department := Department{
				HexId:             d.HexId,
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
	}
	return pcl, nil
}
