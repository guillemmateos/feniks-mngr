package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

func HandlePanic() {
	if r := recover(); r != nil {
		fmt.Println("RECOVER", r)
		debug.PrintStack()
	}
}

var re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

func ParseIp(address string) string {
	return re.FindString(address)
}

func NewId() string {
	return strings.ToLower(shortuuid.New()[:11])
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func DecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func GetQsValue(ctx *gin.Context, key string) string {
	qs := ctx.Request.URL.Query()
	if _, ok := qs[key]; !ok {
		log.Println("invalid qs")
		return ""
	}
	value := qs[key][0]
	if len(value) == 0 {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		log.Println("invalid qs")
		return ""
	}
	return value
}

func ParsePort() int {
	var err error
	defaultPort := 8072
	port := 0
	p := flag.String("port", "", "Web Server Port Number")
	flag.Parse()

	ep := os.Getenv("WEBSERVER_HOST")
	if len(ep) > 0 {
		port, err = strconv.Atoi(ep)
		if err != nil {
			port = defaultPort
			log.Println("An error occurred while converting Redis port value from environment variable: " + err.Error())
		}
	} else if len(*p) > 0 {
		port, err = strconv.Atoi(*p)
		if err != nil {
			port = defaultPort
			log.Println("An error occurred while converting Redis port value from arguments :" + err.Error())
		}
	} else {
		port = defaultPort
	}
	return port
}

func MinMax(array []int) (int, int) {
	if len(array) == 0 {
		return -1, -1
	}
	var max = array[0]
	var min = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
