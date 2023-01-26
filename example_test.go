package stw_test

import (
	"fmt"
	"time"

	"github.com/go-tk/stw"
)

func ExampleSlidingTimeWindow() {
	stw := stw.NewSlidingTimeWindow(300*time.Millisecond /* period */, 3 /* numberOfBuckets */)
	fmt.Println("Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	stw.UpdateWithSample(time.Now(), 1)
	fmt.Println("Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	time.Sleep(151 * time.Millisecond)

	stw.UpdateWithSample(time.Now(), 10)
	fmt.Println("Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	time.Sleep(151 * time.Millisecond)

	stw.UpdateWithSample(time.Now(), 100)
	fmt.Println("Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	time.Sleep(151 * time.Millisecond)

	stw.Update(time.Now())
	fmt.Println("Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	time.Sleep(151 * time.Millisecond)

	stw.Update(time.Now())
	fmt.Println("Sample Count:", stw.Count(), "Sum:", stw.Sum(), "Average:", stw.Average())

	// Output:
	// Sample Count: 0 Sum: 0 Average: NaN
	// Sample Count: 1 Sum: 1 Average: 1
	// Sample Count: 2 Sum: 11 Average: 5.5
	// Sample Count: 2 Sum: 110 Average: 55
	// Sample Count: 1 Sum: 100 Average: 100
	// Sample Count: 0 Sum: 0 Average: NaN
}
