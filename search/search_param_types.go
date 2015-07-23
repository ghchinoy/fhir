package search

import (
	"fmt"
	"math/big"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"github.com/davecgh/go-spew/spew"
)

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

type Date struct {
	Prefix        Prefix
	RangeLowIncl  time.Time
	RangeHighExcl time.Time
	Text          string
}

var DT_REGEX = regexp.MustCompile("([0-9]{4})(-(0[1-9]|1[0-2])(-(0[0-9]|[1-2][0-9]|3[0-1])(T([01][0-9]|2[0-3]):([0-5][0-9])(:([0-5][0-9])(\\.([0-9]+))?)?((Z)|(\\+|-)((0[0-9]|1[0-3]):([0-5][0-9])|(14):(00)))?)?)?)?")

func NewDate(paramStr string) *Date {
	date := &Date{}
	var value string
	date.Prefix, value = ExtractPrefixAndValue(paramStr)
	if date.Prefix != EQ {
		date.Text = date.Prefix.String()
	}

	// TODO: What should we do if the time format is wrong?
	if m := DT_REGEX.FindStringSubmatch(value); m != nil {
		y, mo, d, h, mi, s, ms, tzStr, tzZu, tzOp, tzh, tzm := m[1], m[3], m[5], m[7], m[8], m[10], m[12], m[13], m[14], m[15], m[17], m[18]
		loc := time.Local
		if h != "" {
			// Fix the timezone
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

		layout := "2006-01-02T15:04:05.000"
		if ms != "" {
			// Fix milliseconds (.9 -> .900, .99 -> .990, .999999 -> .999 )
			if ms != "" {
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
			}

			date.RangeLowIncl, _ = time.ParseInLocation(layout, fmt.Sprintf("%s-%s-%sT%s:%s:%s.%s", y, mo, d, h, mi, s, ms), loc)
			date.RangeHighExcl = date.RangeLowIncl.Add(time.Millisecond)
			date.Text += fmt.Sprintf("%s-%s-%sT%s:%s:%s.%s%s", y, mo, d, h, mi, s, ms, tzStr)
		} else if s != "" {
			date.RangeLowIncl, _ = time.ParseInLocation(layout, fmt.Sprintf("%s-%s-%sT%s:%s:%s.000", y, mo, d, h, mi, s), loc)
			date.RangeHighExcl = date.RangeLowIncl.Add(time.Second)
			date.Text += fmt.Sprintf("%s-%s-%sT%s:%s:%s%s", y, mo, d, h, mi, s, tzStr)
		} else if mi != "" {
			date.RangeLowIncl, _ = time.ParseInLocation(layout, fmt.Sprintf("%s-%s-%sT%s:%s:00.000", y, mo, d, h, mi), loc)
			date.RangeHighExcl = date.RangeLowIncl.Add(time.Minute)
			date.Text += fmt.Sprintf("%s-%s-%sT%s:%s%s", y, mo, d, h, mi, tzStr)
		} else if d != "" {
			// NOTE: FHIR spec says that if hours are specified, minutes MUST be specified, so hours-only defaults to days
			date.RangeLowIncl, _ = time.ParseInLocation(layout, fmt.Sprintf("%s-%s-%sT00:00:00.000", y, mo, d), loc)
			date.RangeHighExcl = date.RangeLowIncl.AddDate(0, 0, 1)
			date.Text += fmt.Sprintf("%s-%s-%s", y, mo, d)
		} else if mo != "" {
			date.RangeLowIncl, _ = time.ParseInLocation(layout, fmt.Sprintf("%s-%s-01T00:00:00.000", y, mo), loc)
			date.RangeHighExcl = date.RangeLowIncl.AddDate(0, 1, 0)
			date.Text += fmt.Sprintf("%s-%s", y, mo)
		} else if y != "" {
			date.RangeLowIncl, _ = time.ParseInLocation(layout, fmt.Sprintf("%s-01-01T00:00:00.000", y), loc)
			date.RangeHighExcl = date.RangeLowIncl.AddDate(1, 0, 0)
			date.Text += fmt.Sprintf("%s", y)
		} else {
			// What to do?
		}
	}

	return date
}

/******************************************************************************
 * NUMBER: Searching on a simple numerical value in a resource.
 *
 * http://hl7-fhir.github.io/search.html#number
 ******************************************************************************/

type Number struct {
	Prefix        Prefix
	Number        *big.Rat
	RangeLowIncl  *big.Rat
	RangeHighExcl *big.Rat
	Text          string
}

func NewNumber(paramStr string) *Number {
	n := &Number{Number: new(big.Rat), RangeLowIncl: new(big.Rat), RangeHighExcl: new(big.Rat)}
	var value string
	n.Prefix, value = ExtractPrefixAndValue(paramStr)
	if n.Prefix != EQ {
		n.Text = n.Prefix.String()
	}
	n.Text += value
	n.Number.SetString(value)

	/* FHIR spec defines equality for 100 to be the range [99.5, 100.5) so we must support min/max using rounding semantics.
	 * The basic algorithm for determining low/high is:
	 *   low  (inclusive) = n - 5 / 10^p
	 *   high (exclusive) = n + 5 / 10^p
	 * where n is the number and p is the count of the number's decimal places + 1.
	 */

	p := 1
	i := strings.Index(value, ".")
	if i != -1 {
		p = len(value) - i
	}
	denomInt := new(big.Int).Exp(big.NewInt(int64(10)), big.NewInt(int64(p)), nil)
	denomRat, _ := new(big.Rat).SetString(denomInt.String())
	delta := new(big.Rat).Quo(new(big.Rat).SetInt64(5), denomRat)
	n.RangeLowIncl.Sub(n.Number, delta)
	n.RangeHighExcl.Add(n.Number, delta)

	return n
}

/******************************************************************************
 * QUANTITY: A quantity parameter searches on the Quantity data type.
 *
 * http://hl7-fhir.github.io/search.html#quantity
 ******************************************************************************/

type Quantity struct {
	Prefix Prefix
	Number *Number
	System string
	Code   string
	Text   string
}

func NewQuantity(paramStr string) *Quantity {
	q := &Quantity{}
	var value string
	q.Prefix, value = ExtractPrefixAndValue(paramStr)
	if q.Prefix != EQ {
		q.Text = q.Prefix.String()
	}

	split := escapeFriendlySplit(value, '|')
	q.Number = NewNumber(split[0])
	if len(split) == 3 {
		q.System = unescape(split[1])
		q.Code = unescape(split[2])
	}
	q.Text += fmt.Sprintf("%s|%s|%s", q.Number.Text, escape(q.System), escape(q.Code))

	return q
}

/******************************************************************************
 * REFERENCE: A reference parameter refers to references between resources,
 * e.g. find all Conditions where the subject reference is a particular patient,
 * where the patient is selected by name or identifier.
 *
 * http://hl7-fhir.github.io/search.html#reference
 ******************************************************************************/

type Reference struct {
	Reference string
}

func (r *Reference) IsId() bool {
	return !r.IsUrl()
}

func (r *Reference) IsUrl() bool {
	u, e := url.Parse(r.Reference)
	if e == nil {
		return u.IsAbs()
	}
	return false
}

func NewReference(paramStr string) *Reference {
	return &Reference{paramStr}
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

type String struct {
	String string
}

func NewString(paramStr string) *String {
	return &String{paramStr}
}

/******************************************************************************
 * TOKEN: A token type is a parameter that searches on a pair, a URI and a
 * value. It is used against code or identifier value where the value may have
 * a URI that scopes its meaning. The search is performed against the pair from
 * a Coding or an Identifier.
 *
 * http://hl7-fhir.github.io/search.html#token
 ******************************************************************************/

type Token struct {
	AnySystem bool
	System    string
	Code      string
	Text      string
}

func NewToken(paramStr string) *Token {
	t := &Token{Text: paramStr}
	splitCode := escapeFriendlySplit(paramStr, '|')
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

type URI struct {
	URI string
}

func NewURI(paramStr string) *URI {
	return &URI{paramStr}
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
