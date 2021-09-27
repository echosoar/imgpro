package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/echosoar/imgpro"
)

func imgExec(this js.Value, args []js.Value) interface{} {

	featuresArg := args[0]
	featureLength := featuresArg.Length()
	features := make([]string, featureLength)
	for i := 0; i < featureLength; i++ {
		features[i] = featuresArg.Index(i).String()
	}
	binaryArg := args[1]
	binaryLength := binaryArg.Length()
	binary := make([]byte, binaryLength)
	for i := 0; i < binaryLength; i++ {
		binary[i] = byte(binaryArg.Index(i).Int())
	}

	result := imgpro.RunBinary(binary, features)
	jsonData, _ := json.Marshal(result)
	return js.ValueOf(string(jsonData))
}

func main() {
	done := make(chan int, 0)
	js.Global().Set("imgExec", js.FuncOf(imgExec))
	<-done
}
