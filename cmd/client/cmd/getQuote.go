/*
Copyright © 2022 Alex Rassanov <alexander.rassanov@gmail.com>
*/
package cmd

import (
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	challenge_response "alexander.rassanov/pow-tcp-server/pkg/pow/challenge-response"
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
	"fmt"
	"github.com/spf13/cobra"
	"net"
)

// getQuoteCmd represents the getQuote command
var getQuoteCmd = &cobra.Command{
	Use:   "getQuote",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		protocol.RegisterType(pow.HashCashData{})
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return err
		}
		data, err := challenge_response.GetServiceByStream(conn)
		if err != nil {
			return err
		}
		if quote, ok := data.(string); !ok {
			return protocol.ErrBadPayload
		} else {
			fmt.Println("Received quote:", quote)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getQuoteCmd)

	getQuoteCmd.PersistentFlags().String("host", "localhost", "Client will be connected to this host to get quotes")
	getQuoteCmd.PersistentFlags().Int("port", 1234, "The port of the server")
}
