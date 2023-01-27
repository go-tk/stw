package stw

import (
	"bytes"
	"fmt"
)

func (stw *SlidingTimeWindow) DumpAsString(prefix string) string {
	var buffer bytes.Buffer
	stw.Dump(prefix, &buffer)
	return buffer.String()
}

func (stw *SlidingTimeWindow) Dump(prefix string, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "%sPeriod Per Bucket: %v\n", prefix, stw.periodPerBucket)
	fmt.Fprintf(buffer, "%sPeriod: %v\n", prefix, stw.period)
	for i := range stw.buckets {
		fmt.Fprintf(buffer, "%sBuckets[%d]:\n", prefix, i)
		bucket := &stw.buckets[i]
		stw.dumpBucket(bucket, prefix+"\t", buffer)
	}
	fmt.Fprintf(buffer, "%sTotal Sum: %v\n", prefix, stw.totalSum)
	fmt.Fprintf(buffer, "%sTotal Count: %v\n", prefix, stw.totalCount)
}

func (stw *SlidingTimeWindow) dumpBucket(bucket *bucket, prefix string, buffer *bytes.Buffer) {
	fmt.Fprintf(buffer, "%sNumber: %v\n", prefix, bucket.number)
	fmt.Fprintf(buffer, "%sSum: %v\n", prefix, bucket.sum)
	fmt.Fprintf(buffer, "%sCount: %v\n", prefix, bucket.count)
}
