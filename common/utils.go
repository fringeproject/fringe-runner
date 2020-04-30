package common

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
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

func isAddrInList(addr string, list []string) (bool, error) {
	rawAddr := addr
	seperator := strings.LastIndex(addr, ":")
	if seperator > -1 {
		rawAddr = addr[:seperator]
	}

	ips, err := net.LookupHost(rawAddr)
	if err != nil {
		err = fmt.Errorf("Cannot resolve host: %s", rawAddr)
		return false, err
	}

	for _, ip := range ips {
		for _, x := range list {
			match, _ := regexp.MatchString(x, ip)

			if match {
				return true, nil
			}
		}
	}

	return false, nil
}

func IsValidAddr(addr string) bool {
	isPrivate, err := isAddrInList(addr, privateIpAddresses)
	if err != nil {
		logrus.Debug(err)
		return false
	}

	if isPrivate {
		err := fmt.Errorf("Host is forbiden: %s", addr)
		logrus.Debug(err)
		return false
	}

	return true
}

func GenerateDialer(allowIPAddresses []string) (func(string, string) (net.Conn, error), error) {
	if allowIPAddresses == nil {
		err := fmt.Errorf("Allowed IP addresses cannot be nil.")
		return nil, err
	}

	secureDial := func(network, addr string) (net.Conn, error) {
		isAllowed, err := isAddrInList(addr, allowIPAddresses)
		if err != nil {
			return nil, err
		}

		if !isAllowed {
			isPrivate, err := isAddrInList(addr, privateIpAddresses)
			if err != nil {
				return nil, err
			}

			if isPrivate {
				err := fmt.Errorf("Host is forbiden: %s", addr)
				logrus.Debug(err)
				return nil, err
			}
		}

		// SECURITY: Vuln DNS rebinding
		return dialer.Dial(network, addr)
	}

	return secureDial, nil
}

// func SecureDial(network, addr string) (net.Conn, error) {
// 	if !IsValidAddr(addr) {
// 		err := fmt.Errorf("Host is forbiden: %s", addr)
// 		logrus.Debug(err)
// 		return nil, err
// 	}

// 	return dialer.Dial(network, addr)
// }

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

func GetProxyFromEnv() (string, bool) {
	proxy, exist := os.LookupEnv("HTTP_PROXY")
	if !exist {
		return "", true
	}

	verify, exist := os.LookupEnv("VERIFY_CERT")
	if !exist {
		return proxy, true
	}

	verifyCert, err := strconv.ParseBool(verify)
	if err != nil {
		return proxy, true
	}

	return proxy, verifyCert
}
