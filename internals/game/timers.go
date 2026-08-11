package game

// work in progress
// func TickTimers(Timers []Timer) ([]Timer, []Event) {
// 	timers := []Timer{}
// 	events := []Event{}
// 	for _, timer := range Timers {
// 		timer.CurrFrame -= 1
// 		if timer.CurrFrame >= 0 {
// 			events = append(events, timer.Event)
// 			continue
// 		}
// 		timers = append(timers, timer)
// 	}

// 	return timers, events
// }

// func MakeTimer(setoffFrames int, e Event) Timer {
// 	return Timer{
// 		CurrFrame:    setoffFrames,
// 		SetoffFrames: setoffFrames,
// 		Event:        e,
// 	}
// }
