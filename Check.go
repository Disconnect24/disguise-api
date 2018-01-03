package disguise

import (
	"net/http"
)

// Check handles adding the proper interval for check.cgi along with future
// challenge solving and future mail existence checking.
// BUG(spotlightishere): Challenge solving isn't implemented whatsoever.
func Check(w http.ResponseWriter, r *http.Request, inter int) {

}