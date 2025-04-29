// file: overflowProtector.go
package main

func truncateMessage(message string, maxLength int) string {
	if len(message) > maxLength {
		return message[:maxLength]
	}
	return message
}
