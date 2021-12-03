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
package cmd

import (
	"github.com/piyushjajoo/go-chat/grpc/server"
	"github.com/spf13/cobra"
)

var port string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server starts the chat server for you",
	Long: `server allows you to start your own chat server and that way you can allow others to chat with you.
The server doesn't store any messages so never worry while using this in your terminal'`,
	Run: func(cmd *cobra.Command, args []string) {
		server.StartServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// server port
	serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port on which you want to start the server, default 8080")
}
