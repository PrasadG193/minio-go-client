# minio-go-client

### Build
```
$ go build -o ./minio-go-client .
```

### Usage
```
$ ./minio-go-client --help
Usage of ./minio-go-client:
  -access-key string
        MinIO access key
  -bucket string
        Bucket name
  -end string
        End byte range (exclusive) (default "-1")
  -endpoint string
        MinIO server endpoint with scheme (default "https://localhost:9443")
  -object string
        Object path in the bucket
  -secret-key string
        MinIO secret key
  -secure
        Use secure (HTTPS) connection (default true)
  -skip-tls-verify
        Skip TLS certificate verification
  -start string
        Start byte range (default "0")

```
