package main

import (
	"fmt"
	"net/http"
	"os"

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

func Execute(defaultCmd string) {
	// see: https://github.com/spf13/cobra/issues/823#issuecomment-949732548
	var cmdFound bool
outer:
	for _, cmd := range rootCmd.Commands() {
		for _, arg := range os.Args[1:] {
			if cmd.Name() == arg {
				cmdFound = true
				break outer
			}
		}
	}

	if !cmdFound {
		args := append([]string{defaultCmd}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Uint16P("port", "p", 1337, "the port on which the server listens")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}
