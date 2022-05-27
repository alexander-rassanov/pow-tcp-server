/*
Copyright © 2022 Alexander Rassanov <alexander.rassanov@gmail.com>

*/
package cmd

import (
	"alexander.rassanov/pow-tcp-server/pkg/cache"
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
	"alexander.rassanov/pow-tcp-server/pkg/wordwisdom"
	"context"
	"fmt"
	cache2 "github.com/patrickmn/go-cache"
	"github.com/spf13/cobra"
	"log"
	"net"
)

// BuffSize contains how many bytes messages can contain.
const BuffSize = 1024

// ZeroCount represents difficulty of the required challenges.
const ZeroCount = 1

// runCmd represents the run command.
var runCmd = &cobra.Command{
	Use:   "run",
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
		address := fmt.Sprintf("%s:%d", host, port)
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}
		log.Printf("start server on %s", address)
		localCache := cache2.New(cache2.NoExpiration, cache2.NoExpiration)
		ctx := context.Background()
		for {
			conn, err := listener.Accept()
			log.Printf("%s: accept connection", conn.RemoteAddr().String())
			if err != nil {
				return err
			}
			childCtx, cancelFunc := context.WithCancel(ctx)
			defer cancelFunc()
			go handleIncomingRequest(childCtx, conn, localCache)
		}
	},
}

func quoteService() interface{} {
	return wordwisdom.GetRandQuote()
}

func handleIncomingRequest(ctx context.Context, conn net.Conn, cache cache.Cache) {
	defer conn.Close()
	wordWisdomStream := pow.NewStreamWithHashCash(cache, conn.RemoteAddr().String(), ZeroCount, conn, quoteService)
	if err := wordWisdomStream.ProcessStream(ctx); err != nil {
		log.Printf("%s: error: %s", conn.RemoteAddr().String(), err.Error())
	} else {
		log.Printf("%s: service provided", conn.RemoteAddr().String())
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().String("host", "127.0.0.1", "")
	runCmd.PersistentFlags().Int("port", 1234, "")
}
