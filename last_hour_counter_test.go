package counters

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
	"time"
)

type unixTimerMock struct {
	i             int
	timesToReturn []int64
}

func NewUnixTimerMock(timesToReturn []int64) UnixTimer {
	return &unixTimerMock{
		timesToReturn: timesToReturn,
	}
}

func (m *unixTimerMock) CurrentUnixTime() int64 {
	if m.i >= len(m.timesToReturn) {
		panic("unexpected CurrentUnixTime call")
	}
	res := m.timesToReturn[m.i]
	m.i = m.i + 1
	return res
}

func TestNoMocksBehavior(t *testing.T) {
	var cases = []int{1, 2, 3, 4, 5, 6, 10, 12}

	for _, minutesInBucket := range cases {
		t.Run(fmt.Sprintf("%d minutes in bucket", minutesInBucket), func(t *testing.T) {
			var counter = LastHourCounter{minutesInBucket: minutesInBucket}
			counter.Increment()
			counter.Increment()
			val := counter.Value()
			if val != 2 {
				t.Fatalf("Expected counter value to be 2, but got %d", val)
			}
		})
	}
}

func TestInitializationErrorOnNonDivisorBucket(t *testing.T) {
	assert.PanicsWithValue(t, "Minutes in bucket should be a divisor of 60", func() {
		var counter = LastHourCounter{minutesInBucket: 8}
		counter.Increment()
	})
}

func TestLastHourCounter_SameMinute(t *testing.T) {
	var cases = []int{1, 2, 3, 4, 5, 6, 10, 12}

	for _, minutesInBucket := range cases {
		t.Run(fmt.Sprintf("%d minutes in bucket", minutesInBucket), func(t *testing.T) {
			currentTime := time.Now()

			times := []int64{
				currentTime.Add(-100 * time.Minute).Unix(),
				currentTime.Add(-100 * time.Minute).Unix(),
				currentTime.Add(-100 * time.Minute).Unix(),
			}
			timerMock := NewUnixTimerMock(times)

			counter := LastHourCounter{Timer: timerMock, minutesInBucket: minutesInBucket}
			counter.Increment()
			res := counter.Value()
			if res != 1 {
				t.Fatalf("Expected counter value to be 1, but got %d", res)
			}
		})
	}
}

func TestLastHourCounter_ExactlyHourPassed(t *testing.T) {
	var cases = []int{1, 2, 3, 4, 5, 6, 10, 12}

	for _, minutesInBucket := range cases {
		t.Run(fmt.Sprintf("%d minutes in bucket", minutesInBucket), func(t *testing.T) {
			currentTime := time.Now()

			times := []int64{
				currentTime.Add(-60 * time.Minute).Unix(),
				currentTime.Add(-60*time.Minute + time.Duration(minutesInBucket)*time.Minute - 1*time.Second).Unix(),
				currentTime.Unix(),
			}
			timerMock := NewUnixTimerMock(times)

			counter := LastHourCounter{Timer: timerMock, minutesInBucket: minutesInBucket}
			counter.Increment()
			res := counter.Value()
			if res != 0 {
				t.Fatalf("Expected counter value to be 0, but got %d", res)
			}
		})
	}
}

func TestLastHourCounter_MinuteBeforeHourPassed(t *testing.T) {
	var cases = []int{1, 2, 3, 4, 5, 6, 10, 12}

	for _, minutesInBucket := range cases {
		t.Run(fmt.Sprintf("%d minutes in bucket", minutesInBucket), func(t *testing.T) {
			currentTime := time.Now()

			times := []int64{
				currentTime.Add(-60 * time.Minute).Unix(),
				currentTime.Add(-60*time.Minute + time.Duration(minutesInBucket)*time.Minute).Unix(),
				currentTime.Unix(),
			}
			timerMock := NewUnixTimerMock(times)

			counter := LastHourCounter{Timer: timerMock, minutesInBucket: minutesInBucket}
			counter.Increment()
			res := counter.Value()
			if res != 1 {
				t.Fatalf("Expected counter value to be 1, but got %d", res)
			}
		})
	}
}

func TestLastHourCounter(t *testing.T) {
	var cases = []int{1, 2, 3, 4, 5, 6, 10, 12}

	for _, minutesInBucket := range cases {
		t.Run(fmt.Sprintf("%d minutes in bucket", minutesInBucket), func(t *testing.T) {
			currentTime := time.Now()

			times := []int64{
				currentTime.Add(-90 * time.Minute).Unix(),
				currentTime.Add(-90 * time.Minute).Unix(),
				currentTime.Add(-80 * time.Minute).Unix(),
				currentTime.Add(-70 * time.Minute).Unix(),
				currentTime.Add(-50 * time.Minute).Unix(),
				currentTime.Add(-30 * time.Minute).Unix(),
				currentTime.Unix(),
			}
			timerMock := NewUnixTimerMock(times)

			counter := LastHourCounter{Timer: timerMock, minutesInBucket: minutesInBucket}
			// First call is for timer initialization, last call is for Value()
			for i := 2; i <= len(times)-1; i++ {
				counter.Increment()
			}
			res := counter.Value()
			if res != 2 {
				t.Fatalf("Expected counter value to be 2, but got %d", res)
			}
		})
	}
}

func TestLastHourCounter2(t *testing.T) {
	var cases = []int{1, 2, 3, 4, 5, 6, 10, 12}

	for _, minutesInBucket := range cases {
		t.Run(fmt.Sprintf("%d minutes in bucket", minutesInBucket), func(t *testing.T) {
			currentTime := time.Now()

			times := []int64{
				currentTime.Add(-120 * time.Minute).Unix(),
				currentTime.Add(-90 * time.Minute).Unix(),
				currentTime.Add(-80 * time.Minute).Unix(),
				currentTime.Add(-79 * time.Minute).Unix(),
				currentTime.Add(-75 * time.Minute).Unix(),
				currentTime.Add(-70 * time.Minute).Unix(),
				currentTime.Add(-65 * time.Minute).Unix(),
				currentTime.Add(-64 * time.Minute).Unix(),
				currentTime.Add(-63 * time.Minute).Unix(),
				currentTime.Add(-62 * time.Minute).Unix(),
				currentTime.Add(-61 * time.Minute).Unix(),
				currentTime.Add(-60 * time.Minute).Unix(),
				currentTime.Add(-59 * time.Minute).Unix(),
				currentTime.Add(-50 * time.Minute).Unix(),
				currentTime.Add(-30 * time.Minute).Unix(),
				currentTime.Add(-10 * time.Minute).Unix(),
				currentTime.Add(-2 * time.Second).Unix(),
				currentTime.Unix(),
			}
			timerMock := NewUnixTimerMock(times)

			counter := LastHourCounter{Timer: timerMock, minutesInBucket: minutesInBucket}
			// First call is for timer initialization, last call is for Value()
			expectedCount := 0
			for i := 2; i <= len(times)-1; i++ {
				counter.Increment()
				if i < len(times)-1 && currentTime.Unix()-times[i] < int64(60*60-minutesInBucket*60+1) {
					expectedCount++
				}
			}
			res := counter.Value()
			if res != expectedCount {
				t.Fatalf("Expected counter value to be %d, but got %d", expectedCount, res)
			}
		})
	}
}

func TestLastHourCounter_RandomTimes(t *testing.T) {
	var cases = []int{1, 2, 3, 4, 5, 6, 10, 12}

	for _, minutesInBucket := range cases {
		t.Run(fmt.Sprintf("%d minutes in bucket", minutesInBucket), func(t *testing.T) {
			currentTime := time.Now()
			startTime := currentTime.Add(-4 * time.Hour).Unix()
			finishTime := currentTime.Add(1 * time.Hour).Unix()

			expectedCount := 0
			var timesInMinutes []int
			for i := 1; i <= 500; i++ {
				newTime := rand.Intn(int(finishTime - startTime))
				if int(finishTime-startTime)-newTime < 60*60-minutesInBucket*60+1 {
					expectedCount++
				}
				timesInMinutes = append(timesInMinutes, newTime)
			}
			sort.Ints(timesInMinutes)

			var times []int64
			times = append(times, startTime)
			for _, minutes := range timesInMinutes {
				times = append(times, startTime+int64(minutes))
			}
			times = append(times, finishTime)

			timerMock := NewUnixTimerMock(times)
			counter := LastHourCounter{Timer: timerMock, minutesInBucket: minutesInBucket}
			// First call is for timer initialization, last call is for Value()
			for i := 2; i <= len(times)-1; i++ {
				counter.Increment()
			}
			res := counter.Value()
			if res != expectedCount {
				t.Fatalf("Expected counter value to be %d, but got %d", expectedCount, res)
			}
		})
	}
}
