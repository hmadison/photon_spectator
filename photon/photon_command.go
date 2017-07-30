package photon

const (
	PhotonAcknowledgeType          = 1
	PhotonConnectType              = 2
	PhotonVerifyConnectType        = 3
	PhotonDisconnectType           = 4
	PhotonPingType                 = 5
	PhotonSendReliableType         = 6
	PhotonSendUnreliableType       = 7
	PhotonSendReliableFragmentType = 8
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
