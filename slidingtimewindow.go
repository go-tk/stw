package stw

import "time"

// SlidingTimeWindow represents a sliding time window of samples.
type SlidingTimeWindow struct {
	periodPerBucket time.Duration
	buckets         []bucket
	totalSum        float64
	totalCount      int
}

// NewSlidingTimeWindow creates a sliding time window with the given period and number of buckets.
func NewSlidingTimeWindow(period time.Duration, numberOfBuckets int) *SlidingTimeWindow {
	var stw SlidingTimeWindow
	stw.periodPerBucket = (period + time.Duration(numberOfBuckets-1)) / time.Duration(numberOfBuckets)
	stw.buckets = make([]bucket, numberOfBuckets)
	return &stw
}

// UpdateWithSample advances the window and puts the given sample into a bucket.
func (stw *SlidingTimeWindow) UpdateWithSample(now time.Time, x float64) {
	bucketNumber := stw.doUpdate(now)
	bucket := &stw.buckets[bucketNumber%int64(len(stw.buckets))]
	if bucket.number != bucketNumber {
		// ignore
		return
	}
	bucket.sum += x
	stw.totalSum += x
	bucket.count++
	stw.totalCount++
}

// Update advances the window.
func (stw *SlidingTimeWindow) Update(now time.Time) { stw.doUpdate(now) }

func (stw *SlidingTimeWindow) doUpdate(now time.Time) (bucketNumber0 int64) {
	bucketNumber0 = now.UnixNano() / int64(stw.periodPerBucket)
	i0 := int(bucketNumber0 % int64(len(stw.buckets)))
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
	return
}

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
