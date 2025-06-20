package command

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thinkgos/gin-rest-kit/cmd/gin-gen/util"
)

type RootCmd struct {
	cmd   *cobra.Command
	level string
}

func NewRootCmd() *RootCmd {
	root := &RootCmd{}
	cmd := &cobra.Command{
		Use:           "gin-gen",
		Short:         "gin gen tools",
		Long:          "gin gen tools",
		Version:       util.BuildVersion(),
		SilenceUsage:  false,
		SilenceErrors: false,
		Args:          cobra.NoArgs,
	}
	cobra.OnInitialize(func() {
		textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   false,
			Level:       level(root.level),
			ReplaceAttr: nil,
		})
		slog.SetDefault(slog.New(textHandler))
	})

	cmd.PersistentFlags().StringVarP(&root.level, "level", "l", "info", "log level(debug,info,warn,error)")
	cmd.AddCommand(
		NewDalCmd().Cmd,
		NewApiCmd().Cmd,
	)
	root.cmd = cmd
	return root
}

// Execute adds all child commands to the root command and sets flags appropriately.
func (r *RootCmd) Execute() error {
	return r.cmd.Execute()
}

func level(s string) slog.Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
