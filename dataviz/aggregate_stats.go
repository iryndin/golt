package dataviz

import (
	"cmp"
	"math"
	"sort"
)

type AggregateStats struct {
	simulationStartTime       string
	simulationEndTime         string
	simulationDurationSeconds int64
	total                     int
	min                       int
	max                       int
	mean                      int
	stdDev                    int
	p50                       int
	p75                       int
	p90                       int
	p95                       int
	p99                       int
}

func CalculateAggregateStats(nums []int) AggregateStats {
	minVal, maxVal := findMinMax(nums)
	mean := calculateMean(nums)

	return AggregateStats{
		total:  len(nums),
		min:    minVal,
		max:    maxVal,
		mean:   mean,
		stdDev: calculateStdDev(nums, mean),
		p50:    percentile(nums, 50),
		p75:    percentile(nums, 75),
		p90:    percentile(nums, 90),
		p95:    percentile(nums, 95),
		p99:    percentile(nums, 99),
	}
}

func findMinMax[T cmp.Ordered](nums []T) (T, T) {
	minVal, maxVal := nums[0], nums[0]
	for _, num := range nums {
		if num < minVal {
			minVal = num
		}
		if num > maxVal {
			maxVal = num
		}
	}
	return minVal, maxVal
}

func findMinMaxForI64[T any](x []T, mapper func(i T) int64) (int64, int64) {
	minVal := mapper(x[0])
	maxVal := minVal

	for _, item := range x {
		val := mapper(item)
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}
	return minVal, maxVal
}

func findMinMaxForAny[T cmp.Ordered](x []any, mapper func(i any) T) (T, T) {
	minVal := mapper(x[0])
	maxVal := minVal

	for _, item := range x {
		val := mapper(item)
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}
	return minVal, maxVal
}

func calculateMean(nums []int) int {
	sum := 0.0
	for _, num := range nums {
		sum += float64(num)
	}
	mean := sum / float64(len(nums))
	return int(mean)
}

func calculateStdDev(nums []int, mean int) int {
	varianceSum := 0.0
	for _, num := range nums {
		diff := num - mean
		varianceSum += float64(diff * diff)
	}
	variance := varianceSum / float64(len(nums))
	stdDev := math.Sqrt(variance)
	return int(stdDev)
}

func percentile(nums []int, percent float64) int {
	sorted := append([]int(nil), nums...) // copy the slice
	sort.Ints(sorted)

	k := (percent / 100) * float64(len(sorted)-1)
	f := math.Floor(k)
	c := math.Ceil(k)

	if f == c {
		return sorted[int(k)]
	}

	d0 := float64(sorted[int(f)]) * (c - k)
	d1 := float64(sorted[int(c)]) * (k - f)
	return int(d0 + d1)
}
