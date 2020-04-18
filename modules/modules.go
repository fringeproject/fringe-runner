package modules

import (
	"github.com/fringeproject/fringe-runner/modules/awss3"
	"github.com/fringeproject/fringe-runner/modules/censys"
	"github.com/fringeproject/fringe-runner/modules/certspotter"
	"github.com/fringeproject/fringe-runner/modules/crtsh"
	"github.com/fringeproject/fringe-runner/modules/docker"
	"github.com/fringeproject/fringe-runner/modules/elasticsearch"
	"github.com/fringeproject/fringe-runner/modules/helloworld"
	"github.com/fringeproject/fringe-runner/modules/ipinfo"
	"github.com/fringeproject/fringe-runner/modules/kafka"
	"github.com/fringeproject/fringe-runner/modules/kubernetes"
	"github.com/fringeproject/fringe-runner/modules/nmap"
	"github.com/fringeproject/fringe-runner/modules/securitytrails"
	"github.com/fringeproject/fringe-runner/modules/shodan"
	"github.com/fringeproject/fringe-runner/modules/sslstriper"
	"github.com/fringeproject/fringe-runner/modules/threatcrowd"
	"github.com/fringeproject/fringe-runner/modules/wayback"

	"github.com/fringeproject/fringe-runner/session"
)

func LoadModules(sess *session.Session) {
	sess.RegisterModule(awss3.NewAWSS3())
	sess.RegisterModule(censys.NewCensys())
	sess.RegisterModule(certspotter.NewCertSpotter())
	sess.RegisterModule(crtsh.NewCrtsh())
	sess.RegisterModule(docker.NewDockerAPI())
	sess.RegisterModule(elasticsearch.NewElasticSearchAPI())
	sess.RegisterModule(helloworld.NewHelloWorld())
	sess.RegisterModule(ipinfo.NewIPInfo())
	sess.RegisterModule(kafka.NewKafkaAPI())
	sess.RegisterModule(kubernetes.NewKubernetesAPI())
	sess.RegisterModule(nmap.NewNmap())
	sess.RegisterModule(securitytrails.NewSecurityTrails())
	sess.RegisterModule(shodan.NewShodan())
	sess.RegisterModule(sslstriper.NewSSLStriper())
	sess.RegisterModule(threatcrowd.NewThreatCrowd())
	sess.RegisterModule(wayback.NewWayback())
}
