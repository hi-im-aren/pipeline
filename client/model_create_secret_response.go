/*
 * Pipeline API
 *
 * Pipeline v0.3.0 swagger
 *
 * API version: 0.3.0
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"time"
)

type CreateSecretResponse struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Id        string    `json:"id"`
	Error     string    `json:"error,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	UpdatedBy string    `json:"updatedBy,omitempty"`
	Version   int32     `json:"version,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
}
