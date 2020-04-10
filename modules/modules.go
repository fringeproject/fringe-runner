package modules

import (
	"github.com/fringeproject/fringe-runner/modules/awss3"
	"github.com/fringeproject/fringe-runner/modules/crtsh"
	"github.com/fringeproject/fringe-runner/modules/docker"
	"github.com/fringeproject/fringe-runner/modules/elasticsearch"
	"github.com/fringeproject/fringe-runner/modules/helloworld"
	"github.com/fringeproject/fringe-runner/modules/ipinfo"
	"github.com/fringeproject/fringe-runner/modules/kafka"
	"github.com/fringeproject/fringe-runner/modules/kubernetes"
	"github.com/fringeproject/fringe-runner/modules/nmap"
	"github.com/fringeproject/fringe-runner/modules/wayback"

	"github.com/fringeproject/fringe-runner/session"
)

func LoadModules(sess *session.Session) {
	sess.RegisterModule(awss3.NewAWSS3())
	sess.RegisterModule(crtsh.NewCrtsh())
	sess.RegisterModule(docker.NewDockerAPI())
	sess.RegisterModule(elasticsearch.NewElasticSearchAPI())
	sess.RegisterModule(helloworld.NewHelloWorld())
	sess.RegisterModule(ipinfo.NewIPInfo())
	sess.RegisterModule(kafka.NewKafkaAPI())
	sess.RegisterModule(kubernetes.NewKubernetesAPI())
	sess.RegisterModule(nmap.NewNmap())
	sess.RegisterModule(wayback.NewWayback())
}
