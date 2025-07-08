package importer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/PlakarKorp/integration-caldav/oauth2utils"
	"github.com/PlakarKorp/kloset/objects"
	"github.com/PlakarKorp/kloset/snapshot/importer"
	"github.com/studio-b12/gowebdav"
	"golang.org/x/oauth2"
)

type CaldavImporter struct { //TODO: add more fields as needed
	ctx  context.Context
	opts *importer.Options

	client *gowebdav.Client
	url    string // The URL of the CalDAV server
}

func NewCaldavImporter(appCtx context.Context, opts *importer.Options, name string, config map[string]string) (importer.Importer, error) {

	//TODO: Parse configuration options from `config` map
	url := "https://apidata.googleusercontent.com/caldav/v2/EMAIL@gmail.com/events/"
	username := ""
	password := ""
	isOAuthClient := true

	var client *gowebdav.Client
	if !isOAuthClient {
		client = gowebdav.NewClient(url, username, password)
	} else { // OAuth2 client setup //TODO: make the service used (e.g., Google Calendar) configurable
		googleCal := oauth2utils.OAuthProvider{
			Name: "google",
			Config: &oauth2.Config{
				ClientID:     "ID",     // client ID (provided by the plakar app (production) or by the user directly in a personal app)
				ClientSecret: "SECRET", // client secret (same as above)
				RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
				Scopes:       []string{"SERVICE_SCOPE"}, // e.g., "https://www.googleapis.com/auth/calendar"
				Endpoint:     oauth2.Endpoint{},         // e.g., google.Endpoint for Google Calendar
			},
		}
		client = googleCal.GetClient(url) // maybe not using the url directly... the url could be built from the username
	}

	return &CaldavImporter{
		ctx:  appCtx,
		opts: opts,

		client: client,
		url:    url,
	}, nil
}

func (c *CaldavImporter) Origin() string {
	return c.url
}

func (c *CaldavImporter) Type() string {
	return "caldav"
}

func (c *CaldavImporter) Root() string {
	return "/"
}

func (c *CaldavImporter) Scan() (<-chan *importer.ScanResult, error) {

	results := make(chan *importer.ScanResult, 1000)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		entries, err := c.client.ReadDir("/")
		if err != nil {
			results <- importer.NewScanError("/", fmt.Errorf("error reading directory: %w", err))
			return
		}
		results <- importer.NewScanRecord("/", "", objects.FileInfo{
			Lname:    "/",
			Lsize:    0,
			Lmode:    os.ModeDir | 0755,
			LmodTime: entries[0].ModTime(),
		}, nil, nil)
		if len(entries) == 0 {
			results <- importer.NewScanError("/", fmt.Errorf("no entries found in the root directory"))
			return
		}

		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".ics") {
				data, err := c.client.Read(entry.Name())
				if err != nil {
					results <- importer.NewScanError("/"+entry.Name(), fmt.Errorf("error reading file %s: %w", entry.Name(), err))
					continue
				}

				rd := bytes.NewReader(data)

				results <- importer.NewScanRecord("/"+entry.Name(), "", objects.FileInfo{
					Lname:    entry.Name(),
					Lsize:    entry.Size(),
					Lmode:    entry.Mode(),
					LmodTime: entry.ModTime(),
				}, nil, func() (io.ReadCloser, error) {
					return io.NopCloser(rd), nil
				})
			}
		}
	}()

	go func() {
		wg.Wait()
		defer close(results)
	}()

	return results, nil
}

func (c *CaldavImporter) Close() error {
	return nil
}
