package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type Uq struct {
	UrlString string
	Output    io.Writer
	config    config
}

type config struct {
	Tools map[string]toolSpec
	Urls  map[string]urlSpec
}

func main() {
	uq := &Uq{}

	rootCmd := cobra.Command{
		Use:  "uq",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := context.Background()
			return uq.Run(ctx, args[0], args[1])
		},
	}

	// rootCmd.Flags().StringVar(&uq.UrlString, "url", "-", "the url")
	rootCmd.Execute()
}

func (uq *Uq) Run(ctx context.Context, nickname, urlString string) (err error) {
	uq.Output = os.Stdout

	err = uq.loadConfig(ctx)
	if err != nil {
		return
	}

	originalURL, err := url.Parse(urlString)
	if err != nil {
		return
	}

	u := FromURL(originalURL)

	urlSpec, err := uq.selectUrl(&u)
	if err != nil {
		return
	}

	fmt.Printf("urlSpec = %+v\n", urlSpec)

	toolSpec, err := uq.selectTool(nickname, urlSpec)
	if err != nil {
		return
	}

	fmt.Printf("toolSpec = %+v\n", toolSpec)

	cmdline, err := toolSpec.MergeUrl(u, uq.config)
	if err != nil {
		return
	}

	fmt.Printf("cmdline = %+v\n", cmdline)

	return
}

func (uq *Uq) loadConfig(ctx context.Context) (err error) {
	_, err = toml.DecodeFile("uqcfg.toml", &uq.config)
	if err != nil {
		return
	}

	return
}
