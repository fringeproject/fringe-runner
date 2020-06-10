package modules

import (
	"github.com/fringeproject/fringe-runner/modules/alienvault"
	"github.com/fringeproject/fringe-runner/modules/awss3"
	"github.com/fringeproject/fringe-runner/modules/backup"
	"github.com/fringeproject/fringe-runner/modules/bufferover"
	"github.com/fringeproject/fringe-runner/modules/censys"
	"github.com/fringeproject/fringe-runner/modules/certspotter"
	"github.com/fringeproject/fringe-runner/modules/crtsh"
	"github.com/fringeproject/fringe-runner/modules/dnsdumpster"
	"github.com/fringeproject/fringe-runner/modules/docker"
	"github.com/fringeproject/fringe-runner/modules/elasticsearch"
	"github.com/fringeproject/fringe-runner/modules/hackertarget"
	"github.com/fringeproject/fringe-runner/modules/ipinfo"
	"github.com/fringeproject/fringe-runner/modules/kafka"
	"github.com/fringeproject/fringe-runner/modules/kubernetes"
	"github.com/fringeproject/fringe-runner/modules/nmap"
	"github.com/fringeproject/fringe-runner/modules/securitytrails"
	"github.com/fringeproject/fringe-runner/modules/shodan"
	"github.com/fringeproject/fringe-runner/modules/sslstriper"
	"github.com/fringeproject/fringe-runner/modules/sublist3r"
	"github.com/fringeproject/fringe-runner/modules/takeover"
	"github.com/fringeproject/fringe-runner/modules/threatcrowd"
	"github.com/fringeproject/fringe-runner/modules/threatminer"
	"github.com/fringeproject/fringe-runner/modules/urlscan"
	"github.com/fringeproject/fringe-runner/modules/virustotal"
	"github.com/fringeproject/fringe-runner/modules/wappalyzer"
	"github.com/fringeproject/fringe-runner/modules/wayback"
	"github.com/fringeproject/fringe-runner/modules/yougetsignal"

	"github.com/fringeproject/fringe-runner/session"
)

func LoadModules(sess *session.Session) {
	sess.RegisterModule(alienvault.NewAlienVault())
	sess.RegisterModule(awss3.NewAWSS3())
	sess.RegisterModule(backup.NewBackup())
	sess.RegisterModule(bufferover.NewBufferOver())
	sess.RegisterModule(censys.NewCensys())
	sess.RegisterModule(certspotter.NewCertSpotter())
	sess.RegisterModule(crtsh.NewCrtsh())
	sess.RegisterModule(dnsdumpster.NewDnsdumpster())
	sess.RegisterModule(docker.NewDockerAPI())
	sess.RegisterModule(elasticsearch.NewElasticSearchAPI())
	sess.RegisterModule(hackertarget.NewHackerTarget())
	sess.RegisterModule(ipinfo.NewIPInfo())
	sess.RegisterModule(kafka.NewKafkaAPI())
	sess.RegisterModule(kubernetes.NewKubernetesAPI())
	sess.RegisterModule(nmap.NewNmap())
	sess.RegisterModule(securitytrails.NewSecurityTrails())
	sess.RegisterModule(shodan.NewShodan())
	sess.RegisterModule(sslstriper.NewSSLStriper())
	sess.RegisterModule(sublist3r.NewSublist3r())
	sess.RegisterModule(takeover.NewTakeOver())
	sess.RegisterModule(threatcrowd.NewThreatCrowd())
	sess.RegisterModule(threatminer.NewThreatMiner())
	sess.RegisterModule(urlscan.NewURLScan())
	sess.RegisterModule(virustotal.NewVirusTotal())
	sess.RegisterModule(wappalyzer.NewWappalyzer())
	sess.RegisterModule(wayback.NewWayback())
	sess.RegisterModule(yougetsignal.NewYouGetSignal())
}
