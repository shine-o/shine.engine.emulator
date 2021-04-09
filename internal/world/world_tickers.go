package world

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"time"
)

func worldTimeNotification(np *networking.Parameters) {
	//log.Infof("[world_ticks] worldTimeNotification ticker/worker %v", np.Session)
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			t := time.Now()
			second := t.Second()
			minute := t.Minute()
			hour := t.Hour()
			day := t.Day()
			month := t.Month()
			//year := t.Year()
			weekDay := t.Weekday()
			yearDay := t.YearDay()

			nc := structs.NcMiscServerTimeNotifyCmd{
				Time: structs.TM{
					Seconds:  int32(second),
					Minutes:  int32(minute),
					Hour:     int32(hour),
					MonthDay: int32(day),
					Month:    int32(month),
					//Year:     int32(year),
					Year:    120, // 120 = 2020, why? because potatoes
					WeekDay: int32(weekDay),
					YearDay: int32(yearDay),
					IsDst:   0,
				},
				TimeZone: 11,
			}
			networking.Send(np.OutboundSegments.Send, networking.NC_MISC_SERVER_TIME_NOTIFY_CMD, &nc)
		}
	}
}
