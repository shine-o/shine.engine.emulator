package lib

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"time"
)

var (
	log    *logger.Logger
	xorKey []byte
)

// open TCP socket on port 9010
// for each connection use Context
func Listen(cmd *cobra.Command, args []string) {
	serveConfig()
	log = logger.Init("LoginLogger", true, true, ioutil.Discard)

	if l, err := net.Listen("tcp4", fmt.Sprintf(":%v", viper.GetInt("serve.port"))); err == nil {
		log.Infof("Listening for TCP connections on: %v", l.Addr())
		defer l.Close()
		rand.Seed(time.Now().Unix())
		for {
			if c, err := l.Accept(); err == nil {
				go handleConnection(c)
			} else {
				logger.Fatal(err)
			}
		}
	} else {
		logger.Error(err)
	}
}

func serveConfig() {
	if b, err := hex.DecodeString(viper.GetString("crypt.xorKey")); err != nil {
		log.Error(err)
	} else {
		xorKey = b
	}
}

func handleConnection(c net.Conn) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	log.Infof("Serving %v", c.RemoteAddr().String())

	defer c.Close()
	defer cancel()
	defer log.Infof("Closing connection %v", c.RemoteAddr().String())

	var (
		buf       = make([]byte, 1024)
		r         = bufio.NewReader(c)
		w         = bufio.NewWriter(c)
		xorOffset uint16
		segment   = make(chan []byte)
	)

	xorOffset = randomXorKey(viper.GetInt("crypt.xorLimit"))
	log.Infof("XorKey: %v", xorOffset)
	if b, err := hex.DecodeString("040708f600"); err != nil {
		log.Error(err)
	} else {
		if _, err := w.Write(b); err != nil {
			log.Error(err)
		} else {
			if err = w.Flush(); err != nil {
				log.Error(err)
			}
		}
	}
	xorOffset = 246
	go handleSegments(ctx, segment, &xorOffset)

	var data []byte
	for {
		if n, err := r.Read(buf); err == nil {
			log.Infof("Received %v bytes", n)
			data = append(data, buf[:n]...)
			segment <- buf[:n]
		} else {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
				return
			}
		}
	}
	log.Info(len(data))
}

// Read packet data from segments
func handleSegments(ctx context.Context, segment <-chan []byte, xorOffset *uint16) {
	var (
		data   []byte
		offset int
	)
	offset = 0

	for {
		select {
		case <-ctx.Done():
			return
		case b := <-segment:
			data = append(data, b...)

			if offset > len(data) {
				break
			}

			if offset != len(data) {
				var skipBytes int
				var pLen int
				var pType string
				var pd []byte

				pLen, pType = packetBoundary(offset, data)

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

				xorCipher(pd, xorOffset)
				pc := processPacket(pType, pLen, pd)
				log.Infof("Got one %v", pc.pcb.String())
				offset += skipBytes + pLen
			}
		}
	}
}

// read packet data
// if xorKey is detected in a server flow (packets coming from the server), that is if header == 2055, notify the converse flow
// create PC struct with packet headers + data
func processPacket(pType string, pLen int, packetData []byte) PC {
	var opCode, department, command uint16
	br := bytes.NewReader(packetData)
	binary.Read(br, binary.LittleEndian, &opCode)

	department = opCode >> 10
	command = opCode & 1023

	return PC{
		pcb: ProtocolCommandBase{
			packetType:    pType,
			length:        pLen,
			department:    department,
			command:       command,
			operationCode: opCode,
			data:          packetData,
		},
	}
}

// find out if big or small packet
// return length and type
func packetBoundary(offset int, b []byte) (int, string) {
	if b[offset] == 0 {
		var pLen uint16
		var tempB []byte
		tempB = append(tempB, b[offset:]...)
		br := bytes.NewReader(tempB)
		br.ReadAt(tempB, 1)
		binary.Read(br, binary.LittleEndian, &pLen)
		return int(pLen), "big"
	} else {
		var pLen uint8
		pLen = b[offset]
		return int(pLen), "small"
	}
}

// decrypt encrypted bytes using captured xorKey and xorTable
func xorCipher(eb []byte, xorPos *uint16) {
	xorLimit := uint16(viper.GetInt("crypt.xorLimit"))
	for i, _ := range eb {
		eb[i] ^= xorKey[*xorPos]
		*xorPos++
		if *xorPos >= xorLimit {
			*xorPos = 0
		}
	}
}
