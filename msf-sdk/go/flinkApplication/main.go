package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2"
	"github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2/types"
)

func main() {
	// Command-line argument parsing
	actionPtr := flag.String("action", "", "Action to perform (create or start)")
	applicationNamePtr := flag.String("applicationName", "", "Name of the Flink application")
	codeS3BucketPtr := flag.String("codeS3Bucket", "", "S3 bucket ARN containing application code")
	codeS3FileKeyPtr := flag.String("codeS3FileKey", "", "S3 File Key (Including prefix)")
	applicationRolePtr := flag.String("applicationRole", "", "Service Execution Role")
	snapshotNamePtr := flag.String("snapshotName", "", "Name of the Snapshot")
	appRestoreTypePtr := flag.String("appRestoreType", "", "Restore type when creating application")

	flag.Parse()

	// Load AWS credentials
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		os.Exit(1)
	}

	switch *actionPtr {
	case "create":
		if *applicationNamePtr == "" || *codeS3BucketPtr == "" || *applicationRolePtr == "" || *codeS3FileKeyPtr == "" {
			fmt.Println("Missing required arguments for create application action. Use the --help flag for details.")
			os.Exit(1)
		}
		err = createApplication(*applicationNamePtr, *codeS3BucketPtr, *applicationRolePtr, *codeS3FileKeyPtr, cfg)
		if err != nil {
			fmt.Println("Error creating Flink application:", err)
			os.Exit(1)
		}
	case "start":
		if *applicationNamePtr == "" {
			fmt.Println("Missing required arguments for start application action. Use the --help flag for details.")
			os.Exit(1)
		}
		err = startApplication(*applicationNamePtr, *appRestoreTypePtr, cfg)
		if err != nil {
			fmt.Println("Error starting Flink application:", err)
			os.Exit(1)
		}
	case "stop":
		if *applicationNamePtr == "" {
			fmt.Println("Missing required arguments for stop application action. Use the --help flag for details.")
			os.Exit(1)
		}
		err = stopApplication(*applicationNamePtr, cfg)
		if err != nil {
			fmt.Println("Error stopping Flink application:", err)
			os.Exit(1)
		}
	case "describe":
		if *applicationNamePtr == "" {
			fmt.Println("Missing required arguments for describe application action. Use the --help flag for details.")
		}
		err = describeApplication(*applicationNamePtr, cfg)
		if err != nil {
			fmt.Println("Error describing Flink application:", err)
			os.Exit(1)
		}
	case "create-snapshot":
		if *applicationNamePtr == "" || *snapshotNamePtr == "" {
			fmt.Println("Missing required arguments for sreate-snapshot application action. Use the --help flag for details.")
		}
		err = createApplicationSnapshot(*applicationNamePtr, *snapshotNamePtr, cfg)
		if err != nil {
			fmt.Println("Error creating snapshot", err)
			os.Exit(1)
		}
	}
}

func getARNofCloudWatchLogSteam(appName string, logGrpName string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		return "", err
	}

	cloudwatchClient := cloudwatchlogs.NewFromConfig(cfg)

	// Create msf log group
	createLogGrpInput := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(logGrpName),
	}

	_, err = cloudwatchClient.CreateLogGroup(context.TODO(), createLogGrpInput)
	if err != nil {
		fmt.Printf("Error creating log group: %v", err)
		return "", err
	}

	// Create msf log stream
	createLogStreamInput := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(logGrpName),
		LogStreamName: aws.String(appName),
	}

	_, err = cloudwatchClient.CreateLogStream(context.TODO(), createLogStreamInput)

	if err != nil {
		fmt.Printf("Error creating log stream: %v", err)
		return "", err
	}

	// Describe log stream arn
	describeLogStreamInput := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(logGrpName),
		LogStreamNamePrefix: aws.String(appName),
	}

	describeLogStreamsOutput, err := cloudwatchClient.DescribeLogStreams(context.TODO(), describeLogStreamInput)
	if err != nil {
		fmt.Printf("Error describing the log stream: %v", err)
		return "", err
	}

	return *describeLogStreamsOutput.LogStreams[0].Arn, nil

}

func createApplication(appName string, codeBucketARN string, applicationRole string, codeS3FileKey string, cfg aws.Config) error {
	client := kinesisanalyticsv2.NewFromConfig(cfg)

	logStreamARN, err := getARNofCloudWatchLogSteam(appName, "/zzhe/flinkapps")

	if err != nil {
		fmt.Printf("Error getting arn of logStream: %v", err)
		return err
	}

	// Define application configuration
	applicationCodeConfiguration := &types.ApplicationCodeConfiguration{
		CodeContentType: types.CodeContentTypeZipfile,
		CodeContent: &types.CodeContent{
			S3ContentLocation: &types.S3ContentLocation{
				BucketARN: aws.String(codeBucketARN),
				FileKey:   aws.String(codeS3FileKey),
			},
		},
	}

	createAppInput := &kinesisanalyticsv2.CreateApplicationInput{
		ApplicationName:      aws.String(appName),
		RuntimeEnvironment:   types.RuntimeEnvironmentFlink118,
		ServiceExecutionRole: aws.String(applicationRole),
		ApplicationConfiguration: &types.ApplicationConfiguration{
			FlinkApplicationConfiguration: &types.FlinkApplicationConfiguration{
				MonitoringConfiguration: &types.MonitoringConfiguration{
					ConfigurationType: types.ConfigurationTypeCustom,
					LogLevel:          types.LogLevelInfo,
					MetricsLevel:      types.MetricsLevelOperator,
				},
			},
			ApplicationCodeConfiguration: applicationCodeConfiguration,
			ApplicationSnapshotConfiguration: &types.ApplicationSnapshotConfiguration{
				SnapshotsEnabled: aws.Bool(false),
			},
		},
		CloudWatchLoggingOptions: []types.CloudWatchLoggingOption{
			{
				LogStreamARN: aws.String(logStreamARN),
			},
		},
	}

	createApplicationOutput, err := client.CreateApplication(context.TODO(), createAppInput)
	if err != nil {
		fmt.Println("Error creating Flink application:", err)
		return err
	}

	fmt.Println("Application is created. ARN:", *createApplicationOutput.ApplicationDetail.ApplicationARN)

	return nil
}

func startApplication(appName string, appRestoreType string, cfg aws.Config) error {
	client := kinesisanalyticsv2.NewFromConfig(cfg)

	if appRestoreType == "RESTORE_FROM_LATEST_SNAPSHOT" {
		startApplicationInput := &kinesisanalyticsv2.StartApplicationInput{
			ApplicationName: aws.String(appName),
		}
		_, err := client.StartApplication(context.TODO(), startApplicationInput)
		if err != nil {
			fmt.Printf("Error starting Flink Application %s: %v\n", appName, err)
			return err
		}
	} else {
		startApplicationInput := &kinesisanalyticsv2.StartApplicationInput{
			ApplicationName: aws.String(appName),
		}
		_, err := client.StartApplication(context.TODO(), startApplicationInput)
		if err != nil {
			fmt.Printf("Error starting Flink Application %s: %v\n", appName, err)
			return err
		}
	}

	fmt.Printf("Application %s is started.\n", appName)

	return nil
}

func stopApplication(appName string, cfg aws.Config) error {
	client := kinesisanalyticsv2.NewFromConfig(cfg)

	stopApplicationInput := &kinesisanalyticsv2.StopApplicationInput{
		ApplicationName: aws.String(appName),
		Force:           aws.Bool(false),
	}
	_, err := client.StopApplication(context.TODO(), stopApplicationInput)
	if err != nil {
		fmt.Printf("Error stopping Flink Application %s: %v\n", appName, err)
		return err
	}

	fmt.Printf("Application %s is stopped.\n", appName)

	return nil
}

func describeApplication(appName string, cfg aws.Config) error {
	client := kinesisanalyticsv2.NewFromConfig(cfg)

	describeApplicationInput := &kinesisanalyticsv2.DescribeApplicationInput{
		ApplicationName:          aws.String(appName),
		IncludeAdditionalDetails: aws.Bool(false),
	}
	describeApplicationOutput, err := client.DescribeApplication(context.TODO(), describeApplicationInput)
	if err != nil {
		fmt.Printf("Error describe Flink Application %s: %v\n", appName, err)
		return err
	}

	fmt.Printf("Application %s status is %s\n.", appName, describeApplicationOutput.ApplicationDetail.ApplicationStatus)

	return nil
}

func createApplicationSnapshot(appName string, snapshotName string, cfg aws.Config) error {
	client := kinesisanalyticsv2.NewFromConfig(cfg)

	createApplicationSnapshotInput := &kinesisanalyticsv2.CreateApplicationSnapshotInput{
		ApplicationName: aws.String(appName),
		SnapshotName:    aws.String(snapshotName),
	}

	_, err := client.CreateApplicationSnapshot(context.TODO(), createApplicationSnapshotInput)
	if err != nil {
		fmt.Printf("Error create snapshot %s for Flink Application %s: %v\n", snapshotName, appName, err)
		return err
	}

	fmt.Printf("The snapshot %s of application %s is created\n", snapshotName, appName)
	return nil

}
