package util

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/taxio/errors"
)

type UserClaims struct {
	UserID      uint64
	OfficeID    uint64
	WorkplaceID uint64
	EmployeeID  uint64
	Name        string
	Role        string
}

func GenerateToken(userClaims UserClaims) (string, error) {
	key := os.Getenv("SECRET_KEY")

	claims := jwt.MapClaims{
		"sub":          "AccessToken",
		"user_id":      userClaims.UserID,
		"office_id":    userClaims.OfficeID,
		"workplace_id": userClaims.WorkplaceID,
		"employee_id":  userClaims.EmployeeID,
		"name":         userClaims.Name,
		"role":         userClaims.Role,
		"exp":          time.Now().Add(time.Hour).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		return "", errors.Wrap(err)
	}

	return token, nil
}

func ParseToken(token string) (*jwt.Token, error) {
	key := os.Getenv("SECRET_KEY")

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, errors.Wrap(err)
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return t, nil
}

func GetUserClaims(c *gin.Context) (*UserClaims, error) {
	bearer := c.GetHeader("Authorization")

	re := regexp.MustCompile(`Bearer ([^ ]+)`)
	matches := re.FindStringSubmatch(bearer)
	if len(matches) < 2 {
		return nil, errors.New("Unauthorized")
	}
	token := matches[1]

	t, err := ParseToken(token)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", t.Claims.(jwt.MapClaims)["user_id"].(float64)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	officeID, err := strconv.ParseUint(fmt.Sprintf("%.0f", t.Claims.(jwt.MapClaims)["office_id"].(float64)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	employeeID, err := strconv.ParseUint(fmt.Sprintf("%.0f", t.Claims.(jwt.MapClaims)["employee_id"].(float64)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	workplaceID, err := strconv.ParseUint(fmt.Sprintf("%.0f", t.Claims.(jwt.MapClaims)["workplace_id"].(float64)), 10, 64)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	name := t.Claims.(jwt.MapClaims)["name"].(string)

	role := t.Claims.(jwt.MapClaims)["role"].(string)

	return &UserClaims{
		UserID:      userID,
		OfficeID:    officeID,
		WorkplaceID: workplaceID,
		EmployeeID:  employeeID,
		Name:        name,
		Role:        role,
	}, nil
}
