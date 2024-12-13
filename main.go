package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := flag.String("endpoint", "https://localhost:9443", "MinIO server endpoint with scheme")
	accessKey := flag.String("access-key", "", "MinIO access key")
	secretKey := flag.String("secret-key", "", "MinIO secret key")
	bucket := flag.String("bucket", "", "Bucket name")
	objectPath := flag.String("object", "", "Object path in the bucket")
	readFull := flag.Bool("read-full", false, "read complete object")
	startRange := flag.String("start", "0", "Start byte range")
	endRange := flag.String("end", "-1", "End byte range (exclusive)")
	secure := flag.Bool("secure", true, "Use secure (HTTPS) connection")
	skipTLSVerify := flag.Bool("skip-tls-verify", false, "Skip TLS certificate verification")

	flag.Parse()
	if *accessKey == "" || *secretKey == "" || *bucket == "" || *objectPath == "" {
		log.Fatalln("Access key, secret key, bucket, and object path are required")
	}

	// Initialize MinIO client
	var httpClient *http.Client
	if *skipTLSVerify {
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	minioClient, err := minio.New(*endpoint, &minio.Options{
		Creds:     credentials.NewStaticV4(*accessKey, *secretKey, ""),
		Secure:    *secure,
		Transport: httpClient.Transport,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Parse byte range values
	start, err := strconv.ParseInt(*startRange, 10, 64)
	if err != nil {
		log.Fatalf("Invalid start range: %v", err)
	}
	end, err := strconv.ParseInt(*endRange, 10, 64)
	if err != nil {
		log.Fatalf("Invalid end range: %v", err)
	}
	opts := minio.GetObjectOptions{}
	if !*readFull {
		if end >= 0 {
			opts.SetRange(start, end-1)
		} else {
			end = 0
			opts.SetRange(start, 0) // Use 0 to indicate no end limit
		}
	}

	// Get object metadata
	oi, err := minioClient.StatObject(context.Background(), *bucket, *objectPath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalf("Failed to get metadata: %v", err)
	}
	oiBytes, err := json.MarshalIndent(oi, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal metadata info: %v", err)
	}
	fmt.Printf("Object Metadata:\n%s\n", string(oiBytes))

	// Retrieve the object
	obj, err := minioClient.GetObject(context.Background(), *bucket, *objectPath, opts)
	if err != nil {
		log.Fatalf("Failed to get object: %v", err)
	}
	defer obj.Close()

	// Read and output the object data within the specified range
	var buf bytes.Buffer
	_, err = io.Copy(&buf, obj)
	if err != nil {
		log.Fatalf("Failed to read object data: %v", err)
	}
	if !*readFull {
		fmt.Printf("Read %d bytes from object with withn range %d-%d\n", buf.Len(), start, end)
		return
	}
	fmt.Printf("Read %d bytes from object\n", buf.Len())
}
