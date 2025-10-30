// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: Apache-2.0
/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2019-2020 MinIO, Inc.
 *
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
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/openstor/openstor-go/v7/pkg/s3utils"
)

// objectRetention - object retention specified in
// https://docs.aws.amazon.com/AmazonS3/latest/API/Type_API_ObjectLockConfiguration.html
type objectRetention struct {
	XMLNS           string        `xml:"xmlns,attr,omitempty"`
	XMLName         xml.Name      `xml:"Retention"`
	Mode            RetentionMode `xml:"Mode,omitempty"`
	RetainUntilDate *time.Time    `type:"timestamp" timestampFormat:"iso8601" xml:"RetainUntilDate,omitempty"`
}

func newObjectRetention(mode *RetentionMode, date *time.Time) (*objectRetention, error) {
	objectRetention := &objectRetention{}

	if date != nil && !date.IsZero() {
		objectRetention.RetainUntilDate = date
	}
	if mode != nil {
		if !mode.IsValid() {
			return nil, fmt.Errorf("invalid retention mode `%v`", mode)
		}
		objectRetention.Mode = *mode
	}

	return objectRetention, nil
}

// PutObjectRetentionOptions represents options specified by user for PutObject call
type PutObjectRetentionOptions struct {
	GovernanceBypass bool
	Mode             *RetentionMode
	RetainUntilDate  *time.Time
	VersionID        string
}

// PutObjectRetention sets the retention configuration for an object and specific version.
// Object retention prevents an object version from being deleted or overwritten for a specified period.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bucketName: Name of the bucket
//   - objectName: Name of the object
//   - opts: Options including Mode (GOVERNANCE or COMPLIANCE), RetainUntilDate, optional VersionID, and GovernanceBypass
//
// Returns an error if the operation fails or if the retention settings are invalid.
func (c *Client) PutObjectRetention(ctx context.Context, bucketName, objectName string, opts PutObjectRetentionOptions) error {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return err
	}

	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return err
	}

	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("retention", "")

	if opts.VersionID != "" {
		urlValues.Set("versionId", opts.VersionID)
	}

	retention, err := newObjectRetention(opts.Mode, opts.RetainUntilDate)
	if err != nil {
		return err
	}

	retentionData, err := xml.Marshal(retention)
	if err != nil {
		return err
	}

	// Build headers.
	headers := make(http.Header)

	if opts.GovernanceBypass {
		// Set the bypass goverenance retention header
		headers.Set(amzBypassGovernance, "true")
	}

	reqMetadata := requestMetadata{
		bucketName:       bucketName,
		objectName:       objectName,
		queryValues:      urlValues,
		contentBody:      bytes.NewReader(retentionData),
		contentLength:    int64(len(retentionData)),
		contentMD5Base64: sumMD5Base64(retentionData),
		contentSHA256Hex: sum256Hex(retentionData),
		customHeader:     headers,
	}

	// Execute PUT Object Retention.
	resp, err := c.executeMethod(ctx, http.MethodPut, reqMetadata)
	defer closeResponse(resp)
	if err != nil {
		return err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return httpRespToErrorResponse(resp, bucketName, objectName)
		}
	}
	return nil
}

// GetObjectRetention retrieves the retention configuration for an object and specific version.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bucketName: Name of the bucket
//   - objectName: Name of the object
//   - versionID: Optional version ID to target a specific version (empty string for current version)
//
// Returns the retention mode (GOVERNANCE or COMPLIANCE), retain-until date, and any error.
func (c *Client) GetObjectRetention(ctx context.Context, bucketName, objectName, versionID string) (mode *RetentionMode, retainUntilDate *time.Time, err error) {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return nil, nil, err
	}

	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return nil, nil, err
	}
	urlValues := make(url.Values)
	urlValues.Set("retention", "")
	if versionID != "" {
		urlValues.Set("versionId", versionID)
	}
	// Execute GET on bucket to list objects.
	resp, err := c.executeMethod(ctx, http.MethodGet, requestMetadata{
		bucketName:       bucketName,
		objectName:       objectName,
		queryValues:      urlValues,
		contentSHA256Hex: emptySHA256Hex,
	})
	defer closeResponse(resp)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			return nil, nil, httpRespToErrorResponse(resp, bucketName, objectName)
		}
	}
	retention := &objectRetention{}
	if err = xml.NewDecoder(resp.Body).Decode(retention); err != nil {
		return nil, nil, err
	}

	return &retention.Mode, retention.RetainUntilDate, nil
}
