package cmd

import (
	"encoding/json"
	"github.com/guonaihong/gout"
	"github.com/spf13/cobra"
	"pastebin/server"

	log "github.com/sirupsen/logrus"
)

var (
	cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "run a pastebin serve in local.",
		Long:  "run a pastebin serve in local.",
		Run: func(cmd *cobra.Command, args []string) {
			server.NewRouter()
		},
	}

	cmdPost = &cobra.Command{
		Use:   "post [...]",
		Short: "Post code segment to ?",
		Long:  "Post code segment to ?",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, body := range args {
				var pr PostRequest
				if err := json.Unmarshal([]byte(body), &pr); err != nil {
					log.Errorf("post body: %v\nget err: %v\nskip...", body, err)
					continue
				}

				if pr.Title == "" || pr.Content == "" {
					log.Errorf("post title or content empty")
					continue
				}

				var pb PostResponse
				if err := gout.POST("").Debug(true).SetBody(pr).BindJSON(&pb).Do(); err != nil {
					log.Errorf("post err: %v", err)
					continue
				}
			}
		},
	}
)
