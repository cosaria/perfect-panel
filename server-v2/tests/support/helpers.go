package support

// Ptr 返回任意值的指针，便于测试构造 overlay 输入。
func Ptr[T any](value T) *T {
	return &value
}
