package backup

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fringeproject/fringe-runner/common"
)

type Backup struct {
}

func NewBackup() *Backup {
	mod := &Backup{}

	return mod
}

func (m *Backup) Name() string {
	return "Backup"
}

func (m *Backup) Slug() string {
	return "backup"
}

func (m *Backup) Description() string {
	return "Based on the URL, tries to requests various filenames with backup variations such as .old, .tar, ..."
}

func GenerateBackupFileVariation(filename string) []string {
	// First we need the get the base filename and its extension
	extensionIndex := strings.LastIndex(filename, ".")
	basename := filename
	if extensionIndex > -1 {
		basename = filename[:extensionIndex]
	}

	// List of extensions a devoloper could use to change the name of a file
	// TODO: Add date based on the HTTP header
	extensions := []string{
		// swap vim, nano...
		"txt", "save", "swp", "swo",
		// backup or temp files
		"bk", "bak", "backup", "bkup", "bkp", "temp", "tmp", "old",
		"orig", "original",
		"new", "source",
		// endind
		"~", "1", "2",
		// archive extensions
		"zip", "rar", "tar", "tar.gz", "tar.xz",
	}

	// This list contains the variations based on several rules and some customs
	// examples
	files := []string{
		"#" + filename + "#",
	}

	// Here is some rules to add our custom extensions to the file or basename
	for _, extension := range extensions {
		files = append(files, filename+""+extension)
		files = append(files, filename+"."+extension)
		files = append(files, filename+"_"+extension)
		if extensionIndex > -1 {
			files = append(files, basename+"."+extension)
			files = append(files, basename+"_"+extension)
		}
	}

	return files
}

func ParseURLPath(path string) (string, string) {
	// Here's some examples of how Golang parse a URL, especially the `Path` field.
	// Its quite obvious but it worth to double check before writing some code
	// https://fringeproject.com           -> ''
	// https://fringeproject.com/          -> '/'
	// https://fringeproject.com/uploads   -> '/uploads'
	// https://fringeproject.com/uploads/  -> '/uploads/'
	// https://fringeproject.com/index.php -> '/index.php'

	// We want here the file name, that's mean after the last '/' (if it exists)
	lastSlashIndex := strings.LastIndex(path, "/")

	if lastSlashIndex > -1 {
		return path[:lastSlashIndex+1], path[lastSlashIndex+1:]
	} else {
		return path, ""
	}
}

func (m *Backup) Run(ctx *common.ModuleContext) error {
	rawurl, err := ctx.GetAssetAsURL()
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		logrus.Debug(err)
		err = fmt.Errorf("Cannot parse URL.")
		logrus.Warn(err)
		return err
	}

	path, lastPart := ParseURLPath(parsedURL.Path)
	logrus.Info(parsedURL.Path, " -> ", path, " | ", lastPart)
	// Check if we've a lastPart to perform variation on
	// TODO: Do the test on the last directory name
	if len(lastPart) == 0 {
		logrus.Debug("The last parth of the URL path is empty, nothing to do here.")
		return nil
	}

	// Here the `lastPart` can be a filename (index.php) or a folder (uploads).
	// Then we still use the same variation on the folder because the developper
	// may zip or create a backup of this folder.
	lastPartVariations := GenerateBackupFileVariation(lastPart)

	// Finally, we recreate all the URL from the variations
	// TODO: this is part of a HTTP brute-forcer: add it to a separate module
	statusCodeWhiteList := map[int]bool{
		100: true,
		200: true, 202: true, 204: true,
		301: true, 302: true, 307: true,
		401: true, 403: true,
	}
	for _, part := range lastPartVariations {
		// Change the path so we still use the same others fields from the raw URL
		parsedURL.Path = path + part

		// Test to "GET" the URL and check the status code
		statusCode, _, _, err := ctx.HttpRequest(http.MethodGet, parsedURL.String(), nil, nil)
		if err != nil {
			logrus.Debug(err)
			continue
		}
		logrus.Infof("[%d] %s", *statusCode, parsedURL.String())

		if _, found := statusCodeWhiteList[*statusCode]; found {
			err = ctx.CreateNewAssetAsURL(parsedURL.String())
			if err != nil {
				logrus.Debug(err)
				logrus.Warn("Could not create URL.")
			}
		}
	}

	return nil
}
