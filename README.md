### Net Cat

TCP Chat is a streamlined chat server and client developed in Go, designed to facilitate real-time communication over TCP. It allows multiple users to connect simultaneously, offering unique identifiers for each participant and preserving a log of the conversation.

### Starting the Server

To run the TCP Chat server, use the following command:
   go run .

net-cat
   ./TCPChat [host] [port]

 - [host] (optional): The host to bind the server to. Default is "localhost".
 - [port] (optional): The port number to listen on. Default is "8989".

### Example

Run the server on the default host and port:
   ./TCPChat
  
Run the server on a custom host and port:
   ./TCPChat 0.0.0.0 9999


### Connecting to the Chat

1. Connect to the server using a TCP client, such as Telnet or netcat, or use the provided TCPChatClient binary.
2. Enter your desired username when prompted.
3. Start chatting with other connected users.


# TCP Chat

TCP Chat is a streamlined chat server and client developed in Go, designed to facilitate real-time communication over TCP. It allows multiple users to connect simultaneously, offering unique identifiers for each participant and preserving a log of the conversation.


## Features

- Ability to establish a connection to the server using a designated host and port.
- Each participant is assigned a distinct username.
- Ongoing chat history which new participants can view upon joining.
- Enables live communication among all connected users.
- Supports up to 10 concurrent users.
