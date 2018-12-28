package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"log"
	"os/exec"
	"strings"
)


var rootCmd = &cobra.Command{
	Use:   "readTags",
	Short: "Read EC2 tags into /etc/sshd IAM configuration",
	Run: ReadTags,
}

var accessKeyId = "";
var secretAccessKey = "";
var instanceId = "";
var region = "";
var sshConfigFile = "/tmp/aws-ec2-ssh.conf";

func init() {
	rootCmd.Flags().StringVarP(&accessKeyId, "accessKeyId", "", os.Getenv("AWS_ACCESS_KEY_ID"), "")
	rootCmd.Flags().StringVarP(&secretAccessKey, "secretAccessKey", "", os.Getenv("SECRET_ACCESS_KEY"), "")
	rootCmd.Flags().StringVarP(&instanceId, "instanceId", "", os.Getenv("EC2_INSTANCE_ID"), "")
	rootCmd.Flags().StringVarP(&region, "region", "", os.Getenv("AWS_DEFAULT_REGION"), "")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ReadTags(cmd *cobra.Command, args []string) {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatal(err)
	}
	svc := ec2.New(sess)
	if instanceId == "" {
		out, err := exec.Command("ec2-metadata --instance-id | cut -d ' ' -f 2").Output()
		if err != nil {
			log.Fatal(err)
		}
		instanceId = string(out)
	}
	tagInput := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceId),
				},
			},
			{
				Name: aws.String("key"),
				Values: []*string{
					aws.String("IAM_AUTHORIZED_GROUPS"),
					aws.String("ASSUME_ROLE"),
					aws.String("SUDOERS_GROUPS"),
					aws.String("IAM_AUTHORIZED_GROUP_TAGS"),
					aws.String("SUDOERS_GROUPS_TAGS"),
					aws.String("LOCAL_GROUPS"),
				},
			},
		},
	}
	result, err := svc.DescribeTags(tagInput)
	if err != nil {
		log.Fatal(err)
	}

	lines := []string{};
	for _, elem := range result.Tags {
		key := *elem.Key
		value := *elem.Value
		lines = append(lines, fmt.Sprintf("%s=%s", key, value));
	}
	fmt.Print(strings.Join(lines, "\n"));
}

