package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){
	fmt.Printf("Greet function was invoked with %v", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	response := &greetpb.GreetResponse{
		Result: result,
	}

	return response, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GreetManyTimes was invoked")
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10 ; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{Result: result,}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	response := "Hello ,"
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			fmt.Printf("Received all, now sending response from server and close")
			return stream.SendAndClose(&greetpb.LongGreetResponse{Result: response})

			// Send request
		} else if err != nil {
			log.Fatalf("Error %v", err)
		}

		response += " , " + msg.GetGreeting().GetFirstName()

	}

	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("This function was invoked")
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()

		result := "Hello " + firstName + " !"

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{Result: result})

		if sendErr != nil {
			log.Fatalf("ERR WAS OCCURED %v", sendErr)
			return sendErr
		}


	}
}

func main(){
	fmt.Println("Hello I am server")
	listener, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}


}
