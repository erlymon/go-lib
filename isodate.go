package golib

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/glog"

	"gopkg.in/mgo.v2/bson"
)

type ISODate struct {
	time.Time
}

// ISO8601 format to millis instead of to nanos
const ISO8601Millis = "2006-01-02T15:04:05.000Z0700"

func Now() ISODate {
	return ISODate{time.Now()}
}

func Unix(sec int64, nsec int64) ISODate {
	return ISODate{time.Unix(sec, nsec)}
}

// ParseTimestamp parses a string that represents an ISO8601 time or a unix epoch
func ParseTimestamp(data string) (ISODate, error) {
	d := time.Now().UTC()
	if data != "now" {
		// fmt.Println("we should try to parse")
		dd, err := time.Parse(ISO8601Millis, data)
		fmt.Println("ParseTimestamp", dd)
		if err != nil {
			dd, err = time.Parse(time.RFC3339, data)
			if err != nil {
				dd, err = time.Parse(time.RFC3339Nano, data)
				if err != nil {
					if data == "" {
						data = "0"
					}
					t, err := strconv.ParseInt(data, 10, 64)
					if err != nil {
						return ISODate{}, err
					}
					dd = time.Unix(0, t*int64(time.Millisecond))
				}
			}
		}
		d = dd
	}
	return ISODate{Time: d.UTC()}, nil
}

func (t ISODate) String() string {
	return t.UTC().Format(ISO8601Millis)
}

// UnmarshalJSON implements the json unmarshaller interface
func (t *ISODate) UnmarshalJSON(data []byte) error {
	var value interface{}
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	switch value.(type) {
	case string:
		v := value.(string)
		fmt.Println("UnmarshalJSON TIMESTAMP", v)
		if v == "" {
			return nil
		}
		d, err := ParseTimestamp(v)
		fmt.Println("UnmarshalJSON ParseTimestamp", d, d.UTC())
		if err != nil {
			return err
		}
		*t = d
	case float64:
		*t = ISODate{time.Unix(0, int64(value.(float64))*int64(time.Millisecond)).UTC()}
	default:
		return fmt.Errorf("Couldn't convert json from (%T) %s to a time.Time", value, data)
	}
	return nil
}

func (t *ISODate) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("ISODate.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(ISO8601Millis)+2)
	b = append(b, '"')
	b = t.UTC().AppendFormat(b, ISO8601Millis)
	b = append(b, '"')
	return b, nil
}

// UnmarshalText reads this timestamp from a string value
func (t *ISODate) UnmarshalText(data []byte) error {
	fmt.Println("UnmarshalText", data)
	var value interface{}
	json.Unmarshal(data, &value)

	switch value.(type) {
	case string:
		v := value.(string)
		if v == "" {
			return nil
		}
		d, err := ParseTimestamp(v)
		fmt.Println("UnmarshalText ParseTimestamp", d, d.UTC())
		if err != nil {
			return err
		}
		*t = d
	case float64:
		*t = ISODate{time.Unix(0, int64(value.(float64))*int64(time.Millisecond)).UTC()}
	default:
		return fmt.Errorf("couldn't convert json from (%T) %s to a time.Time", value, data)
	}
	return nil
}

// MarshalText implements the text marshaller interface
func (t ISODate) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// SetBSON customizes the bson serialization for this type
func (t *ISODate) SetBSON(raw bson.Raw) error {
	var ts interface{}
	if err := raw.Unmarshal(&ts); err != nil {
		return err
	}
	glog.Errorln("SetBSON", ts)
	switch ts.(type) {
	case time.Time:
		*t = ISODate{Time: ts.(time.Time).UTC()}
		return nil
	case string:
		tss := ts.(string)
		tt, err := ParseTimestamp(tss)
		fmt.Println("SetBSON ParseTimestamp", tt, tt.UTC())
		if err != nil {
			return err
		}
		*t = tt
		return nil
	case int64:
		*t = ISODate{time.Unix(0, ts.(int64)*int64(time.Millisecond)).UTC()}
		return nil
	case float64:
		*t = ISODate{time.Unix(0, int64(ts.(float64))*int64(time.Millisecond)).UTC()}
		return nil
	}

	return fmt.Errorf("couldn't convert bson data (%T) %s to a Timestamp", ts, ts)
}

// GetBSON customizes the bson serialization for this type
func (t ISODate) GetBSON() (interface{}, error) {
	return t.UTC(), nil
}
