package main

import (
    "log"
    "fmt"
    "context"
    "io"
    "time"
    "github.com/wagnermengue/grpc/pb"
    "google.golang.org/grpc"
)

func main() {
    connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Couldn't connect to gRPC Server: %v", err)
    }
    defer connection.Close()

    client := pb.NewUserServiceClient(connection)
    //AddUser(client)
    //AddUserVerbose(client)
    //AddUsers(client)
    AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
    req := &pb.User{
        Id: "0",
        Name: "João",
        Email: "j@j.com",
    }

    resp, err := client.AddUser(context.Background(), req)
    if err != nil {
        log.Fatalf("Couldn't make gRPC request: %v", err)
    }

    fmt.Println(resp)
}

func AddUserVerbose(client pb.UserServiceClient) {
    req := &pb.User{
        Id: "0",
        Name: "João",
        Email: "j@j.com",
    }

    respStream, err := client.AddUserVerbose(context.Background(), req)
    if err != nil {
        log.Fatalf("Couldn't make gRPC request: %v", err)
    }

    for {
        stream, err := respStream.Recv()

        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatalf("Couldn't receive the message: %v", err)
        }

        fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
    }
}

func AddUsers(client pb.UserServiceClient) {
    reqs := []*pb.User{
        &pb.User{
            Id: "w1",
            Name: "wagner1",
            Email: "w1@w1.com",
        },
        &pb.User{
            Id: "w2",
            Name: "wagner2",
            Email: "w2@w1.com",
        },
        &pb.User{
            Id: "w3",
            Name: "wagner3",
            Email: "w3@w1.com",
        },
    }

    stream, err := client.AddUsers(context.Background())

    if err != nil {
        log.Fatalf("Error creating request: %v", err)
    }

    for _, req := range reqs {
        stream.Send(req)
        time.Sleep(time.Second * 3);
    }

    res, err := stream.CloseAndRecv()
    if err != nil {
        log.Fatalf("Error receiving response: %v", err)
    }

    fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

    stream, err := client.AddUserStreamBoth(context.Background())
    if err != nil {
        log.Fatalf("Error creating request: %v", err)
    }

    reqs := []*pb.User{
        &pb.User{
            Id: "w1",
            Name: "wagner1",
            Email: "w1@w1.com",
        },
        &pb.User{
            Id: "w2",
            Name: "wagner2",
            Email: "w2@w1.com",
        },
        &pb.User{
            Id: "w3",
            Name: "wagner3",
            Email: "w3@w1.com",
        },
    }

    wait := make(chan int)

    go func() {
        for _, req := range reqs {
            fmt.Println("Sending user: ", req.Name)
            stream.Send(req)
            time.Sleep(time.Second * 2)
        }
        stream.CloseSend()
    }()

    go func() {
        for {
            res, err := stream.Recv()
            if err == io.EOF {
                break
            }
            if err != nil {
                log.Fatalf("Error receiving data: %v", err)
                break
            }
            fmt.Printf("Recebendo user %v com status: %v \n", res.GetUser().GetName(), res.GetStatus())
        }
        close(wait)
    }()

    <-wait
}