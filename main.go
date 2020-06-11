// Package main implements a client for Greeter service.
package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/Vinggui/iaut-drivers/go-driver/internal/driverpc"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	token []byte
)

func createDevice() *pb.Device {
	inputOne := &pb.IO{
		Type:  pb.IO_BUTTON,
		Code:  []byte("in3"),
		Name:  "test input 3",
		Value: "true",
	}
	outputFortyTwo := &pb.IO{
		Type:   pb.IO_INT,
		Code:   []byte("out21"),
		Name:   "test output 21",
		Value:  "42",
		Ranges: []string{"0:50"},
	}
	device := &pb.Device{
		Credential: &pb.Credential{
			DriverID: []byte("PMbNixmFx87HsVHS6iAz"),
			Token:    token,
		},
		Code:    []byte("sensor"),
		Address: []byte("127.0.0.1:8080"),
		Name:    "Super sensor",
		Icon:    "$ball",
		Inputs:  []*pb.IO{inputOne},
		Outputs: []*pb.IO{outputFortyTwo},
	}
	return device
}

func makeLog(msg string, log pb.LoggerClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := log.Error(ctx, &pb.LogRequest{
		Credential: &pb.Credential{
			DriverID: []byte("teste"),
			Token:    token,
		},
		Message: msg,
	})
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Server Ans: %v\n", r)
}

func setDev(center pb.CenterAPIClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	device := createDevice()

	_, err := center.SetDevice(ctx, device)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Server Ans: %v\n", r)
}

func report(center pb.CenterAPIClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Report Output
	_, err := center.Report(ctx, &pb.ReportMessage{
		Credential: &pb.Credential{
			DriverID: []byte("PMbNixmFx87HsVHS6iAz"),
			Token:    token,
		},
		DeviceCode: []byte("sensor"),
		OutputCode: []byte("out42"),
		Value:      "23",
	})
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Server Ans: %v\n", r)
}

func getDevs(center pb.CenterAPIClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Get Devices
	r, err := center.GetDevices(ctx, &pb.Credential{
		DriverID: []byte("PMbNixmFx87HsVHS6iAz"),
		Token:    token,
	})
	if err != nil {
		panic(err)
	}
	for {
		_, err := r.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
			// log.Fatalf("%v.ListFeatures(_) = _, %v", center, err)
		}
		// log.Println(dev)
	}
}

func delDev(center pb.CenterAPIClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	device := &pb.Device{
		Credential: &pb.Credential{
			DriverID: []byte("PMbNixmFx87HsVHS6iAz"),
			Token:    token,
		},
		Code: []byte("sensor"),
	}

	// Delete Devices
	_, err := center.DeleteDevice(ctx, device)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Server Ans: %v\n", r)
}

func confirm(center pb.CenterAPIClient, input *pb.InputCommand) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conf := &pb.Confirmation{
		Credential: &pb.Credential{
			DriverID: []byte("PMbNixmFx87HsVHS6iAz"),
			Token:    token,
		},
		Input: input,
	}

	_, err := center.Confirm(ctx, conf)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Server Ans: %v\n", r)
}

func poll(center pb.CenterAPIClient) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		r, err := center.PollRequest(ctx, &pb.Credential{
			DriverID: []byte("PMbNixmFx87HsVHS6iAz"),
			Token:    token,
		})
		if err != nil {
			panic(err)
		}
		for {
			command, err := r.Recv()
			if err == io.EOF ||
				ctx.Err() == context.DeadlineExceeded {
				log.Println("retartarting")
				break
			}
			if err != nil {
				log.Fatalf("%v.ListFeatures(_) = _, %v", center, err)
			}
			confirm(center, command.Input)
			// log.Println(command)
		}
	}
}

func main() {
	token = []byte(os.Args[1])
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log := pb.NewLoggerClient(conn)
	makeLog("Starting", log)
	center := pb.NewCenterAPIClient(conn)
	// setDev(center)
	poll(center)

	// fmt.Printf("Greeting: %v\n", r)
}
