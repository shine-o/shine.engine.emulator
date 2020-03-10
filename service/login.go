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
	ch[3173] = userClientVersionCheckReq
	ch[3175] = userClientVersionCheckAck
	ch[3162] = userUsLoginReq
	ch[3081] = userLoginFailAck
	ch[3082] = userLoginAck
	ch[3076] = userXtrapReq
	ch[3077] = userXtrapAck
	ch[3099] = userWorldStatusReq
	ch[3100] = userWorldStatusAck
	ch[3083] = userWorldSelectReq
	ch[3084] = userWorldSelectAck
	ch[3096] = userNormalLogoutCmd

	hw := protocol.NewHandlerWarden(ch)

	ss := protocol.NewShineService(s, hw)

	ss.Listen(ctx)
}
