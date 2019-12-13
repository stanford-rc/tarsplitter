package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/messiaen/tarsplitter"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tarsplitter",
		Short: "Split of join tar of tar.gz files on file boundaries",
		Long:  ``,
	}

	cmd.AddCommand(NewJoinCmd())
	cmd.AddCommand(NewSplitCmd())

	return cmd
}

func NewSplitCmd() *cobra.Command {
	splitSize := int64(1000)
	mib := int64(1048576)
	cmd := &cobra.Command{
		Use:   "split <tar file> <destination direction>",
		Short: "Split tar or tar.gz files on file boundaries",
		Long:  ``,
		Args:  cobra.ExactArgs(2),
		RunE: func(ccmd *cobra.Command, args []string) error {
			sourceFn := args[0]
			dest := args[1]
			reader, err := os.Open(sourceFn)
			if err != nil {
				return err
			}

			useGzip, err := tarsplitter.IsGzip(sourceFn)
			if err != nil {
				return err
			}
			return tarsplitter.SplitTar(dest, filepath.Base(sourceFn), reader, useGzip, splitSize*mib)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&splitSize, "split-size", "s", splitSize, "max size of split in MiB")

	return cmd
}

func NewJoinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join <tar files ...>",
		Short: "Join multiple tar or tar.gz files into one",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(ccmd *cobra.Command, args []string) error {
			fmt.Println("Wouldn't that be nice!")
			return nil
		},
	}
	return cmd
}

func main() {
	cmd := NewCmd()
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
