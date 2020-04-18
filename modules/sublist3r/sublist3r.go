package sublist3r

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type Sublist3r struct {
}

func NewSublist3r() *Sublist3r {
	mod := &Sublist3r{}

	return mod
}

func (m *Sublist3r) Name() string {
	return "Sublist3r"
}

func (m *Sublist3r) Slug() string {
	return "sublist3r"
}

func (m *Sublist3r) Description() string {
	return "Requests Sublist3r API. Ref: https://github.com/aboul3la/Sublist3r"
}

func (m *Sublist3r) Run(ctx *common.ModuleContext) error {
	hostname, err := ctx.GetAssetAsHostname()
	if err != nil {
		return err
	}

	url := "https://api.sublist3r.com/search.php?domain=" + hostname
	var subResponse []string
	_, _, _, err = ctx.HTTPRequestJson(http.MethodGet, url, nil, &subResponse, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot request Sublist3r API.")
		logrus.Warn(err)
		return err
	}

	for _, host := range subResponse {
		err = ctx.CreateNewAssetAsHostname(host)
		if err != nil {
			logrus.Info("Something went wrong creating the new asset.")
			logrus.Debug(err)
		}
	}

	return nil
}
