package stw_test

import (
	"testing"
	"time"

	"github.com/go-tk/stw"
	"github.com/go-tk/testcase"
	"github.com/stretchr/testify/assert"
)

func TestNewSlidingTimeWindow(t *testing.T) {
	type C struct {
		period          time.Duration
		numberOfBuckets int
		expectedState   string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		testcase.DoCallback(0, t, c)

		slidingTimeWindow := stw.NewSlidingTimeWindow(c.period, c.numberOfBuckets)
		state := slidingTimeWindow.DumpAsString("")
		assert.Equal(t, c.expectedState, state)
	})

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.period = 3 * time.Second
		c.numberOfBuckets = 2
		c.expectedState = `
Period Per Bucket: 1.5s
Period: 3s
Buckets[0]:
	Number: 0
	Sum: 0
	Count: 0
Buckets[1]:
	Number: 0
	Sum: 0
	Count: 0
Total Sum: 0
Total Count: 0
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.period = 10 * time.Nanosecond
		c.numberOfBuckets = 3
		c.expectedState = `
Period Per Bucket: 4ns
Period: 12ns
Buckets[0]:
	Number: 0
	Sum: 0
	Count: 0
Buckets[1]:
	Number: 0
	Sum: 0
	Count: 0
Buckets[2]:
	Number: 0
	Sum: 0
	Count: 0
Total Sum: 0
Total Count: 0
`[1:]
	}).Run(t)
}

func TestSlidingTimeWindow_UpdateWithSample(t *testing.T) {
	type C struct {
		slidingTimeWindow *stw.SlidingTimeWindow
		expectedState     string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		slidingTimeWindow := stw.NewSlidingTimeWindow(9*time.Second, 3)
		slidingTimeWindow.UpdateWithSample(time.Unix(12, 1234567), 33)
		c.slidingTimeWindow = slidingTimeWindow

		testcase.DoCallback(0, t, c)

		state := slidingTimeWindow.DumpAsString("")
		assert.Equal(t, c.expectedState, state)
	})

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.UpdateWithSample(time.Unix(11, 1234567), 22)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(1, 1234567), 99)
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 3
	Sum: 22
	Count: 1
Buckets[1]:
	Number: 4
	Sum: 77
	Count: 2
Buckets[2]:
	Number: 2
	Sum: 0
	Count: 0
Total Sum: 99
Total Count: 3
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.UpdateWithSample(time.Unix(11, 1234567), 22)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(16, 1234567), 11)
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 3
	Sum: 22
	Count: 1
Buckets[1]:
	Number: 4
	Sum: 77
	Count: 2
Buckets[2]:
	Number: 5
	Sum: 11
	Count: 1
Total Sum: 110
Total Count: 4
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.UpdateWithSample(time.Unix(11, 1234567), 22)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(16, 1234567), 11)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(21, 1234567), 33)
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 6
	Sum: 0
	Count: 0
Buckets[1]:
	Number: 7
	Sum: 33
	Count: 1
Buckets[2]:
	Number: 5
	Sum: 11
	Count: 1
Total Sum: 44
Total Count: 2
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.UpdateWithSample(time.Unix(11, 1234567), 22)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow.UpdateWithSample(time.Unix(23, 1234567), 99)
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 6
	Sum: 0
	Count: 0
Buckets[1]:
	Number: 7
	Sum: 99
	Count: 1
Buckets[2]:
	Number: 5
	Sum: 0
	Count: 0
Total Sum: 99
Total Count: 1
`[1:]
	}).Run(t)
}

func TestSlidingTimeWindow_Update(t *testing.T) {
	type C struct {
		slidingTimeWindow *stw.SlidingTimeWindow
		expectedState     string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		slidingTimeWindow := stw.NewSlidingTimeWindow(9*time.Second, 3)
		slidingTimeWindow.UpdateWithSample(time.Unix(11, 1234567), 22)
		slidingTimeWindow.UpdateWithSample(time.Unix(12, 1234567), 33)
		slidingTimeWindow.UpdateWithSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow = slidingTimeWindow

		testcase.DoCallback(0, t, c)

		state := slidingTimeWindow.DumpAsString("")
		assert.Equal(t, c.expectedState, state)
	})

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.Update(time.Unix(18, 1234567))
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 6
	Sum: 0
	Count: 0
Buckets[1]:
	Number: 4
	Sum: 77
	Count: 2
Buckets[2]:
	Number: 5
	Sum: 0
	Count: 0
Total Sum: 77
Total Count: 2
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.Update(time.Unix(23, 1234567))
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 6
	Sum: 0
	Count: 0
Buckets[1]:
	Number: 7
	Sum: 0
	Count: 0
Buckets[2]:
	Number: 5
	Sum: 0
	Count: 0
Total Sum: 0
Total Count: 0
`[1:]
	}).Run(t)
}

func TestSlidingTimeWindow_Period_Average_Sum_Count(t *testing.T) {
	slidingTimeWindow := stw.NewSlidingTimeWindow(11*time.Second, 22)
	slidingTimeWindow.UpdateWithSample(time.Now(), 1)
	slidingTimeWindow.UpdateWithSample(time.Now(), 2)
	slidingTimeWindow.UpdateWithSample(time.Now(), 3)
	assert.Equal(t, 11*time.Second, slidingTimeWindow.Period())
	assert.Equal(t, 2.0, slidingTimeWindow.Average())
	assert.Equal(t, 6.0, slidingTimeWindow.Sum())
	assert.Equal(t, 6.0, slidingTimeWindow.Sum())
	assert.Equal(t, 3, slidingTimeWindow.Count())
}
