package duration_test

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"go.farcloser.world/core/duration"
)

func ExampleHumanDuration() {
	fmt.Println(duration.HumanDuration(450 * time.Millisecond))
	fmt.Println(duration.HumanDuration(47 * time.Second))
	fmt.Println(duration.HumanDuration(1 * time.Minute))
	fmt.Println(duration.HumanDuration(3 * time.Minute))
	fmt.Println(duration.HumanDuration(35 * time.Minute))
	fmt.Println(duration.HumanDuration(35*time.Minute + 40*time.Second))
	fmt.Println(duration.HumanDuration(1 * time.Hour))
	fmt.Println(duration.HumanDuration(1*time.Hour + 45*time.Minute))
	fmt.Println(duration.HumanDuration(3 * time.Hour))
	fmt.Println(duration.HumanDuration(3*time.Hour + 59*time.Minute))
	fmt.Println(duration.HumanDuration(3*time.Hour + 60*time.Minute))
	fmt.Println(duration.HumanDuration(24 * time.Hour))
	fmt.Println(duration.HumanDuration(24*time.Hour + 12*time.Hour))
	fmt.Println(duration.HumanDuration(2 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(7 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(13*24*time.Hour + 5*time.Hour))
	fmt.Println(duration.HumanDuration(2 * 7 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(2*7*24*time.Hour + 4*24*time.Hour))
	fmt.Println(duration.HumanDuration(3 * 7 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(4 * 7 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(4*7*24*time.Hour + 3*24*time.Hour))
	fmt.Println(duration.HumanDuration(1 * 30 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(1*30*24*time.Hour + 2*7*24*time.Hour))
	fmt.Println(duration.HumanDuration(2 * 30 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(3*30*24*time.Hour + 1*7*24*time.Hour))
	fmt.Println(duration.HumanDuration(5*30*24*time.Hour + 2*7*24*time.Hour))
	fmt.Println(duration.HumanDuration(13 * 30 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(23 * 30 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(24 * 30 * 24 * time.Hour))
	fmt.Println(duration.HumanDuration(24*30*24*time.Hour + 2*7*24*time.Hour))
	fmt.Println(duration.HumanDuration(3*365*24*time.Hour + 2*30*24*time.Hour))
	// Output:
	// Less than a second
	// 47 seconds
	// About a minute
	// 3 minutes
	// 35 minutes
	// 35 minutes
	// About an hour
	// 2 hours
	// 3 hours
	// 4 hours
	// 4 hours
	// 24 hours
	// 36 hours
	// 2 days
	// 7 days
	// 13 days
	// 2 weeks
	// 2 weeks
	// 3 weeks
	// 4 weeks
	// 4 weeks
	// 4 weeks
	// 6 weeks
	// 2 months
	// 3 months
	// 5 months
	// 13 months
	// 23 months
	// 24 months
	// 2 years
	// 3 years
}

func TestHumanDuration(t *testing.T) {
	t.Parallel()

	// Useful duration abstractions
	day := 24 * time.Hour
	week := 7 * day
	month := 30 * day
	year := 365 * day

	assert.Equal(t, "Less than a second", duration.HumanDuration(450*time.Millisecond))
	assert.Equal(t, "1 second", duration.HumanDuration(1*time.Second))
	assert.Equal(t, "45 seconds", duration.HumanDuration(45*time.Second))
	assert.Equal(t, "46 seconds", duration.HumanDuration(46*time.Second))
	assert.Equal(t, "59 seconds", duration.HumanDuration(59*time.Second))
	assert.Equal(t, "About a minute", duration.HumanDuration(60*time.Second))
	assert.Equal(t, "About a minute", duration.HumanDuration(1*time.Minute))
	assert.Equal(t, "3 minutes", duration.HumanDuration(3*time.Minute))
	assert.Equal(t, "35 minutes", duration.HumanDuration(35*time.Minute))
	assert.Equal(t, "35 minutes", duration.HumanDuration(35*time.Minute+40*time.Second))
	assert.Equal(t, "45 minutes", duration.HumanDuration(45*time.Minute))
	assert.Equal(t, "45 minutes", duration.HumanDuration(45*time.Minute+40*time.Second))
	assert.Equal(t, "46 minutes", duration.HumanDuration(46*time.Minute))
	assert.Equal(t, "59 minutes", duration.HumanDuration(59*time.Minute))
	assert.Equal(t, "About an hour", duration.HumanDuration(1*time.Hour))
	assert.Equal(t, "About an hour", duration.HumanDuration(1*time.Hour+29*time.Minute))
	assert.Equal(t, "2 hours", duration.HumanDuration(1*time.Hour+31*time.Minute))
	assert.Equal(t, "2 hours", duration.HumanDuration(1*time.Hour+59*time.Minute))
	assert.Equal(t, "3 hours", duration.HumanDuration(3*time.Hour))
	assert.Equal(t, "3 hours", duration.HumanDuration(3*time.Hour+29*time.Minute))
	assert.Equal(t, "4 hours", duration.HumanDuration(3*time.Hour+31*time.Minute))
	assert.Equal(t, "4 hours", duration.HumanDuration(3*time.Hour+59*time.Minute))
	assert.Equal(t, "4 hours", duration.HumanDuration(3*time.Hour+60*time.Minute))
	assert.Equal(t, "24 hours", duration.HumanDuration(24*time.Hour))
	assert.Equal(t, "36 hours", duration.HumanDuration(1*day+12*time.Hour))
	assert.Equal(t, "2 days", duration.HumanDuration(2*day))
	assert.Equal(t, "7 days", duration.HumanDuration(7*day))
	assert.Equal(t, "13 days", duration.HumanDuration(13*day+5*time.Hour))
	assert.Equal(t, "2 weeks", duration.HumanDuration(2*week))
	assert.Equal(t, "2 weeks", duration.HumanDuration(2*week+4*day))
	assert.Equal(t, "3 weeks", duration.HumanDuration(3*week))
	assert.Equal(t, "4 weeks", duration.HumanDuration(4*week))
	assert.Equal(t, "4 weeks", duration.HumanDuration(4*week+3*day))
	assert.Equal(t, "4 weeks", duration.HumanDuration(1*month))
	assert.Equal(t, "6 weeks", duration.HumanDuration(1*month+2*week))
	assert.Equal(t, "2 months", duration.HumanDuration(2*month))
	assert.Equal(t, "2 months", duration.HumanDuration(2*month+2*week))
	assert.Equal(t, "3 months", duration.HumanDuration(3*month))
	assert.Equal(t, "3 months", duration.HumanDuration(3*month+1*week))
	assert.Equal(t, "5 months", duration.HumanDuration(5*month+2*week))
	assert.Equal(t, "13 months", duration.HumanDuration(13*month))
	assert.Equal(t, "23 months", duration.HumanDuration(23*month))
	assert.Equal(t, "24 months", duration.HumanDuration(24*month))
	assert.Equal(t, "2 years", duration.HumanDuration(24*month+2*week))
	assert.Equal(t, "3 years", duration.HumanDuration(3*year+2*month))
}
