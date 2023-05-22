package tools

import "time"

func GetReleaseDates(period string) (time.Time, time.Time) {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	startDate := now
	endDate := now

	if period == "past" {
		startDate = startDate.AddDate(0, 0, -14)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, -7)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "week" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "month" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, 28)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "extended" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, 84)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	}

	return startDate, endDate
}
