package handler

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

	"H4sIAAAAAAAC/7RSTY/TMBD9K9bAMWoK7Cl3JPYAQrCcVj2Y5DVx5djemUmlqsp/R+O20JU4wsmO52Xe",
	"x8yZ+jyXnJBUqDuT9BNmX68fmTN/g5ScBPZQOBewBtTyDBE/1oKeCqgjUQ5ppHVtiPGyBMZA3fNv4K65",
	"AfPPA3qltaFPiDH/Vw5DhrTP1iOGHlee5GdDfX58MhkaNNrnDwG77+Bj6EENHcEScqKO3m22m60hc0Hy",
	"JVBHH+pTQ8XrVMW2k5mx2wi1w5x4DTk9DtRdrFY8+xkKFuqezxSs/csCPlFzUxUGurenvKC5TuYuipAU",
	"I5jWdWfoS4ZVyfvt1o4+J0WqUnwpMfRVTHsQs3S+a/iWsaeO3rR/dqG9LkL7ekI1zgHScyh6ieYJoo6h",
	"CycL6GH78M+4X2/gX7i/ZHX7vKShLoQs8+z5ZJqmIC6IOyyizjs1iUhDySGp0+xGqDvlxYl6Vgwb9zXC",
	"C9yACIVT+/2G31yIBXy8zWzhSB1NqqVr25h7H6csSutu/RUAAP//x+wl0k8DAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
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
	res := make(map[string]func() ([]byte, error))
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
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
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
