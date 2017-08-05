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

type ReliableMessageParamaters map[string]interface{}

func DecodeReliableMessage(msg ReliableMessage) (ReliableMessageParamaters, error) {
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
			params[paramsKey] = decodeInt8Type(buf)
		case Float32Type:
			params[paramsKey] = decodeFloat32Type(buf)
		case Int32Type:
			params[paramsKey] = decodeInt32Type(buf)
		case Int16Type, 7:
			params[paramsKey] = decodeInt16Type(buf)
		case Int64Type:
			params[paramsKey] = decodeInt64Type(buf)
		case StringType:
			params[paramsKey] = decodeStringType(buf)
		case BooleanType:
			result, err := decodeBooleanType(buf)

			if err != nil {
				return nil, err
			}

			params[paramsKey] = result
		case SliceInt8Type:
			params[paramsKey] = decodeSliceInt8Type(buf)
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
			array[j] = decodeFloat32Type(buf)
		}

		return array, nil
	case Int32Type:
		array := make([]int32, length)

		for j := 0; j < int(length); j++ {
			array[j] = decodeInt32Type(buf)
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
			array[j] = decodeInt64Type(buf)
		}

		return array, nil
	case StringType:
		array := make([]string, length)

		for j := 0; j < int(length); j++ {
			array[j] = decodeStringType(buf)
		}

		return array, nil
	case BooleanType:
		array := make([]bool, length)

		for j := 0; j < int(length); j++ {
			result, err := decodeBooleanType(buf)

			if err != nil {
				return array, err
			}

			array[j] = result
		}

		return array, nil
	case SliceInt8Type:
		array := make([][]int8, length)

		for j := 0; j < int(length); j++ {
			array[j] = decodeSliceInt8Type(buf)
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

func decodeInt8Type(buf *bytes.Buffer) (temp int8) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeFloat32Type(buf *bytes.Buffer) (temp float32) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeInt16Type(buf *bytes.Buffer) (temp int16) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeInt32Type(buf *bytes.Buffer) (temp int32) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeInt64Type(buf *bytes.Buffer) (temp int64) {
	binary.Read(buf, binary.BigEndian, &temp)
	return
}

func decodeStringType(buf *bytes.Buffer) string {
	var length uint16

	binary.Read(buf, binary.BigEndian, &length)

	strBytes := make([]byte, length)
	buf.Read(strBytes)

	return string(strBytes[:])
}

func decodeBooleanType(buf *bytes.Buffer) (bool, error) {
	var value uint8

	binary.Read(buf, binary.BigEndian, &value)

	if value == 0 {
		return false, nil
	} else if value == 1 {
		return true, nil
	} else {
		return false, fmt.Errorf("Invalid value for boolean of %d", value)
	}

}

func decodeSliceInt8Type(buf *bytes.Buffer) []int8 {
	var length uint32

	binary.Read(buf, binary.BigEndian, &length)

	array := make([]int8, length)

	for j := 0; j < int(length); j++ {
		var temp int8
		binary.Read(buf, binary.BigEndian, &temp)
		array[j] = temp
	}

	return array
}
