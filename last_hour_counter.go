package counters

import (
	"time"
)

type LastHourCounter struct {
	lastHourUnixTimestamp int64
	countPerBucket        []int
	minutesInBucket       int
	Timer                 UnixTimer
}

type UnixTimer interface {
	CurrentUnixTime() int64
}

type unixTimeGetter struct{}

func (unixTimeGetter) CurrentUnixTime() int64 {
	return time.Now().Unix()
}

const defaultMinutesInBucket = 4

func (c *LastHourCounter) initialize() {
	if c.Timer == nil {
		c.Timer = unixTimeGetter{}
	}
	if c.minutesInBucket == 0 {
		c.minutesInBucket = defaultMinutesInBucket
	}
	if 60%c.minutesInBucket != 0 {
		panic("Minutes in bucket should be a divisor of 60")
	}
	c.countPerBucket = make([]int, 60/c.minutesInBucket)
	c.lastHourUnixTimestamp = c.Timer.CurrentUnixTime() - 60*60
}

func (c *LastHourCounter) Value() int {
	if c.lastHourUnixTimestamp == 0 {
		c.initialize()
	}
	currUnixTimestamp := c.Timer.CurrentUnixTime()
	c.adjustStateToCurrentTime(currUnixTimestamp)

	var total int
	for _, count := range c.countPerBucket {
		total += count
	}
	return total
}

func (c *LastHourCounter) Increment() {
	if c.lastHourUnixTimestamp == 0 {
		c.initialize()
	}
	currUnixTimestamp := c.Timer.CurrentUnixTime()
	c.adjustStateToCurrentTime(currUnixTimestamp)

	c.countPerBucket[len(c.countPerBucket)-1] = c.countPerBucket[len(c.countPerBucket)-1] + 1
}

func (c *LastHourCounter) adjustStateToCurrentTime(currUnixTimestamp int64) {
	hourAgo := currUnixTimestamp - 60*60
	if hourAgo < c.lastHourUnixTimestamp {
		// nothing to do
		return
	}
	minutesPassedToHourAgo := int((hourAgo - c.lastHourUnixTimestamp) / 60)

	bucketsToHourAgo := (minutesPassedToHourAgo / c.minutesInBucket) + 1
	c.lastHourUnixTimestamp = c.lastHourUnixTimestamp + int64(bucketsToHourAgo*c.minutesInBucket*60)

	i := 0
	for ; i < len(c.countPerBucket)-bucketsToHourAgo; i++ {
		c.countPerBucket[i] = c.countPerBucket[i+bucketsToHourAgo]
	}
	for ; i <= len(c.countPerBucket)-1; i++ {
		c.countPerBucket[i] = 0
	}
}
