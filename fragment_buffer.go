package photon_spectator

import (	
	lru "github.com/hashicorp/golang-lru"
)

// Provides a LRU backed buffer which will assemble ReliableFragments
// into a single PhotonCommand with type ReliableMessage
type FragmentBuffer struct {
	cache *lru.Cache
}

// Offers a message to the buffer. Returns nil when no new commands could be assembled from the
// buffer's contents.
func (buf *FragmentBuffer) Offer(msg ReliableFragment) *PhotonCommand {
	var entry fragmentBufferEntry
	
	if buf.cache.Contains(msg.SequenceNumber) {
		obj, _ := buf.cache.Get(msg.SequenceNumber)
		entry = obj.(fragmentBufferEntry)
		entry.Fragments[int(msg.FragmentNumber)] = msg.Data
		
	} else {
		entry.FragmentsNeeded = int(msg.FragmentCount)
		entry.Fragments = make(map[int][]byte)
		entry.Fragments[int(msg.FragmentNumber)] = msg.Data
	}
	
	if entry.Finished() {
		command := entry.Make()
		buf.cache.Remove(msg.SequenceNumber)
		return &command
	} else {
		buf.cache.Add(msg.SequenceNumber, entry)		
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

// Makes a new instance of a FragmentBuffer
func NewFragmentBuffer() *FragmentBuffer {
	var f FragmentBuffer
	f.cache, _ = lru.New(128)
	return &f
}
