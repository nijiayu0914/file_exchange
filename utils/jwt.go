package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"time"
)

var (
	sigKey = []byte("Xdw+Ry%zTp+K1OiG779_DgklyH_tSfs4")
	encKey = []byte("MTPEx#6lnF9eivyIMFXYLdFj2#V36=oR")
)


func GenerateToken(userName string, password string) (token string, err error){
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name":  userName,
		"password": password,
		"exp":  time.Now().Add(48 * time.Hour).Unix(),
	})
	token, err = at.SignedString(sigKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(token string) (iris.Map, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return sigKey, nil
	})
	if err != nil {
		return iris.Map{}, err
	}
	res := iris.Map{
		"user_name": claim.Claims.(jwt.MapClaims)["user_name"].(string),
		"exp": claim.Claims.(jwt.MapClaims)["exp"].(float64),
	}
	return res, nil
}

func Verification(userName string, token string) (checked bool, err error){
	res, err := ParseToken(token)
	if err != nil{
		return false, errors.New("token解析失败")
	}else{
		if res["user_name"].(string) != userName{
			return false, nil
		}
		timeDiff := res["exp"].(float64) - float64(time.Now().Unix())
		if timeDiff > 0{
			return true, nil
		}else {
			return false, nil
		}
	}
}
