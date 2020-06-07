package sslstriper

import (
	"crypto/tls"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type SSLStriper struct {
}

func NewSSLStriper() *SSLStriper {
	mod := &SSLStriper{}

	return mod
}

func (m *SSLStriper) Name() string {
	return "SSL/TLS striper"
}

func (m *SSLStriper) Slug() string {
	return "ssl-striper"
}

func (m *SSLStriper) Description() string {
	return "Get information from X.509 certificat of a hostname (default port 443)."
}

func (m *SSLStriper) ResourceURLs() []common.ModuleResource {
	return nil
}

func (m *SSLStriper) Run(ctx *common.ModuleContext) error {
	host, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	// Connect to the hostname on port 443 to get the certificat
	cfg := tls.Config{}
	conn, err := tls.Dial("tcp", host+":443", &cfg)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("SSL/TLS connection failed")
		logrus.Warn(err)
		return err
	}

	// Grab the last certificate in the chain
	certChain := conn.ConnectionState().PeerCertificates
	cert := certChain[0]

	for _, hostname := range cert.DNSNames {
		hostname = common.CleanHostname(hostname)
		err = ctx.CreateNewAssetAsHostname(hostname)
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	for _, ip := range cert.IPAddresses {
		err = ctx.CreateNewAssetAsIP(ip.String())
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	for _, uri := range cert.URIs {
		err = ctx.CreateNewAssetAsURL(uri.String())
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	return nil
}
