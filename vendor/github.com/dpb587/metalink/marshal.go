package metalink

import (
	"encoding/xml"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func Unmarshal(data []byte, meta4 *Metalink) error {
	err := xml.Unmarshal(data, meta4)
	if err != nil {
		return bosherr.WrapError(err, "Unmarshaling XML")
	}

	return nil
}

func Marshal(r Metalink) ([]byte, error) {
	data, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, bosherr.WrapError(err, "Marshaling XML")
	}

	return append(data, '\n'), nil
}
