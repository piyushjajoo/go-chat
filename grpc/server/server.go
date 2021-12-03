/*
Copyright © 2021 Piyush Jajoo piyush.jajoo1991@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package server implements a server for Chat service.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	pb "github.com/piyushjajoo/go-chat/grpc/gochat"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

// Connection stores the streams for active userss
type Connection struct {
	stream       pb.Broadcast_CreateStreamServer
	id           string
	displayName  string
	chattingWith []string
	active       bool
	error        chan error
}

type Server struct {
	pb.UnimplementedBroadcastServer
	Connection map[string]*Connection
}

func (s *Server) CreateStream(pconn *pb.Connect, stream pb.Broadcast_CreateStreamServer) error {
	grpcLog.Infof("handing connection request for %s", pconn.User.GetDisplayName())

	conn := &Connection{
		stream:       stream,
		id:           pconn.User.GetId(),
		displayName:  pconn.User.GetDisplayName(),
		chattingWith: pconn.GetChattingWith(),
		active:       true,
		error:        make(chan error),
	}

	// check to see if we already have a user with the display name in the connection, if yes
	// drop the connection
	// FIXME: allow multiple users with same name
	if _, ok := s.Connection[conn.displayName]; ok {
		grpcLog.Error("cannot accept this connection as there is already a user with same name")
		return errors.New("try using a different user name")
	} else {
		// accept the connection
		s.Connection[conn.displayName] = conn
	}

	return <-conn.error
}

// canChatWith validates if the sender can chat with provided user
func canChatWith(senderChattingWith, currentUserChattingWith []string, senderUserDisplayName, currentUserDisplayName string) bool {
	// self talking is always healthy
	if senderUserDisplayName == currentUserDisplayName {
		return true
	}
	// sending doesn't want to talk to anyone
	if len(senderChattingWith) == 0 {
		return false
	}

	// sender wants to chat with all, let's check if current user wants to chat or not
	if len(senderChattingWith) == 1 && strings.EqualFold(senderChattingWith[0], "all") {
		grpcLog.Infof("sender %s, current user %s, current user chatting with %v", senderUserDisplayName, currentUserDisplayName, currentUserChattingWith)
		if len(currentUserChattingWith) == 0 { // current user not chatting with anyone
			return false
		} else if len(currentUserChattingWith) == 1 && strings.EqualFold(currentUserChattingWith[0], "all") {
			// current user chatting with all
			return true
		}
		// current user only chatting with selective users
		for _, userName := range currentUserChattingWith {
			if senderUserDisplayName == userName {
				return true
			}
		}
		// current user not chatting with sender
		return false
	}

	// check if sender is chatting with current user
	for _, userName := range senderChattingWith {
		if currentUserDisplayName == userName {
			return true
		}
	}

	return false
}

func (s *Server) BroadcastMessage(_ context.Context, msg *pb.Message) (*pb.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for sendingTo, conn := range s.Connection {
		wait.Add(1)

		go func(sendingTo string, msg *pb.Message, conn *Connection) {
			defer wait.Done()

			senderConn := s.Connection[msg.User.GetDisplayName()]

			if conn.active && canChatWith(senderConn.chattingWith, conn.chattingWith, senderConn.displayName, conn.displayName) {
				grpcLog.Infof("sending message to %s: %v", sendingTo, conn.stream)
				err := conn.stream.Send(msg)
				if err != nil {
					grpcLog.Errorf("error with stream: %v - error: %v; try re-connecting..", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(sendingTo, msg, conn)

	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &pb.Close{}, nil
}

// StartServer starts the server
func StartServer(port string) {

	server := &Server{Connection: make(map[string]*Connection)}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		grpcLog.Fatalf("error creating the server %v", err)
	}

	grpcLog.Infof("starting server at %s", listener.Addr().String())

	pb.RegisterBroadcastServer(grpcServer, server)
	err = grpcServer.Serve(listener)
	if err != nil {
		grpcLog.Fatalf("error starting server at %s: %v", listener.Addr().String(), err)
	}
}
