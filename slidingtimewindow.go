package stw

import (
	"math"
	"time"
)

var (
	pinf = math.Inf(+1)
	ninf = math.Inf(-1)
)

// SlidingTimeWindow represents a sliding time window of samples.
type SlidingTimeWindow struct {
	periodPerBucket time.Duration
	period          time.Duration
	buckets         []bucket
	totalSum        float64
	totalCount      int
}

// NewSlidingTimeWindow creates a sliding time window with the given period and the given number of buckets.
func NewSlidingTimeWindow(period time.Duration, numberOfBuckets int) *SlidingTimeWindow {
	var stw SlidingTimeWindow
	stw.periodPerBucket = (period + time.Duration(numberOfBuckets-1)) / time.Duration(numberOfBuckets)
	stw.period = stw.periodPerBucket * time.Duration(numberOfBuckets)
	stw.buckets = make([]bucket, numberOfBuckets)
	for i := range stw.buckets {
		bucket := &stw.buckets[i]
		bucket.Min = pinf
		bucket.Max = ninf
	}
	return &stw
}

// AddSample puts the given sample into a bucket and removes outdated samples.
func (stw *SlidingTimeWindow) AddSample(now time.Time, x float64) {
	bucketNumber := stw.doAdvance(now.UnixNano())
	bucket := &stw.buckets[bucketNumber%int64(len(stw.buckets))]
	if bucket.number != bucketNumber {
		// Sample x is outdated, ignore it.
		return
	}
	bucket.Sum += x
	stw.totalSum += x
	bucket.Count++
	stw.totalCount++
	bucket.Min = math.Min(bucket.Min, x)
	bucket.Max = math.Max(bucket.Max, x)
}

// Advance removes outdated samples.
func (stw *SlidingTimeWindow) Advance(now time.Time) { stw.doAdvance(now.UnixNano()) }

func (stw *SlidingTimeWindow) doAdvance(now int64) (bucketNumber0 int64) {
	bucketNumber0 = now / int64(stw.periodPerBucket)
	i0 := int(bucketNumber0 % int64(len(stw.buckets)))
	// Reset buckets with outdated samples.
	bucketNumber := bucketNumber0
	for i := i0; i >= 0; i-- {
		bucket := &stw.buckets[i]
		if bucket.number >= bucketNumber {
			return
		}
		bucket.number = bucketNumber
		stw.totalSum -= bucket.Sum
		bucket.Sum = 0
		stw.totalCount -= bucket.Count
		bucket.Count = 0
		bucket.Min = pinf
		bucket.Max = ninf
		bucketNumber--
	}
	for i := len(stw.buckets) - 1; i > i0; i-- {
		bucket := &stw.buckets[i]
		if bucket.number >= bucketNumber {
			return
		}
		bucket.number = bucketNumber
		stw.totalSum -= bucket.Sum
		bucket.Sum = 0
		stw.totalCount -= bucket.Count
		bucket.Count = 0
		bucket.Min = pinf
		bucket.Max = ninf
		bucketNumber--
	}
	// If all buckets are reset, totalSum should be zero, otherwise it's caused by floating-point errors.
	stw.totalSum = 0
	return
}

// Period returns the period of the time window.
func (stw *SlidingTimeWindow) Period() time.Duration { return stw.period }

// Average returns the average of samples. If there is no sample, NaN is returned.
func (stw *SlidingTimeWindow) Average() float64 { return stw.totalSum / float64(stw.totalCount) }

// Sum returns the sum of samples.
func (stw *SlidingTimeWindow) Sum() float64 { return stw.totalSum }

// Count returns the count of samples.
func (stw *SlidingTimeWindow) Count() int { return stw.totalCount }

// Min returns the minimum of samples. If there is no sample, +Inf is returned.
func (stw *SlidingTimeWindow) Min() float64 {
	min := pinf
	for i := range stw.buckets {
		min = math.Min(min, stw.buckets[i].Min)
	}
	return min
}

// Max returns the maximum of samples. If there is no sample, -Inf is returned.
func (stw *SlidingTimeWindow) Max() float64 {
	max := ninf
	for i := range stw.buckets {
		max = math.Max(max, stw.buckets[i].Max)
	}
	return max
}

// func (stw *SlidingTimeWindow) Reduce(x0 float64, f func(x float64, bucket Bucket) (y float64)) (x float64) {
//	x = x0
//	for i := range stw.buckets {
//		x = f(x, stw.buckets[i].Bucket)
//	}
//	return
// }

type bucket struct {
	number int64

	Bucket
}

// Bucket represents a subset of samples.
type Bucket struct {
	Sum   float64
	Count int
	Min   float64
	Max   float64
}
