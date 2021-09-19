package services

import (
	"github.com/wasmerio/wasmer-go/wasmer"
)

func getString(instance *wasmer.Instance, ptr wasmer.Value, ptrLen wasmer.Value) string {
	mem, merr := instance.Exports.GetMemory("memory")
	if merr != nil {
		return ""
	}
	memData := mem.Data()
	strPtr := ptr.Unwrap().(int32)
	strLen := ptrLen.Unwrap().(int32)
	str := memData[strPtr : strPtr+strLen]
	return string(str)
}

func setString(instance *wasmer.Instance, ptr wasmer.Value, value string) {
	mem, merr := instance.Exports.GetMemory("memory")
	if merr != nil {
		return
	}
	strPtr := ptr.Unwrap().(int32)
	memData := mem.Data()
	end := len(value)
	for nth := 0; nth < end; nth++ {
		memData[int(strPtr)+nth] = value[nth]
	}
}
