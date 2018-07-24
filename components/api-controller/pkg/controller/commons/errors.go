package commons

import (
	"fmt"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// Logs root cause and returns new error to hide implementation details
func HandleError(rootCause error, msg string) error {
	errId := uuid.NewV4()
	log.Errorf("[Error '%s']: %v", errId, rootCause)
	return fmt.Errorf("%s (error code = '%s')", msg, errId)
}
