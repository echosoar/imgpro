package imgtype

// Core 内核
type Core struct {
	FilePath string
	Result   map[string]Value

	processorMap  map[string]*Processor
	features      []string
	originFeature []string
}

// Result 结果
type Result map[string]Value

// Value 值
type Value struct {
	Type   ValueType
	Int    int
	String string
	Bytes  []byte
	Rgba   [][]RGBA
	Values map[string]Value
}

// RGBA 值
type RGBA struct {
	R int
	G int
	B int
	A int
}

// ValueType 值类型
type ValueType int

const (
	// ValueTypeInt int
	ValueTypeInt ValueType = 0
	// ValueTypeString string
	ValueTypeString ValueType = 1
	// ValueTypeBytes bytes
	ValueTypeBytes ValueType = 2
	// ValueTypeRGBA rgba
	ValueTypeRGBA ValueType = 3
	// ValueTypeMap map[key] ValueType
	ValueTypeMap ValueType = 4
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
		core.features = append(core.features, condition)
	}
}

// Run 运行
func (core *Core) Run(filePath string) {
	core.FilePath = filePath
	core.Result = make(map[string]Value)
	core.features = removeDuplicateStringValues(core.features)
	for _, feature := range core.features {
		core.runProcessor(feature)
	}
}

// GetResult 获取结果
func (core *Core) GetResult() Result {
	result := make(map[string]Value)
	for _, feature := range core.originFeature {
		result[feature] = core.Result[feature]
	}
	return result
}

func (core *Core) runProcessor(feature string) {
	_, valueExists := core.Result[feature]
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
		core.Result[feat] = value
	}
}

// New 实例化core
func New(features []string) *Core {
	return &Core{
		features:      features,
		originFeature: features,
		processorMap:  make(map[string]*Processor),
	}
}
