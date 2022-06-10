package core

import (
	"bufio"
	"bytes"
	"os"

	utils "github.com/echosoar/imgpro/utils"
)

// Core 内核
type Core struct {
	FileBinary []byte
	// FilePath   string
	Result map[string]Value

	processorMap  map[string]*Processor
	features      []string
	originFeature []string
	ioReader      *bytes.Reader
}

// Result 结果
type Result map[string]Value

// Value 值
type Value struct {
	Type   ValueType
	Int    int
	String string
	Bytes  []byte
	Rgba   []RGBA
	Rect   []ValuePosition
	Values map[string]Value
	List   []Value
	Frames []Value
}

type ValuePosition struct {
	X int
	Y int
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
	// ValueTypeList []value
	ValueTypeList ValueType = 5
	// ValueTypeRect []ValuePosition
	ValueTypeRect ValueType = 6
	// ValueTypeFrame []ValueP
	ValueTypeFrames ValueType = 7
)

// Processor 处理器
type Processor struct {
	Keys          []string
	PreConditions []string
	Runner        func(*Core) map[string]Value
}

// Bind 绑定处理器
func (core *Core) Bind(processor *Processor) {
	for _, feature := range processor.Keys {
		core.processorMap[feature] = processor
	}
}

// Run 运行
func (core *Core) Run(filePath string) {
	// core.FilePath = filePath

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			panic("File \"" + filePath + "\" not exists")
		}
		panic("stat error")
	}
	// get the size
	size := fileInfo.Size()

	fileHandler, err := os.Open(filePath)
	if err != nil {
		panic("open error")
	}
	defer fileHandler.Close()
	fileBytes := make([]byte, size)
	reader := bufio.NewReader(fileHandler)
	_, readErr := reader.Read(fileBytes)
	if readErr != nil {
		panic("file read error")
	}
	core.RunBinary(fileBytes)
}

// RunBinary 二进制运行
func (core *Core) RunBinary(binary []byte) {
	core.FileBinary = binary
	core.Result = make(map[string]Value)
	core.features = utils.RemoveDuplicateStringValues(core.features)
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

	for _, preFeature := range processor.PreConditions {
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
