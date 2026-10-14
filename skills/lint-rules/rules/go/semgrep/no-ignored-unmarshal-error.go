package example

import "encoding/json"

func parse(b []byte) (v map[string]any, err error) {
	// ruleid: no-ignored-unmarshal-error
	_ = json.Unmarshal(b, &v)
	// ruleid: no-ignored-unmarshal-error
	_ = json.NewDecoder(nil).Decode(&v)
	// ok: no-ignored-unmarshal-error
	err = json.Unmarshal(b, &v)
	return v, err
}
