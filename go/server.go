package derrit

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
)

type Protocol struct {
	Conn  net.Conn
	proto *textproto.Conn
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
	err := p.sendIHaz()
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

func (p *Protocol) sendIHaz() error {
	_, err := p.proto.Cmd("IHAZ blablabla")
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
	case "IHAZ":
		log.Println("Got ICANHAZ")
		p.handleIHaz(line[cmddel+1:])
		break
	case "UCANHAZ":
		log.Println("Got UCANHAZ")
		break
	default:
		log.Println("WTF, what's that?", line)
	}
	return nil
}

func (p *Protocol) handleIHaz(line string) error {
	for {
		go handleIHazLine(line)
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
