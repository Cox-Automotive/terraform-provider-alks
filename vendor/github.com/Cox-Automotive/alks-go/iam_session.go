package alks

import (
	"log"
)

// CreateIamSession creates a new IAM STS session. If no error is returned
// then you will received a IamSessionResponse object containing your session
// keys.
func (c *Client) CreateIamSession() (*SessionResponse, error) {
	log.Println("[INFO] Creating IAM session")

	return c.CreateSession(1, true)
}
