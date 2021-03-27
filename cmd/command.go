package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/spf13/cobra"
	"os"
	"pastebin/server"
	"pastebin/server/driver"

	log "github.com/sirupsen/logrus"
)

var (
	cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "run a pastebin serve in local",
		Long:  "run a pastebin serve in local",
		Run: func(cmd *cobra.Command, args []string) {
			server.NewRouter()
		},
	}

	cmdPost = &cobra.Command{
		Use:   "post [...]",
		Short: "Post code segment to https://cs.xuthus.cc",
		Long:  "Post code segment to https://cs.xuthus.cc",
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
)
