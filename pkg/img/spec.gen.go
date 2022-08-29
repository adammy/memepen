// Package img provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package img

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RWS2/cNhD+KwTbo7xSHgUCHRM36bZFDCQtcih8oMVZaQzxEXK0jmHsfy+GlNar1doN",
	"3DZoTxa58/rm+2boO9k4450FS1HWdzI2HRiVPt8EUARro1p464LhKw2xCegJnZX1sYHQsEELUShxpSII",
	"5B8EOTH43im9koX0wXkIhJASJINl2BRQYBTUnQyyccEokrW8QqvCrSwk3XqQtYwU0LZytytkgM8DBtCy",
	"/mNMc7k3c1fX0JDcFYcAfv548f5RhGywR3jTQTioao54E5xZgh1Cv4z/+4dfJ6Cboe8Fn8ml833AGeYh",
	"4AJwIb+cte7MKgM55qIFnPxUA9aPMTCBNUBKK1Ji48IM6RJkk/qlL+xDrdQXdh+XUWpFUBIaOMIsblQU",
	"Y7QZfnY4Y4cl7YXsANuOlrl/SvdTp7OVcJvjPotPnNeSoA6j8AEa1TdDzzUIZbWI5AJoJkhtHWrROEto",
	"B9X3t2IyRdsmdy7aqC9oBiPr51VVFdKgzcdn+9LRErQQuHbUJ3g455qVGCx+HkCgBku4QQiJCq6ewHgu",
	"cK6RAfVfiGR9Pl0smvhthVrIG9TULTN+4uspZ7L5bzF2vGS449y64mAIJnB7ZS5nkOOg3bhlA37rQORB",
	"/Ahhiw0IH9wWNUTBI6fYLCYhtECMohi3Ufpk9Bp6SPDuO5ZQElIP+zEfo8tCbiHEnPvZqlpVTI7zYJVH",
	"WcsXq2r1ggdeUZdmvcwB+dO7SI8uz3GUWcsWbjIsrmQPZK3nDjJ3FyK9dvo2rRZnCWxKo7zvsUl+5XV0",
	"9v7d4q/vA2xkLb8r7x+2cnzVyuOFn7rPp/LaQzuP8zXPzOT7ZFdvn+Bqhp7Qq0AlW5/xbn5SC9Krvpsr",
	"mcIA6SJ6Z2Om93lV/WMMZHZT0rlaLn5hvb3MqeY/vVZafMhqYJsfTtmsLUGwqk9yhiB+DMHlKY2DMdzJ",
	"SV9iEhipNvLo5vMl246aLu/S3/X5LifqgU68kufpfnor+Tsvq1G4oB/S+YFjGqigDBAELuZOIofmIZOF",
	"HBflWIw85qk46PmJBzy/HlzR+nzanSdqO9bX5YL9l8sE7514M8ohsXbShMRbN1j9tzjLvXqQs0K2cGL1",
	"vAPKvASggLD9WmYmv/8FLd9wKP81et8BPTyPbJlcMwXpHxPZEfm6LHvXqL5zkepX1auq3D6Tu+LQJNZl",
	"iaZdKY8rAwY82FXjTDK83P0ZAAD//ziNi032DAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}