package packet

import (
	"bytes"
	"fmt"
)

type Packet interface {
	Decode([]byte) error     // []byte -> struct
	Encode() ([]byte, error) // struct -> []byte
}

const (
	CommandConn      = 0x01
	CommandConnAck   = 0x02
	CommandSubmit    = 0x03
	CommandSubmitAck = 0x04
)

func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]
	pktBody := packet[1:]

	switch commandID {
	case CommandConn:
		return nil, nil
	case CommandConnAck:
		return nil, nil
	case CommandSubmit:
		s := Submit{}
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	case CommandSubmitAck:
		s := SubmitAck{}
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, fmt.Errorf("unknown commandID [%d]", commandID)
	}
}

func Encode(p Packet) ([]byte, error) {
	var commandID uint8
	var pktBody []byte
	var err error

	switch t := p.(type) {
	case *Submit:
		commandID = CommandSubmit
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{[]byte{commandID}, pktBody}, nil), nil
}

type Submit struct {
	ID      string
	Payload []byte
}

func (s *Submit) Decode(pktBody []byte) error {
	s.ID = string(pktBody[:8])
	s.Payload = pktBody[8:]
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), s.Payload}, nil), nil
}

type SubmitAck struct {
	ID     string
	Result uint8
}

func (sa *SubmitAck) Decode(pktBody []byte) error {
	sa.ID = string(pktBody[:8])
	sa.Result = pktBody[8]
	return nil
}

func (sa *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(sa.ID[:8]), []byte{sa.Result}}, nil), nil
}
