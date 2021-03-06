package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	iktest "github.com/dhaval070/imagekit-go/test"
	"github.com/google/go-cmp/cmp"
)

var ctx = context.Background()
var metadataApi *API

func TestMain(m *testing.M) {
	var err error
	metadataApi, err = NewFromConfiguration(iktest.Cfg)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
func getHandler(statusCode int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, body)
	}
}

func TestMetadata_FromFile(t *testing.T) {
	var respBody = `{"height":801,"width":597,"size":59718,"format":"jpg","hasColorProfile":true,"quality":0,"density":72,"hasTransparency":false,"exif":{},"pHash":"85d07f1fe4ae8be2"}`

	var err error
	var respObj = &Metadata{}

	if err = json.Unmarshal([]byte(respBody), respObj); err != nil {
		t.Error(err)
	}

	handler := getHandler(200, respBody)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	metadataApi.Config.API.Prefix = ts.URL + "/"

	resp, err := metadataApi.FromAsset(ctx, "3325344545345")

	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(resp.Data, *respObj) {
		t.Errorf("\n%v\n%v\n", resp.Data, *respObj)
	}
}

func TestMetadata_FromUrl(t *testing.T) {
	var respBody = `{"height":801,"width":597,"size":59718,"format":"jpg","hasColorProfile":true,"quality":0,"density":72,"hasTransparency":false,"exif":{},"pHash":"85d07f1fe4ae8be2"}`

	var err error
	var respObj = &Metadata{}

	if err = json.Unmarshal([]byte(respBody), respObj); err != nil {
		t.Error(err)
	}

	handler := getHandler(200, respBody)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	metadataApi.Config.API.Prefix = ts.URL + "/"

	resp, err := metadataApi.FromUrl(ctx, "https://ik.imagekit.io/xk1m7xkgi/default-image.jpg")

	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(resp.Data, *respObj) {
		t.Errorf("\n%v\n%v\n", resp.Data, *respObj)
	}
}

func TestMetadata_CreateCustomField(t *testing.T) {
	var respBody = `{"id":"62a8966b663ef736f841fe28","name":"speed","label":"Speed","schema":{"type":"Number","defaultValue":100,"minValue":1,"maxValue":120}}`

	var err error
	var expected = &CustomField{}

	if err = json.Unmarshal([]byte(respBody), expected); err != nil {
		t.Error(err)
	}

	handler := getHandler(201, respBody)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	metadataApi.Config.API.Prefix = ts.URL + "/"

	resp, err := metadataApi.CreateCustomField(ctx, CreateFieldParam{
		Name:  "speed",
		Label: "Speed",
		Schema: Schema{
			Type:         "Number",
			DefaultValue: 100,
			MinValue:     1,
			MaxValue:     120,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(resp.Data, *expected) {
		t.Errorf("\n%v\n%v\n", resp.Data, *expected)
	}
}

func TestMetadata_CustomFields(t *testing.T) {
	var respBody = `[{"id":"629f6b437eb0fe6f1b66d864","name":"price","label":"Price","schema":{"type":"Number","isValueRequired":false,"minValue":1,"maxValue":1000}},{"id":"629f6b6d7eb0fe344f66e1b6","name":"country","label":"Country","schema":{"type":"SingleSelect","isValueRequired":false,"selectOptions":["USA","Canada"]}},{"id":"62a8764d663ef721e93f4ea9","name":"clearance","label":"Clearance","schema":{"type":"MultiSelect","selectOptions":["one","two"]}},{"id":"62a876b1663ef7728f3f5348","name":"mileage","label":"Mileage","schema":{"type":"Number"}},{"id":"62a8966b663ef736f841fe28","name":"speed","label":"Speed","schema":{"type":"Number","defaultValue":100,"minValue":1,"maxValue":120}}]`

	var err error
	var expected = []CustomField{}

	if err = json.Unmarshal([]byte(respBody), &expected); err != nil {
		t.Error(err)
	}

	handler := getHandler(200, respBody)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	metadataApi.Config.API.Prefix = ts.URL + "/"
	resp, err := metadataApi.CustomFields(ctx, false)

	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(resp.Data, expected) {
		t.Errorf("%v\n%v", resp.Data, expected)
	}
}

func TestMetadata_UpdateCustomField(t *testing.T) {
	var respBody = `{"id":"629f6b437eb0fe6f1b66d864","name":"price","label":"Cost","schema":{"type":"Number"}}`

	var err error
	var expected = CustomField{}

	if err = json.Unmarshal([]byte(respBody), &expected); err != nil {
		t.Error(err)
	}

	handler := getHandler(200, respBody)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	metadataApi.Config.API.Prefix = ts.URL + "/"
	resp, err := metadataApi.UpdateCustomField(ctx, UpdateCustomFieldParam{
		FieldId: "629f6b437eb0fe6f1b66d864",
		Label:   "Cost",
	})

	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(resp.Data, expected) {
		t.Errorf("%v\n%v", resp.Data, expected)
	}
}

func TestMetadata_DeleteCustomField(t *testing.T) {
	var respBody = ``
	var err error

	handler := getHandler(204, respBody)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	metadataApi.Config.API.Prefix = ts.URL + "/"
	_, err = metadataApi.DeleteCustomField(ctx, "62a8966b663ef736f841fe28")
	if err != nil {
		log.Println("got error")
		t.Error(err)
	}
}
