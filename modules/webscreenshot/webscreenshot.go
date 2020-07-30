package webscreenshot

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type WebScreenshot struct {
}

func NewWebScreenshot() *WebScreenshot {
	mod := &WebScreenshot{}

	return mod
}

func (m *WebScreenshot) Name() string {
	return "Web Screenshot"
}

func (m *WebScreenshot) Slug() string {
	return "webscreenshot"
}

func (m *WebScreenshot) Description() string {
	return "Take a screenshot of a website."
}

func (m *WebScreenshot) ResourceURLs() []common.ModuleResource {
	return nil
}

func takeScreenshotCmd(renderer string, args []string) error {
	logrus.Debugf("Execute renderer: %s %s", renderer, args)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, renderer, args...)
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("The renderer timed out.")
		}

		return err
	}

	return nil
}

func getFilenameFromURL(url string) string {
	filename := strings.ReplaceAll(url, "://", "_")
	re := regexp.MustCompile(`[^\w\-_\. ]`)
	filename = re.ReplaceAllString(filename, "_")

	return filename + ".png"
}

func (m *WebScreenshot) Run(ctx *common.ModuleContext) error {
	url, err := ctx.GetAssetAsURL()
	if err != nil {
		return err
	}

	renderer, err := ctx.GetConfigurationValue("webscreenshot_renderer")
	if err != nil {
		err := fmt.Errorf("You must provide a render using \"webscreenshot_renderer\".")
		return err
	}

	rendererPath, err := ctx.GetConfigurationValue("webscreenshot_renderer_path")
	if err != nil {
		err := fmt.Errorf("You must provide the render path using \"webscreenshot_renderer_path\".")
		return err
	}

	outputDirectory, err := ctx.GetConfigurationValue("webscreenshot_output")
	if err != nil {
		err := fmt.Errorf("You must provide a output directory for screenshots using \"webscreenshot_output\".")
		return err
	}

	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		err := fmt.Errorf("The webscreenshots output directory does not exist.")
		return err
	}

	screenshot := filepath.Join(outputDirectory, getFilenameFromURL(url))
	windowSize := "1200,800"
	renderer = strings.ToLower(renderer)
	args := []string{}

	if renderer == "chrome" || renderer == "chromium" {
		// https://developers.google.com/web/updates/2017/04/headless-chrome
		args = append(args, []string{
			"--headless",
			"--disable-gpu",
			"--hide-scrollbars",
			"--incognito",
			"--allow-running-insecure-content",
			"--ignore-certificate-errors",
			"--ignore-urlfetcher-cert-requests",
			"--reduce-security-for-testing",
			"--no-sandbox",
			"--disable-crash-reporter",
		}...)
	} else if renderer == "firefox" {
		// https://developer.mozilla.org/en-US/docs/Mozilla/Firefox/Headless_mode
		args = append(args, []string{
			// "--new-instance",
			// "--headless" // You can omit -headless when using --screenshot
		}...)
	} else {
		return fmt.Errorf("The renderer is invalid, please use \"chrome\", \"chromium\" or \"firefox\".")
	}

	args = append(args, []string{
		"--screenshot=" + screenshot,
		"--window-size=" + windowSize,
		url,
	}...)

	err = takeScreenshotCmd(rendererPath, args)
	if err != nil {
		return err
	}

	logrus.Infof("Screenshot written to: %s", screenshot)

	return nil
}
