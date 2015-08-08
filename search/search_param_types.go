package search

import (
	"math/big"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"github.com/davecgh/go-spew/spew"
)

type Query struct {
	Resource string
	Query    string
}

func (q *Query) Params() []SearchParam {
	var results []SearchParam
	queryMap, _ := url.ParseQuery(q.Query)
	for param, values := range queryMap {
		info, ok := SearchParameterDictionary[q.Resource][param]
		if ok {
			for _, value := range values {
				results = append(results, info.CreateSearchParam(value))
			}
		}
	}
	return results
}

type SearchParam interface {
	getInfo() SearchParamInfo
}

/******************************************************************************
 * SearchParamInfo contains information about a FHIR search parameter,
 * including its name, type, and paths.  Paths are represented as a simple
 * string map that maps a search path to a data type, for example:
 * {
 *   "valueDateTime" : "dateTime" ,
 *   "valuePeriod"   : "Period"   ,
 * }
 *
 * This map is used to allow the FHIR query parameter to be translated into a
 * mongo query object.
 ******************************************************************************/
type SearchParamInfo struct {
	Name  string
	Type  string
	Paths map[string]string
}

func (s SearchParamInfo) CreateSearchParam(paramStr string) SearchParam {
	switch s.Type {
	case "composite":
		return nil
	case "date":
		return ParseDateParam(paramStr, s)
	case "number":
		return ParseNumberParam(paramStr, s)
	case "quantity":
		return ParseQuantityParam(paramStr, s)
	case "reference":
		return ParseReferenceParam(paramStr, s)
	case "string":
		return ParseStringParam(paramStr, s)
	case "token":
		return ParseTokenParam(paramStr, s)
	case "uri":
		return ParseURIParam(paramStr, s)
	}
	return nil
}

/******************************************************************************
 * COMPOSITE: a resource may also specify composite parameters that take
 * sequences of single values that match other defined parameters as an
 * argument. The matching parameter of each component in such a sequence is
 * documented in the definition of the parameter. These sequences are formed by
 * joining the single values with a "$". Note that this sequence is a single
 * value and itself can be composed into a set of values, so that, for example,
 * multiple matching state-on-date parameters can be specified as
 * state-on-date=new$2013-05-04,active$2013-05-05.
 *
 * http://hl7-fhir.github.io/search.html#composite
 ******************************************************************************/

// TODO: Implement composite parameters

/******************************************************************************
 * DATE: A date parameter searches on a date/time or period. As is usual for
 * date/time related functionality, while the concepts are relatively
 * straight-forward, there are a number of subtleties involved in ensuring
 * consistent behavior.
 *
 * http://hl7-fhir.github.io/search.html#date
 ******************************************************************************/

type DateParam struct {
	SearchParamInfo
	Prefix Prefix
	Date   *Date
}

func (d *DateParam) getInfo() SearchParamInfo {
	return d.SearchParamInfo
}

func ParseDateParam(paramStr string, info SearchParamInfo) *DateParam {
	date := &DateParam{SearchParamInfo: info}

	var value string
	date.Prefix, value = ExtractPrefixAndValue(paramStr)
	date.Date = ParseDate(value)

	return date
}

var DT_REGEX = regexp.MustCompile("([0-9]{4})(-(0[1-9]|1[0-2])(-(0[0-9]|[1-2][0-9]|3[0-1])(T([01][0-9]|2[0-3]):([0-5][0-9])(:([0-5][0-9])(\\.([0-9]+))?)?((Z)|(\\+|-)((0[0-9]|1[0-3]):([0-5][0-9])|(14):(00)))?)?)?)?")

type Date struct {
	Value     time.Time
	Precision DatePrecision
}

func (d *Date) String() string {
	return d.Value.Format(d.Precision.Layout())
}

func (d *Date) RangeLowIncl() time.Time {
	return d.Value
}

func (d *Date) RangeHighExcl() time.Time {
	switch d.Precision {
	case YEAR:
		return d.Value.AddDate(1, 0, 0)
	case MONTH:
		return d.Value.AddDate(0, 1, 0)
	case DAY:
		return d.Value.AddDate(0, 0, 1)
	case MINUTE:
		return d.Value.Add(time.Minute)
	case SECOND:
		return d.Value.Add(time.Second)
	case MILLISECOND:
		return d.Value.Add(time.Millisecond)
	default:
		return d.Value.Add(time.Millisecond)
	}
}

func ParseDate(dateStr string) *Date {
	dt := &Date{}

	dateStr = strings.TrimSpace(dateStr)
	if m := DT_REGEX.FindStringSubmatch(dateStr); m != nil {
		y, mo, d, h, mi, s, ms, tzZu, tzOp, tzh, tzm := m[1], m[3], m[5], m[7], m[8], m[10], m[12], m[14], m[15], m[17], m[18]

		switch {
		case ms != "":
			dt.Precision = MILLISECOND

			// Fix milliseconds (.9 -> .900, .99 -> .990, .999999 -> .999 )
			switch len(ms) {
			case 1:
				ms += "00"
			case 2:
				ms += "0"
			case 3:
				// do nothing
			default:
				ms = ms[:3]
			}
		case s != "":
			dt.Precision = SECOND
		case mi != "":
			dt.Precision = MINUTE
		// NOTE: Skip hour precision since FHIR specification disallows it
		case d != "":
			dt.Precision = DAY
		case mo != "":
			dt.Precision = MONTH
		case y != "":
			dt.Precision = YEAR
		default:
			dt.Precision = MILLISECOND
		}

		// Get the location (if no time components or no location, use local)
		loc := time.Local
		if h != "" {
			if tzZu == "Z" {
				loc, _ = time.LoadLocation("UTC")
			} else if tzOp != "" && tzh != "" && tzm != "" {
				tzhi, _ := strconv.Atoi(tzh)
				tzmi, _ := strconv.Atoi(tzm)
				offset := tzhi*60*60 + tzmi*60
				if tzOp == "-" {
					offset *= -1
				}
				loc = time.FixedZone(tzOp+tzh+tzm, offset)
			}
		}

		// Convert to a time.Time
		yInt, _ := strconv.Atoi(y)
		moInt, err := strconv.Atoi(mo)
		if err != nil {
			moInt = 1
		}
		dInt, err := strconv.Atoi(d)
		if err != nil {
			dInt = 1
		}
		hInt, _ := strconv.Atoi(h)
		miInt, _ := strconv.Atoi(mi)
		sInt, _ := strconv.Atoi(s)
		msInt, _ := strconv.Atoi(ms)

		dt.Value = time.Date(yInt, time.Month(moInt), dInt, hInt, miInt, sInt, msInt*1000*1000, loc)
	} else {
		// TODO: What should we do if the time format is wrong?  Right now, we default to NOW
		dt.Precision = MILLISECOND
		dt.Value = time.Now()
	}

	return dt
}

type DatePrecision int

const (
	YEAR DatePrecision = iota
	MONTH
	DAY
	MINUTE
	SECOND
	MILLISECOND
)

func (p DatePrecision) Layout() string {
	switch p {
	case YEAR:
		return "2006"
	case MONTH:
		return "2006-01"
	case DAY:
		return "2006-01-02"
	case MINUTE:
		return "2006-01-02T15:04-07:00"
	case SECOND:
		return "2006-01-02T15:04:05-07:00"
	case MILLISECOND:
		return "2006-01-02T15:04:05.000-07:00"
	default:
		return "2006-01-02T15:04:05.000-07:00"
	}
}

/******************************************************************************
 * NUMBER: Searching on a simple numerical value in a resource.
 *
 * http://hl7-fhir.github.io/search.html#number
 ******************************************************************************/

type NumberParam struct {
	SearchParamInfo
	Prefix Prefix
	Number *Number
}

func (n *NumberParam) getInfo() SearchParamInfo {
	return n.SearchParamInfo
}

func ParseNumberParam(paramStr string, info SearchParamInfo) *NumberParam {
	n := &NumberParam{SearchParamInfo: info}

	var value string
	n.Prefix, value = ExtractPrefixAndValue(paramStr)
	n.Number = ParseNumber(value)

	return n
}

type Number struct {
	Value     *big.Rat
	Precision int
}

func (n *Number) String() string {
	return n.Value.FloatString(n.Precision)
}

func (n *Number) RangeLowIncl() *big.Rat {
	return new(big.Rat).Sub(n.Value, n.rangeDelta())
}

func (n *Number) RangeHighExcl() *big.Rat {
	return new(big.Rat).Add(n.Value, n.rangeDelta())
}

/* FHIR spec defines equality for 100 to be the range [99.5, 100.5) so we must support min/max
 * using rounding semantics. The basic algorithm for determining low/high is:
 *   low  (inclusive) = n - 5 / 10^p
 *   high (exclusive) = n + 5 / 10^p
 * where n is the number and p is the count of the number's decimal places + 1.
 *
 * This function returns the delta ( 5 / 10^p )
 */
func (n *Number) rangeDelta() *big.Rat {
	p := n.Precision + 1
	denomInt := new(big.Int).Exp(big.NewInt(int64(10)), big.NewInt(int64(p)), nil)
	denomRat, _ := new(big.Rat).SetString(denomInt.String())
	return new(big.Rat).Quo(new(big.Rat).SetInt64(5), denomRat)
}

func ParseNumber(numStr string) *Number {
	n := &Number{}

	numStr = strings.TrimSpace(numStr)
	n.Value, _ = new(big.Rat).SetString(numStr)
	i := strings.Index(numStr, ".")
	if i != -1 {
		n.Precision = len(numStr) - i - 1
	} else {
		n.Precision = 0
	}

	return n
}

/******************************************************************************
 * QUANTITY: A quantity parameter searches on the Quantity data type.
 *
 * http://hl7-fhir.github.io/search.html#quantity
 ******************************************************************************/

type QuantityParam struct {
	SearchParamInfo
	Prefix Prefix
	Number *Number
	System string
	Code   string
}

func (q *QuantityParam) getInfo() SearchParamInfo {
	return q.SearchParamInfo
}

func ParseQuantityParam(paramStr string, info SearchParamInfo) *QuantityParam {
	q := &QuantityParam{SearchParamInfo: info}

	var value string
	q.Prefix, value = ExtractPrefixAndValue(paramStr)

	split := escapeFriendlySplit(value, '|')
	q.Number = ParseNumber(split[0])
	if len(split) == 3 {
		q.System = unescape(split[1])
		q.Code = unescape(split[2])
	}

	return q
}

/******************************************************************************
 * REFERENCE: A reference parameter refers to references between resources,
 * e.g. find all Conditions where the subject reference is a particular patient,
 * where the patient is selected by name or identifier.
 *
 * http://hl7-fhir.github.io/search.html#reference
 ******************************************************************************/

type ReferenceParam struct {
	SearchParamInfo
	Reference string
}

func (r *ReferenceParam) getInfo() SearchParamInfo {
	return r.SearchParamInfo
}

func (r *ReferenceParam) GetId() (id string, ok bool) {
	if r.Reference == "" {
		ok = false
	} else if !r.isUrl() {
		i := strings.LastIndex(r.Reference, "/")
		id = r.Reference[i+1:]
		ok = true
	} else {
		ok = false
	}
	return
}

func (r *ReferenceParam) GetType() (typ string, ok bool) {
	if !r.isUrl() {
		i := strings.LastIndex(r.Reference, "/")
		if i == -1 {
			ok = false
		} else {
			typ = r.Reference[:i]
			ok = true
		}
	} else {
		ok = false
	}
	return
}

func (r *ReferenceParam) GetUrl() (url string, ok bool) {
	if r.isUrl() {
		url = r.Reference
		ok = true
	} else {
		ok = false
	}
	return
}

func (r *ReferenceParam) isUrl() bool {
	u, e := url.Parse(r.Reference)
	if e == nil {
		return u.IsAbs()
	}
	return false
}

func (r *ReferenceParam) IsId() bool {
	return !r.isUrl()
}

func (r *ReferenceParam) IsUrl() bool {
	return r.isUrl()
}

func ParseReferenceParam(paramStr string, info SearchParamInfo) *ReferenceParam {
	//r := unescape(paramStr)
	return &ReferenceParam{info, unescape(paramStr)}
}

/******************************************************************************
 * STRING: The string parameter refers to simple string searches against
 * sequences of characters. Matches are case- and accent- insensitive. By
 * default, a field matches a string query if the value of the field equals or
 * starts with the supplied parameter value, after both have been normalized by
 * case and accent.
 *
 * http://hl7-fhir.github.io/search.html#string
 ******************************************************************************/

type StringParam struct {
	SearchParamInfo
	String string
}

func (s *StringParam) getInfo() SearchParamInfo {
	return s.SearchParamInfo
}

func ParseStringParam(paramString string, info SearchParamInfo) *StringParam {
	return &StringParam{info, unescape(paramString)}
}

/******************************************************************************
 * TOKEN: A token type is a parameter that searches on a pair, a URI and a
 * value. It is used against code or identifier value where the value may have
 * a URI that scopes its meaning. The search is performed against the pair from
 * a Coding or an Identifier.
 *
 * http://hl7-fhir.github.io/search.html#token
 ******************************************************************************/

type TokenParam struct {
	SearchParamInfo
	System    string
	Code      string
	AnySystem bool
}

func (t *TokenParam) getInfo() SearchParamInfo {
	return t.SearchParamInfo
}

func ParseTokenParam(paramString string, info SearchParamInfo) *TokenParam {
	t := &TokenParam{SearchParamInfo: info}
	splitCode := escapeFriendlySplit(paramString, '|')
	if len(splitCode) > 1 {
		t.System = unescape(splitCode[0])
		t.Code = unescape(splitCode[1])
	} else {
		t.AnySystem = true
		t.Code = unescape(splitCode[0])
	}
	return t
}

/******************************************************************************
 * URI: The uri parameter refers to an element which is URI (RFC 3986). Matches
 * are precise (e.g. case, accent, and escape) sensitive, and the entire URI
 * must match.
 *
 * http://hl7-fhir.github.io/search.html#uri
 ******************************************************************************/

type URIParam struct {
	SearchParamInfo
	URI string
}

func (u *URIParam) getInfo() SearchParamInfo {
	return u.SearchParamInfo
}

func ParseURIParam(paramStr string, info SearchParamInfo) *URIParam {
	return &URIParam{info, unescape(paramStr)}
}

/******************************************************************************
 * PREFIX: For the ordered parameter types number, date, and quantity, a prefix
 * to the parameter value may be used to control the nature of the matching.
 *
 * http://hl7-fhir.github.io/search.html#prefix
 ******************************************************************************/

type Prefix string

const (
	EQ Prefix = "eq"
	NE Prefix = "ne"
	GT Prefix = "gt"
	LT Prefix = "lt"
	GE Prefix = "ge"
	LE Prefix = "le"
	AP Prefix = "ap"
)

func (p Prefix) String() string {
	return string(p)
}

func ExtractPrefixAndValue(s string) (Prefix, string) {
	prefix := EQ
	for _, p := range []Prefix{EQ, NE, GT, LT, GE, LE, AP} {
		if strings.HasPrefix(s, p.String()) {
			prefix = p
			break
		}
	}
	return prefix, strings.TrimPrefix(s, prefix.String())
}

/******************************************************************************
 * ESCAPING SEARCH PARAMETERS: In the rules above, special rules are defined
 * for the characters "$", ",", and "|". As a consequence, if these characters
 * appear in an actual parameter value, they must be differentiated from their
 * use as separator characters. When any of these characters appear in an
 * actual parameter value, they must be prepended by the character "\" (which
 * also must be used to prepend itself).
 *
 * http://hl7-fhir.github.io/search.html#escaping
 ******************************************************************************/

func escapeFriendlySplit(s string, sep byte) []string {
	var result []string

	start := 0
	for i := range s {
		if s[i] == sep {
			// Count the preceding backslashes to see if it is escaped
			numBS := 0
			for j := i - 1; j >= 0; j-- {
				if s[j] == '\\' {
					numBS++
				} else {
					break
				}
			}
			// If number of preceding backslashes are even, it is not escaped
			if numBS%2 == 0 {
				result = append(result, s[start:i])
				start = i + 1
			}
		}
	}
	result = append(result, s[start:])

	return result
}

func unescape(s string) string {
	// A little hacky, but... otherwise there's a lot of annoying lookbacks/lookaheads
	s = strings.Replace(s, "\\\\", "```ie.bs```", -1)
	s = strings.Replace(s, "\\|", "|", -1)
	s = strings.Replace(s, "\\$", "$", -1)
	s = strings.Replace(s, "\\,", ",", -1)
	return strings.Replace(s, "```ie.bs```", "\\", -1)
}

func escape(s string) string {
	// A little hacky, but... otherwise there's a lot of annoying lookbacks/lookaheads
	s = strings.Replace(s, "\\", "```ie.bs```", -1)
	s = strings.Replace(s, "|", "\\|", -1)
	s = strings.Replace(s, "$", "\\$", -1)
	s = strings.Replace(s, ",", "\\,", -1)
	return strings.Replace(s, "```ie.bs```", "\\\\", -1)
}
