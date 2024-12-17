// Forked from github.com/docker/go-units
package duration

import (
	"fmt"
	"time"
)

const (
	hoursPerDay         = 24
	daysPerWeek         = 7
	daysPerMonth        = 30
	threeSixtyFive      = 365
	roundingMarginToOne = 0.5

	minSeconds = 1
	maxSeconds = 60
	maxMinutes = 60
	maxHours   = 48
	maxDays    = 24 * 7 * 2
	maxWeeks   = 24 * 30 * 2
	maxYears   = 24 * 365 * 2
)

// HumanDuration returns a human-readable approximation of a duration
// (eg. "About a minute", "4 hours ago", etc.).
func HumanDuration(duration time.Duration) string {
	switch seconds := int(duration.Seconds()); {
	case seconds < minSeconds:
		return "Less than a second"
	case seconds == 1:
		return "1 second"
	case seconds < maxSeconds:
		return fmt.Sprintf("%d seconds", seconds)
	}

	switch minutes := int(duration.Minutes()); {
	case minutes == 1:
		return "About a minute"
	case minutes < maxMinutes:
		return fmt.Sprintf("%d minutes", minutes)
	}

	switch hours := int(duration.Hours() + roundingMarginToOne); {
	case hours == 1:
		return "About an hour"
	case hours < maxHours:
		return fmt.Sprintf("%d hours", hours)
	case hours < maxDays:
		return fmt.Sprintf("%d days", hours/hoursPerDay)
	case hours < maxWeeks:
		return fmt.Sprintf("%d weeks", hours/hoursPerDay/daysPerWeek)
	case hours < maxYears:
		return fmt.Sprintf("%d months", hours/hoursPerDay/daysPerMonth)
	}

	return fmt.Sprintf("%d years", int(duration.Hours())/hoursPerDay/threeSixtyFive)
}
