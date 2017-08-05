package photon_spectator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

const (
	NilType       = 42
	Int8Type      = 98
	Float32Type   = 102
	Int32Type     = 105
	Int16Type     = 107
	Int64Type     = 108
	StringType    = 115
	BooleanType   = 111
	SliceInt8Type = 120
	SliceType     = 121
)

func DecodeReliableMessage(msg ReliableMessage) (map[string]interface{}, error) {
	buf := bytes.NewBuffer(msg.Data)
	params := make(map[string]interface{})

	for i := 0; i < int(msg.ParamaterCount); i++ {
		var paramID uint8
		var paramType uint8

		binary.Read(buf, binary.BigEndian, &paramID)
		binary.Read(buf, binary.BigEndian, &paramType)

		paramsKey := strconv.Itoa(int(paramID))

		switch paramType {
		case NilType, 0:
			// Do nothing
		case Int8Type:
			var temp int8

			binary.Read(buf, binary.BigEndian, &temp)

			params[paramsKey] = temp
		case Float32Type:
			var temp float32

			binary.Read(buf, binary.BigEndian, &temp)

			params[paramsKey] = temp
		case Int32Type:
			var temp int32

			binary.Read(buf, binary.BigEndian, &temp)

			params[paramsKey] = temp
		case Int16Type, 7:
			var temp int16

			binary.Read(buf, binary.BigEndian, &temp)

			params[paramsKey] = temp
		case Int64Type:
			var temp int64

			binary.Read(buf, binary.BigEndian, &temp)

			params[paramsKey] = temp
		case StringType:
			var length uint16

			binary.Read(buf, binary.BigEndian, &length)

			strBytes := make([]byte, length)
			buf.Read(strBytes)

			params[paramsKey] = string(strBytes[:])
		case BooleanType:
			var value uint8

			binary.Read(buf, binary.BigEndian, &value)

			if value == 0 {
				params[paramsKey] = false
			} else if value == 1 {
				params[paramsKey] = true
			} else {
				return nil, fmt.Errorf("Invalid value for boolean of %d", value)
			}
		case SliceInt8Type:
			var length uint32

			binary.Read(buf, binary.BigEndian, &length)

			array := make([]int8, length)

			for j := 0; j < int(length); j++ {
				var temp int8
				binary.Read(buf, binary.BigEndian, &temp)
				array[j] = temp
			}

			params[paramsKey] = array
		case SliceType:
			array, error := decodeSlice(buf)
			if error != nil {
				return nil, fmt.Errorf("Slice Error: %s; Current Params: %+v", error.Error(), params)
			}
			params[paramsKey] = array
		default:
			return nil, fmt.Errorf("Invalid type of %d; Current Params: %+v", paramType, params)
		}
	}

	return params, nil
}

func decodeSlice(buf *bytes.Buffer) (interface{}, error) {
	var length uint16
	var sliceType uint8

	binary.Read(buf, binary.BigEndian, &length)
	binary.Read(buf, binary.BigEndian, &sliceType)

	switch sliceType {
	case Float32Type:
		array := make([]float32, length)

		for j := 0; j < int(length); j++ {
			var temp float32
			binary.Read(buf, binary.BigEndian, &temp)
			array[j] = temp
		}

		return array, nil
	case Int32Type:
		array := make([]int32, length)

		for j := 0; j < int(length); j++ {
			var temp int32
			binary.Read(buf, binary.BigEndian, &temp)
			array[j] = temp
		}

		return array, nil
	case Int16Type:
		array := make([]int16, length)

		for j := 0; j < int(length); j++ {
			var temp int16
			binary.Read(buf, binary.BigEndian, &temp)
			array[j] = temp
		}

		return array, nil
	case Int64Type:
		array := make([]int64, length)

		for j := 0; j < int(length); j++ {
			var temp int64
			binary.Read(buf, binary.BigEndian, &temp)
			array[j] = temp
		}

		return array, nil
	case StringType:
		array := make([]string, length)

		for j := 0; j < int(length); j++ {
			var strLen uint16

			binary.Read(buf, binary.BigEndian, &strLen)

			strBytes := make([]byte, strLen)
			buf.Read(strBytes)

			array[j] = string(strBytes[:])
		}

		return array, nil
	case BooleanType:
		array := make([]bool, length)

		for j := 0; j < int(length); j++ {
			var value uint8

			binary.Read(buf, binary.BigEndian, &value)

			if value == 0 {
				array[j] = false
			} else if value == 1 {
				array[j] = true
			} else {
				return nil, fmt.Errorf("Invalid value for boolean of %d", value)
			}
		}

		return array, nil
	case SliceInt8Type:
		array := make([][]int8, length)

		for j := 0; j < int(length); j++ {
			var sliceLength uint32

			binary.Read(buf, binary.BigEndian, &sliceLength)

			sliceArray := make([]int8, length)

			for k := 0; k < int(length); k++ {
				var temp int8
				binary.Read(buf, binary.BigEndian, &temp)
				sliceArray[j] = temp
			}

			array[j] = sliceArray
		}

		return array, nil
	case SliceType:
		array := make([]interface{}, length)

		for j := 0; j < int(length); j++ {
			subArray, error := decodeSlice(buf)

			if error != nil {
				return nil, error
			}

			array[j] = subArray
		}

		return array, nil
	default:
		return nil, fmt.Errorf("Invalid slice type of %d", sliceType)
	}
}
