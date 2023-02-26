package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "showdown [file]",
	Short: "Live markdown previewer",
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

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cmd.Flags().Uint16P("port", "p", 1337, "the port on which the server will listen")
}
