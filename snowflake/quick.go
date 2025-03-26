package snowflake

//
// import (
// 	"time"
// )
//
// var f1, f2 *Sonyflake
// var startTime time.Time
//
// func init() {
// 	startTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
// 	f1 = NewSonyflake(Settings{
// 		StartTime: startTime, MachineID: func() (uint16, error) {
// 			return 1, nil
// 		},
// 	})
// 	f2 = NewSonyflake(Settings{
// 		StartTime: startTime, MachineID: func() (uint16, error) {
// 			return 2, nil
// 		},
// 	})
// }
//
// func QuickID() uint64 {
// 	id, _ := f1.NextID()
// 	return id
// }
//
// func QuickID2() uint64 {
// 	id, _ := f2.NextID()
// 	return id
// }
