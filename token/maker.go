package token

import "time"

type Maker interface {
	//CreateToken creates a token for the given username with the given duration.
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
