package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	flag.String("har", "", "*.har file")
	flag.String("dir", "content", "savedir")
	flag.Parse()
}

func main() {
	var downloadDir = flag.Lookup("dir").Value.String()
	os.Mkdir(downloadDir, os.ModePerm)
	parseHAR(downloadDir)
}

func parseHAR(downloadDir string) {
	var harFile = flag.Lookup("har").Value.String()
	file, err := os.Open(harFile)
	if err != nil {
		fmt.Fprint(os.Stderr, err, `: `, `"`+harFile+`"`)
		os.Exit(1)
	}
	var har HARlog
	var decoder = json.NewDecoder(file)
	if err := decoder.Decode(&har); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	file.Close()
	for _, e := range har.Log.Entries {
		if e.Response.Content.Size == 0 {
			continue
		}
		var URL *url.URL
		URL, err = url.Parse(e.Request.URL)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		var filename = URL.Path[strings.LastIndex(URL.Path, `/`)+1:]
		var b []byte
		var err error
		switch e.Response.Content.Encoding {
		case "base64":
			b, err = base64.StdEncoding.DecodeString(e.Response.Content.Text)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		default:
			b = []byte(e.Response.Content.Text)
		}
		if err := ioutil.WriteFile(filepath.Join(downloadDir, filename), b, os.ModePerm); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// HARlog - har файл запросов
type HARlog struct {
	Log struct {
		Version string `json:"version"`
		Creator struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"creator"`
		Pages   []interface{} `json:"pages"`
		Entries []struct {
			StartedDateTime time.Time `json:"startedDateTime"`
			Time            float64   `json:"time"`
			Request         struct {
				Method      string `json:"method"`
				URL         string `json:"url"`
				HTTPVersion string `json:"httpVersion"`
				Headers     []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"headers"`
				QueryString []interface{} `json:"queryString"`
				Cookies     []struct {
					Name     string      `json:"name"`
					Value    string      `json:"value"`
					Expires  interface{} `json:"expires"`
					HTTPOnly bool        `json:"httpOnly"`
					Secure   bool        `json:"secure"`
				} `json:"cookies"`
				HeadersSize float64 `json:"headersSize"`
				BodySize    float64 `json:"bodySize"`
			} `json:"request"`
			Response struct {
				Status      int    `json:"status"`
				StatusText  string `json:"statusText"`
				HTTPVersion string `json:"httpVersion"`
				Headers     []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"headers"`
				Cookies []interface{} `json:"cookies"`
				Content struct {
					Size     int    `json:"size"`
					MimeType string `json:"mimeType"`
					Text     string `json:"text"`
					Encoding string `json:"encoding"`
				} `json:"content"`
				RedirectURL  string `json:"redirectURL"`
				HeadersSize  int    `json:"headersSize"`
				BodySize     int    `json:"bodySize"`
				TransferSize int    `json:"_transferSize"`
			} `json:"response"`
			Cache struct {
			} `json:"cache"`
			Timings struct {
				Blocked         float64 `json:"blocked"`
				DNS             float64 `json:"dns"`
				Ssl             float64 `json:"ssl"`
				Connect         float64 `json:"connect"`
				Send            float64 `json:"send"`
				Wait            float64 `json:"wait"`
				Receive         float64 `json:"receive"`
				BlockedQueueing float64 `json:"_blocked_queueing"`
			} `json:"timings"`
			ServerIPAddress string `json:"serverIPAddress"`
			Initiator       struct {
				Type string `json:"type"`
			} `json:"_initiator"`
			Priority     string `json:"_priority"`
			ResourceType string `json:"_resourceType"`
			Connection   string `json:"connection"`
		} `json:"entries"`
	} `json:"log"`
}
