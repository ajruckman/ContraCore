package main

import (
    "context"
    "net"

    . "github.com/ajruckman/xlib"
    "google.golang.org/grpc"

    query "github.com/ajruckman/ContraCore/internal/rpc"
)

type server struct{}

func (s *server) Alert(context.Context, *query.QueryNotification) (*query.QueryNotificationResponse, error) {
    return &query.QueryNotificationResponse{Status: 1}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    Err(err)

    s := grpc.NewServer()
    query.RegisterQueryServer(s, &server{})
    err = s.Serve(lis)
    Err(err)
}
