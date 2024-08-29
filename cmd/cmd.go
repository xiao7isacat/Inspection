package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var AddCommand = &cobra.Command{
	Use: "add",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Add(); err != nil {
			os.Exit(1)
		}
	},
}

var UpdateCommand = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Update(); err != nil {
			os.Exit(1)
		}
	},
}

var GetCommand = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Get(); err != nil {
			os.Exit(1)
		}
	},
}

var DeleteCommand = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		if err := Delete(); err != nil {
			os.Exit(1)
		}
	},
}

func Get() error {
	return nil
}

func Delete() error {
	return nil
}

func Update() error {
	return nil
}

func Add() error {
	return nil
}
