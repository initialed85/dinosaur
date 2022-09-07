package sessions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type SupportedLanguage struct {
	Name         string `json:"name,omitempty"`
	FriendlyName string `json:"friendly_name,omitempty"`
	FolderPath   string `json:"folder_path,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	BuildCmd     string `json:"build_cmd,omitempty"`
	RunCmd       string `json:"run_cmd,omitempty"`
	Code         string `json:"code,omitempty"`
}

var supportedLanguageByName = make(map[string]SupportedLanguage)

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	listPath := path.Join(cwd, "pkg/sessions/languages")

	files, err := ioutil.ReadDir(listPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, fileInfo := range files {
		if !fileInfo.IsDir() {
			continue
		}

		manifestFilePath := path.Join(listPath, fileInfo.Name(), "manifest.json")

		manifestFileInfo, err := os.Stat(manifestFilePath)
		if err != nil {
			continue
		}

		if manifestFileInfo.IsDir() {
			log.Fatalf("expected %v to be a JSON file but it was a folder", manifestFilePath)
		}

		manifestFileJSON, err := ioutil.ReadFile(manifestFilePath)
		if err != nil {
			log.Fatal(err)
		}

		supportedLanguage := SupportedLanguage{
			Name: fileInfo.Name(),
		}

		err = json.Unmarshal(manifestFileJSON, &supportedLanguage)
		if err != nil {
			log.Fatal(err)
		}

		codePath := path.Join(listPath, fileInfo.Name(), supportedLanguage.FolderPath, supportedLanguage.FileName)

		code, err := ioutil.ReadFile(codePath)
		if err != nil {
			log.Fatal(err)
		}

		supportedLanguage.Code = string(code)

		supportedLanguageByName[supportedLanguage.Name] = supportedLanguage
	}
}
