package photon_spectator

type FragmentBuffer struct {
	entries map[int]fragmentBufferEntry
}

func (buf *FragmentBuffer) Offer(msg ReliableFragment) *PhotonCommand {
	entry, ok := buf.entries[int(msg.SequenceNumber)]

	if ok {
		entry.Fragments[int(msg.FragmentNumber)] = msg.Data
	} else {
		entry.FragmentsNeeded = int(msg.FragmentCount)
		entry.Fragments = make(map[int][]byte)
		entry.Fragments[int(msg.FragmentNumber)] = msg.Data
	}

	if entry.Finished() {
		command := entry.Make()
		delete(buf.entries, int(msg.SequenceNumber))
		return &command
	} else {
		buf.entries[int(msg.SequenceNumber)] = entry
		return nil
	}
}

type fragmentBufferEntry struct {
	FragmentsNeeded int
	Fragments       map[int][]byte
}

func (buf fragmentBufferEntry) Finished() bool {
	return len(buf.Fragments) == buf.FragmentsNeeded
}

func (buf fragmentBufferEntry) Make() PhotonCommand {
	var data []byte

	for i := 0; i < buf.FragmentsNeeded; i++ {
		data = append(data, buf.Fragments[i]...)
	}

	return PhotonCommand{Type: SendReliableType, Data: data}
}

func NewFragmentBuffer() *FragmentBuffer {
	var f FragmentBuffer
	f.entries = make(map[int]fragmentBufferEntry)
	return &f
}
