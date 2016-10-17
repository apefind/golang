package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/apefind/golang/shellutil"
)

var AppConfig = `
Graphics:
    - Gimp
Multimedia:
    - VLC
    - Vox
Network:
    - FileZilla
    - Firefox
    - Thunderbird
    - TorBrowser
Office:
    - LibreOffice
Programming:
    - Atom
    - PyCharm CE
    - Racket v6.5
    - TextMate
Tools:
    - FoundImages
    - iTerm
`

func usage() {
	fmt.Fprintf(os.Stderr, "\n%s: rearrange application folder\n\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func getConfig(config string) (string, error) {
	if config == "" {
		return AppConfig, nil
	}
	data, err := ioutil.ReadFile(config)
	return string(data), err
}

func getAppCategories(applications string) (map[string]string, error) {
	config, appCategories := make(map[string][]string), make(map[string]string)
	err := yaml.Unmarshal([]byte(applications), &config)
	if err != nil {
		return appCategories, err
	}
	for category, apps := range config {
		for _, app := range apps {
			appCategories[app] = category
		}
	}
	return appCategories, nil
}

func getAppFilenames(appDir string) ([]string, error) {
	return filepath.Glob(appDir + string(filepath.Separator) + "*.app")
}

func getAppAndCategory(filename string, appCategories map[string]string) (string, string) {
	app := shellutil.GetFileBasename(filename)
	category, _ := appCategories[app]
	return app, category
}

func moveApp(path, appDir, app, category string) error {
	dir := appDir + string(filepath.Separator) + category
	if !shellutil.IsDir(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	appPath := dir + string(filepath.Separator) + app + ".app"
	if shellutil.IsSymlink(appPath) {
		if err := os.Remove(appPath); err != nil {
			return err
		}
	} else if shellutil.IsDir(appPath) {
		if err := os.RemoveAll(appPath); err != nil {
			return err
		}
	}
	return os.Rename(path, appPath)
}

func SyncAppLocations(appDir, config string) error {
	config, err := getConfig(config)
	if err != nil {
		return err
	}
	appCategories, err := getAppCategories(config)
	if err != nil {
		return err
	}
	apps, err := getAppFilenames(appDir)
	if err != nil {
		return err
	}
	for _, path := range apps {
		app, category := getAppAndCategory(path, appCategories)
		if category != "" {
			log.Printf("%s -> %s\n", app, category)
			if err := moveApp(path, appDir, app, category); err != nil {
				return err
			}
		} else {
			log.Printf("%s -> no category found\n", app)
		}
	}
	return nil
}

func main() {
	log.SetFlags(0)
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	var appDir, config string
	flag.Usage = usage
	flag.StringVar(&appDir, "appdir", user.HomeDir+"/Applications", "Application directory")
	flag.StringVar(&config, "config", "", "YAML configuration")
	flag.Parse()
	if err := SyncAppLocations(appDir, config); err != nil {
		log.Fatal(err)
	}
}
