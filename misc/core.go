package misc

import (
	"strconv"
	"sync"
)

// ParallelForeach 并行处理 value 分段任务
// value: 总数值
// segments: 分段数（协程数）
// fn: 每个分段的处理函数，参数为起始 start 和结束 end（不含）
func ParallelForeach(value, segments int, fn func(start, end int)) {
	if value <= 0 || segments <= 0 {
		return
	}

	// 限制分段数不超过高度
	if segments > value {
		segments = value
	}

	segmentHeight := value / segments
	remainder := value % segments

	var wg sync.WaitGroup
	for i := 0; i < segments; i++ {
		wg.Add(1)
		go func(segIndex int) {
			defer wg.Done()
			start := segIndex * segmentHeight
			if segIndex > 0 {
				start += min(segIndex, remainder)
			}
			var end int
			if segIndex < remainder {
				end = start + segmentHeight + 1
			} else {
				end = start + segmentHeight
			}
			fn(start, end)
		}(i)
	}
	wg.Wait()
}

func MapGetFloat(m map[string]string, key string) (float64, bool) {
	if m == nil {
		return 0, false
	}

	val, ok := m[key]
	valFloat, err := strconv.ParseFloat(val, 64)
	return valFloat, ok && err == nil
}

func MapGetInt(m map[string]string, key string) (int64, bool) {
	if m == nil {
		return 0, false
	}

	val, ok := m[key]
	valInt, err := strconv.ParseInt(val, 10, 64)
	return valInt, ok && err == nil
}

func MapMultiGetFloat(m map[string]string, keys ...string) ([]float64, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	var vals []float64
	var ok bool = true
	for _, key := range keys {
		val, ok1 := MapGetFloat(m, key)
		vals = append(vals, val)
		if !ok1 {
			ok = false
		}
	}

	return vals, ok
}

func MapMultiGetInt(m map[string]string, keys ...string) ([]int64, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	var vals []int64
	var ok bool = true
	for _, key := range keys {
		val, ok1 := MapGetInt(m, key)
		vals = append(vals, val)
		if !ok1 {
			ok = false
		}
	}

	return vals, ok
}
