/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alexander.rassanov/pow-tcp-server/pkg/pow"
	"alexander.rassanov/pow-tcp-server/pkg/protocol"
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/spf13/cobra"
	"log"
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
		gob.Register(pow.HashCashData{})
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

		_, err = conn.Write(protocol.NewMessage(protocol.RequestChallenge, "").Encode())
		if err != nil {
			return err
		}
		conn.Write([]byte{'\n'})
		connReader := bufio.NewReader(conn)
		b, err := connReader.ReadBytes('\n')
		if err != nil {
			return err
		}
		log.Printf("received bytes %v", b)
		m, err := protocol.ParseMessage(b)
		fmt.Println("Parse", m, err)
		if err != nil {
			return err
		}
		hd, ok := m.Payload.(pow.HashCashData)
		if !ok {
			return protocol.ErrBadPayload
		}
		resolvedhd, err := hd.Resolve()
		log.Printf("resolved %v", resolvedhd)
		if err != nil {
			return err
		}
		_, err = conn.Write(protocol.NewMessage(protocol.RequestService, resolvedhd).Encode())
		conn.Write([]byte{'\n'})
		log.Printf("to server hash %v", resolvedhd)
		if err != nil {
			return err
		}
		bytes, err := connReader.ReadBytes('\n')
		log.Printf("from requestServer %x", bytes)
		if err != nil {
			return err
		}
		m, err = protocol.ParseMessage(bytes)
		if err != nil {
			return err
		}
		quote, ok := m.Payload.(string)
		if !ok {
			return protocol.ErrBadPayload
		}
		fmt.Println(quote)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getQuoteCmd)

	getQuoteCmd.PersistentFlags().String("host", "localhost", "Client will be connected to this host to get quotes")
	getQuoteCmd.PersistentFlags().Int("port", 1234, "The port of the server")
}
