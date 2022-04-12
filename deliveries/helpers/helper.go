package helpers

import (
	web "tupulung/entities/web"
)

func MakeErrorResponse(status string, code int, err string, links map[string]string) web.ErrorResponse {
	return web.ErrorResponse{
		Status: status,
		Code:   code,
		Error:  err,
		Links:  links,
	}
}
