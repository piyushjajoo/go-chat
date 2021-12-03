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
	"fmt"
	"github.com/piyushjajoo/go-chat/grpc/client"
	"github.com/spf13/cobra"
)

var remoteServerHost, whoDoYouWantToChatWith, yourName string

// withCmd represents the with command
var withCmd = &cobra.Command{
	Use:   "with",
	Short: "with sub-command let's you specify who you want to chat with",
	Long: `with sub-command let's you specify who you chat with`,
	Run: func(cmd *cobra.Command, args []string) {
		if yourName == "" {
			cobra.CheckErr(fmt.Sprint("please provider your name with `name` argument"))
		}
		client.StartClient(yourName, whoDoYouWantToChatWith, remoteServerHost)
	},
}

func init() {
	rootCmd.AddCommand(withCmd)

	withCmd.Flags().StringVarP(&remoteServerHost, "remove-server-host", "s", "localhost:8080", "Remote server host where you want to join chat e.g 10.11.12.13:8080, default is localhost")
	withCmd.Flags().StringVarP(&whoDoYouWantToChatWith, "chatting-with", "c", "all", "comma separated list of users names on the remote host you want to chat with e.g. A,B,C, default is you can chat with all")
	withCmd.Flags().StringVarP(&yourName, "name", "n", "", "your display name you want users to see e.g. Piyush")
}
