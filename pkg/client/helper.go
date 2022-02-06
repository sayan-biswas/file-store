package client

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func checkFile(fileName string) (bool, error) {
	storeURL, err := url.Parse(storeConfig.URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, err
	}
	storeURL.Path = "store/check/file"
	values := storeURL.Query()
	values.Add("file", fileName)
	storeURL.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodGet, storeURL.String(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, err
	}
	if res.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func checkSHA(SHA []byte) (bool, error) {
	storeURL, err := url.Parse(storeConfig.URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, err
	}
	storeURL.Path = "store/check/sha"
	values := storeURL.Query()
	values.Add("sha", hex.EncodeToString(SHA))
	storeURL.RawQuery = values.Encode()
	req, err := http.NewRequest(http.MethodGet, storeURL.String(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, err
	}
	if res.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}
