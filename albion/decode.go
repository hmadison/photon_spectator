package albion

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/hmadison/ao_spectator/photon"
)

const (
	Int32Type       = 105
	Int16Type       = 107
	Int64Type       = 108
	StringType      = 115
	BooleanType     = 111
	SliceUInt16Type = 121
)

func DecodeReliableMessage(msg photon.ReliableMessage) (map[string]interface{}, error) {
	buf := bytes.NewBuffer(msg.Data)
	params := make(map[string]interface{})

	for i := 0; i < int(msg.ParamaterCount); i++ {
		var paramID uint8
		var paramType uint8

		binary.Read(buf, binary.BigEndian, &paramID)
		binary.Read(buf, binary.BigEndian, &paramType)

		paramsKey := strconv.Itoa(int(paramID))

		switch paramType {
		case Int32Type:
			var temp int32

			binary.Read(buf, binary.BigEndian, &temp)

			params[paramsKey] = temp
		case Int16Type:
			var temp int16

			binary.Read(buf, binary.BigEndian, &temp)

			params[paramsKey] = temp
		case Int64Type:
			var temp int64

			binary.Read(buf, binary.BigEndian, &temp)

			params[paramsKey] = temp
		case StringType:
			var unknown uint8
			var length uint8

			binary.Read(buf, binary.BigEndian, &unknown)
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
		case SliceUInt16Type:
			var unknown uint8
			var length uint16

			binary.Read(buf, binary.BigEndian, &length)
			binary.Read(buf, binary.BigEndian, &unknown)

			array := make([]int16, length)

			for j := 0; j < int(length); j++ {
				var temp int16
				binary.Read(buf, binary.BigEndian, &temp)
				array[j] = temp
			}

			params[paramsKey] = array
		default:
			return nil, fmt.Errorf("Invalid type of %d", paramType)
		}
	}

	return params, nil
}
