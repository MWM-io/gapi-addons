package spanner

import (
	"strconv"
	"time"

	"cloud.google.com/go/spanner"
)

// Time overrides the spanner.Time to produce a unix milli timestamp when marshalled to json
type Time time.Time

// MarshalJSON implements the json.Marshaler interface to produce a unix milli timestamp
func (nt Time) MarshalJSON() ([]byte, error) {
	casted := time.Time(nt)

	return []byte(strconv.FormatInt(casted.UnixMilli(), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle a unix milli timestamp and convert it to a Time
func (nt *Time) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err == nil {
		*nt = Time(time.UnixMilli(i))
		return nil
	}

	// if it's not a unix timestamp try decoding as time this is useful for json object created by spanner TO_JSON()
	if b[0] == '"' {
		b = b[1 : len(b)-1]
	}

	date, errT := time.Parse(time.RFC3339Nano, string(b))
	if errT != nil {
		return errT
	}

	*nt = Time(date)
	return nil
}

// CommitTimestamp returns a spanner.CommitTimestamp cast to a models.Time
func CommitTimestamp() Time {
	return Time(spanner.CommitTimestamp)
}

// NullTime overrides the spanner.NullTime to produce a unix milli timestamp when marshalled to json
type NullTime spanner.NullTime

// MarshalJSON implements the json.Marshaler interface to produce a unix milli timestamp
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return []byte(strconv.FormatInt(nt.Time.UnixMilli(), 10)), nil
	}

	return []byte("0"), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle a unix milli timestamp and convert it to a NullTime
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "0" {
		nt.Valid = false
		return nil
	}

	i, err := strconv.ParseInt(string(b), 10, 64)
	if err == nil {
		*nt = NullTime(spanner.NullTime{Valid: true, Time: time.UnixMilli(i)})
		return nil
	}

	// if it's not a unix timestamp try decoding as time this is useful for json object created by spanner TO_JSON()
	if b[0] == '"' {
		b = b[1 : len(b)-1]
	}
	date, errT := time.Parse(time.RFC3339Nano, string(b))
	if errT != nil {
		return errT
	}

	*nt = NullTime{
		Time:  date,
		Valid: true,
	}
	return nil
}

// NullCommitTimestamp returns a spanner.CommitTimestamp casted to a models.NullTime
func NullCommitTimestamp() NullTime {
	return NullTime{
		Time:  spanner.CommitTimestamp,
		Valid: true,
	}
}
