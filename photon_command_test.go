package photon_spectator

import (
	"reflect"
	"testing"
)

func TestPhotonCommand_ReliableMessage_InvalidType(t *testing.T) {
	var cmd PhotonCommand
	_, err := cmd.ReliableMessage()

	if err == nil {
		t.Fail()
	}
}

func TestPhotonCommand_ReliableMessage_OperationRequest(t *testing.T) {
	var cmd PhotonCommand
	cmd.Type = SendReliableType
	cmd.Data = []byte{0x00, OperationRequest, 0x01, 0x00, 0x01}

	msg, err := cmd.ReliableMessage()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if msg.Type != OperationRequest {
		t.Errorf("Type invalid")
	}

	if msg.OperationCode != uint8(1) {
		t.Errorf("OperationCode invalid")
	}

	if msg.ParamaterCount != int16(1) {
		t.Errorf("ParamaterCount invalid")
	}
}

func TestPhotonCommand_ReliableMessage_EventData(t *testing.T) {
	var cmd PhotonCommand
	cmd.Type = SendReliableType
	cmd.Data = []byte{0x00, EventDataType, 0x01, 0x00, 0x01}

	msg, err := cmd.ReliableMessage()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if msg.Type != EventDataType {
		t.Errorf("Type invalid")
	}

	if msg.EventCode != uint8(1) {
		t.Errorf("EventCode invalid")
	}

	if msg.ParamaterCount != int16(1) {
		t.Errorf("ParamaterCount invalid")
	}
}

func TestPhotonCommand_ReliableMessage_OperationResponse(t *testing.T) {
	var cmd PhotonCommand
	cmd.Type = SendReliableType
	cmd.Data = []byte{0x00, otherOperationResponse, 0x1, 0x00, 0x01, 0x01, 0x00, 0x01}

	msg, err := cmd.ReliableMessage()

	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if msg.Type != OperationResponse {
		t.Errorf("Type invalid")
	}

	if msg.OperationCode != uint8(1) {
		t.Errorf("OperationCode invalid")
	}

	if msg.OperationResponseCode != uint16(1) {
		t.Errorf("OperationResponseCode invalid")
	}

	if msg.OperationDebugByte != uint8(1) {
		t.Errorf("OperationDebugByte invalid")
	}
}

func TestPhotonCommand_ReliableFragment_InvalidType(t *testing.T) {
	var cmd PhotonCommand
	_, err := cmd.ReliableFragment()

	if err == nil {
		t.Fail()
	}
}

func TestPhotonCommand_ReliableFragment(t *testing.T) {
	expected := ReliableFragment{
		SequenceNumber: 1,
		FragmentCount:  1,
		FragmentNumber: 1,
		TotalLength:    1,
		FragmentOffset: 1,
		Data:           []uint8{},
	}

	var cmd PhotonCommand
	cmd.Type = SendReliableFragmentType
	cmd.Data = []byte{
		0x0, 0x0, 0x0, 0x1, // SequenceNumber
		0x0, 0x0, 0x0, 0x1, // FragmentCount
		0x0, 0x0, 0x0, 0x1, // FragmentNumber
		0x0, 0x0, 0x0, 0x1, // TotalLength
		0x0, 0x0, 0x0, 0x1, // FragmentOffset
	}

	fragment, _ := cmd.ReliableFragment()

	if !reflect.DeepEqual(expected, fragment) {
		t.Errorf("Expected %#v but got %#v", expected, fragment)
	}

}
