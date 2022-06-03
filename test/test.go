package test

import (
	"path"
	"path/filepath"
	"runtime"

	"github.com/dhaval070/imagekit-go/config"
)

var Cfg *config.Configuration

const (
	PrivateKey  = "private_XxZH+I8BfOoIsY0M9CQtS4nyNSk="
	PublicKey   = "public_fGfgv45RjwmkbzGMRy1gKTcHf4M="
	UrlEndpoint = "https://ik.imagekit.io/dk1m7xkgi/"
)

func init() {
	Cfg = config.NewFromParams(PrivateKey, PublicKey, UrlEndpoint)
}

// LogoURL is the URL of the publicly available logo.
const LogoURL = "https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png"

// Base64Image us a base64 encoded test image.
const Base64Image = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

// ImageFilePath is a full path to the test image file.
var ImageFilePath = TestDataDir() + "new-york-cityscape-buildings.jpg"

// TestDataDir returns the full path to the directory with test files.
func TestDataDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	return filepath.Dir(d) + "/test/testdata/"
}
