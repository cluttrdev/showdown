package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "showdown",
	Short: "Live markdown previewer",
}

var startCmd = &cobra.Command{
	Use:   "start [file]",
	Short: "Start previewing a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		port, err := cmd.Flags().GetUint16("port")
		if err != nil {
			return err
		}

		app := NewApplication(file)
		return app.run(port)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop previewing",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetUint16("port")
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://127.0.0.1:%d/shutdown", port)

		_, err = http.Get(url)
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().Uint16P("port", "p", 1337, "the port on which the server listens")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}
