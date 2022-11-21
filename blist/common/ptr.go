package common

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

func StringPtr(v string) *string {
	return &v
}

func StringUnptr(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func StringsPtr(v []string) *[]string {
	return &v
}

func StringsUnptr(v *[]string) []string {
	if v != nil {
		return *v
	}
	return make([]string, 0)
}

func IntPtr(v int) *int {
	return &v
}

func IntUnptr(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

func TimePtr(v time.Time) *time.Time {
	return &v
}

func TimeUnptr(v *time.Time) time.Time {
	if v != nil {
		return *v
	}
	return time.Time{}
}

func Uint64Ptr(v uint64) *uint64 {
	return &v
}

func Uint64Unptr(v *uint64) uint64 {
	if v != nil {
		return *v
	}
	return 0
}

func DurationPtr(v time.Duration) *time.Duration {
	return &v
}

func DurationUnptr(v *time.Duration) time.Duration {
	if v != nil {
		return *v
	}
	return 0
}

func ObjectIdPtr(v bson.ObjectId) *bson.ObjectId {
	return &v
}

func ObjectIdsPtr(v []bson.ObjectId) *[]bson.ObjectId {
	return &v
}

func BoolPtr(v bool) *bool {
	return &v
}

func BoolUnptr(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}
