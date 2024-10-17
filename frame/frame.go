package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

type FramePayload []byte

type StreamFrameCodec interface {
	Encode(io.Writer, FramePayload) error   // 将FramePayload添加totalLength变为Frame并发包
	Decode(io.Reader) (FramePayload, error) // 提取数据填入framepayload并返回
}

var ErrShortWrite = errors.New("short write")
var ErrShortRead = errors.New("short read")

type myFrameCodec struct{}

func NewMyFrameCodec() StreamFrameCodec {
	return &myFrameCodec{}
}

func (p *myFrameCodec) Encode(w io.Writer, framePayLoad FramePayload) error {
	var f = framePayLoad
	var totalLength int32 = int32(len(framePayLoad)) + 4

	err := binary.Write(w, binary.BigEndian, &totalLength)
	if err != nil {
		return err
	}

	n, err := w.Write([]byte(f))
	if err != nil {
		return err
	}

	if n != len(framePayLoad) {
		return ErrShortWrite
	}

	return nil
}

func (p *myFrameCodec) Decode(r io.Reader) (FramePayload, error) {
	var totalLength int32
	err := binary.Read(r, binary.BigEndian, &totalLength)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLength-4)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	if n != int(totalLength-4) {
		return nil, ErrShortRead
	}

	return FramePayload(buf), nil
}
