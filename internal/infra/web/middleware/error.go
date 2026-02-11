// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package middleware

type UnauthorizedError struct {
	Message string `json:"message"`
}
