package utils

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// StandardizeSpaces -> Join long query to one line query
func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// GenerateSlug  -> Replace space with dash, lower it, and trim the space
// DEPRECATED never use this no more
func GenerateSlug(inputStr string) string {
	re := regexp.MustCompile("[$#<|>{}~%`\\[\\]'^]")
	inputStr = re.ReplaceAllString(inputStr, "")
	return strings.Trim(strings.ToLower(strings.Replace(inputStr, " ", "-", -1)), " ")
}

// UnescapeString UTF-8 string
// e.g. convert "\u0e27\u0e23\u0e0d\u0e32" to "วรญา"
func UnescapeString(str string) (ustr string) {
	json.Unmarshal([]byte(`"`+str+`"`), &ustr)
	return
}

// String2Bool :nodoc:
func String2Bool(s string) bool {
	if s != "" {
		i, err := strconv.ParseBool(s)
		if err == nil {
			return i
		}
	}
	return false
}

// String2Int64 :nodoc:
func String2Int64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// String2Int64WithDefault :nodoc:
func String2Int64WithDefault(s string, d int64) int64 {
	i := String2Int64(s)
	if i == 0 {
		return d
	}
	return i
}

// GetIDFromSlug :nodoc:
func GetIDFromSlug(s string) int64 {
	ss := strings.Split(s, "-")
	id, err := strconv.ParseInt(ss[len(ss)-1], 10, 64)
	if err != nil {
		return 0
	}

	return id
}

// StringPointer2String :nodoc:
func StringPointer2String(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// StringPointer2Float64 :nodoc:
func StringPointer2Float64(s *string) float64 {
	if s != nil {
		f, err := strconv.ParseFloat(*s, 64)
		if err != nil {
			return float64(0)
		}
		return f
	}
	return float64(0)
}

// StringPointer2Int64 :nodoc:
func StringPointer2Int64(s *string) int64 {
	if s == nil {
		return int64(0)
	}

	return String2Int64(*s)
}

// ArrayStringPointer2ArrayInt64 :nodoc:
func ArrayStringPointer2ArrayInt64(s *[]*string) []int64 {
	var i []int64
	if s != nil {
		for _, val := range *s {
			i = append(i, StringPointer2Int64(val))
		}
		return i
	}
	return nil
}
