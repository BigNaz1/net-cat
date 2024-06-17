package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var welcome string = `
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
`

const (
	tcp        string = "tcp"
	maxOfUsers int    = 10
)

type Server struct {
	ln           net.Listener
	clients      map[net.Conn]*client
	names        map[string]bool
	history      []string
	countOfUsers int
	mx           *sync.Mutex
	rwMx         *sync.RWMutex
}

type client struct {
	Name   string
	Writer *bufio.Writer
}

func NewServer() *Server {
	return &Server{
		clients: make(map[net.Conn]*client),
		names:   make(map[string]bool),
		mx:      &sync.Mutex{},
		rwMx:    &sync.RWMutex{},
	}
}

func (serve *Server) Listen(host, port string) error {
	addr := fmt.Sprint(host + ":" + port)
	listener, err := net.Listen(tcp, addr)
	if err == nil {
		serve.ln = listener
	}
	fmt.Printf("Listening on the : [%s]\n", addr)
	return err
}

func (serve *Server) Run() error {
	for {
		conn, err := serve.ln.Accept()
		if err != nil {
			return err
		}
		serve.mx.Lock()
		if serve.countOfUsers <= maxOfUsers {
			serve.countOfUsers++
			go serve.handleRequest(conn)
		} else {
			conn.Close()
			return errors.New("Max count of connections")
		}
		serve.mx.Unlock()
	}
}

func (serve *Server) Stop() {
	serve.ln.Close()
}

func (serve *Server) handleRequest(conn net.Conn) {

	serve.welcomeMsg(conn)

	client, _ := serve.accept(conn)
	// how handle error?

	serve.broadcast(fmt.Sprintf("%s has joined our chat...\n", client.Name))

	serve.addClient(conn, client)

	serve.uploadHistory(client)

	serve.sendMessage(conn)

	serve.deleteClient(conn)

	serve.broadcast(fmt.Sprintf("%s has left our chat...\n", client.Name))

	defer conn.Close()
}

func (serve *Server) welcomeMsg(conn net.Conn) {
	conn.Write([]byte(welcome))
}

func (serve *Server) accept(conn net.Conn) (*client, error) {
	writer := bufio.NewWriter(conn)
	client := &client{Writer: writer}
	for {
		conn.Write([]byte("Please, enter your name: "))

		name, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return nil, err
		}
		name = strings.Trim(name, " \r\n")

		if name != "" {

			serve.rwMx.RLock()
			flagofName := serve.names[name]
			serve.rwMx.RUnlock()

			if !flagofName {
				client.Name = name
				serve.addName(name)
				break
			} else {
				conn.Write([]byte("This name is already used, Kindly write a unique name?\n"))
			}
		} else {
			conn.Write([]byte("Name Cannot be empty\n"))
		}

	}
	return client, nil
}

func (serve *Server) addName(name string) {
	serve.rwMx.Lock()
	defer serve.rwMx.Unlock()
	serve.names[name] = true
}

func (serve *Server) broadcast(servMsg string) {
	serve.rwMx.RLock()
	defer serve.rwMx.RUnlock()
	for _, client := range serve.clients {
		client.Writer.WriteString(servMsg)
		client.Writer.Flush()
	}
}

func (serve *Server) addClient(conn net.Conn, client *client) {
	serve.mx.Lock()
	serve.clients[conn] = client
	serve.mx.Unlock()
}

func (serve *Server) uploadHistory(client *client) {
	serve.rwMx.RLock()
	defer serve.rwMx.RUnlock()
	for _, msg := range serve.history {
		client.Writer.WriteString(msg)
		msg += "\n"
		client.Writer.Flush()

	}
}

func (serve *Server) sendMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.Trim(message, " \r\n")
		if message != "" {
			currTime := time.Now().Format("2006-01-02 15:04:05")
			fullMsg := fmt.Sprintf("[%s][%s]:%s\n", currTime, serve.clients[conn].Name, message)
			serve.broadcast(fullMsg)
			serve.rwMx.Lock()
			serve.history = append(serve.history, fullMsg)
			serve.rwMx.Unlock()
		}
	}
}

func (serve *Server) deleteClient(conn net.Conn) {
	serve.rwMx.Lock()
	defer serve.rwMx.Unlock()
	delete(serve.clients, conn)
}
