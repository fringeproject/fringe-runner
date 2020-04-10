package common

import (
	"encoding/base64"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	dialer = net.Dialer{
		Timeout:   20 * time.Second,
		KeepAlive: 20 * time.Second,
	}

	privateIpAddresses []string = []string{
		"^0\\.0\\.0\\.0",
		"^240\\.0\\.0\\.0",
		"^203\\.0\\.113\\.0",
		"^198\\.51\\.100\\.0",
		"^198\\.18\\.0\\.0",
		"^192\\.88\\.99\\.0",
		"^192\\.0\\.2\\.0",
		"^100\\.64\\.0\\.0",
		"^255\\.255\\.255\\.255",
		"^192\\.168\\.",
		"^172\\.(?:1[6-9]|2[0-9]|3[0-1])\\.",
		"^10\\.",
		"^127\\.",
		"^169\\.254\\.",
	}
)

func AppendIfMissing(slice []string, e string) []string {
	if !StringInSlice(slice, e) {
		return append(slice, e)
	}

	return slice
}

func StringInSlice(slice []string, e string) bool {
	for _, v := range slice {
		if v == e {
			return true
		}
	}
	return false
}

func IsValidAddr(addr string) bool {
	rawIP := addr
	seperator := strings.LastIndex(addr, ":")
	if seperator > -1 {
		rawIP = addr[:seperator]
	}

	ips, err := net.LookupHost(rawIP)
	if err != nil {
		logrus.Warnf("Cannot resolve host %s", rawIP)
		return false
	}

	for _, ip := range ips {
		for _, x := range privateIpAddresses {
			match, _ := regexp.MatchString(x, ip)

			if match {
				logrus.Debugf("Match forbiden host on: [%s] [%s] [%s]", x, ip, addr)
				return false
			}
		}
	}

	return true
}

func SecureDial(network, addr string) (net.Conn, error) {
	if !IsValidAddr(addr) {
		errorMsg := "Host is forbiden: " + addr
		logrus.Info(errorMsg)
		return nil, fmt.Errorf(errorMsg)
	}

	return dialer.Dial(network, addr)
}

func IsIPv4(rawString string) bool {
	parts := strings.Split(rawString, ".")

	if len(parts) < 4 {
		return false
	}

	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}

	}
	return true
}

func CleanHostname(hostname string) string {
	hostname = strings.ToLower(hostname)

	hostname = strings.TrimPrefix(hostname, "*.")
	hostname = strings.TrimPrefix(hostname, ".")
	hostname = strings.TrimSuffix(hostname, ".")

	return hostname
}

func IsHostname(host string) bool {
	host = strings.Trim(host, " ")

	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	return re.MatchString(host)
}

func IsURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func GetBasicAuthHeader(username string, password string) HTTPHeader {
	auth := username + ":" + password
	value := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	return HTTPHeader{
		Name:  "Authorization",
		Value: value,
	}
}
