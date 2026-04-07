package notify

import (
	stderrors "errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errInvalidNotification = stderrors.New("invalid notification")

func markInvalidNotification(err error) error {
	if err == nil {
		return nil
	}
	return stderrors.Join(errInvalidNotification, err)
}

func isInvalidNotification(err error) bool {
	return stderrors.Is(err, errInvalidNotification)
}

func writePlainText(c *gin.Context, status int, body string) {
	c.String(status, "%s", body)
}

func writeEmptyStatus(c *gin.Context, status int) {
	c.Status(status)
}

func writeProtocolFailure(c *gin.Context, err error, invalidStatus int, invalidBody string, internalBody string, emptyBody bool) {
	if isInvalidNotification(err) {
		if emptyBody {
			writeEmptyStatus(c, invalidStatus)
			return
		}
		writePlainText(c, invalidStatus, invalidBody)
		return
	}

	if emptyBody {
		writeEmptyStatus(c, http.StatusInternalServerError)
		return
	}
	writePlainText(c, http.StatusInternalServerError, internalBody)
}
