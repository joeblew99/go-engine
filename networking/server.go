package networking

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net"
)

const serverPacketBufferSize = 1000
const serverSessionBufferSize = 20
const sessionTimeout = 10 * 60 // 10 minutes

type Server struct {
	conn           *net.UDPConn
	sessions       map[string]*Session
	newSessions    chan *Session
	packets        chan Packet
	onClientJoined func(clientId string)
	bytesSent      int64
	bytesReceived  int64
}

func NewServer() *Server {
	var server *Server
	server = &Server{
		sessions:    make(map[string]*Session),
		newSessions: make(chan *Session, serverSessionBufferSize),
		packets:     make(chan Packet, clientPacketBufferSize),
		onClientJoined: func(clientId string) {
			server.WriteMessage("", clientId, []byte{})
		},
	}
	return server
}

func (s *Server) Listen(port int) {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println("Error Resolving udp address: ", err)
		return
	}

	s.conn, err = net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error listening on udp address: ", err)
		return
	}

	data := make([]byte, 65500)
	go func() {
		for s.conn != nil {
			n, addr, err := s.conn.ReadFromUDP(data)
			if err != nil {
				fmt.Println("Error reading udp packet: ", err)
				continue
			}

			s.bytesReceived += int64(n)
			dataBuf := bytes.NewBuffer(data[0:n])
			gzipReader, err := gzip.NewReader(dataBuf)
			if err != nil {
				fmt.Println("Error creating gzip Reader for udp packet: ", err)
				continue
			}

			data, err := ioutil.ReadAll(gzipReader)
			if err != nil {
				fmt.Println("Error unzipping udp packet: ", err)
				continue
			}

			packet, err := Decode(data)
			if err != nil {
				fmt.Println("Error Decoding udp packet: ", err)
				continue
			}

			if len(packet.Token) == 0 {
				s.newSessions <- NewSession(addr)
			}

			s.packets <- packet
		}
	}()
}

func (s *Server) WriteMessage(command, token string, data []byte) {
	session, ok := s.sessions[token]
	packet := Packet{
		Token:   token,
		Command: command,
		Data:    data,
	}

	if ok {
		data := Encode(packet)

		var gzipBuf bytes.Buffer
		gzipWriter := gzip.NewWriter(&gzipBuf)
		_, err := gzipWriter.Write(data)
		if err != nil {
			fmt.Println("Error Gzip compressing udp message: ", err)
			return
		}

		if err := gzipWriter.Flush(); err != nil {
			fmt.Println("Error Flushing Gzip writer for udp message: ", err)
			return
		}

		if err := gzipWriter.Close(); err != nil {
			fmt.Println("Error Closing Gzip writer for udp message: ", err)
			return
		}

		gzipData := gzipBuf.Bytes()
		s.bytesSent += int64(len(gzipData))
		s.conn.WriteToUDP(gzipData, session.addr)
		if err != nil {
			fmt.Println("Error Writing udp message: ", err)
		}
	}
}

func (s *Server) BroadcastMessage(command string, data []byte) {
	for token, _ := range s.sessions {
		s.WriteMessage(command, token, data)
	}
}

func (s *Server) ClientJoinedEvent(callback func(clientId string)) {
	s.onClientJoined = callback
}

func (s *Server) Update(dt float64) {
	// check for new sessions
	select {
	case newSession := <-s.newSessions:
		s.sessions[newSession.token] = newSession
		s.onClientJoined(newSession.token)
	default:
	}

	// check for session timeouts
	for token, session := range s.sessions {
		if session.idleTime > sessionTimeout {
			delete(s.sessions, token)
		}
		session.idleTime = session.idleTime + dt
	}
}

func (s *Server) GetNextMessage() (Packet, bool) {
	select {
	case packet := <-s.packets:
		if session, ok := s.sessions[packet.Token]; ok {
			session.idleTime = 0
		}
		return packet, true
	default:
	}
	return Packet{}, false
}

func (s *Server) Close() {
	s.conn.Close()
	s.conn = nil
}
