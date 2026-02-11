// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package handler

type ValidationError struct {
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

type BadRequestError struct {
	Message string `json:"message"`
}

type InternalServerError struct {
	Message string `json:"message"`
}
