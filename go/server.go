package derrit

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
	"encoding/json"
)

type Protocol struct {
	Conn  net.Conn
	proto *textproto.Conn
	ohai  bool
}

func Dial(address string) (*Protocol, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &Protocol{Conn: conn}, nil
}

func (p *Protocol) Start() error {
	p.proto = textproto.NewConn(p.Conn)
	err := p.sendOhai()
	if err != nil {
		return err
	}
	for {
		err = p.handleRequest()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Protocol) sendIhaz() error {
	asset, err := json.Marshal(NewRegister("Hello World"))
	if err != nil {
		return err
	}
	_, err = p.proto.Cmd(fmt.Sprintf("IHAZ %s", asset))
	if err != nil {
		return err
	}
	return nil
}

func (p *Protocol) handleRequest() error {
	line, err := p.proto.ReadLine()
	if err != nil {
		return err
	}
	cmddel := strings.Index(line, " ")
	if cmddel <= 0 {
		return fmt.Errorf("No Command given")
	}
	cmd := line[:cmddel]
	log.Printf("Got %q command", cmd)
	switch cmd {
	case "OHAI":
		return p.handleOhai(line[cmddel+1:])
	case "IHAZ":
		return p.handleIHaz(line[cmddel+1:])
	case "UCANHAZ":
		break
	default:
		log.Println("WTF, what's that?", line)
	}
	return nil
}

func (p *Protocol) handleIHaz(line string) error {
	for {
		handleIHazLine(line)
		line, err := p.proto.ReadLine()
		if err != nil {
			return err
		}
		if line == "" {
			break
		}
	}
	return nil
}

func (p *Protocol) sendOhai() error {
	_, err := p.proto.Cmd("OHAI 4d1e7ca8-a6d6-4b0c-9f2b-2ade9f2269ab")
	return err
}

func (p *Protocol) handleOhai(line string) error {
	p.ohai = true
	log.Printf("Got OHAI from %s", line)
	return p.sendIhaz()
}

func handleIHazLine(line string) {
	log.Println("Got line:", line)
}

func Listen(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Uarggh errrooooor:", err)
			continue
		}
		go handleConn(conn)

	}
}

func handleConn(conn net.Conn) {
	p := &Protocol{Conn: conn}
	p.Start()
}
