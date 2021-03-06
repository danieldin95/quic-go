package ackhandler

import (
	"fmt"

	"github.com/lucas-clemente/quic-go/internal/protocol"
)

type sentPacketHistory struct {
	packetList *PacketList
	packetMap  map[protocol.PacketNumber]*PacketElement
}

func newSentPacketHistory() *sentPacketHistory {
	return &sentPacketHistory{
		packetList: NewPacketList(),
		packetMap:  make(map[protocol.PacketNumber]*PacketElement),
	}
}

func (h *sentPacketHistory) SentPacket(p *Packet) {
	el := h.packetList.PushBack(*p)
	h.packetMap[p.PacketNumber] = el
}

// Iterate iterates through all packets.
func (h *sentPacketHistory) Iterate(cb func(*Packet) (cont bool, err error)) error {
	cont := true
	var next *PacketElement
	for el := h.packetList.Front(); cont && el != nil; el = next {
		var err error
		next = el.Next()
		cont, err = cb(&el.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// FirstOutStanding returns the first outstanding packet.
func (h *sentPacketHistory) FirstOutstanding() *Packet {
	if !h.HasOutstandingPackets() {
		return nil
	}
	return &h.packetList.Front().Value
}

func (h *sentPacketHistory) Len() int {
	return len(h.packetMap)
}

func (h *sentPacketHistory) Remove(p protocol.PacketNumber) error {
	el, ok := h.packetMap[p]
	if !ok {
		return fmt.Errorf("packet %d not found in sent packet history", p)
	}
	h.packetList.Remove(el)
	delete(h.packetMap, p)
	return nil
}

func (h *sentPacketHistory) HasOutstandingPackets() bool {
	return h.packetList.Len() > 0
}
