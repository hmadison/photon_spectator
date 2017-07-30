package photon

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	// Command types
	AcknowledgeType          = 1
	ConnectType              = 2
	VerifyConnectType        = 3
	DisconnectType           = 4
	PingType                 = 5
	SendReliableType         = 6
	SendUnreliableType       = 7
	SendReliableFragmentType = 8
	// Message types
	OperationRequest       = 2
	otherOperationResponse = 3
	EventDataType          = 4
	OperationResponse      = 7
)

type PhotonCommand struct {
	// Header
	Type                   uint8
	ChannelID              uint8
	Flags                  uint8
	ReservedByte           uint8
	Length                 int32
	ReliableSequenceNumber int32

	// Body
	Data []byte
}

type ReliableMessage struct {
	// Header
	Signature uint8
	Type      uint8

	// OperationRequest
	OperationCode uint8

	// EventData
	EventCode uint8

	// OperationResponse
	OperationResponseCode uint16
	OperationDebugByte    uint8

	ParamaterCount int16
	Data           []byte
}

func (c PhotonCommand) ReliableMessage() (msg ReliableMessage, err error) {
	if c.Type != SendReliableType {
		return msg, fmt.Errorf("Command can't be converted")
	}

	buf := bytes.NewBuffer(c.Data)

	binary.Read(buf, binary.BigEndian, &msg.Signature)
	binary.Read(buf, binary.BigEndian, &msg.Type)

	if msg.Type == otherOperationResponse {
		msg.Type = OperationResponse
	}

	switch msg.Type {
	case OperationRequest:
		binary.Read(buf, binary.BigEndian, &msg.OperationCode)
	case EventDataType:
		binary.Read(buf, binary.BigEndian, &msg.EventCode)
	case OperationResponse:
		binary.Read(buf, binary.BigEndian, &msg.OperationCode)
		binary.Read(buf, binary.BigEndian, &msg.OperationResponseCode)
		binary.Read(buf, binary.BigEndian, &msg.OperationDebugByte)
	}

	binary.Read(buf, binary.BigEndian, &msg.ParamaterCount)
	msg.Data = buf.Bytes()

	return
}
