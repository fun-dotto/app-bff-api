package domain

import "time"

// PersonalCalendarItem 個人カレンダーアイテム
type PersonalCalendarItem struct {
	Date          time.Time
	Slot          TimetableSlot
	TimetableItem TimetableItem
}
