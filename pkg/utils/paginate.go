package utils

func Paginate[T any](items []T, page int, pageSize int) ([]T, int) {
	totalItems := len(items)

	// 处理特殊情况
	if totalItems == 0 || pageSize <= 0 {
		return []T{}, totalItems
	}

	// 计算总页数
	totalPages := (totalItems + pageSize - 1) / pageSize

	// 确保页码在有效范围内
	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	// 计算起始索引
	startIndex := (page - 1) * pageSize

	// 计算结束索引
	endIndex := startIndex + pageSize
	if endIndex > totalItems {
		endIndex = totalItems
	}

	// 返回分页后的结果和总条目数
	return items[startIndex:endIndex], totalItems
}
