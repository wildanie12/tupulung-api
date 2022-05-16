package auth

import "tupulung/entities"

type AuthServiceInterface interface {
	Login(AuthReq entities.AuthRequest) (interface{}, error)
	Me(ID int, token interface{}) (interface{}, error)
}
