package stw_test

import (
	"math"
	"testing"
	"time"

	. "github.com/go-tk/stw"
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

		slidingTimeWindow := NewSlidingTimeWindow(c.period, c.numberOfBuckets)
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
	Min: +Inf
	Max: -Inf
Buckets[1]:
	Number: 0
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
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
	Min: +Inf
	Max: -Inf
Buckets[1]:
	Number: 0
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Buckets[2]:
	Number: 0
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Total Sum: 0
Total Count: 0
`[1:]
	}).Run(t)
}

func TestSlidingTimeWindow_AddSample(t *testing.T) {
	type C struct {
		slidingTimeWindow *SlidingTimeWindow
		expectedState     string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		slidingTimeWindow := NewSlidingTimeWindow(9*time.Second, 3)
		slidingTimeWindow.AddSample(time.Unix(12, 1234567), 33)
		c.slidingTimeWindow = slidingTimeWindow

		testcase.DoCallback(0, t, c)

		state := slidingTimeWindow.DumpAsString("")
		assert.Equal(t, c.expectedState, state)
	})

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.AddSample(time.Unix(11, 1234567), 22)
		c.slidingTimeWindow.AddSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow.AddSample(time.Unix(1, 1234567), 99)
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 3
	Sum: 22
	Count: 1
	Min: 22
	Max: 22
Buckets[1]:
	Number: 4
	Sum: 77
	Count: 2
	Min: 33
	Max: 44
Buckets[2]:
	Number: 2
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Total Sum: 99
Total Count: 3
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.AddSample(time.Unix(11, 1234567), 22)
		c.slidingTimeWindow.AddSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow.AddSample(time.Unix(16, 1234567), 11)
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 3
	Sum: 22
	Count: 1
	Min: 22
	Max: 22
Buckets[1]:
	Number: 4
	Sum: 77
	Count: 2
	Min: 33
	Max: 44
Buckets[2]:
	Number: 5
	Sum: 11
	Count: 1
	Min: 11
	Max: 11
Total Sum: 110
Total Count: 4
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.AddSample(time.Unix(11, 1234567), 22)
		c.slidingTimeWindow.AddSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow.AddSample(time.Unix(16, 1234567), 11)
		c.slidingTimeWindow.AddSample(time.Unix(21, 1234567), 33)
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 6
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Buckets[1]:
	Number: 7
	Sum: 33
	Count: 1
	Min: 33
	Max: 33
Buckets[2]:
	Number: 5
	Sum: 11
	Count: 1
	Min: 11
	Max: 11
Total Sum: 44
Total Count: 2
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.AddSample(time.Unix(11, 1234567), 22)
		c.slidingTimeWindow.AddSample(time.Unix(13, 1234567), 44)
		c.slidingTimeWindow.AddSample(time.Unix(23, 1234567), 99)
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 6
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Buckets[1]:
	Number: 7
	Sum: 99
	Count: 1
	Min: 99
	Max: 99
Buckets[2]:
	Number: 5
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Total Sum: 99
Total Count: 1
`[1:]
	}).Run(t)
}

func TestSlidingTimeWindow_Advance(t *testing.T) {
	type C struct {
		slidingTimeWindow *SlidingTimeWindow
		expectedState     string
	}
	tc := testcase.New(func(t *testing.T, c *C) {
		slidingTimeWindow := NewSlidingTimeWindow(9*time.Second, 3)
		slidingTimeWindow.AddSample(time.Unix(11, 1234567), 22)
		slidingTimeWindow.AddSample(time.Unix(12, 1234567), 33)
		slidingTimeWindow.AddSample(time.Unix(13, 1234567), 44)
		slidingTimeWindow.AddSample(time.Unix(14, 1234567), 55)
		c.slidingTimeWindow = slidingTimeWindow

		testcase.DoCallback(0, t, c)

		state := slidingTimeWindow.DumpAsString("")
		assert.Equal(t, c.expectedState, state)
	})

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.Advance(time.Unix(18, 1234567))
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 6
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Buckets[1]:
	Number: 4
	Sum: 132
	Count: 3
	Min: 33
	Max: 55
Buckets[2]:
	Number: 5
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Total Sum: 132
Total Count: 3
`[1:]
	}).Run(t)

	tc.Copy().SetCallback(0, func(t *testing.T, c *C) {
		c.slidingTimeWindow.Advance(time.Unix(23, 1234567))
		c.expectedState = `
Period Per Bucket: 3s
Period: 9s
Buckets[0]:
	Number: 6
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Buckets[1]:
	Number: 7
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Buckets[2]:
	Number: 5
	Sum: 0
	Count: 0
	Min: +Inf
	Max: -Inf
Total Sum: 0
Total Count: 0
`[1:]
	}).Run(t)
}

func TestSlidingTimeWindow_Period_Average_Sum_Count(t *testing.T) {
	slidingTimeWindow := NewSlidingTimeWindow(11*time.Second, 22)
	assert.True(t, math.IsNaN(slidingTimeWindow.Average()))
	slidingTimeWindow.AddSample(time.Now(), 1)
	slidingTimeWindow.AddSample(time.Now(), 2)
	slidingTimeWindow.AddSample(time.Now(), 3)
	assert.Equal(t, 11*time.Second, slidingTimeWindow.Period())
	assert.Equal(t, 2.0, slidingTimeWindow.Average())
	assert.Equal(t, 6.0, slidingTimeWindow.Sum())
	assert.Equal(t, 3, slidingTimeWindow.Count())
}

func TestSlidingTimeWindow_Min_Max(t *testing.T) {
	slidingTimeWindow := NewSlidingTimeWindow(9*time.Second, 3)
	assert.True(t, math.IsInf(slidingTimeWindow.Min(), +1))
	assert.True(t, math.IsInf(slidingTimeWindow.Max(), -1))
	slidingTimeWindow.AddSample(time.Unix(10, 1234567), -11)
	slidingTimeWindow.AddSample(time.Unix(11, 1234567), 222)
	slidingTimeWindow.AddSample(time.Unix(12, 1234567), 33)
	slidingTimeWindow.AddSample(time.Unix(13, 1234567), -444)
	slidingTimeWindow.AddSample(time.Unix(14, 1234567), 55)
	slidingTimeWindow.AddSample(time.Unix(15, 1234567), 666)
	slidingTimeWindow.AddSample(time.Unix(16, 1234567), -77)
	assert.Equal(t, -444.0, slidingTimeWindow.Min())
	assert.Equal(t, 666.0, slidingTimeWindow.Max())
}

func TestSlidingTimeWindow_Reduce(t *testing.T) {
	slidingTimeWindow := NewSlidingTimeWindow(9*time.Second, 3)
	slidingTimeWindow.AddSample(time.Unix(11, 1234567), 22)
	slidingTimeWindow.AddSample(time.Unix(12, 1234567), 33)
	slidingTimeWindow.AddSample(time.Unix(13, 1234567), 44)

	{
		x := slidingTimeWindow.Reduce(0, func(x float64, bucket Bucket) (y float64) {
			return x + float64(bucket.Count)
		})
		assert.Equal(t, x, 3.0)
	}

	{
		x := slidingTimeWindow.Reduce(math.Inf(-1), func(x float64, bucket Bucket) (y float64) {
			if bucket.Count == 0 {
				return x
			}
			return math.Max(x, bucket.Sum/float64(bucket.Count))
		})
		assert.Equal(t, x, 38.5)
	}

	{
		x := slidingTimeWindow.Reduce(math.Inf(+1), func(x float64, bucket Bucket) (y float64) {
			if bucket.Count == 0 {
				return x
			}
			return math.Min(x, bucket.Sum/float64(bucket.Count))
		})
		assert.Equal(t, x, 22.0)
	}
}
