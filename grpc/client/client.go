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

// Package client implements a client for Chat service.
package client

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	pb "github.com/piyushjajoo/go-chat/grpc/gochat"
	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2
var client pb.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

// connect connects the user remote server
func connect(user *pb.User, chattingWith []string) error {
	var streamerror error

	stream, err := client.CreateStream(context.Background(), &pb.Connect{
		User:         user,
		Active:       true,
		ChattingWith: chattingWith,
	})
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	wait.Add(1)
	go func(str pb.Broadcast_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("error reading message: %v", err)
				break
			}
			grpcLog.Infof("%s : %s\n", msg.User.DisplayName, msg.Message)
		}
	}(stream)

	return streamerror
}

// StartClient starts the client and connects to the remote server
func StartClient(name, whoDoYouWantToChatWith, remoteServerHost string) {
	grpcLog.Infof("starting client for %s", name)

	timestamp := time.Now()
	done := make(chan int)

	// get the remote host connection
	conn, err := grpc.Dial(fmt.Sprintf("%s", remoteServerHost), grpc.WithInsecure())
	if err != nil {
		grpcLog.Fatalf("couldn't connect to host %s: %v", remoteServerHost, err)
	}

	// create a broadcast client
	client = pb.NewBroadcastClient(conn)

	// initialize the user details
	id := sha256.Sum256([]byte(timestamp.String() + name))
	user := &pb.User{
		Id:          hex.EncodeToString(id[:]),
		DisplayName: name,
	}

	// region with the remote host
	err = connect(user, strings.Split(whoDoYouWantToChatWith, ","))
	if err != nil {
		grpcLog.Fatalf("error while creating user stream: %v", err)
	}

	wait.Add(1)
	go func() {
		defer wait.Done()

		// take input from command line
		scanner := bufio.NewScanner(os.Stdin)


		for scanner.Scan() {

			// initialize the message
			msg := &pb.Message{
				Id:        user.Id,
				User:      user,
				Message:   scanner.Text(),
				Timestamp: timestamp.String(),
			}

			// send message
			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				grpcLog.Errorf("error sending message: %v", err)
				break
			}

		}

	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	grpcLog.Infof("started client for %s", name)

	<-done
}
