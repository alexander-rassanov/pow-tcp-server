/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net"
	"log"

	"github.com/spf13/cobra"
)

// BuffSize contains how much bytes messages can contain.
const BuffSize = 1024
// ZeroCount represents diffucilty of the required challenges.
const ZeroCount = 20

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
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			return err
		}
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return err
		}
		
		for {
			conn, err := listener.Accept()
			if err != nil {
				return err
			}
			go handleIncomingRequest(conn)
		}
	},
}

func handleIncomingRequest(conn net.Conn) {
	defer conn.Close()
	
	
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().String("host", "127.0.0.1", "")
	runCmd.PersistentFlags().Int("port", 1234, "")
}
