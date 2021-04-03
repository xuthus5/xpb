package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"pastebin/common"
	"pastebin/server"
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
		Use:   "post",
		Short: "Post code segment to https://xpb.xuthus.cc",
		Args: func(cmd *cobra.Command, args []string) error {
			_, err := cmd.Flags().GetStringSlice("post")
			if err != nil {
				log.Errorf("get flags err: %+v", err)
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			//var cs driver.CodeSegmentRecord
			//err := json.Unmarshal([]byte(record), &cs)
			//if err != nil {
			//	log.Errorf("err: %v", err)
			//	return
			//}
			//
			//if cs.Title == "" || cs.Content == "" {
			//	log.Errorf("post title or content empty")
			//	return
			//}
			//
			//var pb server.Response
			//if err := gout.POST("").Debug(true).SetBody(cs).BindJSON(&pb).Do(); err != nil {
			//	log.Errorf("post err: %v", err)
			//	return
			//}
			//
			//_, _ = fmt.Fprintf(os.Stdout, "url: https://cs.xuthus.cc/%v", pb.Data.(driver.CodeSegmentRecord).ShortKey)
			//params, _ := cmd.Flags().GetStringSlice("post")
			//
			//for i, param := range params {
			//	log.Infof("idx: %v, param: %+v", i, param)
			//}

			//for _, body := range params {
			//	var cs driver.CodeSegmentRecord
			//	if err := json.Unmarshal([]byte(body), &cs); err != nil {
			//		log.Errorf("post body: %v\nget err: %v\nskip...", body, err)
			//		continue
			//	}
			//
			//	if cs.Title == "" || cs.Content == "" {
			//		log.Errorf("post title or content empty")
			//		continue
			//	}
			//
			//	var pb server.Response
			//	if err := gout.POST("").Debug(true).SetBody(cs).BindJSON(&pb).Do(); err != nil {
			//		log.Errorf("post err: %v", err)
			//		continue
			//	}
			//
			//	_, _ = fmt.Fprintf(os.Stdout, "url: https://cs.xuthus.cc/%v", pb.Data.(driver.CodeSegmentRecord).ShortKey)
			//}
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
