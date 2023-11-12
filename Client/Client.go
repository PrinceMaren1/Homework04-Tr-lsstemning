package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"

	gRPC "Homework04/Proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	gRPC.UnimplementedClientConnectionServer
}

var lamportTime int64
var peers = flag.String("peers", "", "ports of other clients") 	// example usage -ports "1000,1001,1002"
var port = flag.Int64("port", 1000, "port to use for this") // port is also used as id this clients ID when communicating with other client

var portToPeerClient map[int64]gRPC.ClientConnectionClient = make(map[int64]gRPC.ClientConnectionClient)
// var ServerConn *grpc.ClientConn maybe not needed?

func main() {
	flag.Parse()
	fmt.Println("New client")
	launchServer()
	
	ports := strings.Split(*peers, ",")
	

	for i := 0; i < len(ports); i++ {
		p, err := strconv.Atoi(ports[i])

		if err != nil {
			fmt.Println(("PANIC"))
		}

		connect(int64(p))
		_, _ = portToPeerClient[int64(p)].Connection(context.Background(), &gRPC.Greeting{Port: *port})
	}
}

func launchServer() {
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Printf("Failed to listen on port %v: %v\n", port, err)
		return
	}
	grpcServer := grpc.NewServer()

	server := &Server{
		// clientStreams: make(map[string]gRPC.ServerConnection_SendMessagesServer),
	}

	gRPC.RegisterClientConnectionServer(grpcServer, server)

	if err := grpcServer.Serve(list); err != nil {
		fmt.Printf("Failed to serve gRPC server over port %v %v\n", port, err)
	}

	fmt.Println("Started listening for incoming messages")
}

func request() {
	
}

func connect(dialPort int64) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	fmt.Printf("Client %d: Attemps to dial on port %v", port, dialPort)

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

func (s *Server) Connection(ctx context.Context, msg *gRPC.Greeting) (*gRPC.Empty, error) {
	connect(msg.Port)
	return &gRPC.Empty{}, nil
}
