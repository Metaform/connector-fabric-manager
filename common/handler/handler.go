//  Copyright (c) 2025 Metaform Systems, Inc
//
//  This program and the accompanying materials are made available under the
//  terms of the Apache License, Version 2.0 which is available at
//  https://www.apache.org/licenses/LICENSE-2.0
//
//  SPDX-License-Identifier: Apache-2.0
//
//  Contributors:
//       Metaform Systems, Inc. - initial API and implementation
//

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

const contentType = "application/json"

// ErrorResponse represents a generic JSON error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
	ID      string `json:"id,omitempty"`
}

type HttpHandler struct {
	Monitor system.LogMonitor
}

func (h HttpHandler) WriteError(w http.ResponseWriter, message string, statusCode int) {
	h.WriteErrorWithID(w, message, statusCode, "")
}

// WriteErrorWithID writes a JSON error response to the response writer
func (h HttpHandler) WriteErrorWithID(w http.ResponseWriter, message string, statusCode int, errorID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	}

	if errorID != "" {
		response.ID = errorID
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Monitor.Infow("Error encoding JSON error response: %v", err)
	}
}

func (h HttpHandler) InvalidMethod(w http.ResponseWriter, req *http.Request, expectedMethod string) bool {
	if req.Method != expectedMethod {
		h.WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return true
	}
	return false
}

func (h HttpHandler) ReadPayload(w http.ResponseWriter, req *http.Request, definition any) bool {
	// Read the request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.WriteError(w, "Failed to read request body", http.StatusBadRequest)
		return false
	}

	defer req.Body.Close()

	if err := json.Unmarshal(body, definition); err != nil {
		h.WriteError(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return false
	}

	if err := api.Validator.Struct(definition); err != nil {
		h.WriteError(w, "Invalid definition: "+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

func (h HttpHandler) HandleError(w http.ResponseWriter, err error) {
	switch {
	case types.IsClientError(err):
		var clientErr types.ClientError
		errors.As(err, &clientErr)
		if badReq, ok := clientErr.(types.BadRequestError); ok {
			h.WriteError(w, fmt.Sprintf("Bad request: %s", badReq.Message), http.StatusBadRequest)
		} else {
			h.WriteError(w, fmt.Sprintf("Client error: %v", clientErr), http.StatusBadRequest)
		}
	case types.IsFatal(err):
		id := uuid.New().String()
		h.Monitor.Infow("Internal Error [%s]: %v", id, err)
		h.WriteError(w, fmt.Sprintf("Internal server error occurred [%s]", id), http.StatusInternalServerError)
	default:
		h.WriteError(w, fmt.Sprintf("Operation failed: %s", err.Error()), http.StatusInternalServerError)
	}
}

func (h HttpHandler) Created(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", contentType)
}

func (h HttpHandler) Accepted(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", contentType)
}

func (h HttpHandler) ResponseAccepted(w http.ResponseWriter, response any) {
	h.Accepted(w)
	h.write(w, response)
}

func (h HttpHandler) OK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", contentType)
}

func (h HttpHandler) ResponseOK(w http.ResponseWriter, response any) {
	h.OK(w)
	h.write(w, response)
}

func (h HttpHandler) ResponseCreated(w http.ResponseWriter, response any) {
	h.Created(w)
	h.write(w, response)
}

func (h HttpHandler) write(w http.ResponseWriter, response any) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Monitor.Infow("Error encoding response: %v", err)
	}
}

func (h HttpHandler) ExtractPathVariable(w http.ResponseWriter, req *http.Request, key string) (string, bool) {
	value := chi.URLParam(req, key)
	if value == "" {
		h.WriteError(w, fmt.Sprintf("Missing %s parameter", key), http.StatusBadRequest)
		return "", false
	}
	return value, true
}
