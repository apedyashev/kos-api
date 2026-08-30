// Package logger centralizes error logging so every domain logs consistently.
package logger

import "log"

// Error logs a contextual message together with the underlying error.
func Error(message string, err error) {
	log.Printf("%s: %v", message, err)
}
