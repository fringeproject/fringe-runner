package nmap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/common/assets"
)

type Nmap struct {
}

func NewNmap() *Nmap {
	return &Nmap{}
}

func (m *Nmap) Name() string {
	return "nmap"
}

func (m *Nmap) Slug() string {
	return "nmap"
}

func (m *Nmap) Description() string {
	return "Run the network mapper nmap."
}

func RunNmapScan(nmapPath string, nmapArgs []string) (*NmapExecution, error) {
	var (
		out, errs bytes.Buffer
	)

	cmd := exec.Command(nmapPath, nmapArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &errs
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	if errs.Len() > 0 {
		// If there is something in stderr then we check if it end with "QUTTING!"
		errorMsg := errs.String()

		if strings.HasSuffix(errorMsg, "QUITTING!\r\n") {
			return nil, fmt.Errorf("There was somethign wrong running nmap:\n%s", errorMsg)
		} else {
			// It's a warning, log it and continue the parsing
			logrus.Warn(errorMsg)
		}
	}

	// Parse the result
	result := &NmapExecution{}
	err = xml.Unmarshal(out.Bytes(), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Nmap) Run(ctx *common.ModuleContext) error {
	target, err := ctx.GetAssetAsIP()
	if err != nil {
		return err
	}

	if !common.IsValidAddr(target) {
		err = fmt.Errorf("nmap cannot scan local address %s", target)
		return err
	}

	// Get the nmap path from the configuration. If it's not set, then use the
	// nmap command from the PATH
	nmapPath, err := ctx.GetConfigurationValue("NMAP_PATH")
	if err != nil {
		nmapPath = "nmap"
	}

	args := []string{
		// "--top-ports", "20", // Use the "Top 20 most scanned ports"
		// Use the top 20 ports and add some customs
		// 2375: docker api
		// 8083: kafka
		// 9200, 9300: elasticsearch
		// 10250: Kubernetes API
		"-p", "21,22,23,25,53,80,110,111,135,139,143,443,445,993,995,1723,2375,3306,3389,5900,8080,8083,9200,9300,10250",
		"--unprivileged", // Assume the user lacks raw socket privileges
		"-sV",            // Probe open ports to determine service/version info
		"-n",             //  Never do DNS resolution
		// Disable OS detection as we may need a privileged user
		// "-O",              // Enable OS detection
		"-T4",             // Set timing template (higher is faster)
		"--open",          //  Only show open (or possibly open) ports
		"--min-rate=1500", // Send packets no slower than 1500 per second
		"-oX", "-",        // Use the XML in stdout output for parsing
		target, // Set the scan target
	}

	// Run nmap with the previous arguments
	result, err := RunNmapScan(nmapPath, args)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Nmap scan failed with error: %s", err)
		logrus.Warn(err)
		return err
	}

	for _, host := range result.Hosts {
		if host.Status.State == "up" {
			for _, port := range host.Ports {
				// Check if the port is open (should be already done with --open arg)
				if port.State.State == "open" {
					portMsg := fmt.Sprintf("Port %d is open with service (%s)", port.PortId, port.Service.Name)

					err = ctx.CreateNewAsset(portMsg, assets.AssetTypes["raw"])
					if err != nil {
						logrus.Debug(err)
						logrus.Warn("Could not create vulnerability.")
					}
				}
			}
		}
	}

	return nil
}
