package util

// GetTotalPageIndexes 获取总的分页索引列表，总数小于pageSize的不分页
func GetTotalPageIndexes(startPageIndex, pageSize, total int64) []int64 {
	if pageSize < total {
		maxPageIndex := total / pageSize
		remainder := total % pageSize

		if remainder > 0 {
			maxPageIndex += 1
		}

		totalPageIndexes := make([]int64, 0, maxPageIndex)
		for i := startPageIndex; i <= maxPageIndex; i++ {
			totalPageIndexes = append(totalPageIndexes, i)
		}

		return totalPageIndexes
	}

	return []int64{}
}
