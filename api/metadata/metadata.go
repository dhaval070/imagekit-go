package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	neturl "net/url"

	"github.com/dhaval070/imagekit-go/api"
	"github.com/dhaval070/imagekit-go/config"
	"github.com/dhaval070/imagekit-go/logger"
)

// API is the main struct for media
type API struct {
	Config config.Configuration
	Logger *logger.Logger
	Client api.HttpClient
}

// New creates a new Media API instance from the environment variable.
func New() (*API, error) {
	c, err := config.New()
	if err != nil {
		return nil, err
	}

	return NewFromConfiguration(c)
}

// NewFromConfiguration a new Media API instance with the given Configuration.
func NewFromConfiguration(c *config.Configuration) (*API, error) {
	return &API{
		Config: *c,
		Client: &http.Client{},
		Logger: logger.New(),
	}, nil
}

// MetadataResponse represents main struct of metadata response of the sdk
type MetadataResponse struct {
	Data Metadata
	api.Response
}

// Metadata represents struct of metadata response from api
type Metadata struct {
	Height          int
	Width           int
	Size            int64
	Format          string
	HasColorProfile bool
	Quality         int
	Density         int
	HasTransparency bool
	PHash           string
	Exif            Mexif
}

type Mexif struct {
	Image            ImageExif
	Thumbnail        ThumbnailExif
	Exif             Exif
	Gps              Gps
	Interoperability Interoperability
	Makernote        map[string]interface{}
}

type ImageExif struct {
	Make             string
	Model            string
	Orientation      string
	XResolution      int
	YResolution      int
	ResolutionUnit   int
	Software         string
	ModifyDate       time.Time
	YCbCrPositioning int
	ExifOffset       int
	GPSInfo          int
}

type ThumbnailExif struct {
	Compression     int
	XResolution     int
	YResolution     int
	ResolutionUnit  int
	ThumbnailOffset int
	ThumbnailLength int
}

type Exif struct {
	ExposureTime             time.Time
	FNumber                  float32
	ExposureProgram          int
	ISO                      int
	ExifVersion              string
	DateTimeOriginal         time.Time
	CreateDate               time.Time
	ShutterSpeedValue        float32
	ApertureValue            float32
	ExposureCompensation     int
	MeteringMode             int
	Flash                    int
	FocalLength              int
	SubSEcTime               string
	SubSecTimeOriginal       string
	FlashpixVersion          string
	ColorSpace               int
	ExifImageWidth           int
	ExifImageHeight          int
	InteropOffset            int
	FocalPlaneXResolution    float32
	FocalPlaneYResolution    float32
	FocalPlaneResolutionUnit int
	CustomRendered           int
	ExposureMode             int
	WhiteBalance             int
	SceneCaptutureType       int
}

type Gps struct {
	GPSVersionID []int
}

type Interoperability struct {
	InteropIndex   string
	InteropVersion string
}

type Schema struct {
	Type            string        `json:"type"`
	SelectOptions   []interface{} `json:"selectOptions"`
	DefaultValue    string        `json:"defaultValue"`
	IsValueRequired bool          `json:"isValueRequired"`
	MinValue        string        `json:"minValue"`
	MaxValue        string        `json:"maxValue"`
	MinLength       int           `json:"minLength"`
	MaxLength       int           `json:"maxLength"`
}

func (m *API) get(ctx context.Context, url string, query map[string]string) (*http.Response, error) {
	var err error
	urlObj, err := neturl.Parse(api.BuildPath(m.Config.API.Prefix, url))
	if err != nil {
		return nil, err
	}

	values := urlObj.Query()
	for k, v := range query {
		values.Set(k, v)
	}

	sUrl := urlObj.String() + "?" + values.Encode()
	req, err := http.NewRequest(http.MethodGet, sUrl, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(m.Config.Cloud.PrivateKey, "")

	return m.Client.Do(req.WithContext(ctx))
}

// FromAsset fetches metadata of media library file
func (m *API) FromAsset(ctx context.Context, fileId string) (*MetadataResponse, error) {
	if fileId == "" {
		return nil, errors.New("fileId can not be blank")
	}

	var response = &MetadataResponse{}

	resp, err := m.get(ctx, fmt.Sprintf("files/%s/metadata", fileId), nil)
	defer api.DeferredBodyClose(resp)

	api.SetResponseMeta(resp, response)

	if err != nil {
		return response, err
	}

	if resp.StatusCode != 200 {
		err = response.ParseError()
	} else {
		err = json.Unmarshal(response.Body(), &response.Data)
	}

	return response, err
}

func (m *API) FromUrl(ctx context.Context, url string) (*MetadataResponse, error) {
	var err error
	if url == "" {
		return nil, errors.New("url can not be blank")
	}

	var response = &MetadataResponse{}

	if err != nil {
		return nil, err
	}

	resp, err := m.get(ctx, "metadata", map[string]string{"url": url})
	defer api.DeferredBodyClose(resp)

	api.SetResponseMeta(resp, response)

	if err != nil {
		return response, err
	}

	if resp.StatusCode != 200 {
		err = response.ParseError()
	} else {
		err = json.Unmarshal(response.Body(), &response.Data)
	}

	return response, err
}
