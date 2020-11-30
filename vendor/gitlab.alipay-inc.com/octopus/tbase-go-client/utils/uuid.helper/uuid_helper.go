package uuid_helper

import (
	"fmt"
	"math/rand"

	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"
)

func GenerateUUID() (string, error) {
	b := make([]byte, 16)
	var _, err = rand.Read(b)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[UUID_HELPER] generate UUID error. error: %v", err)
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}
