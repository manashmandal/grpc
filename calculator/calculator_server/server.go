package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"net"
)

type server struct {}

func (*server) Sum(ctx context.Context, request* calculatorpb.SumRequest) (*calculatorpb.SumResponse, error){
	firstNumber := request.GetFirstNumber()
	secondNumber := request.GetSecondNumber()
	result := firstNumber + secondNumber
	response := &calculatorpb.SumResponse{SumResult: result}
	return response, nil
}


func (*server) PrimeNumberDecomposition(request* calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received primenumberdecomposition rpc %v\n", request)
	number := request.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number % divisor == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{PrimeFactor: divisor,}
			stream.Send(res)
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased %v\n", divisor)
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	result := int32(0)
	times := int32(0)

	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			result = result / times
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{Result: result })
		} else if err != nil {
			log.Fatalf("ERROR OCCURED %v", err)
		}

		result += msg.GetNumber()
		times += 1
	}

	return nil
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("Findmaximum is invoked")
	// Maximum
	maximum := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while receiving %v", err)
			return err
		}
		number := req.GetNumber()

		if number > maximum {
			maximum = number
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{Result: maximum})
			if sendErr != nil {
				log.Fatalf("Error while sending %v", sendErr)
				return sendErr
			}
		}
	}
	return nil
}

func main(){
	listener, error := net.Listen("tcp", "0.0.0.0:50051")

	if error != nil {
		log.Fatalf("Can't listen because %v", error)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(listener) ; err != nil {
		log.Fatalf("Can't listen because %v", err)
	}
}
