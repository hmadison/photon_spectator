package photon_spectator

import (
	"testing"

	"github.com/google/gopacket"
)

func TestPhotonLayer(t *testing.T) {
	photonHeader := []byte{
		0x00, 0x01, // PeerIdx
		0x01,                   // CrcEnabled
		0x01,                   // CommandCount
		0x00, 0x00, 0x00, 0x01, // Timestamp
		0x00, 0x00, 0x00, 0x01, // Challenge
	}

	photonCommand := []byte{
		AcknowledgeType,        // Type
		0x01,                   // ChannelID
		0x01,                   // Flags
		0x04,                   // ReservedByte
		0x00, 0x00, 0x00, 0x0c, // Length
		0x00, 0x00, 0x00, 0x01, // ReliableSequenceNumber
	}

	data := append(photonHeader, photonCommand...)
	packet := gopacket.NewPacket(data, PhotonLayerType, gopacket.Default)

	photonLayer := packet.Layer(PhotonLayerType)

	if photonLayer == nil {
		t.Errorf("Photon layer should be present")
	}

	packetContent, _ := photonLayer.(PhotonLayer)

	if packetContent.PeerId != uint16(1) {
		t.Errorf("PeerId invalid")
	}

	if packetContent.CrcEnabled != uint8(1) {
		t.Errorf("CrcEnabled invalid")
	}

	if packetContent.CommandCount != uint8(1) {
		t.Errorf("CommandCount invalid")
	}

	if packetContent.Timestamp != uint32(1) {
		t.Errorf("Timestamp invalid")
	}

	if packetContent.Challenge != 1 {
		t.Errorf("Challenge invalid")
	}

	if len(packetContent.Commands) != 1 {
		t.Errorf("Commands length invalid")
	}

	command := packetContent.Commands[0]

	if command.Type != AcknowledgeType {
		t.Errorf("Type invalid")
	}

	if command.ChannelID != uint8(1) {
		t.Errorf("ChannelID invalid")
	}

	if command.Flags != uint8(1) {
		t.Errorf("Flags invalid")
	}

	if command.ReservedByte != uint8(4) {
		t.Errorf("ReservedByte invalid")
	}

	if command.Length != PhotonCommandHeaderLength {
		t.Errorf("Length invalid")
	}

	if command.ReliableSequenceNumber != 1 {
		t.Errorf("ReliableSequenceNumber invalid")
	}

	if packetContent.LayerContents() == nil {
		t.Errorf("LayerContents invalid")
	}

}

func TestMalformedCommand(t *testing.T) {
	photonHeader := []byte{
		0x00, 0x01, // PeerIdx
		0x01,                   // CrcEnabled
		0x01,                   // CommandCount
		0x00, 0x00, 0x00, 0x01, // Timestamp
		0x00, 0x00, 0x00, 0x01, // Challenge
	}

	photonCommand := []byte{
		AcknowledgeType,        // Type
		0x01,                   // ChannelID
		0x01,                   // Flags
		0x04,                   // ReservedByte
		0x00, 0x0c, 0x0c, 0x0c, // Length
		0x00, 0x00, 0x00, 0x01, // ReliableSequenceNumber
	}

	data := append(photonHeader, photonCommand...)
	packet := gopacket.NewPacket(data, PhotonLayerType, gopacket.Default)

	photonLayer := packet.Layer(PhotonLayerType)

	if photonLayer != nil {
		t.Errorf("Photon layer should be absent")
	}
}
