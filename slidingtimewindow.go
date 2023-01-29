package stw

import "time"

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
	return &stw
}

// UpdateWithSample removes outdated samples and puts the given sample into a bucket.
func (stw *SlidingTimeWindow) UpdateWithSample(now time.Time, x float64) {
	bucketNumber := stw.doUpdate(now.UnixNano())
	bucket := &stw.buckets[bucketNumber%int64(len(stw.buckets))]
	if bucket.number != bucketNumber {
		// Sample x is outdated, ignore it.
		return
	}
	bucket.sum += x
	stw.totalSum += x
	bucket.count++
	stw.totalCount++
}

// Update removes outdated samples.
func (stw *SlidingTimeWindow) Update(now time.Time) { stw.doUpdate(now.UnixNano()) }

func (stw *SlidingTimeWindow) doUpdate(now int64) (bucketNumber0 int64) {
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
		stw.totalSum -= bucket.sum
		bucket.sum = 0
		stw.totalCount -= bucket.count
		bucket.count = 0
		bucketNumber--
	}
	for i := len(stw.buckets) - 1; i > i0; i-- {
		bucket := &stw.buckets[i]
		if bucket.number >= bucketNumber {
			return
		}
		bucket.number = bucketNumber
		stw.totalSum -= bucket.sum
		bucket.sum = 0
		stw.totalCount -= bucket.count
		bucket.count = 0
		bucketNumber--
	}
	// If all buckets are reset, totalSum should be zero, otherwise it's caused by floating-point errors.
	stw.totalSum = 0
	return
}

// Period returns the period of the time window.
func (stw *SlidingTimeWindow) Period() time.Duration { return stw.period }

// Average returns the average of samples.
func (stw *SlidingTimeWindow) Average() float64 { return stw.totalSum / float64(stw.totalCount) }

// Sum returns the sum of samples.
func (stw *SlidingTimeWindow) Sum() float64 { return stw.totalSum }

// Count returns the count of samples.
func (stw *SlidingTimeWindow) Count() int { return stw.totalCount }

type bucket struct {
	number int64
	sum    float64
	count  int
}
