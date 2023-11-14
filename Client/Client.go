package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
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

var sequenceNumber int64
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
	log.Println("Starting client")

	f, err := os.OpenFile(fmt.Sprintf("log_%v", *port), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	go launchServer()
	ports := strings.Split(*peers, ",")

	for i := 0; i < len(ports); i++ {
		if (ports[i] == "") {
			continue
		}
		p, err := strconv.Atoi(ports[i])

		if err != nil {
			log.Println(("PANIC"))
		}

		connect(int64(p))
		_, _ = portToPeerClient[int64(p)].Connection(context.Background(), &gRPC.Greeting{Port: *port})
	}

	for {
		time.Sleep(time.Duration((rand.Intn(10) + 2)) * time.Second)

		if state == "RELEASED" {
			initiateAccess(portToPeerClient)
		}
	}
}

func launchServer() {
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		log.Printf("Failed to listen on port %v: %v\n", *port, err)
		return
	}
	grpcServer := grpc.NewServer()

	server := &Server{}

	gRPC.RegisterClientConnectionServer(grpcServer, server)

	log.Println("Started listening for incoming messages")
	if err := grpcServer.Serve(list); err != nil {
		log.Printf("Failed to serve gRPC server over port %v %v\n", port, err)
	}
}

func initiateAccess(portToPeer map[int64]gRPC.ClientConnectionClient) {
	log.Printf("Client %v is initiating access request with sequence number %v \n", *port, sequenceNumber)
	state = "WANTED"
	outstandingResponses = len(portToPeer)
	requestTime = sequenceNumber

	for _, val := range portToPeer {
		_, _ = val.RequestAccess(context.Background(), &gRPC.Request{Id: *port, Time: sequenceNumber})
	}

	for outstandingResponses > 0 {

	}
	state = "HELD"

	PerformCriticalSection()
	log.Println("Sending responses to queued clients now that we are releasing access...")
	for (len(queue) > 0) {
		// pop element from queue
		el := queue[0]
		queue = queue[1:]

		log.Printf("Sending response to %v \n", el)
		_, _ = portToPeerClient[el].Receive(context.Background(), &gRPC.Response{Id: *port, Time: sequenceNumber})
	}
	state = "RELEASED"
}

func connect(dialPort int64) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	log.Printf("Client %d: Attemps to dial to peer on port %v", *port, dialPort)

	var conn *grpc.ClientConn

	conn, err := grpc.Dial(fmt.Sprintf(":%v", dialPort), opts...)

	if err != nil {
		log.Printf("Failed to Dial : %v\n", err)
		return
	}

	server := gRPC.NewClientConnectionClient(conn)
	portToPeerClient[dialPort] = server
	// ServerConn = conn
	log.Println("The connection is: ", conn.GetState().String())
}

func UpdateTime(time int64) {
	sequenceNumber = max(sequenceNumber, time) + 1
}

// Performs critical section
func PerformCriticalSection() {
	log.Printf("Client%vDoingCriticalStuff.HelpSpacebarNotSupported\n", *port)
	time.Sleep(time.Duration(5) * time.Second)
}

func (s *Server) Connection(ctx context.Context, msg *gRPC.Greeting) (*gRPC.Empty, error) {
	connect(msg.Port)
	return &gRPC.Empty{}, nil // error handling missing
}

func (s *Server) RequestAccess(ctx context.Context, msg *gRPC.Request) (*gRPC.Empty, error) {
	UpdateTime(msg.Time)
	log.Printf("Received request for access from client %v with sequence number %v \n", msg.Id, msg.Time)
	log.Printf("Comparing my sequence number: %v to the request: %v \n", sequenceNumber, msg.Time)
	// If we are holding access, or want access with higher priority, we queue the response until we are finished
	// Priority is determined by request time, with Id (represented here by port number) as a tiebreaker
	if state == "HELD" || state == "WANTED" && (requestTime < msg.Time || (requestTime == msg.Time && msg.Id < *port)) {
		// queue response
		log.Print(": Responding with NO\n")
		queue = append(queue, msg.Id)
	} else {
		log.Print(": Responding with YES\n")
		_, _ = portToPeerClient[int64(msg.Id)].Receive(context.Background(), &gRPC.Response{Id: *port, Time: sequenceNumber})
		return &gRPC.Empty{}, nil
	}
	return &gRPC.Empty{},nil
}

func (s *Server) Receive(ctx context.Context, msg *gRPC.Response) (*gRPC.Empty, error) {
	log.Printf("Received access permission from client %v. Waiting for %v more responses\n", msg.Id, outstandingResponses - 1)
	outstandingResponses = outstandingResponses - 1
	return &gRPC.Empty{}, nil
}
