package misc

import "sync"

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
