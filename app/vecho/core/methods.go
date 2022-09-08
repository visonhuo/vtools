package core

import "math/rand"

const (
	defaultMinPort = 30000
	defaultMaxPort = 65535
)

func randomPort() int {
	return defaultMinPort + rand.Intn(defaultMaxPort-defaultMinPort)
}
