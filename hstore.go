package enthstore

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
)

// ErrTypeMustBeObject is the error returned by UnmarshalGQL when
// the value provided cannot be decoded.
var ErrTypeMustBeObject = errors.New("type Hstore must be a object")

// Hstore represents the hstore type of Postgres.
type Hstore map[string]*string

// FromMap creates a new Hstore from a map.
func FromMap(m map[string]string) Hstore {
	hs := Hstore{}

	for k := range m {
		v := m[k]
		hs[k] = &v
	}

	return hs
}

// Has check if the key exists.
func (h Hstore) Has(key string) bool {
	_, ok := h[key]
	return ok
}

// Set defines a value for the provided key.
func (h Hstore) Set(key string, val *string) {
	h[key] = val
}

// SetString defines a value for the provided key.
func (h Hstore) SetString(key string, val string) {
	h[key] = &val
}

// Get return the value from the provided key.
func (h Hstore) Get(key string) *string {
	if h == nil {
		return nil
	}

	val, found := h[key]
	if !found {
		return nil
	}

	return val
}

// GetString return the value from the provided key
// or an empty string if the key is not found
// or the value is nil.
func (h Hstore) GetString(key string) string {
	if h == nil {
		return ""
	}

	val, found := h[key]
	if !found {
		return ""
	}

	if val == nil {
		return ""
	}

	return *val
}

// Del deletes a key=>value pair.
func (h Hstore) Del(key string) {
	delete(h, key)
}

// String returns the string representation of the Hstore.
func (h Hstore) String() string {
	data, err := h.Value()
	if err != nil {
		return err.Error()
	}

	switch v := data.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	}

	return fmt.Sprint(data)
}

// Scan implements the interface Scanner.
func (h *Hstore) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if *h == nil {
		*h = Hstore{}
	}

	var input string

	switch v := value.(type) {
	case string:
		input = v
	case []byte:
		input = string(v)
	default:
		return fmt.Errorf("invalid input type: %T", v)
	}

	input = strings.TrimSpace(input)

	if len(input) == 0 {
		return nil
	}

	record := [][]byte{{}, {}}
	cur := 0
	isEscaping := false
	insideQuote := false
	lastQuoted := false
	ignore := []byte{' ', '\t', '\n', '\r', '='}

	savePair := func() {
		if len(record) == 0 {
			return
		}

		key := string(record[0])
		value := string(record[1])

		if !lastQuoted && strings.ToUpper(value) == "NULL" {
			h.Set(key, nil)
		} else {
			h.SetString(key, value)
		}

		record[0] = []byte{}
		record[1] = []byte{}
		cur = 0
	}

	for _, c := range []byte(input) {
		if isEscaping {
			isEscaping = false
			record[cur] = append(record[cur], c)
			continue
		}

		if c == '\\' {
			isEscaping = true
			continue
		}

		if c == '"' {
			insideQuote = !insideQuote
			if !insideQuote {
				lastQuoted = true
			}
			continue
		}

		if insideQuote {
			record[cur] = append(record[cur], c)
			continue
		}

		if bytes.IndexByte(ignore, c) != -1 {
			continue
		}

		if c == '>' {
			cur++
			lastQuoted = false
			continue
		}

		if c == ',' {
			savePair()
			continue
		}

		record[cur] = append(record[cur], c)
	}

	savePair()
	return nil
}

// Value implements the interface driver.Valuer.
func (h Hstore) Value() (driver.Value, error) {
	if h == nil {
		return nil, nil
	}

	parts := make([]string, 0, len(h))
	for key, val := range h {
		var part string
		if val == nil {
			part = quoteValue(key) + "=>NULL"
		} else {
			part = quoteValue(key) + "=>" + quoteValue(*val)
		}

		parts = append(parts, part)
	}
	return strings.Join(parts, ","), nil
}

// FormatParam defines how format the placeholder.
func (h *Hstore) FormatParam(param string, info *sql.StmtInfo) string {
	return param + "::hstore"
}

// UnmarshalGQL implements the interface graphql.Unmarshaler.
func (h *Hstore) UnmarshalGQL(v interface{}) error {
	val, ok := v.(map[string]interface{})
	if ok {
		hs := Hstore{}
		for key, val := range val {
			if val == nil {
				hs[key] = nil
				continue
			}
			v := fmt.Sprint(val)
			hs[key] = &v
		}
		*h = hs
		return nil
	}

	return ErrTypeMustBeObject
}

// MarshalGQL implements the interface graphql.Marshaler.
func (h Hstore) MarshalGQL(w io.Writer) {
	_ = json.NewEncoder(w).Encode(h)
}

// Equals check if two Hstore are equals.
func (h Hstore) Equals(other Hstore) bool {
	if len(h) != len(other) {
		return false
	}

	for k1, v1 := range h {
		v2, found := other[k1]
		if !found {
			return false
		}

		if (v1 == nil && v2 != nil) || (v2 == nil && v1 != nil) {
			return false
		}

		if v1 == nil && v2 == nil {
			continue
		}

		if *v1 != *v2 {
			return false
		}
	}

	return true
}

// SchemaType defines the schema-type of the Hstore object.
func (Hstore) SchemaType() map[string]string {
	return map[string]string{
		dialect.Postgres: "hstore",
	}
}
