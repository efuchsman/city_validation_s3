package main

import (
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/efuchsman/city_validation_s3/internal/citiesapi"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {

	cities := citiesapi.NewClient()
	validElements, err := cities.ReturnValidTmpElements()
	if err != nil {
		log.Fatalf("Error returing valid tmp elements: %v", err)
	}
	invalidElements, err := cities.ReturnInvalidTmpElements()
	if err != nil {
		log.Fatalf("Error returning invalid tmp elements: %v", err)
	}
	badFiles := cities.ReturnUnprocessableFiles()

	// Concurrent processing through the use of goroutines to executes the build of the
	// two JSON files and the CSV files and waits to check for errors until all concurrent goroutines are finished
	var wg sync.WaitGroup

	// Channel for catching errors within the goroutines
	errChan := make(chan error, 3)
	wg.Add(3)
	go func() {
		defer wg.Done()
		log.Infof("Creating valid elements JSON")
		if err := cities.CreateValidElementsJSON(validElements, "valid_elements.json"); err != nil {
			errChan <- err
		}
		log.Infof("Finished creating valid elements JSON")
	}()
	go func() {
		defer wg.Done()
		log.Infof("Creating invalid elements JSON")
		if err := cities.CreateInvalidElementsJSON(invalidElements, "invalid_elements.json"); err != nil {
			errChan <- err
		}
		log.Infof("Finished creating invalid elements JSON")
	}()
	go func() {
		defer wg.Done()
		log.Infof("Creating unprocessable files CSV")
		if err := cities.CreateUnprocessableFilesCSV(badFiles.Files, "unprocessable_files.csv"); err != nil {
			errChan <- err
		}
		log.Infof("Finished creating unprocessable files CSV")
	}()

	// Waits for goroutines to execute and then closes the error channel
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for errors caught within the goroutines
	for err := range errChan {
		log.Fatalf("Received error from goroutine: %v", err)
	}

	// Upload files to S3
	if err := uploadToS3("results/valid_elements.json", "city-validation-s3", "results/valid_elements.json"); err != nil {
		log.Fatalf("Error uploading valid_elements.json to S3: %v", err)
	}
	if err := uploadToS3("results/invalid_elements.json", "city-validation-s3", "results/invalid_elements.json"); err != nil {
		log.Fatalf("Error uploading invalid_elements.json to S3: %v", err)
	}
	if err := uploadToS3("results/unprocessable_files.csv", "city-validation-s3", "results/unprocessable_files.csv"); err != nil {
		log.Fatalf("Error uploading unprocessable_files.csv to S3: %v", err)
	}
}

func uploadToS3(filePath string, bucketName string, objectKey string) error {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	awsRegion := os.Getenv("AWS_REGION")
	awsId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsId, awsKey, ""),
	})

	if err != nil {
		return err
	}

	svc := s3.New(sess)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	return err
}
