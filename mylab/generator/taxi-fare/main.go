package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/google/uuid"
)

type TaxiFare struct {
	ID             string  `json:"id"`
	StartTime      string  `json:"starttime"`
	EndTime        string  `json:"endtime"`
	StartLocation  string  `json:"startlocation"`
	EndLocation    string  `json:"endlocation"`
	Distance       float64 `json:"distance"`
	PassengerCount int     `json:"passengercount"`
	FareAmount     float64 `json:"fareamount"`
}

func generateTaxiFare() TaxiFare {
	endTime := time.Now()

	// Get a random Start Time within 24 hours of End Time
	maxOffset := time.Hour * 24
	startTime := endTime.Add(-time.Duration(rand.Int63n(int64(maxOffset))))
	return TaxiFare{
		ID:             uuid.NewString(),
		StartTime:      startTime.Format("2006-01-02 15:04:05.999"), // Within the 24 hours
		EndTime:        endTime.Format("2006-01-02 15:04:05.999"),
		StartLocation:  fmt.Sprintf("%f, %f", rand.Float64()*180-90, rand.Float64()*360-180),
		EndLocation:    fmt.Sprintf("%f, %f", rand.Float64()*180-90, rand.Float64()*360-180),
		Distance:       rand.Float64() * 50,  // Up to 50 km
		PassengerCount: rand.Intn(5) + 1,     // 1 to 5 passengers
		FareAmount:     rand.Float64() * 100, // Up to $100
	}
}

func sendToKinesis(streamName string, awsRegion string, fareData []byte) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return fmt.Errorf("error loading AWS config: %w", err)
	}

	client := kinesis.NewFromConfig(cfg)

	_, err = client.PutRecord(context.TODO(), &kinesis.PutRecordInput{
		Data:         fareData,
		StreamName:   &streamName,
		PartitionKey: aws.String(uuid.NewString()),
	})

	if err != nil {
		return fmt.Errorf("error sending to Kinesis: %w", err)
	}

	return nil
}

func printDebugFile(filename string, fareData []byte) error {
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer logFile.Close()

	_, err = logFile.WriteString(string(fareData) + "\n")
	if err != nil {
		return fmt.Errorf("error writing to log file: %w", err)
	}
	return nil
}

func main() {
	streamName := "FlinkLabTaxiFareStream"
	awsRegion := "ap-southeast-1"

	eventsPerSecond := 5
	reportingInterval := 5 * time.Second

	eventsSentCount := 0
	ticker := time.NewTicker(reportingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("Sent %d events to Kinesis in the last %v\n", eventsSentCount, reportingInterval)
			eventsSentCount = 0 // Reset count for next interval

		default:
			taxiFare := generateTaxiFare()
			fareJSON, _ := json.Marshal(taxiFare)

			err := sendToKinesis(streamName, awsRegion, fareJSON)
			if err != nil {
				fmt.Println(err)
			} else {
				eventsSentCount++
			}

			// Print Debug log
			err = printDebugFile("taxifare.log", fareJSON)

			if err != nil {
				fmt.Println(err)
			}

			// Calculate sleep time to maintain the desired rate
			elapsed := time.Since(time.Now())
			sleepDuration := time.Second/time.Duration(eventsPerSecond) - elapsed
			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}
		}
	}
}
