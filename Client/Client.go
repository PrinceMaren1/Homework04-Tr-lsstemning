package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	gRPC "Homework04/Proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	gRPC.UnimplementedClientConnectionServer
}

var lamportTime int64
var requestTime int64
var outstandingResponses int
var queue = []int64{}
var peers = flag.String("peers", "", "ports of other clients") // example usage -ports "1000,1001,1002"
var port = flag.Int64("port", 1000, "port to use for this")    // port is also used as id this clients ID when communicating with other client

var portToPeerClient map[int64]gRPC.ClientConnectionClient = make(map[int64]gRPC.ClientConnectionClient)

// valid state values are RELEASED, WANTED, HELD
var state = "RELEASED"

// var ServerConn *grpc.ClientConn maybe not needed?

func main() {
	flag.Parse()
	fmt.Println("Starting client")
	go launchServer()
	ports := strings.Split(*peers, ",")

	for i := 0; i < len(ports); i++ {
		if (ports[i] == "") {
			continue
		}
		p, err := strconv.Atoi(ports[i])

		if err != nil {
			fmt.Println(("PANIC"))
		}

		connect(int64(p))
		_, _ = portToPeerClient[int64(p)].Connection(context.Background(), &gRPC.Greeting{Port: *port})
	}

	for {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

		if state == "RELEASED" {
			initiateAccess(portToPeerClient)
		}
	}
}

func launchServer() {
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		fmt.Printf("Failed to listen on port %v: %v\n", *port, err)
		return
	}
	grpcServer := grpc.NewServer()

	server := &Server{
		// clientStreams: make(map[string]gRPC.ServerConnection_SendMessagesServer),
	}

	gRPC.RegisterClientConnectionServer(grpcServer, server)

	fmt.Println("Started listening for incoming messages")
	if err := grpcServer.Serve(list); err != nil {
		fmt.Printf("Failed to serve gRPC server over port %v %v\n", port, err)
	}
}

func initiateAccess(portToPeer map[int64]gRPC.ClientConnectionClient) {
	state = "WANTED"
	outstandingResponses = len(portToPeer)
	requestTime = lamportTime

	for _, val := range portToPeer {
		_, _ = val.RequestAccess(context.Background(), &gRPC.Request{Id: *port, Time: lamportTime})
	}

	for outstandingResponses > 0 {

	}
	state = "HELD"

	WriteToFile()
	for (len(queue) > 0) {
		// pop element from queue
		el := queue[0]
		queue = queue[1:]

		_, _ = portToPeerClient[el].Receive(context.Background(), &gRPC.Response{Id: *port, Time: lamportTime})
	}
	state = "RELEASED"
}

func connect(dialPort int64) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	fmt.Printf("Client %d: Attemps to dial to peer on port %v", *port, dialPort)

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(fmt.Sprintf(":%v", dialPort), opts...)

	if err != nil {
		fmt.Printf("Failed to Dial : %v\n", err)
		return
	}

	server := gRPC.NewClientConnectionClient(conn)
	portToPeerClient[dialPort] = server
	// ServerConn = conn
	fmt.Println("The connection is: ", conn.GetState().String())
}

func UpdateTime(time int64) {
	lamportTime = max(lamportTime, time) + 1
}

// Performs critical section
func WriteToFile() {
	fmt.Printf("Client %v doing critical stuff\n", *port)
}

func (s *Server) Connection(ctx context.Context, msg *gRPC.Greeting) (*gRPC.Empty, error) {
	connect(msg.Port)
	return &gRPC.Empty{}, nil // error handling missing
}

func (s *Server) requestAccess(ctx context.Context, msg *gRPC.Request) (*gRPC.Empty, error) {
	if state == "HELD" || state == "WANTED" && requestTime < msg.Time {
		// queue response
		queue = append(queue, msg.Id)
	} else {
		_, _ = portToPeerClient[int64(msg.Id)].Receive(context.Background(), &gRPC.Response{Id: *port, Time: lamportTime})
		return &gRPC.Empty{}, nil
	}
	return nil,nil
}

func (s *Server) receive(ctx context.Context, msg *gRPC.Response) (*gRPC.Empty, error) {
	outstandingResponses = outstandingResponses - 1
	return &gRPC.Empty{}, nil
}
