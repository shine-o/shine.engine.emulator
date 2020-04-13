package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	)

// 4097
func charLoginReq(ctx context.Context, pc *networking.Command)  {
	// get character where user_id and slot match
	// check which map the character is at and position
	// charLoginAck
}

//4099
func charLoginAck(ctx context.Context, pc *networking.Command)  {
	// query the zone master for connection info for the map
	// send it to the client
}