package imgtype

// Core 内核
type Core struct {
	Features     []string
	FilePath     string
	processorMap map[string]*Processor
	result       map[string]Value
}

// Result 结果
type Result map[string]Value

// Value 值
type Value struct {
	Type ValueType
	Int  int64
}

// ValueType 值类型
type ValueType int

const (
	// ValueTypeInt int
	ValueTypeInt ValueType = 0
)

// Processor 处理器
type Processor struct {
	Keys         []string
	Precondition []string
	Runner       func(*Core) map[string]Value
}

// Bind 绑定处理器
func (core *Core) Bind(processor *Processor) {
	for _, feature := range processor.Keys {
		core.processorMap[feature] = processor
	}
	for _, condition := range processor.Precondition {
		core.Features = append(core.Features, condition)
	}
}

// Run 运行
func (core *Core) Run(filePath string) Result {
	core.FilePath = filePath
	core.Features = removeDuplicateStringValues(core.Features)
	core.result = make(map[string]Value)
	for _, feature := range core.Features {
		core.runProcessor(feature)
	}

	return core.result
}

func (core *Core) runProcessor(feature string) {
	_, valueExists := core.result[feature]
	if valueExists {
		return
	}
	processor, exists := core.processorMap[feature]
	if !exists {
		panic("Feature \"" + feature + "\" not found processor")
	}

	for _, preFeature := range processor.Precondition {
		core.runProcessor(preFeature)
	}

	result := processor.Runner(core)
	for feat, value := range result {
		core.result[feat] = value
	}
}

// New 实例化core
func New() *Core {
	return &Core{
		processorMap: make(map[string]*Processor),
	}
}
