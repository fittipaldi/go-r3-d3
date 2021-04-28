package middlewaresv1

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fittipaldi/go-r3-d3/app/apiv1/modelsv1"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	authHeader = "Authorization"
	authPrefix = "Bearer"
)

func getHeader(name string, r *http.Request) string {
	hdrs := r.Header[name]
	if len(hdrs) > 0 {
		return hdrs[0]
	}
	return ""
}

func getSecretSessionKeys(authorization string) string {
	token := ""

	authorization = strings.TrimSpace(authorization)
	if authorization != "" {
		hasBearer := false
		re := regexp.MustCompile(`(?i)` + authPrefix)
		for _, match := range re.FindAllString(authorization, -1) {
			if match == authPrefix {
				hasBearer = true
			}
		}
		if hasBearer {
			var re = regexp.MustCompile(`^` + authPrefix)
			token = strings.TrimSpace(re.ReplaceAllString(authorization, `$2`))
		}
	}

	return token
}

// Check if a request is authorized
func CheckTheAuthorization(r *http.Request, gormDB *gorm.DB) (bool, int) {
	authorization := getHeader(authHeader, r)
	log.Info("Header Authorization: " + authorization)
	headerToken := getSecretSessionKeys(authorization)
	log.Info(fmt.Sprintf("Endpoint [%s] Method [%s] Token [%s]", r.URL.Path, r.Method, headerToken))

	if headerToken != "" {
		var token modelsv1.Token
		gormDB.First(&token, "status = ? AND token = ?", "1", headerToken)
		if token.ID <= 0 {
			return false, http.StatusUnauthorized
		}
	} else {
		return false, http.StatusUnauthorized
	}

	return true, http.StatusOK
}
