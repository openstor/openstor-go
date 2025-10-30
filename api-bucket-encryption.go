// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: Apache-2.0
/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2020 MinIO, Inc.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openstor

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"net/url"

	"github.com/openstor/openstor-go/v7/pkg/s3utils"
	"github.com/openstor/openstor-go/v7/pkg/sse"
)

// SetBucketEncryption sets the default encryption configuration on an existing bucket.
// The encryption configuration specifies the default encryption behavior for objects uploaded to the bucket.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bucketName: Name of the bucket
//   - config: Server-side encryption configuration to apply
//
// Returns an error if the operation fails or if config is nil.
func (c *Client) SetBucketEncryption(ctx context.Context, bucketName string, config *sse.Configuration) error {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}

	if config == nil {
		return errInvalidArgument("configuration cannot be empty")
	}

	buf, err := xml.Marshal(config)
	if err != nil {
		return err
	}

	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("encryption", "")

	// Content-length is mandatory to set a default encryption configuration
	reqMetadata := requestMetadata{
		bucketName:       bucketName,
		queryValues:      urlValues,
		contentBody:      bytes.NewReader(buf),
		contentLength:    int64(len(buf)),
		contentMD5Base64: sumMD5Base64(buf),
	}

	// Execute PUT to upload a new bucket default encryption configuration.
	resp, err := c.executeMethod(ctx, http.MethodPut, reqMetadata)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp, bucketName, "")
	}
	return nil
}

// RemoveBucketEncryption removes the default encryption configuration from a bucket.
// After removal, the bucket will no longer apply default encryption to new objects.
// It uses the provided context to control cancellations and timeouts.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bucketName: Name of the bucket
//
// Returns an error if the operation fails.
func (c *Client) RemoveBucketEncryption(ctx context.Context, bucketName string) error {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}

	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("encryption", "")

	// DELETE default encryption configuration on a bucket.
	resp, err := c.executeMethod(ctx, http.MethodDelete, requestMetadata{
		bucketName:       bucketName,
		queryValues:      urlValues,
		contentSHA256Hex: emptySHA256Hex,
	})
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return httpRespToErrorResponse(resp, bucketName, "")
	}
	return nil
}

// GetBucketEncryption retrieves the default encryption configuration from a bucket.
// It uses the provided context to control cancellations and timeouts.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bucketName: Name of the bucket
//
// Returns the bucket's encryption configuration or an error if the operation fails.
func (c *Client) GetBucketEncryption(ctx context.Context, bucketName string) (*sse.Configuration, error) {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return nil, err
	}

	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("encryption", "")

	// Execute GET on bucket to get the default encryption configuration.
	resp, err := c.executeMethod(ctx, http.MethodGet, requestMetadata{
		bucketName:  bucketName,
		queryValues: urlValues,
	})

	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp, bucketName, "")
	}

	encryptionConfig := &sse.Configuration{}
	if err = xmlDecoder(resp.Body, encryptionConfig); err != nil {
		return nil, err
	}

	return encryptionConfig, nil
}
