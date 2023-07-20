/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go/parser"
	"io"
	"os"
	"time"

	"github.com/AnthonyHewins/unbabel/internal/cmdline"
	"github.com/marhaupe/json2struct/pkg/generator"
	"github.com/spf13/cobra"
)

// build vars
var (
	version string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "unbabel",
	Short: "Take JSON, pbf, SQL or Go structs and turn them into any of the aforementioned",
	Long: `unbabel INPUTTYPE [FILE instead of stdin]
Take JSON and turn it into PBF, Go structs, or SQL`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Println(version)
			return nil
		}

		n := len(args)
		if n > 3 {
			return fmt.Errorf("incorrect number of args")
		}

		if args[0] == "version" {
			fmt.Println(version)
			return nil
		}

		reader := os.Stdin
		if n == 2 {
			var err error
			reader, err = os.Open(args[1])
			if err != nil {
				return err
			}
		}

		switch args[0] {
		case "json":
			buf, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			golang, err := generator.GenerateOutputFromString(string(buf))
			if err != nil {
				return err
			}

			syntaxTree, err := parser.ParseExpr(golang)
			if err != nil {
				return err
			}
		case "sql":
			return fmt.Errorf("unimplemented")
		case "go":
			return fmt.Errorf("unimplemented")
		case "pbf":
			return fmt.Errorf("unimplemented")
		default:
			return fmt.Errorf("unknown format %s", args[0])
		}

		return nil
	},
}

func init() {
	f := rootCmd.Flags()
	f.BoolP("toggle", "t", false, "Help message for toggle")
	f.BoolP("version", "v", false, "Print version")

	pf := rootCmd.PersistentFlags()

	pf.String(cmdline.LogLevel, "", "Log level to use. None for no logs, or debug, warn/warning, info, error/err")
	pf.String(cmdline.LogExporter, "", "Log exporter to use. By default, it goes off log level: info/debug go to STDOUT, warn/error to STDERR. Specify 'stderr' to write to stderr, and anything else opens a file")
	pf.String(cmdline.LogFmt, "", "Log format to use. Blank or 'json' will create a json logger, or you can use logfmt/text")
	pf.Bool(cmdline.LogSource, false, "Make all logging show where the log occurred")

	pf.String("trace-exporter", "", "Export data using this exporter. Options are stdout (can be configured to go to a file using trace-exporter-arg), otlp with gRPC, jaegar. Use 'none' or leave blank to skip tracing")
	pf.String("trace-exporter-arg", "", "Export data using this URI. For otlp and jaegar this will point to the collector of tracing, for stdout this will point to a file rather than stdout")
	pf.Duration("trace-exporter-timeout", time.Second*5, "How long the tracer will try to export before it abandons the whole process (not supported for all trace exporters)")

	pf.StringP("json-out", "j", "", "Return JSON output to this filename. Blank to not create output. If the file exists, it will be appended with a newline; if it doesn't, it will create the file")
	pf.StringP("sql-out", "s", "", "Return SQL output to this filename. Blank to not create output. If the file exists, it will be appended with a newline; if it doesn't, it will create the file")
	pf.StringP("go-out", "g", "", "Return Go struct output to this filename. Blank to not create output. If the file exists, it will be appended with a newline; if it doesn't, it will create the file")
	pf.StringP("pbf-out", "p", "", "Return protobuf output to this filename. Blank to not create output. If the file exists, it will be appended with a newline; if it doesn't, it will create the file")
}
