package bloomfilter

import (
	"github.com/bits-and-blooms/bloom/v3"
)

// CountingBloomFilter 封装计数布隆过滤器
type CountingBloomFilter struct {
	filter *bloom.BloomFilter
}

// NewCountingBloomFilter 创建一个新的计数布隆过滤器
func NewCountingBloomFilter() *CountingBloomFilter {
	n := uint(1_000_000)
	p := 0.01
	return &CountingBloomFilter{
		filter: bloom.NewWithEstimates(n, p),
	}
}

// Add 添加元素到计数布隆过滤器
func (cbf *CountingBloomFilter) Add(element string) {
	cbf.filter.Add([]byte(element))
}

// Contains 检查元素是否存在于计数布隆过滤器中
func (cbf *CountingBloomFilter) Contains(element string) bool {
	return cbf.filter.Test([]byte(element))
}
