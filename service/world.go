package service

import (
	"context"
	"encoding/hex"
	"github.com/google/logger"
	protocol "github.com/shine-o/shine.engine.protocol"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var log logger.Logger

func Start(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	s := &protocol.Settings{}

	if xk, err := hex.DecodeString(viper.GetString("crypt.xorKey")); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		s.XorKey = xk
	}

	s.XorLimit = uint16(viper.GetInt("crypt.xorLimit"))

	if path, err := filepath.Abs(viper.GetString("protocol.nc-data")); err != nil {
		log.Error(err)
	} else {
		s.CommandsFilePath = path
	}

	ch := make(map[uint16]func(ctx context.Context, pc *protocol.Command))

	// for now, nada

	hw := protocol.NewHandlerWarden(ch)

	ss := protocol.NewShineService(s, hw)

	ss.Listen(ctx)
}
