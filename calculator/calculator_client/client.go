package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"
)

func main(){
	calculatorClient, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't connect %v", err)
	}
	defer calculatorClient.Close()

	calculatorService := calculatorpb.NewCalculatorServiceClient(calculatorClient)

	//doUnary(calculatorService)
	//doServerStreaming(calculatorService)
	//doClientStreaming(calculatorService)
	doBidiStreaming(calculatorService)
}



func doUnary(calculatorServiceClient calculatorpb.CalculatorServiceClient){
	req := &calculatorpb.SumRequest{
		FirstNumber: 10,
		SecondNumber: 20,
	}

	response, err := calculatorServiceClient.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Calculatorserviceclient %v", err)
	}

	fmt.Printf("%v", response)


}

func doServerStreaming(client calculatorpb.CalculatorServiceClient){
	req := &calculatorpb.PrimeNumberDecompositionRequest{Number: 12}

	response, err := client.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error %v", err)
	}

	for {
		msg, err := response.Recv()

		if err == io.EOF {
			log.Printf("Received all")
			break
		} else if err != nil {
			log.Fatalf("ERROR %v", err)
		}
		fmt.Printf("RECEIVED %v\n", msg)
	}
}

func doBidiStreaming(client calculatorpb.CalculatorServiceClient) {
	stream, err := client.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream %v", err)
	}

	waitc := make(chan struct{})

	go func(){
		numbers := []int32{4, 7, 2, 18, 100}
		for _, num := range numbers {
			stream.Send(&calculatorpb.FindMaximumRequest{Number: num})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func(){
		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				close(waitc)
				break

			}

			if err != nil {
				log.Fatalf("Error while receiving %v", err)
			}

			maximum := msg.GetResult()

			fmt.Printf("Received a new maximum %v", maximum)
		}
	}()

	<- waitc
}

func doClientStreaming(client calculatorpb.CalculatorServiceClient){

	numbers := [10]int32{}
	for i := 0; i < 10; i++ {
		numbers[i] = int32(i * 100)
	}

	stream,  err := client.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("ERROR COMPUTE  AVG %v\n", err)
	}

	for _, num := range numbers {
		req := &calculatorpb.ComputeAverageRequest{Number: num}
		fmt.Printf("SeNDING %v\n", num)
		stream.Send(req)
	}

	response, error := stream.CloseAndRecv()

	if error != nil {
		log.Fatalf("ERROR %v\n", error)
	}

	fmt.Printf("%v\n", response)
}