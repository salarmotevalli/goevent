package httpmsg

import (
	"event-manager/pkg/errmsg"
	"event-manager/pkg/richerror"
	"net/http"
)

func Error(err error) (message string, code int) {
	switch err.(type) {
	case richerror.RichError:
		re := err.(richerror.RichError)
		msg := re.Message()

		code := mapKindToHTTPStatusCode(re.Kind())

		// we should not expose unexpected error messages
		if code >= 500 {
			msg = errmsg.ErrorMsgSomethingWentWrong
		}

		return msg, code
	default:
		return err.Error(), http.StatusBadRequest
	}

}

func mapKindToHTTPStatusCode(k richerror.Kind) int {
	switch k {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
