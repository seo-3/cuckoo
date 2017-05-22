package cmd

import (
	"fmt"
	"time"
	"os"
	"github.com/Songmu/retry"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

// cpCmd represents the cp command
var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a svc library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		err := execute()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func execute() error {

		file, err := os.Open(in)
    if err != nil {
			return fmt.Errorf("err opening file: %s", err)
    }
    defer file.Close()

    cre := credentials.NewStaticCredentials(
        access,
        secret,
        "")

    svc := s3.New(session.New(), &aws.Config{
        Credentials: cre,
        Region:      aws.String(region),
    })

		err = retry.Retry(uint(count), time.Duration(interval)*time.Second, func() error {
    	// return error once in a while
			_, err := svc.PutObject(&s3.PutObjectInput{
	        Bucket: aws.String(bucket),
	        Key:    aws.String(out),
	        Body:   file,
	    })
			return err
		})
		if err != nil {
    	return fmt.Errorf("faild to retire the host: %s", err)
		}
		return nil
}

var count, interval int
var region, access, secret, bucket, in, out string

func init() {

	s3Cmd.AddCommand(cpCmd)

	cpCmd.Flags().StringVar(&region, "region", "ap-northeast-1", "s3 region")
	cpCmd.Flags().StringVar(&access, "access", "", "aws access key id")
	cpCmd.Flags().StringVar(&secret, "secret", "", "aws secret key id")

	cpCmd.Flags().IntVar(&count, "count", 3, "retry max count")
	cpCmd.Flags().IntVar(&interval, "interval", 1, "retry interval second")

	cpCmd.Flags().StringVar(&bucket, "bucket", "", "s3 bucket")
	cpCmd.Flags().StringVar(&in, "in", "", "input file name")
	cpCmd.Flags().StringVar(&out, "out", "", "output file")
}
