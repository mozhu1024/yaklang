package ssaapi

type OperationConfig struct {
	// 限制递归深度，每一次递归核心函数，计数器都会加一
	// 上下文计数器受到这个限制
	MaxDepth int
}

type OperationOption func(*OperationConfig)

func WithMaxDepth(maxDepth int) OperationOption {
	return func(operationConfig *OperationConfig) {
		operationConfig.MaxDepth = maxDepth
	}
}

func NewOperations(opt ...OperationOption) *OperationConfig {
	config := &OperationConfig{
		MaxDepth: -1,
	}

	for _, o := range opt {
		o(config)
	}
	return config
}