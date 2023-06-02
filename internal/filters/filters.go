package filters

import (
	"math/rand"
	"time"

	"github.com/flosch/pongo2"
)

func init() {

	rand.Seed(time.Now().UTC().UnixNano())
	pongo2.RegisterFilter("namelocationiduri", filterNameLocationIdURI)
}

func filterNameLocationIdURI(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {

	return pongo2.AsValue("something"), nil
}
