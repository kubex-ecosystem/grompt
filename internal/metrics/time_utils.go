// Package metrics - Time utilities for timezone and period handling
package metrics

import (
	"fmt"
	"time"
)

// TimeUtils provides utilities for timezone-aware time handling
type TimeUtils struct {
	defaultTimezone string
}

// NewTimeUtils creates a new TimeUtils instance
func NewTimeUtils(defaultTimezone string) *TimeUtils {
	if defaultTimezone == "" {
		defaultTimezone = "UTC"
	}
	return &TimeUtils{
		defaultTimezone: defaultTimezone,
	}
}

// ParseTimeRange parses a time range with timezone support
func (tu *TimeUtils) ParseTimeRange(start, end time.Time, timezone string) (TimeRange, error) {
	if timezone == "" {
		timezone = tu.defaultTimezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return TimeRange{}, fmt.Errorf("invalid timezone %s: %w", timezone, err)
	}

	// Convert times to the specified timezone
	startInTZ := start.In(loc)
	endInTZ := end.In(loc)

	if startInTZ.After(endInTZ) {
		return TimeRange{}, fmt.Errorf("start time %v is after end time %v", startInTZ, endInTZ)
	}

	return TimeRange{
		Start:    startInTZ,
		End:      endInTZ,
		Timezone: timezone,
	}, nil
}

// GetPeriodBoundaries calculates period boundaries based on granularity
func (tu *TimeUtils) GetPeriodBoundaries(baseTime time.Time, granularity string, timezone string, periodsBack int) ([]TimeRange, error) {
	if timezone == "" {
		timezone = tu.defaultTimezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %s: %w", timezone, err)
	}

	baseTimeInTZ := baseTime.In(loc)
	var periods []TimeRange

	for i := 0; i < periodsBack; i++ {
		var start, end time.Time

		switch granularity {
		case "hour":
			end = baseTimeInTZ.Add(time.Duration(-i) * time.Hour)
			start = end.Add(-time.Hour)
			// Round to hour boundaries
			end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), 0, 0, 0, loc)
			start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), 0, 0, 0, loc)

		case "day":
			end = baseTimeInTZ.AddDate(0, 0, -i)
			start = end.AddDate(0, 0, -1)
			// Round to day boundaries
			end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, loc)
			start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)

		case "week":
			// Calculate start of week (Monday)
			end = baseTimeInTZ.AddDate(0, 0, -i*7)
			weekday := int(end.Weekday())
			if weekday == 0 { // Sunday
				weekday = 7
			}
			end = end.AddDate(0, 0, -(weekday-1)) // Move to Monday
			end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, loc)
			start = end.AddDate(0, 0, -7)

		case "month":
			end = baseTimeInTZ.AddDate(0, -i, 0)
			start = end.AddDate(0, -1, 0)
			// Round to month boundaries
			end = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, loc)
			start = time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, loc)

		case "quarter":
			quarterMonth := ((baseTimeInTZ.Month()-1)/3)*3 + 1
			end = time.Date(baseTimeInTZ.Year(), quarterMonth, 1, 0, 0, 0, 0, loc)
			end = end.AddDate(0, -i*3, 0)
			start = end.AddDate(0, -3, 0)

		case "year":
			end = time.Date(baseTimeInTZ.Year()-i, 1, 1, 0, 0, 0, 0, loc)
			start = end.AddDate(-1, 0, 0)

		default:
			return nil, fmt.Errorf("unsupported granularity: %s", granularity)
		}

		periods = append(periods, TimeRange{
			Start:    start,
			End:      end,
			Timezone: timezone,
		})
	}

	return periods, nil
}

// GetBusinessHours filters time ranges to business hours only
func (tu *TimeUtils) GetBusinessHours(timeRange TimeRange, startHour, endHour int, excludeWeekends bool) ([]TimeRange, error) {
	if startHour < 0 || startHour > 23 || endHour < 0 || endHour > 23 {
		return nil, fmt.Errorf("invalid hour range: %d-%d", startHour, endHour)
	}

	loc, err := time.LoadLocation(timeRange.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %s: %w", timeRange.Timezone, err)
	}

	var businessPeriods []TimeRange
	current := timeRange.Start.In(loc)
	end := timeRange.End.In(loc)

	for current.Before(end) {
		// Skip weekends if requested
		if excludeWeekends && (current.Weekday() == time.Saturday || current.Weekday() == time.Sunday) {
			current = current.AddDate(0, 0, 1)
			current = time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, loc)
			continue
		}

		// Calculate business hours for this day
		dayStart := time.Date(current.Year(), current.Month(), current.Day(), startHour, 0, 0, 0, loc)
		dayEnd := time.Date(current.Year(), current.Month(), current.Day(), endHour, 0, 0, 0, loc)

		// Ensure we don't go beyond the original time range
		if dayStart.Before(timeRange.Start) {
			dayStart = timeRange.Start
		}
		if dayEnd.After(timeRange.End) {
			dayEnd = timeRange.End
		}

		if dayStart.Before(dayEnd) {
			businessPeriods = append(businessPeriods, TimeRange{
				Start:    dayStart,
				End:      dayEnd,
				Timezone: timeRange.Timezone,
			})
		}

		// Move to next day
		current = current.AddDate(0, 0, 1)
		current = time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, loc)
	}

	return businessPeriods, nil
}

// CalculateWorkingHours calculates working hours between two times
func (tu *TimeUtils) CalculateWorkingHours(start, end time.Time, timezone string, startHour, endHour int, excludeWeekends bool) (float64, error) {
	timeRange, err := tu.ParseTimeRange(start, end, timezone)
	if err != nil {
		return 0, err
	}

	businessPeriods, err := tu.GetBusinessHours(timeRange, startHour, endHour, excludeWeekends)
	if err != nil {
		return 0, err
	}

	var totalHours float64
	for _, period := range businessPeriods {
		duration := period.End.Sub(period.Start)
		totalHours += duration.Hours()
	}

	return totalHours, nil
}

// GetTimezoneOffset returns the UTC offset for a timezone at a specific time
func (tu *TimeUtils) GetTimezoneOffset(t time.Time, timezone string) (int, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, fmt.Errorf("invalid timezone %s: %w", timezone, err)
	}

	_, offset := t.In(loc).Zone()
	return offset, nil
}

// ConvertToUTC converts a time to UTC
func (tu *TimeUtils) ConvertToUTC(t time.Time, fromTimezone string) (time.Time, error) {
	if fromTimezone == "" {
		fromTimezone = tu.defaultTimezone
	}

	loc, err := time.LoadLocation(fromTimezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone %s: %w", fromTimezone, err)
	}

	// If time doesn't have timezone info, assume it's in fromTimezone
	if t.Location() == time.UTC {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	}

	return t.UTC(), nil
}

// ConvertFromUTC converts a UTC time to a specific timezone
func (tu *TimeUtils) ConvertFromUTC(utcTime time.Time, toTimezone string) (time.Time, error) {
	if toTimezone == "" {
		toTimezone = tu.defaultTimezone
	}

	loc, err := time.LoadLocation(toTimezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone %s: %w", toTimezone, err)
	}

	return utcTime.In(loc), nil
}

// GetPeriodDuration returns the duration for a granularity period
func (tu *TimeUtils) GetPeriodDuration(granularity string) (time.Duration, error) {
	switch granularity {
	case "hour":
		return time.Hour, nil
	case "day":
		return 24 * time.Hour, nil
	case "week":
		return 7 * 24 * time.Hour, nil
	case "month":
		return 30 * 24 * time.Hour, nil // Approximate
	case "quarter":
		return 90 * 24 * time.Hour, nil // Approximate
	case "year":
		return 365 * 24 * time.Hour, nil // Approximate
	default:
		return 0, fmt.Errorf("unsupported granularity: %s", granularity)
	}
}

// IsBusinessDay checks if a date is a business day (Monday-Friday)
func (tu *TimeUtils) IsBusinessDay(t time.Time) bool {
	weekday := t.Weekday()
	return weekday != time.Saturday && weekday != time.Sunday
}

// GetNextBusinessDay returns the next business day
func (tu *TimeUtils) GetNextBusinessDay(t time.Time, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone %s: %w", timezone, err)
	}

	next := t.In(loc).AddDate(0, 0, 1)
	for !tu.IsBusinessDay(next) {
		next = next.AddDate(0, 0, 1)
	}

	return next, nil
}

// CalculateBusinessDays returns the number of business days between two dates
func (tu *TimeUtils) CalculateBusinessDays(start, end time.Time, timezone string) (int, error) {
	timeRange, err := tu.ParseTimeRange(start, end, timezone)
	if err != nil {
		return 0, err
	}

	current := timeRange.Start
	days := 0

	for current.Before(timeRange.End) {
		if tu.IsBusinessDay(current) {
			days++
		}
		current = current.AddDate(0, 0, 1)
	}

	return days, nil
}

// FormatTimeForTimezone formats a time for display in a specific timezone
func (tu *TimeUtils) FormatTimeForTimezone(t time.Time, timezone, format string) (string, error) {
	if timezone == "" {
		timezone = tu.defaultTimezone
	}

	convertedTime, err := tu.ConvertFromUTC(t, timezone)
	if err != nil {
		return "", err
	}

	if format == "" {
		format = time.RFC3339
	}

	return convertedTime.Format(format), nil
}

// ValidateTimezone checks if a timezone is valid
func (tu *TimeUtils) ValidateTimezone(timezone string) error {
	_, err := time.LoadLocation(timezone)
	if err != nil {
		return fmt.Errorf("invalid timezone %s: %w", timezone, err)
	}
	return nil
}

// GetCommonTimezones returns a list of common timezone identifiers
func (tu *TimeUtils) GetCommonTimezones() []string {
	return []string{
		"UTC",
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"America/Sao_Paulo",
		"Europe/London",
		"Europe/Paris",
		"Europe/Berlin",
		"Europe/Amsterdam",
		"Europe/Zurich",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Asia/Seoul",
		"Asia/Kolkata",
		"Asia/Dubai",
		"Australia/Sydney",
		"Australia/Melbourne",
		"Pacific/Auckland",
	}
}

// TimeRangeOverlaps checks if two time ranges overlap
func (tr TimeRange) Overlaps(other TimeRange) bool {
	return tr.Start.Before(other.End) && other.Start.Before(tr.End)
}

// Duration returns the duration of the time range
func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// Contains checks if a time is within the time range
func (tr TimeRange) Contains(t time.Time) bool {
	return !t.Before(tr.Start) && t.Before(tr.End)
}

// Split splits a time range into smaller periods based on granularity
func (tr TimeRange) Split(granularity string) ([]TimeRange, error) {
	tu := NewTimeUtils(tr.Timezone)
	return tu.GetPeriodBoundaries(tr.End, granularity, tr.Timezone, int(tr.Duration().Hours()/24)+1)
}