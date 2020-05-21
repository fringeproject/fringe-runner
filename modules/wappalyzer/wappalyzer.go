package wappalyzer

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"

	"github.com/fringeproject/fringe-runner/common"
)

type Wappalyzer struct {
	Apps map[string]WappalyzerApp
	Tags []string
}

type WappalyzerFile struct {
	Apps map[string]WappalyzerApp `json:"apps"`
}

type WappalyzerApp struct {
	Cookies map[string]string `json:"cookies"`
	// This field can be a string or an array of string
	HTML    interface{}       `json:"html"`
	Headers map[string]string `json:"headers"`
	Meta    map[string]string `json:"meta"`
	// This field can be a string or an array of string
	Implies interface{} `json:"implies"`
}

func NewWappalyzer() *Wappalyzer {
	mod := &Wappalyzer{}

	return mod
}

func (m *Wappalyzer) Name() string {
	return "Wappalyzer"
}

func (m *Wappalyzer) Slug() string {
	return "wappalyzer"
}

func (m *Wappalyzer) Description() string {
	return "Analyzer HTTP response using Wappalyzer. Ref: https://www.wappalyzer.com/"
}

func (w *Wappalyzer) Parse(body *[]byte, headers *http.Header) error {
	w.Tags = make([]string, 0)

	err := w.ParseMetas(body)
	if err != nil {
		logrus.Warn("Error parsing metadata")
	}

	err = w.ParseCookies(headers)
	if err != nil {
		logrus.Warn("Error parsing cookies")
	}

	err = w.ParseHeaders(headers)
	if err != nil {
		logrus.Warn("Error parsing headers")
	}

	err = w.ParseHTML(body)
	if err != nil {
		logrus.Warn("Error parsing html")
	}

	return nil
}

func (w *Wappalyzer) addTag(tag string) {
	w.Tags = common.AppendIfMissing(w.Tags, tag)

	app := w.Apps[tag]
	implies := app.Implies
	v := reflect.ValueOf(implies)

	switch v.Kind() {
	case reflect.String:
		w.addTag(implies.(string))

	case reflect.Slice:
		tagsSlice, ok := v.Interface().([]interface{})
		if !ok {
			logrus.Warn("Wappalyzer: cannot convert to slice: ", v.Interface())
		} else {
			for _, tagString := range tagsSlice {
				w.addTag(tagString.(string))
			}
		}

	case reflect.Invalid:
		// nothing to do

	default:
		logrus.Warn("Wappalyzer: unknow implies type: ", v)
	}
}

func (w *Wappalyzer) ParseHeaders(headers *http.Header) error {
	for respHeaderName, respHeaderValue := range *headers {
		respHeaderName = strings.ToLower(respHeaderName)

		for appName, appValue := range w.Apps {
			for appHeaderName, appHeaderValue := range appValue.Headers {
				if strings.ToLower(appHeaderName) == respHeaderName {

					headerRegexp, err := regexp.Compile(strings.Split(appHeaderValue, "\\;")[0])
					if err != nil {
						logrus.Warnf("Cannot compile header regexp: %s", strings.Split(appHeaderValue, "\\;")[0])
					} else {
						if headerRegexp.MatchString(respHeaderValue[0]) {
							if strings.Contains(appHeaderValue, "version") {
								headerMatches := headerRegexp.FindStringSubmatch(respHeaderValue[0])
								if len(headerMatches[1]) > 0 {
									w.addTag(appName + "/" + headerMatches[1])
									w.Tags = common.AppendIfMissing(w.Tags, appName+"/"+headerMatches[1])
								} else {
									w.addTag(appName)
								}
							} else {
								w.addTag(appName)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (w *Wappalyzer) ParseMetas(body *[]byte) error {
	z := html.NewTokenizer(bytes.NewReader(*body))

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return nil
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "meta" {
				for appName, appValue := range w.Apps {
					for appMetaName, appMetaValue := range appValue.Meta {
						content, ok := extractMetaProperty(t, appMetaName)

						if ok {
							headerRegexp, err := regexp.Compile(strings.Split(appMetaValue, "\\;")[0])
							if err != nil {
								logrus.Warn("Cannot compile meta regexp: ", strings.Split(appMetaValue, "\\;")[0])
							} else {
								if headerRegexp.MatchString(content) {
									if strings.Contains(appMetaValue, "version") {
										headerMatches := headerRegexp.FindStringSubmatch(content)
										w.addTag(appName + "/" + headerMatches[1])
									} else {
										w.addTag(appName)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func extractMetaProperty(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}

func (w *Wappalyzer) ParseCookies(headers *http.Header) error {
	cookieCount := len((*headers)["Set-Cookie"])
	if cookieCount == 0 {
		return nil
	}

	for _, line := range (*headers)["Set-Cookie"] {
		for appName, appValue := range w.Apps {
			for cookieName := range appValue.Cookies {
				if cookieName+"=" == line {
					w.addTag(appName)
				}
			}
		}
	}

	return nil
}

func (w *Wappalyzer) ParseHTML(body *[]byte) error {
	html := string(*body)

	for appName, appValue := range w.Apps {
		v := reflect.ValueOf(appValue.HTML)

		switch v.Kind() {
		case reflect.String:
			htmlValue := appValue.HTML.(string)
			htmlRegex, err := regexp.Compile(strings.Split(htmlValue, "\\;")[0])
			if err != nil {
				logrus.Debug("Cannot compile html regexp: ", strings.Split(htmlValue, "\\;")[0])
			} else {
				if htmlRegex.MatchString(html) {
					w.addTag(appName)
				}
			}

		case reflect.Slice:
			htmlValues := appValue.HTML.([]interface{})
			for _, htmlValue := range htmlValues {
				htmlValueString := htmlValue.(string)
				htmlRegex, err := regexp.Compile(strings.Split(htmlValueString, "\\;")[0])
				if err != nil {
					logrus.Debug("Cannot compile html regexp: ", strings.Split(htmlValueString, "\\;")[0])
				} else {
					if htmlRegex.MatchString(html) {
						w.addTag(appName)
					}
				}
			}

		case reflect.Invalid:
			// nothing to do

		default:
			logrus.Warnf("Wappalyzer: unknow implies type: %s", v)
		}
	}

	return nil
}

func (m *Wappalyzer) Run(ctx *common.ModuleContext) error {
	url, err := ctx.GetAssetAsURL()
	if err != nil {
		return err
	}

	wappalyzer := &WappalyzerFile{}
	err = common.ReadJSONFile("./lists/wappalyzer.json", wappalyzer)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot read wappalyzer JSON file. Please, check \"lists/wappalyzer.json\" file.")
		logrus.Warn(err)
		return err
	}
	m.Apps = wappalyzer.Apps

	_, body, headers, err := ctx.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Error fetching URL")
		logrus.Warn(err)
		return err
	}

	err = m.Parse(body, headers)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Could not parse the HTTP response with Wappalyzer.")
		logrus.Warn(err)
		return err
	}

	for _, tag := range m.Tags {
		err = ctx.AddTag(tag)
		if err != nil {
			logrus.Debug(err)
			logrus.Warn("Could not create tag.")
		}
	}

	return nil
}
