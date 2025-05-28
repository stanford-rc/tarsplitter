package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/stanford-rc/tarsplitter"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tarsplitter",
		Short: "Split tar or tar.gz files on file boundaries",
		Long:  ``,
	}

	cmd.AddCommand(NewSplitCmd())

	return cmd
}

func NewSplitCmd() *cobra.Command {
	splitSize := int64(1000)
	mib := int64(1048576)
	cmd := &cobra.Command{
		Use:   "split <tar file> <destination directory>",
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

			baseName := filepath.Base(sourceFn)
			for {
				ext := filepath.Ext(baseName)
				if ext == ".tar" || ext == ".gz" {
					baseName = baseName[0:len(baseName)-len(ext)]
				} else {
					break
				}
			}
		
			ts, err := tarsplitter.NewTarSplitter(
				dest, baseName, splitSize * mib, useGzip)
			if err != nil {
				return err
			}

			return ts.Split(reader)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&splitSize, "split-size", "s", splitSize, "max size of split in MiB")

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
