package photon_spectator

import (
	"bytes"
	"encoding/binary"

	"github.com/google/gopacket"
)

const (
	PhotonCommandHeaderLength = 12
)

var PhotonLayerType = gopacket.RegisterLayerType(5056,
	gopacket.LayerTypeMetadata{
		Name:    "PhotonLayerType",
		Decoder: gopacket.DecodeFunc(decodePhotonPacket)})

type PhotonLayer struct {
	// Header
	PeerID       uint16
	CrcEnabled   uint8
	CommandCount uint8
	Timestamp    uint32
	Challenge    int32

	// Commands
	Commands []PhotonCommand

	// Interface stuff
	contents []byte
	payload  []byte
}

func (p PhotonLayer) LayerType() gopacket.LayerType { return PhotonLayerType }
func (p PhotonLayer) LayerContents() []byte         { return p.contents }
func (p PhotonLayer) LayerPayload() []byte          { return p.payload }

func decodePhotonPacket(data []byte, p gopacket.PacketBuilder) error {
	layer := PhotonLayer{}
	buf := bytes.NewBuffer(data)

	// Read the header
	binary.Read(buf, binary.BigEndian, &layer.PeerID)
	binary.Read(buf, binary.BigEndian, &layer.CrcEnabled)
	binary.Read(buf, binary.BigEndian, &layer.CommandCount)
	binary.Read(buf, binary.BigEndian, &layer.Timestamp)
	binary.Read(buf, binary.BigEndian, &layer.Challenge)

	var commands []PhotonCommand

	// Read each command
	for i := 0; i < int(layer.CommandCount); i++ {
		var command PhotonCommand

		// Command header
		binary.Read(buf, binary.BigEndian, &command.Type)
		binary.Read(buf, binary.BigEndian, &command.ChannelID)
		binary.Read(buf, binary.BigEndian, &command.Flags)
		binary.Read(buf, binary.BigEndian, &command.ReservedByte)
		binary.Read(buf, binary.BigEndian, &command.Length)
		binary.Read(buf, binary.BigEndian, &command.ReliableSequenceNumber)

		// Command data
		dataLength := int(command.Length) - PhotonCommandHeaderLength

		// Ensure we don't try to read more than we have
		if dataLength > buf.Len() {
			panic("Data is malformed")
		}

		command.Data = make([]byte, dataLength)
		buf.Read(command.Data)

		commands = append(commands, command)
	}

	layer.Commands = commands

	// Split and store the read and unread data
	dataUsed := len(data) - buf.Len()
	layer.contents = data[0:dataUsed]
	layer.payload = buf.Bytes()

	p.AddLayer(layer)
	return p.NextDecoder(gopacket.LayerTypePayload)
}
