package vmWareGo

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// getEnvString returns string from environment variable.
func getEnvString(v string, def string) string {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	return r
}

// getEnvBool returns boolean from environment variable.
func getEnvBool(v string, def bool) bool {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	switch strings.ToLower(r[0:1]) {
	case "t", "y", "1":
		return true
	}

	return false
}

//func processOverride(u *url.URL, c ClientParams) {
//	//envUsername := os.Getenv(envUserName)
//	//envPassword := os.Getenv(envPassword)
//
//	// Override username if provided
//	if envUsername != "" {
//		var password string
//		var ok bool
//
//		if u.User != nil {
//			password, ok = u.User.Password()
//		}
//
//		if ok {
//			u.User = url.UserPassword(envUsername, password)
//		} else {
//			u.User = url.User(envUsername)
//		}
//	}
//
//	// Override password if provided
//	if envPassword != "" {
//		var username string
//
//		if u.User != nil {
//			username = u.User.Username()
//		}
//
//		u.User = url.UserPassword(username, envPassword)
//	}
//}

func StringInSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
