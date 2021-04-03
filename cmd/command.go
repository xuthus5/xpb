package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/spf13/cobra"
	"os"
	"pastebin/common"
	"pastebin/server"
	"pastebin/server/driver"

	log "github.com/sirupsen/logrus"
)

var (
	cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "run a pbx serve in local",
		Run: func(cmd *cobra.Command, args []string) {
			server.NewRouter()
		},
	}

	cmdPost = &cobra.Command{
		Use:   "post [...]",
		Short: "Post code segment to https://xpb.xuthus.cc",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, body := range args {
				var cs driver.CodeSegmentRecord
				if err := json.Unmarshal([]byte(body), &cs); err != nil {
					log.Errorf("post body: %v\nget err: %v\nskip...", body, err)
					continue
				}

				if cs.Title == "" || cs.Content == "" {
					log.Errorf("post title or content empty")
					continue
				}

				var pb server.Response
				if err := gout.POST("").Debug(true).SetBody(cs).BindJSON(&pb).Do(); err != nil {
					log.Errorf("post err: %v", err)
					continue
				}

				_, _ = fmt.Fprintf(os.Stdout, "url: https://cs.xuthus.cc/%v", pb.Data.(driver.CodeSegmentRecord).ShortKey)
			}
		},
	}

	Version string

	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Display version information.",
		Long:  "Display version information.",
		Run: func(cmd *cobra.Command, args []string) {
			common.CrudeOutput(fmt.Sprintf("version info:\npbx-server: %s\npbx-webui: %s", Version, Version))
		},
	}
)
