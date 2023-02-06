package stw_test

import (
	"fmt"
	"time"

	"github.com/go-tk/stw"
)

func ExampleSlidingTimeWindow() {
	stw := stw.NewSlidingTimeWindow(300*time.Millisecond /* period */, 3 /* numberOfBuckets */)
	fmt.Println("1) Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	stw.AddSample(time.Now(), 1)
	fmt.Println("2) Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	time.Sleep(151 * time.Millisecond)

	stw.AddSample(time.Now(), 10)
	fmt.Println("3) Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	time.Sleep(151 * time.Millisecond)

	stw.AddSample(time.Now(), 100)
	fmt.Println("4) Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	time.Sleep(151 * time.Millisecond)

	stw.Update(time.Now())
	fmt.Println("5) Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	time.Sleep(151 * time.Millisecond)

	stw.Update(time.Now())
	fmt.Println("6) Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	// Output:
	// 1) Sample Count: 0 Sum: 0 Average: NaN
	// 2) Sample Count: 1 Sum: 1 Average: 1
	// 3) Sample Count: 2 Sum: 11 Average: 5.5
	// 4) Sample Count: 2 Sum: 110 Average: 55
	// 5) Sample Count: 1 Sum: 100 Average: 100
	// 6) Sample Count: 0 Sum: 0 Average: NaN
}
