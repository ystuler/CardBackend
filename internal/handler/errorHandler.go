package handler

import (
	"back/internal/exceptions"
	"back/internal/util"
	"errors"
	"net/http"
)

func writeAppError(w http.ResponseWriter, err error) {
	var appErr *exceptions.AppError
	if errors.As(err, &appErr) {
		util.WriteError(w, appErr.StatusCode, appErr.Message)
		return
	}

	util.WriteError(w, http.StatusInternalServerError, exceptions.ErrInternalServer)
}
