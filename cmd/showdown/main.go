package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "showdown [file]",
	Short: "Live markdown previewer",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetUint16("port")
		if err != nil {
			return err
		}

		stop, err := cmd.Flags().GetBool("stop")
		if err != nil {
			return err
		}

		if stop {
			return StopServer(port)
		}

		if len(args) == 0 {
			return fmt.Errorf("Missing file argument")
		}
		file := args[0]

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
	cmd.Flags().Uint16P("port", "p", 1337, "the port the server listens on")
	cmd.Flags().Bool("stop", false, "stop a running server")
}
