package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{conn: conn, msgCh: msgCh, delCh: delCh}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		value, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var cmd Command
		if value.Type() == resp.Array {
			rawCMD := value.Array()[0]
			switch rawCMD.String() {
			case CommandClient:
				cmd = ClientCommand{
					value: value.Array()[1].String(),
				}
			case CommandGET:
				cmd = GetCommand{
					key: value.Array()[1].Bytes(),
				}
			case CommandSET:
				cmd = SetCommand{
					key: value.Array()[1].Bytes(),
					val: value.Array()[2].Bytes(),
				}
			case CommandHELLO:
				cmd = HelloCommand{
					value: value.Array()[1].String(),
				}
			default:
				fmt.Println("unhandled command", rawCMD)
			}
			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,
			}
		}
	}
	return nil
}
