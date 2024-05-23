package broker

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ServerOption struct {
	ServerName        string
	ServerDescription string
	ServerAddress     string
	ServerPort        int
	ServerUtilization int
	ServerCapacity    int
	ServerEnabled     bool
}

type BrokerServer struct {
	host          string
	port          int
	serverOptions []ServerOption
	worldSession  []net.Conn
}

// NewBrokerServer creates a new BrokerServer instance with the given host, port, server options, and world session.
//
// Parameters:
// - host: the host address for the server.
// - port: the port number for the server.
// - options: a slice of ServerOption structs representing the server options.
// - worldSession: a slice of net.Conn representing the world session.
//
// Returns:
// - *BrokerServer: a pointer to the newly created BrokerServer instance.
func NewBrokerServer(host string, port int, options []ServerOption, worldSession []net.Conn) *BrokerServer {
	return &BrokerServer{
		host:          host,
		port:          port,
		serverOptions: options,
		worldSession:  worldSession,
	}
}

// Listen starts a TCP server on the specified host and port. It listens for incoming
// connections and handles them in separate goroutines.
//
// No parameters.
// No return value.
func (bs *BrokerServer) Listen() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bs.host, bs.port))
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Broker TCP Bound on %s:%d\n", bs.host, bs.port)
	for _, option := range bs.serverOptions {
		fmt.Printf("Server: %s - %s on port %d\n", option.ServerName, option.ServerDescription, option.ServerPort)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go bs.handleClient(client)
	}
}

// handleClient handles the client connection and reads incoming data from the client.
//
// It takes a net.Conn as a parameter, which represents the client connection.
// The function reads data from the client using the Read method of the net.Conn.
// It checks the length of the received data and performs different actions based on the received command.
// If the command is 0x1013, it generates a login packet and sends it to the client.
// If the command is 0x1100, it generates a directory packet and sends it to the client.
// If the command is unknown, it prints a message indicating the unknown command.
// The function continues reading data from the client until an error occurs or the client connection is closed.
//
// Parameters:
// - client: a net.Conn representing the client connection.
//
// Return type: None.
func (bs *BrokerServer) handleClient(client net.Conn) {
	defer client.Close()
	fmt.Printf("New connection from %s\n", client.RemoteAddr())

	buf := make([]byte, 1024)
	socketRxSum := 0

	for {
		n, err := client.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		if n < 6 {
			fmt.Println("Invalid Packet (length < 6)")
			fmt.Println(buf[:n])
			continue
		}

		payloadSize := int(binary.LittleEndian.Uint16(buf[:2]))
		clientCommand := int(binary.LittleEndian.Uint16(buf[4:6]))
		socketRxSum += payloadSize

		switch clientCommand {
		case 0x1013:
			fmt.Println("BS: Authentication Request")
			loginPacket := bs.generatePacket(-1, 0x1312, intToBytes(0x0000, 2, true))
			client.Write(loginPacket)
		case 0x1100:
			fmt.Println("BS: Server Directory Request")
			directoryPacket := bs.createDirectoryPacket()
			directoryPacket = bs.generatePacket(0, 0x1102, directoryPacket)
			client.Write(directoryPacket)
		default:
			fmt.Printf("Unknown command: 0x%04X\n", clientCommand)
		}
	}
}

// createDirectoryPacket generates a directory packet for the BrokerServer.
//
// It returns a byte slice containing the directory packet.
func (bs *BrokerServer) createDirectoryPacket() []byte {
	directoryPacket := []byte{0x00, 0x00, 0x01, byte(len(bs.serverOptions))}
	for i, option := range bs.serverOptions {
		option.ServerUtilization = len(bs.worldSession)
		directoryPacket = append(directoryPacket, bs.getIndividualServer(option, i)...)
	}
	return directoryPacket
}

// generatePacket generates a packet with the given parameters.
//
// Parameters:
// - sentPacketLength: the length of the previously sent packet.
// - command: the command of the packet.
// - dataBytes: the data to be included in the packet.
//
// Returns:
// - []byte: the generated packet.
func (bs *BrokerServer) generatePacket(sentPacketLength, command int, dataBytes []byte) []byte {
	packetExpectedLength := len(dataBytes) + 6
	packetSequence := getSequence(sentPacketLength + packetExpectedLength)

	if sentPacketLength == -1 {
		packetSequence = 0xCBEB
	}

	response := make([]byte, 0, packetExpectedLength)
	response = append(response, intToBytes(packetExpectedLength, 2)...)
	response = append(response, intToBytes(packetSequence, 2)...)
	response = append(response, intToBytes(command, 2)...)
	response = append(response, dataBytes...)

	return response
}

// boolToInt converts a boolean value to an integer.
//
// Parameters:
// - b: a boolean value to be converted.
//
// Returns:
// - int: the integer representation of the boolean value. If the boolean value is true, it returns 1. Otherwise, it returns 0.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// getIndividualServer generates a byte slice containing the individual server information.
//
// Parameters:
// - option: the ServerOption struct containing the server details.
// - position: the position of the server in the list.
//
// Returns:
// - []byte: the byte slice containing the individual server information.
func (bs *BrokerServer) getIndividualServer(option ServerOption, position int) []byte {
	extendedDescription := fmt.Sprintf("%s\r\n[%d/%d] players online", option.ServerDescription, option.ServerUtilization, option.ServerCapacity)
	response := make([]byte, 0)
	response = append(response, byte(position), 0x00, 0x00)
	response = append(response, byte(len(option.ServerName)))
	response = append(response, option.ServerName...)
	response = append(response, byte(len(extendedDescription)))
	response = append(response, extendedDescription...)
	for _, part := range strings.Split(option.ServerAddress, ".") {
		ip, _ := strconv.Atoi(part)
		response = append(response, byte(ip))
	}
	response = append(response, intToBytes(option.ServerPort, 2, true)...)
	response = append(response, intToBytes(option.ServerUtilization, 2)...)
	response = append(response, intToBytes(option.ServerCapacity, 2)...)
	response = append(response, byte(boolToInt(option.ServerEnabled))) // Convert bool to byte
	return response
}

// intToBytes converts an integer value to a byte slice of the specified size.
//
// Parameters:
// - value: the integer value to be converted.
// - size: the size of the byte slice to be returned.
// - bigEndian (optional): a boolean flag indicating whether the byte slice should be in big-endian format.
//
// Returns:
// - []byte: the byte slice representation of the integer value.
func intToBytes(value int, size int, bigEndian ...bool) []byte {
	result := make([]byte, size)
	if len(bigEndian) > 0 && bigEndian[0] {
		for i := 0; i < size; i++ {
			result[size-i-1] = byte(value & 0xFF)
			value >>= 8
		}
	} else {
		for i := 0; i < size; i++ {
			result[i] = byte(value & 0xFF)
			value >>= 8
		}
	}
	return result
}

// getSequence calculates the sequence value based on the sum of packet lengths.
//
// Parameters:
// - sumPacketLength: the sum of packet lengths.
//
// Returns:
// - int: the calculated sequence value.
func getSequence(sumPacketLength int) int {
	return (((sumPacketLength * 0x43FD) & 0xFFFF) - 0x53FD) & 0xFFFF
}
