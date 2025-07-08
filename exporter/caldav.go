package exporter

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"path"

	"github.com/PlakarKorp/integration-caldav/oauth2utils"
	"github.com/PlakarKorp/kloset/objects"
	"github.com/PlakarKorp/kloset/snapshot/exporter"
	"github.com/studio-b12/gowebdav"
)

type CaldavExporter struct { //TODO: add more fields as needed
	ctx  context.Context
	opts *exporter.Options

	client *gowebdav.Client
	url    string // The URL of the CalDAV server
}

func NewCaldavExporter(appCtx context.Context, opts *exporter.Options, name string, config map[string]string) (exporter.Exporter, error) {

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

	return &CaldavExporter{
		ctx:  appCtx,
		opts: opts,

		client: client,
		url:    url,
	}, nil
}

func (c *CaldavExporter) Root() string {
	return "/"
}

func (c *CaldavExporter) CreateDirectory(pathname string) error {
	return nil
}

func (c *CaldavExporter) StoreFile(pathname string, fp io.Reader, size int64) error {
	pathname = path.Base(pathname)

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", pathname, err)
	}

	//TODO: look at it returns an error, even if the file is written successfully
	if c.client.Write(pathname, data, 0644) != nil {
		return fmt.Errorf("error writing %s: %w", pathname, err)
	}
	return nil
}

func (c *CaldavExporter) SetPermissions(pathname string, fileinfo *objects.FileInfo) error {
	return nil
}

func (c *CaldavExporter) Close() error {
	return nil
}
