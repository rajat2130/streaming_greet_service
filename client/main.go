package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"streaming_greet_service/greetpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// handler function to stream the greetings from gRPC service to HTTP client
func greetHandler(c *gin.Context) {
	// Set up the connection to the gRPC server
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Could not connect to server: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to gRPC server"})
		return
	}
	defer conn.Close()

	// Create a client for the GreetService
	client := greetpb.NewGreetServiceClient(conn)

	// Prepare the request with greeting details
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rajat",
			LastName:  "Kashyap",
		},
	}

	// Set up the gRPC stream for receiving the greeting messages
	stream, err := client.GreetManyTimes(c, req)
	if err != nil {
		log.Printf("Error while calling GreetManyTimes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while calling GreetManyTimes"})
		return
	}

	// Streaming response from gRPC to HTTP client
	c.Stream(func(w io.Writer) bool {
		// Receive the next message from the gRPC stream
		res, err := stream.Recv()
		if err == io.EOF {
			// End of stream
			return false
		}
		if err != nil {
			log.Printf("Error receiving stream response: %v", err)
			return false
		}

		// Send the greeting to the HTTP client
		_, err = fmt.Fprintf(w, "%s\n", res.GetResult())
		if err != nil {
			log.Printf("Error sending data to client: %v", err)
			return false
		}

		// Simulate a slight delay between sending each greeting (optional)
		time.Sleep(1 * time.Second)
		return true
	})
}

func main() {
	// Set up Gin router
	r := gin.Default()

	// Define the route for streaming greetings
	r.GET("/greet", greetHandler)

	// Start the Gin server on port 8080
	r.Run(":8080")
}
