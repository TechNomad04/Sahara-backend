package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sahara/internal/store"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	EntityType string `json:"entity_type"`
	jwt.RegisteredClaims
}

type Token struct {
	Access     string
	Refresh    string
	JTIAcc     string
	JTIRef     string
	ExpAcc     time.Time
	ExpRef     time.Time
	UserID     string
	EntityType string
	Issuer     string
	Audience   string
}

func IssueToken(userID, entityType string) (*Token, error) {
	now := time.Now()

	t := &Token{
		UserID:     userID,
		EntityType: entityType,
		JTIAcc:     uuid.NewString(),
		JTIRef:     uuid.NewString(),
		ExpAcc:     now.Add(15 * time.Minute),
		ExpRef:     now.Add(7 * 24 * time.Hour),
		Issuer:     "sahara",
		Audience:   "sahara-client",
	}

	accClaims := Claims{
		EntityType: entityType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        t.JTIAcc,
			Issuer:    t.Issuer,
			Audience:  jwt.ClaimStrings{t.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
		},
	}

	refClaims := Claims{
		EntityType: entityType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        t.JTIRef,
			Issuer:    t.Issuer,
			Audience:  jwt.ClaimStrings{t.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(t.ExpRef),
		},
	}

	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, accClaims)
	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims)

	var err error

	t.Access, err = acc.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	t.Refresh, err = ref.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func Persist(ctx context.Context, r *store.Redis, t *Token) error {
	if err := r.SetJTI(ctx, "access:"+t.JTIAcc, t.UserID, t.ExpAcc); err != nil {
		return err
	}

	if err := r.SetJTI(ctx, "refresh:"+t.JTIRef, t.UserID, t.ExpRef); err != nil {
		return err
	}

	return nil
}

func SendTokens(c *gin.Context, t *Token, entity any) {
	c.JSON(http.StatusOK, gin.H{
		"access_token":  t.Access,
		"refresh_token": t.Refresh,
		"access_exp":    t.ExpAcc,
		"refresh_exp":   t.ExpRef,
		"user_id":       t.UserID,
		"entity_type":   t.EntityType,
		"entity":        entity,
	})
}

func ParseAccess(tokenStr string) (*Claims, error) {
	return parseWithSecret(
		tokenStr,
		os.Getenv("ACCESS_SECRET"),
	)
}

func ParseRefresh(tokenStr string) (*Claims, error) {
	return parseWithSecret(
		tokenStr,
		os.Getenv("REFRESH_SECRET"),
	)
}

func parseWithSecret(tokenStr, secret string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods(
			[]string{jwt.SigningMethodHS256.Alg()},
		),
	)

	token, err := parser.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil &&
		time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}