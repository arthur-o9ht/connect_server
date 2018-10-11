package tcp

import (
	"github.com/connect_server/lib"
	"log"
	"net"
	"sync"
	"time"
)

type HandleFunc func() error

type Server struct {
	Ip           string
	Port         string
	Listener     net.Listener
	ConnArray    map[string]net.Conn
	HeartBeatMap map[string]chan byte
	Handle       HandleFunc
	MaxConnNum   int
	TimeOut      time.Duration
	Lock         sync.RWMutex
}

func (s *Server) NewServer(i string, p string, m int, t int) {
	s.Lock = sync.RWMutex{}
	s.Lock.Lock()
	s.Ip = i
	s.Port = p
	s.MaxConnNum = m
	s.TimeOut = time.Second * time.Duration(t)
	s.ConnArray = make(map[string]net.Conn, s.MaxConnNum)
	s.HeartBeatMap = make(map[string]chan byte, s.MaxConnNum)
}

func (s *Server) UnRegister(name string) error {
	s.Lock.Lock()
	err := s.ConnArray[name].Close()
	if err != nil {
		return err
	}
	delete(s.ConnArray, name)
	s.Lock.Unlock()
	return nil
}

func (s *Server) Register(name string, c net.Conn) {
	s.Lock.Lock()
	s.ConnArray[name] = c
	s.Lock.Unlock()
}

//心跳计时，根据GravelChannel判断Client是否在设定时间内发来信息
func HeartBeating(conn net.Conn, readerChannel chan byte, timeout int) {
	select {
	case fk := <-readerChannel:
		lib.Info(conn.RemoteAddr().String() + "receive data string:" + string(fk))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	case <-time.After(time.Second * time.Duration(timeout)):
		lib.Warn("It's really weird to get Nothing!!!")
		conn.Close()
		break
	}
}

func GravelChannel(n []byte, mess chan byte) {
	for _, v := range n {
		mess <- v
	}
	close(mess)
}

func (s *Server) Start() {
	addr, err := net.ResolveTCPAddr("tcp", s.Ip+":"+s.Port)
	if err != nil {
		panic(err)
	}
	s.Listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		c, err := s.Listener.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		s.Register(c.RemoteAddr().String(), c)
		for {
			buffer := make([]byte, 1024)
			n, err := c.Read(buffer)
			if err != nil {

			}
		}
	}
}
