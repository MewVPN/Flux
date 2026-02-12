package util

import (
	"crypto/rand"
	"math/big"
)

func AgentName() string {
	adjectives := []string{
		"goofy", "graceful", "brave", "sleepy",
		"curious", "gentle", "wild", "silent",
	}

	animals := []string{
		"shark", "otter", "falcon", "orca",
		"fox", "lynx", "turtle", "raven",
	}

	return pick(adjectives) + " " + pick(animals)
}

func pick(list []string) string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(list))))
	return list[n.Int64()]
}
