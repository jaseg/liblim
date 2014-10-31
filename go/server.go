package liblim

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

var InstanceId string

func init() {
	InstanceId = os.Getenv("LIBLIM_INSTANCE_ID")
	tmpUuid := uuid.Parse(InstanceId)
	if tmpUuid == nil {
		InstanceId = uuid.NewRandom().String()
	}
}

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

func (p *Protocol) Close() error {
	return p.proto.Close()
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
	payload := line[cmddel+1:]
	log.Printf("Got %q command", cmd)
	switch cmd {
	case "ØHAI":
		return p.handleOhai(payload)
	case "IHAZ":
		return p.handleIHaz(payload)
	case "UCANHAZ":
		return p.handleUcanhaz(payload)
	default:
		log.Println("WTF dat shit?", line)
	}
	return nil
}

func (p *Protocol) handleIHaz(line string) error {
	for {
		line, err := p.proto.ReadLine()
		if err != nil {
			return err
		}
		if line == "" {
			break
		}
		log.Println("Got IHAZ:", line)
	}
	return p.sendUcanhaz()
}

func (p *Protocol) sendOhai() error {
	_, err := p.proto.Cmd("ØHAI " + InstanceId)
	return err
}

func (p *Protocol) handleOhai(line string) error {
	p.ohai = true
	log.Printf("Got ØHAI from %s", line)
	return p.sendIhaz()
}

func (p *Protocol) handleUcanhaz(line string) error {
	log.Println("Got UCANHAZ:", line)
	return nil
}

func (p *Protocol) sendUcanhaz() error {
	_, err := p.proto.Cmd("UCANHAZ bsfgsdafgdsf")
	return err
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
	log.Println(p.Start())
}
