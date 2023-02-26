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
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		app := NewApplication(file)
		app.run()
	},
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cmd.Flags().IntP("port", "p", 1337, "the port on which the server will listen")
}
