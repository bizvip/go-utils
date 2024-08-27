/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package snowflake

import (
	"time"

	"github.com/prometheus/common/log"
)

var f1, f2 *Sonyflake
var startTime time.Time

func init() {
	startTime = time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC)
	f1 = NewSonyflake(Settings{StartTime: startTime, MachineID: func() (uint16, error) {
		return 1, nil
	}})
	f2 = NewSonyflake(Settings{StartTime: startTime, MachineID: func() (uint16, error) {
		return 1, nil
	}})
}

func QuickID() uint64 {
	id, err := f1.NextID()
	if err != nil {
		log.Fatalf("Error getting ID: %s", err)
		return 0
	}
	return id
}

func QuickID2() uint64 {
	id, err := f2.NextID()
	if err != nil {
		log.Errorf("Error getting ID: %s", err)

	}
	return id
}
