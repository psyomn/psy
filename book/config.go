package book

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	grSecret = ""
	grKey    = ""
)

func baseConfigPath() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".config", "goodreads")
}

func goodReadsSecretFilePath() string {
	return filepath.Join(baseConfigPath(), "secret")
}

func goodReadsKeyFilePath() string {
	return filepath.Join(baseConfigPath(), "key")
}

func missingInfo() error {
	msg := `missing files with information:
  %s
  %s
You can find the key information here:
  https://www.goodreads.com/api/keys`

	return fmt.Errorf(msg,
		goodReadsSecretFilePath(),
		goodReadsKeyFilePath())
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ret, err := ioutil.ReadAll(f)
	return ret, err
}

func loadConfig() error {
	grkfp := goodReadsKeyFilePath()
	grsfp := goodReadsSecretFilePath()
	if !fileExists(grkfp) || !fileExists(grsfp) {
		return missingInfo()
	}

	secret, err := readFile(grkfp)
	if err != nil {
		return err
	}

	key, err := readFile(grsfp)
	if err != nil {
		return err
	}

	grKey = string(secret)
	grSecret = string(key)

	return nil
}
