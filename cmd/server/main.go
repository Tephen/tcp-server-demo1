package main

import (
	"fmt"
	"net"

	"github.com/Tephen/tcp-server-demo1/frame"
	"github.com/Tephen/tcp-server-demo1/packet"
)

func handlePacket(framePayload []byte) (ackFramePayload []byte, err error) {
	var p packet.Packet
	p, err = packet.Decode(framePayload)
	if err != nil {
		fmt.Println("handleConn: packet decode error:", err)
		return
	}
	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		fmt.Printf("recv submit: id = %s, payload=%s\n", submit.ID, string(submit.Payload))
		// 返回ACK
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0,
		}
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			fmt.Println("handleConn: packet encode error:", err)
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	frameCodec := frame.NewMyFrameCodec()

	framePayload, err := frameCodec.Decode(c)
	if err != nil {
		fmt.Println("handleConn: frame decode error:", err)
		return
	}

	// 处理收到的包并生成ackFrame包
	ackFramePayload, err := handlePacket(framePayload)
	if err != nil {
		fmt.Println("handleConn: handle packet error:", err)
		return
	}

	// write ack frame to the connection
	err = frameCodec.Encode(c, ackFramePayload)
	if err != nil {
		fmt.Println("handleConn: frame encode error:", err)
		return
	}
}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go handleConn(c)
	}
}
