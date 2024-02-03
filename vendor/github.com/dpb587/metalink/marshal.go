package metalink

import (
	"encoding/json"
	"encoding/xml"

	"github.com/pkg/errors"
)

func Unmarshal(data []byte, meta4 *Metalink) error {
	if len(data) > 0 && data[0] == '{' {
		return UnmarshalJSON(data, meta4)
	}

	return UnmarshalXML(data, meta4)
}

func UnmarshalXML(data []byte, meta4 *Metalink) error {
	err := xml.Unmarshal(data, meta4)
	if err != nil {
		return errors.Wrap(err, "Unmarshaling XML")
	}

	return nil
}

func UnmarshalJSON(data []byte, meta4 *Metalink) error {
	err := json.Unmarshal(data, meta4)
	if err != nil {
		return errors.Wrap(err, "Unmarshaling JSON")
	}

	return nil
}

func MarshalXML(r Metalink) ([]byte, error) {
	data, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "Marshaling XML")
	}

	return append(data, '\n'), nil
}

func MarshalJSON(r Metalink) ([]byte, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "Marshaling JSON")
	}

	return append(data, '\n'), nil
}
