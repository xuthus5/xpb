package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path/filepath"
	"runtime"

	nested "github.com/antonfisher/nested-logrus-formatter"
)

func RunCommand() {
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		NoColors:        true,
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(" [%v:%v]", filepath.Base(fullPath), line)
		},
	})

	var rootCmd = &cobra.Command{
		Use:   "xpb",
		Short: "xpb is a pastebin cli program.the official website is ?",
	}

	rootCmd.AddCommand(cmdServe, cmdPost)
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("exec err: %v\n", err)
		return
	}
}
