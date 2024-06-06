package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	defaultPort    = "8989"
	maxConnections = 10
)

type client struct {
	conn     net.Conn
	name     string
	messages chan string
}

type server struct {
	listener net.Listener
	clients  map[*client]bool
	join     chan *client
	leave    chan *client
	messages chan string
	mutex    sync.Mutex
}

func newServer() *server {
	return &server{
		clients:  make(map[*client]bool),
		join:     make(chan *client),
		leave:    make(chan *client),
		messages: make(chan string),
	}
}

func (s *server) start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	s.listener = listener

	go s.run()

	log.Printf("Server started on port %s", port)
	return nil
}

func (s *server) run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		if len(s.clients) >= maxConnections {
			conn.Write([]byte("Server is full. Please try again later.\n"))
			conn.Close()
			continue
		}

		client := &client{
			conn:     conn,
			messages: make(chan string),
		}

		go s.handleClient(client)
	}
}

func (s *server) handleClient(c *client) {
	defer func() {
		s.leave <- c
		c.conn.Close()
	}()

	s.sendWelcomeMessage(c.conn)

	// Read the client's name
	reader := bufio.NewReader(c.conn)
	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading client name: %v", err)
			return
		}
		c.name = strings.TrimSpace(name)

		if c.name != "" {
			break
		}

		c.conn.Write([]byte("Name cannot be empty. Please try again.\n"))
	}

	s.join <- c

	go s.sendMessages(c)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Client closed the connection
				return
			}
			log.Printf("Error reading message from client %s: %v", c.name, err)
			return
		}

		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		s.messages <- fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), c.name, msg)
	}
}

func (s *server) sendWelcomeMessage(conn net.Conn) {
	// Read the logo from logo.txt
	logoContent, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Printf("Error reading logo file: %v", err)
		// Send a default welcome message if the logo file is not found
		conn.Write([]byte("Welcome to TCP-Chat!\n"))
		return
	}

	welcomeMsg := fmt.Sprintf("Welcome to TCP-Chat!\n%s\n", logoContent)
	conn.Write([]byte(welcomeMsg))
}

func (s *server) sendMessages(c *client) {
	for msg := range c.messages {
		_, err := c.conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Printf("Error sending message to client %s: %v", c.name, err)
			return
		}
	}
}

func (s *server) broadcast(msg string, sender *client) {
	for c := range s.clients {
		if c != sender {
			c.messages <- msg
		}
	}
}

func (s *server) sendPreviousMessages(c *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// No need to store past messages separately anymore
}

func main() {
	port := defaultPort
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	server := newServer()
	err := server.start(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	for {
		select {
		case client := <-server.join:
			server.mutex.Lock()
			server.clients[client] = true
			server.mutex.Unlock()
			server.broadcast(fmt.Sprintf("%s joined the chat\n", client.name), client)

		case client := <-server.leave:
			server.mutex.Lock()
			delete(server.clients, client)
			server.mutex.Unlock()
			close(client.messages)
			server.broadcast(fmt.Sprintf("%s left the chat\n", client.name), client)

		case msg := <-server.messages:
			server.broadcast(msg, nil)
		}
	}
}
