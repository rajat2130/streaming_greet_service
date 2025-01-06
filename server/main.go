package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"streaming_greet_service/greetpb"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

// StreamingGreet is the function to stream multiple greetings
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet Function was Invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	// Send 10 greetings
	for i := 0; i < 10; i++ {
		result := "Hey " + firstName + " " + lastName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		// Send the response to the client
		if err := stream.Send(res); err != nil {
			return fmt.Errorf("failed to send response: %v", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func main() {
	fmt.Println("Starting GRPC Server on port 50052")

	// Set up the listener
	listener, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new GRPC server
	s := grpc.NewServer()

	// Register the server with the GRPC service
	greetpb.RegisterGreetServiceServer(s, &server{})

	// Start serving
	log.Println("Server is ready to serve...")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
