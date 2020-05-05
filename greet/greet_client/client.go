package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"
)

func main(){
	fmt.Println("Hello I'm client")

	clientConnection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldnt connect %v", err)
	}

	defer clientConnection.Close()

	client := greetpb.NewGreetServiceClient(clientConnection)

	//request := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{FirstName: "Manash", LastName: "Mandal"}}

	//doUnary(client)
	//doServerStreaming(client)
	//doClientStreaming(client)
	doBidiStreaming(client)
}



func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Manash",
			LastName: "Mandal",
		},
	}

	resp, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling greet RPC: %v", err)
	}

	log.Printf("Response from Greet %v", resp.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient){
	fmt.Println("Streaming rpc")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Manash",
			LastName: "mandal",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Cant greet %v", err)
	}

	//msg, err := resStream.Recv()

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading stream %v", err)
		}
		log.Printf("Response from greet many times %v", msg.GetResult())
	}

}

func doBidiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting bidi streaming")

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{
			FirstName: "Manash",
		},
		},

		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{
			FirstName: "Mandal",
		},
		},

		&greetpb.GreetEveryoneRequest{Greeting: &greetpb.Greeting{
			FirstName: "HEY HEY HEY",
		},
		},
	}

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while sending %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message %v", req )
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				stream.CloseSend()
				break
			}
			if err != nil {
				log.Fatalf("ERROR %v while receiving", err)
			}
			fmt.Printf("RESULT WAS %v", res.GetResult())
		}

	}()

	<- waitc
}

func doClientStreaming(c greetpb.GreetServiceClient){
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{
			FirstName: "Manash",
		},
		},

		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{
			FirstName: "Mandal",
		},
		},

		&greetpb.LongGreetRequest{Greeting: &greetpb.Greeting{
			FirstName: "HEY HEY HEY",
		},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Some error occured %v", err)
	}

	for _, req := range requests {
		fmt.Printf("SENDING REQUEST %v\n", req)
		stream.Send(req)
		time.Sleep(1000* time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("ERROR HAPPENED %v", err)
	}

	fmt.Printf("RESPONSE %v", res)

}