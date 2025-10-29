OpenStor Go Client SDK for Amazon S3 Compatible Cloud Storage [![Apache V2 License](https://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/openstor/openstor-go/blob/master/LICENSE)
==================================================================================================================================================================================================================================================================================================================================================================================================================

The OpenStor Go Client SDK provides straightforward APIs to access any Amazon S3 compatible object storage.

This Quickstart Guide covers how to install the OpenStor client SDK and create a sample file uploader. For a complete list of APIs and examples, see the [godoc documentation](https://pkg.go.dev/github.com/openstor/openstor-go/v7).

These examples presume a working [Go development environment](https://golang.org/doc/install).

Download from Github
--------------------

From your project directory:

```sh
go get github.com/openstor/openstor-go/v7
```

Initialize an OpenStor Client Object
--------------------------------

The OpenStor client requires the following parameters to connect to an Amazon S3 compatible object storage:

| Parameter         | Description                                                |
|-------------------|------------------------------------------------------------|
| `endpoint`        | URL to object storage service.                             |
| `_openstor.Options_` | All the options such as credentials, custom transport etc. |

```go
package main

import (
	"log"

	"github.com/openstor/openstor-go/v7"
	"github.com/openstor/openstor-go/v7/pkg/credentials"
)

func main() {
	endpoint := "play.min.io"
	accessKeyID := "Q3AM3UQ867SPQQA43P2F"
	secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
	useSSL := true

	// Initialize openstor client object.
	client, err := openstor.New(endpoint, &openstor.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", client) // client is now set up
}
```

Example - File Uploader
-----------------------

This sample code connects to an object storage server, creates a bucket, and uploads a file to the bucket.

### FileUploader.go

This example does the following:

-	Connects to the MinIO `play` server using the provided credentials.
-	Creates a bucket named `testbucket`.
-	Uploads a file named `testdata` from `/tmp`.
-	Verifies the file was created using `mc ls`.

	```go
	// FileUploader.go OpenStor example
	package main

	import (
		"context"
		"log"

		"github.com/openstor/openstor-go/v7"
		"github.com/openstor/openstor-go/v7/pkg/credentials"
	)

	func main() {
		ctx := context.Background()
		endpoint := "your-s3-endpoint.example.com"
		accessKeyID := "YOUR-ACCESSKEYID"
		secretAccessKey := "YOUR-SECRETACCESSKEY"
		useSSL := true

		// Initialize openstor client object.
		client, err := openstor.New(endpoint, &openstor.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}

		// Make a new bucket called testbucket.
		bucketName := "testbucket"
		location := "us-east-1"

		err = client.MakeBucket(ctx, bucketName, openstor.MakeBucketOptions{Region: location})
		if err != nil {
			// Check to see if we already own this bucket (which happens if you run this twice)
			exists, errBucketExists := client.BucketExists(ctx, bucketName)
			if errBucketExists == nil && exists {
				log.Printf("We already own %s\n", bucketName)
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Printf("Successfully created %s\n", bucketName)
		}

		// Upload the test file
		// Change the value of filePath if the file is in another location
		objectName := "testdata"
		filePath := "/tmp/testdata"
		contentType := "application/octet-stream"

		// Upload the test file with FPutObject
		info, err := client.FPutObject(ctx, bucketName, objectName, filePath, openstor.PutObjectOptions{ContentType: contentType})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	}
	```

**1. Create a test file containing data:**

Note: The example uses public S3-compatible endpoint credentials for demonstration purposes. In production, use your own endpoint and credentials.

You can do this with `dd` on Linux or macOS systems:

```sh
dd if=/dev/urandom of=/tmp/testdata bs=2048 count=10
```

or `fsutil` on Windows:

```sh
fsutil file createnew "C:\Users\<username>\Desktop\sample.txt" 20480
```

**2. Run FileUploader with the following commands:**

```sh
go mod init example/FileUploader
go get github.com/openstor/openstor-go/v7
go get github.com/openstor/openstor-go/v7/pkg/credentials
go run FileUploader.go
```

The output resembles the following:

```sh
2023/11/01 14:27:55 Successfully created testbucket
2023/11/01 14:27:55 Successfully uploaded testdata of size 20480
```

**3. Verify the Uploaded File With `mc ls`:**

```sh
# Use your S3-compatible client to list objects in testbucket
[2023-11-01 14:27:55 UTC]  20KiB STANDARD TestDataFile
```

API Reference
-------------

The full API Reference is available here.


### API Reference : Bucket Operations

-	[`MakeBucket`](https://min.io/docs/minio/linux/developers/go/API.html#MakeBucket)
-	[`ListBuckets`](https://min.io/docs/minio/linux/developers/go/API.html#ListBuckets)
-	[`BucketExists`](https://min.io/docs/minio/linux/developers/go/API.html#BucketExists)
-	[`RemoveBucket`](https://min.io/docs/minio/linux/developers/go/API.html#RemoveBucket)
-	[`ListObjects`](https://min.io/docs/minio/linux/developers/go/API.html#ListObjects)
-	[`ListIncompleteUploads`](https://min.io/docs/minio/linux/developers/go/API.html#ListIncompleteUploads)

### API Reference : Bucket policy Operations

-	[`SetBucketPolicy`](https://min.io/docs/minio/linux/developers/go/API.html#SetBucketPolicy)
-	[`GetBucketPolicy`](https://min.io/docs/minio/linux/developers/go/API.html#GetBucketPolicy)

### API Reference : Bucket notification Operations

-	[`SetBucketNotification`](https://min.io/docs/minio/linux/developers/go/API.html#SetBucketNotification)
-	[`GetBucketNotification`](https://min.io/docs/minio/linux/developers/go/API.html#GetBucketNotification)
-	[`RemoveAllBucketNotification`](https://min.io/docs/minio/linux/developers/go/API.html#RemoveAllBucketNotification)
-	[`ListenBucketNotification`](https://min.io/docs/minio/linux/developers/go/API.html#ListenBucketNotification) (MinIO Extension)
-	[`ListenNotification`](https://min.io/docs/minio/linux/developers/go/API.html#ListenNotification) (MinIO Extension)

### API Reference : File Object Operations

-	[`FPutObject`](https://min.io/docs/minio/linux/developers/go/API.html#FPutObject)
-	[`FGetObject`](https://min.io/docs/minio/linux/developers/go/API.html#FGetObject)

### API Reference : Object Operations

-	[`GetObject`](https://min.io/docs/minio/linux/developers/go/API.html#GetObject)
-	[`PutObject`](https://min.io/docs/minio/linux/developers/go/API.html#PutObject)
-	[`PutObjectStreaming`](https://min.io/docs/minio/linux/developers/go/API.html#PutObjectStreaming)
-	[`StatObject`](https://min.io/docs/minio/linux/developers/go/API.html#StatObject)
-	[`CopyObject`](https://min.io/docs/minio/linux/developers/go/API.html#CopyObject)
-	[`RemoveObject`](https://min.io/docs/minio/linux/developers/go/API.html#RemoveObject)
-	[`RemoveObjects`](https://min.io/docs/minio/linux/developers/go/API.html#RemoveObjects)
-	[`RemoveIncompleteUpload`](https://min.io/docs/minio/linux/developers/go/API.html#RemoveIncompleteUpload)
-	[`SelectObjectContent`](https://min.io/docs/minio/linux/developers/go/API.html#SelectObjectContent)

### API Reference : Presigned Operations

-	[`PresignedGetObject`](https://min.io/docs/minio/linux/developers/go/API.html#PresignedGetObject)
-	[`PresignedPutObject`](https://min.io/docs/minio/linux/developers/go/API.html#PresignedPutObject)
-	[`PresignedHeadObject`](https://min.io/docs/minio/linux/developers/go/API.html#PresignedHeadObject)
-	[`PresignedPostPolicy`](https://min.io/docs/minio/linux/developers/go/API.html#PresignedPostPolicy)

### API Reference : Client custom settings

-	[`SetAppInfo`](https://min.io/docs/minio/linux/developers/go/API.html#SetAppInfo)
-	[`TraceOn`](https://min.io/docs/minio/linux/developers/go/API.html#TraceOn)
-	[`TraceOff`](https://min.io/docs/minio/linux/developers/go/API.html#TraceOff)

Explore Further
---------------

-	[Godoc Documentation](https://pkg.go.dev/github.com/openstor/openstor-go/v7)

Contribute
----------

[Contributors Guide](https://github.com/openstor/openstor-go/blob/master/CONTRIBUTING.md)

License
-------

This SDK is distributed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0), see [LICENSE](https://github.com/openstor/openstor-go/blob/master/LICENSE) and [NOTICE](https://github.com/openstor/openstor-go/blob/master/NOTICE) for more information.
