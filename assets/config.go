// +build dev

package assets

import "net/http"

// Config contains project configuration.
var Config http.FileSystem = http.Dir("config")
