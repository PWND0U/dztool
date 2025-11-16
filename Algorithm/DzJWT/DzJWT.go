package DzJWT

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func jwtErrParse(err error) error {
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return fmt.Errorf("令牌签名无效")
	}
	if errors.Is(err, jwt.ErrTokenExpired) {
		return fmt.Errorf("令牌已过期")
	}
	if errors.Is(err, jwt.ErrTokenNotValidYet) {
		return fmt.Errorf("令牌未生效")
	}
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return fmt.Errorf("令牌格式错误")
	}
	return err
}

type JWTClaims struct {
	jwt.RegisteredClaims
	ExtraData map[string]interface{} `json:"extra_data,omitempty"`
}

func NewJWTClaims() *JWTClaims {
	return &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  nil,
			ExpiresAt: nil,
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
		ExtraData: make(map[string]interface{}),
	}
}

func NewJWTClaimsByToken(jwtString, secret string) *JWTClaims {
	claims, err := ParseJWTToEntity(jwtString, secret)
	if err != nil {
		return nil
	}
	return claims
}

func (jc *JWTClaims) SetIssuer(issuer string) {
	jc.Issuer = issuer
}

func (jc *JWTClaims) SetSubject(subject string) {
	jc.Subject = subject
}

func (jc *JWTClaims) SetAudience(audience []string) {
	jc.Audience = audience
}

func (jc *JWTClaims) SetExpiresAt(t time.Time) {
	jc.ExpiresAt = jwt.NewNumericDate(t)
}

func (jc *JWTClaims) SetNotBefore(t time.Time) {
	jc.NotBefore = jwt.NewNumericDate(t)
}

func (jc *JWTClaims) SetIssuedAt(t time.Time) {
	jc.IssuedAt = jwt.NewNumericDate(t)
}

func (jc *JWTClaims) SetID(id string) {
	jc.ID = id
}

func (jc *JWTClaims) GenJWT(secret string) (string, error) {
	return GenJWT(*jc, secret)
}

func (jc *JWTClaims) SetExtraData(extraData map[string]interface{}) {
	jc.ExtraData = extraData
}

func (jc *JWTClaims) AddExtraData(extraData map[string]interface{}) {
	maps.Insert(jc.ExtraData, maps.All(extraData))
}

func (jc *JWTClaims) AddExtraDataByKey(key string, value interface{}) {
	jc.ExtraData[key] = value
}

func genJwtClaims(claims JWTClaims) (jwt.MapClaims, error) {
	jsonMap := make(map[string]interface{})
	marshal, err := json.Marshal(claims.RegisteredClaims)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &jsonMap)
	if err != nil {
		return nil, err
	}
	maps.Insert(jsonMap, maps.All(claims.ExtraData))
	return jsonMap, nil
}

func GenJWT(claims JWTClaims, secret string) (string, error) {
	jsonMap, err := genJwtClaims(claims)
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jsonMap).SignedString([]byte(secret))
}

func VerifyJWT(jwtString, secret string) bool {
	_, err := ParseJWT(jwtString, secret)
	return err == nil
}

func ParseJWT(jwtString, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, jwtErrParse(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("令牌格式错误")
}

func ParseJWTToEntity(jwtString, secret string) (*JWTClaims, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, jwtErrParse(err)
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("令牌格式错误")
}
