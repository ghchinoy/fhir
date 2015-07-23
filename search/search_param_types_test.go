package search

import (
	"math/big"
	"time"

	. "gopkg.in/check.v1"
)

type SearchPTSuite struct {
	MDT *time.Location
}

var _ = Suite(&SearchPTSuite{})

func (s *SearchPTSuite) SetUpSuite(c *C) {
	s.MDT = time.FixedZone("MDT", -7*60*60)
}

/******************************************************************************
 * DATE
 ******************************************************************************/

func (s *SearchPTSuite) TestDatesToMilliseconds(c *C) {

	d := NewDate("2013-01-02T12:13:14.999-07:00")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14.999-07:00")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 999000000, s.MDT).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, s.MDT).UnixNano())

	d = NewDate("2013-01-02T12:13:14.999Z")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14.999Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 999000000, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("2013-01-02T12:13:14.999")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14.999")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 999000000, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.Local).UnixNano())

	// Test different levels of precision
	d = NewDate("2013-01-02T12:13:14.9")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14.900")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 900000000, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 901000000, time.Local).UnixNano())

	d = NewDate("2013-01-02T12:13:14.09")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14.090")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 90000000, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 91000000, time.Local).UnixNano())

	d = NewDate("2013-01-02T12:13:14.009")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14.009")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 9000000, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 10000000, time.Local).UnixNano())

	d = NewDate("2013-01-02T12:13:14.987654321")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14.987")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 987000000, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 988000000, time.Local).UnixNano())
}

func (s *SearchPTSuite) TestDatesToSeconds(c *C) {

	d := NewDate("2013-01-02T12:13:14-07:00")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14-07:00")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, s.MDT).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, s.MDT).UnixNano())

	d = NewDate("2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("2013-01-02T12:13:14")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.Local).UnixNano())
}

func (s *SearchPTSuite) TestDatesToMinutes(c *C) {

	d := NewDate("2013-01-02T12:13-07:00")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13-07:00")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 0, 0, s.MDT).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 14, 0, 0, s.MDT).UnixNano())

	d = NewDate("2013-01-02T12:13Z")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 0, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 14, 0, 0, time.UTC).UnixNano())

	d = NewDate("2013-01-02T12:13")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 14, 0, 0, time.Local).UnixNano())
}

// NOTE: FHIR spec says that if hours are specified, minutes MUST be specified, so hours-only is invalid

func (s *SearchPTSuite) TestDatesToDays(c *C) {

	// Timezone should be ignored when no time components are included
	d := NewDate("2013-01-02T-07:00")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 3, 0, 0, 0, 0, time.Local).UnixNano())

	d = NewDate("2013-01-02Z")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 3, 0, 0, 0, 0, time.Local).UnixNano())

	d = NewDate("2013-01-02")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 3, 0, 0, 0, 0, time.Local).UnixNano())
}

func (s *SearchPTSuite) TestDatesToMonths(c *C) {

	// Timezone should be ignored when no time components are included
	d := NewDate("2013-01T-07:00")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.February, 1, 0, 0, 0, 0, time.Local).UnixNano())

	d = NewDate("2013-01Z")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.February, 1, 0, 0, 0, 0, time.Local).UnixNano())

	d = NewDate("2013-01")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.February, 1, 0, 0, 0, 0, time.Local).UnixNano())
}

func (s *SearchPTSuite) TestDatesToYears(c *C) {

	// Timezone should be ignored when no time components are included
	d := NewDate("2013T-07:00")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2014, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())

	d = NewDate("2013Z")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2014, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())

	d = NewDate("2013")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2014, time.January, 1, 0, 0, 0, 0, time.Local).UnixNano())
}

func (s *SearchPTSuite) TestLeapAndNonLeapYears(c *C) {

	// Non-Leap Year
	d := NewDate("1995-02-28")
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(1995, time.March, 1, 0, 0, 0, 0, time.Local).UnixNano())

	// Leap Year
	d = NewDate("1996-02-28")
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(1996, time.February, 29, 0, 0, 0, 0, time.Local).UnixNano())

	// Centurial Non-Leap Year (divisible by 4, but centuries are not leap years unless they are divisible by 400)
	d = NewDate("1900-02-28")
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(1900, time.March, 1, 0, 0, 0, 0, time.Local).UnixNano())

	// Centurial Leap Year (divisible by 4, and a century, but also divisible by 400-- so it IS a leap year)
	d = NewDate("2000-02-28")
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2000, time.February, 29, 0, 0, 0, 0, time.Local).UnixNano())
}

func (s *SearchPTSuite) TestDatePrefixes(c *C) {

	// Test prefixes
	d := NewDate("2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("eq2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, EQ)
	c.Assert(d.Text, Equals, "2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("ne2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, NE)
	c.Assert(d.Text, Equals, "ne2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("gt2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, GT)
	c.Assert(d.Text, Equals, "gt2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("lt2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, LT)
	c.Assert(d.Text, Equals, "lt2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("ge2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, GE)
	c.Assert(d.Text, Equals, "ge2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("le2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, LE)
	c.Assert(d.Text, Equals, "le2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())

	d = NewDate("ap2013-01-02T12:13:14Z")
	c.Assert(d.Prefix, Equals, AP)
	c.Assert(d.Text, Equals, "ap2013-01-02T12:13:14Z")
	c.Assert(d.RangeLowIncl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 14, 0, time.UTC).UnixNano())
	c.Assert(d.RangeHighExcl.UnixNano(), Equals, time.Date(2013, time.January, 2, 12, 13, 15, 0, time.UTC).UnixNano())
}

/******************************************************************************
 * NUMBER
 ******************************************************************************/

func (s *SearchPTSuite) TestNumbersThatAreInts(c *C) {
	n := NewNumber("100")

	c.Assert(n.Prefix, Equals, EQ)
	c.Assert(n.Text, Equals, "100")
	f, _ := n.Number.Float64()
	c.Assert(f, Equals, float64(100))
	f, _ = n.RangeLowIncl.Float64()
	c.Assert(f, Equals, float64(99.5))
	f, _ = n.RangeHighExcl.Float64()
	c.Assert(f, Equals, float64(100.5))
}

func (s *SearchPTSuite) TestNumbersThatAreNegativeInts(c *C) {
	n := NewNumber("-100")

	c.Assert(n.Prefix, Equals, EQ)
	c.Assert(n.Text, Equals, "-100")
	c.Assert(n.Number.IsInt(), Equals, true)
	f, _ := n.Number.Float64()
	c.Assert(f, Equals, float64(-100))
	f, _ = n.RangeLowIncl.Float64()
	c.Assert(f, Equals, float64(-100.5))
	f, _ = n.RangeHighExcl.Float64()
	c.Assert(f, Equals, float64(-99.5))
}

func (s *SearchPTSuite) TestNumbersThatAreDecimals(c *C) {
	n := NewNumber("100.00")

	c.Assert(n.Prefix, Equals, EQ)
	c.Assert(n.Text, Equals, "100.00")
	f, _ := n.Number.Float64()
	c.Assert(f, Equals, float64(100))
	f, _ = n.RangeLowIncl.Float64()
	c.Assert(f, Equals, float64(99.995))
	f, _ = n.RangeHighExcl.Float64()
	c.Assert(f, Equals, float64(100.005))

	n = NewNumber("100.0000000000000000000000000000000000000000")

	c.Assert(n.Prefix, Equals, EQ)
	c.Assert(n.Text, Equals, "100.0000000000000000000000000000000000000000")
	f, _ = n.Number.Float64()
	c.Assert(f, Equals, float64(100))
	f, _ = n.RangeLowIncl.Float64()
	c.Assert(f, Equals, float64(99.99999999999999999999999999999999999999995))
	f, _ = n.RangeHighExcl.Float64()
	c.Assert(f, Equals, float64(100.00000000000000000000000000000000000000005))
}

func (s *SearchPTSuite) TestNegativeNumbersThatAreDecimals(c *C) {
	n := NewNumber("-100.00")

	c.Assert(n.Prefix, Equals, EQ)
	c.Assert(n.Text, Equals, "-100.00")
	f, _ := n.Number.Float64()
	c.Assert(f, Equals, float64(-100))
	f, _ = n.RangeLowIncl.Float64()
	c.Assert(f, Equals, float64(-100.005))
	f, _ = n.RangeHighExcl.Float64()
	c.Assert(f, Equals, float64(-99.995))

	n = NewNumber("-100.0000000000000000000000000000000000000000")

	c.Assert(n.Prefix, Equals, EQ)
	c.Assert(n.Text, Equals, "-100.0000000000000000000000000000000000000000")
	f, _ = n.Number.Float64()
	c.Assert(f, Equals, float64(-100))
	f, _ = n.RangeLowIncl.Float64()
	c.Assert(f, Equals, float64(-100.00000000000000000000000000000000000000005))
	f, _ = n.RangeHighExcl.Float64()
	c.Assert(f, Equals, float64(-99.99999999999999999999999999999999999999995))
}

func (s *SearchPTSuite) TestNumberPrefixes(c *C) {
	b100, _ := new(big.Rat).SetString("100")
	b99_5, _ := new(big.Rat).SetString("99.5")
	b100_5, _ := new(big.Rat).SetString("100.5")

	n := NewNumber("100")
	c.Assert(n.Prefix, Equals, EQ)
	c.Assert(n.Text, Equals, "100")
	c.Assert(n.Number, DeepEquals, b100)
	c.Assert(n.RangeLowIncl, DeepEquals, b99_5)
	c.Assert(n.RangeHighExcl, DeepEquals, b100_5)

	n = NewNumber("eq100")
	c.Assert(n.Prefix, Equals, EQ)
	c.Assert(n.Text, Equals, "100")
	c.Assert(n.Number, DeepEquals, b100)
	c.Assert(n.RangeLowIncl, DeepEquals, b99_5)
	c.Assert(n.RangeHighExcl, DeepEquals, b100_5)

	n = NewNumber("ne100")
	c.Assert(n.Prefix, Equals, NE)
	c.Assert(n.Text, Equals, "ne100")
	c.Assert(n.Number, DeepEquals, b100)
	c.Assert(n.RangeLowIncl, DeepEquals, b99_5)
	c.Assert(n.RangeHighExcl, DeepEquals, b100_5)

	n = NewNumber("gt100")
	c.Assert(n.Prefix, Equals, GT)
	c.Assert(n.Text, Equals, "gt100")
	c.Assert(n.Number, DeepEquals, b100)
	c.Assert(n.RangeLowIncl, DeepEquals, b99_5)
	c.Assert(n.RangeHighExcl, DeepEquals, b100_5)

	n = NewNumber("lt100")
	c.Assert(n.Prefix, Equals, LT)
	c.Assert(n.Text, Equals, "lt100")
	c.Assert(n.Number, DeepEquals, b100)
	c.Assert(n.RangeLowIncl, DeepEquals, b99_5)
	c.Assert(n.RangeHighExcl, DeepEquals, b100_5)

	n = NewNumber("ge100")
	c.Assert(n.Prefix, Equals, GE)
	c.Assert(n.Text, Equals, "ge100")
	c.Assert(n.Number, DeepEquals, b100)
	c.Assert(n.RangeLowIncl, DeepEquals, b99_5)
	c.Assert(n.RangeHighExcl, DeepEquals, b100_5)

	n = NewNumber("le100")
	c.Assert(n.Prefix, Equals, LE)
	c.Assert(n.Text, Equals, "le100")
	c.Assert(n.Number, DeepEquals, b100)
	c.Assert(n.RangeLowIncl, DeepEquals, b99_5)
	c.Assert(n.RangeHighExcl, DeepEquals, b100_5)

	n = NewNumber("ap100")
	c.Assert(n.Prefix, Equals, AP)
	c.Assert(n.Text, Equals, "ap100")
	c.Assert(n.Number, DeepEquals, b100)
	c.Assert(n.RangeLowIncl, DeepEquals, b99_5)
	c.Assert(n.RangeHighExcl, DeepEquals, b100_5)
}

/******************************************************************************
 * QUANTITY
 ******************************************************************************/

func (s *SearchPTSuite) TestQuantitiesWithSystemsAndUnits(c *C) {
	q := NewQuantity("5.4|http://unitsofmeasure.org|mg")

	c.Assert(q.Prefix, Equals, EQ)
	c.Assert(q.Text, Equals, "5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")
}

func (s *SearchPTSuite) TestQuantitiesWithOnlyUnits(c *C) {
	q := NewQuantity("5.4||mg")

	c.Assert(q.Prefix, Equals, EQ)
	c.Assert(q.Text, Equals, "5.4||mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "")
	c.Assert(q.Code, Equals, "mg")
}

func (s *SearchPTSuite) TestQuantitiesWithNoUnits(c *C) {
	q := NewQuantity("5.4")

	c.Assert(q.Prefix, Equals, EQ)
	c.Assert(q.Text, Equals, "5.4||")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "")
	c.Assert(q.Code, Equals, "")
}

func (s *SearchPTSuite) TestNegativeQuantities(c *C) {
	q := NewQuantity("-10|http://unitsofmeasure.org|mg")

	c.Assert(q.Prefix, Equals, EQ)
	c.Assert(q.Text, Equals, "-10|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "-10")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")
}

func (s *SearchPTSuite) TestQuantitiesWithEscapedPipesAndSlashes(c *C) {
	q := NewQuantity("5.4|foo\\|bar|foo\\\\\\|baz")

	c.Assert(q.Prefix, Equals, EQ)
	c.Assert(q.Text, Equals, "5.4|foo\\|bar|foo\\\\\\|baz")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "foo|bar")
	c.Assert(q.Code, Equals, "foo\\|baz")
}

func (s *SearchPTSuite) TestQuantityPrefixes(c *C) {
	q := NewQuantity("5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Prefix, Equals, EQ)
	c.Assert(q.Text, Equals, "5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")

	q = NewQuantity("eq5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Prefix, Equals, EQ)
	c.Assert(q.Text, Equals, "5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")

	q = NewQuantity("ne5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Prefix, Equals, NE)
	c.Assert(q.Text, Equals, "ne5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")

	q = NewQuantity("gt5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Prefix, Equals, GT)
	c.Assert(q.Text, Equals, "gt5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")

	q = NewQuantity("lt5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Prefix, Equals, LT)
	c.Assert(q.Text, Equals, "lt5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")

	q = NewQuantity("ge5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Prefix, Equals, GE)
	c.Assert(q.Text, Equals, "ge5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")

	q = NewQuantity("le5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Prefix, Equals, LE)
	c.Assert(q.Text, Equals, "le5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")

	q = NewQuantity("ap5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Prefix, Equals, AP)
	c.Assert(q.Text, Equals, "ap5.4|http://unitsofmeasure.org|mg")
	c.Assert(q.Number.Text, Equals, "5.4")
	c.Assert(q.System, Equals, "http://unitsofmeasure.org")
	c.Assert(q.Code, Equals, "mg")
}

/******************************************************************************
 * REFERENCE
 ******************************************************************************/

func (s *SearchPTSuite) TestReferenceID(c *C) {
	r := NewReference("23")

	c.Assert(r.Reference, Equals, "23")
	c.Assert(r.IsId(), Equals, true)
	c.Assert(r.IsUrl(), Equals, false)
}

func (s *SearchPTSuite) TestReferenceAbsoluteURL(c *C) {
	r := NewReference("http://acme.org/fhir/Patient/23")

	c.Assert(r.Reference, Equals, "http://acme.org/fhir/Patient/23")
	c.Assert(r.IsId(), Equals, false)
	c.Assert(r.IsUrl(), Equals, true)
}

func (s *SearchPTSuite) TestReferenceRelativeURL(c *C) {
	r := NewReference("Patient/23")

	// According to FHIR spec, URLs must be absolute, so this is interpreted as ID
	c.Assert(r.Reference, Equals, "Patient/23")
	c.Assert(r.IsId(), Equals, true)
	c.Assert(r.IsUrl(), Equals, false)
}

/******************************************************************************
 * STRING
 ******************************************************************************/

func (s *SearchPTSuite) TestStrings(c *C) {
	n := NewString("Hello World")

	c.Assert(n.String, Equals, "Hello World")
}

/******************************************************************************
 * TOKEN
 ******************************************************************************/

func (s *SearchPTSuite) TestTokenCode(c *C) {
	t := NewToken("M")

	c.Assert(t.Text, Equals, "M")
	c.Assert(t.AnySystem, Equals, true)
	c.Assert(t.Code, Equals, "M")
	c.Assert(t.System, Equals, "")
}

func (s *SearchPTSuite) TestTokenSystemAndCode(c *C) {
	t := NewToken("http://hl7.org/fhir/v2/0001|M")

	c.Assert(t.Text, Equals, "http://hl7.org/fhir/v2/0001|M")
	c.Assert(t.AnySystem, Equals, false)
	c.Assert(t.Code, Equals, "M")
	c.Assert(t.System, Equals, "http://hl7.org/fhir/v2/0001")
}

func (s *SearchPTSuite) TestTokenSystemlessCode(c *C) {
	t := NewToken("|M")

	c.Assert(t.Text, Equals, "|M")
	c.Assert(t.AnySystem, Equals, false)
	c.Assert(t.Code, Equals, "M")
	c.Assert(t.System, Equals, "")
}

func (s *SearchPTSuite) TestTokensWithEscapedPipesAndSlashes(c *C) {
	t := NewToken("foo\\|bar")

	c.Assert(t.Text, Equals, "foo\\|bar")
	c.Assert(t.AnySystem, Equals, true)
	c.Assert(t.Code, Equals, "foo|bar")
	c.Assert(t.System, Equals, "")

	t = NewToken("foo\\|bar|foo\\\\\\|baz")

	c.Assert(t.Text, Equals, "foo\\|bar|foo\\\\\\|baz")
	c.Assert(t.AnySystem, Equals, false)
	c.Assert(t.Code, Equals, "foo\\|baz")
	c.Assert(t.System, Equals, "foo|bar")
}

/******************************************************************************
 * URI
 ******************************************************************************/

func (s *SearchPTSuite) TestURIs(c *C) {
	n := NewURI("http://acme.org/fhir/ValueSet/123")

	c.Assert(n.URI, Equals, "http://acme.org/fhir/ValueSet/123")
}

/******************************************************************************
 * PREFIX
 ******************************************************************************/

func (s *SearchPTSuite) TestPrefixes(c *C) {
	x, y := ExtractPrefixAndValue("eq10")
	c.Assert(x, Equals, EQ)
	c.Assert(y, Equals, "10")

	x, y = ExtractPrefixAndValue("ne10")
	c.Assert(x, Equals, NE)
	c.Assert(y, Equals, "10")

	x, y = ExtractPrefixAndValue("gt10")
	c.Assert(x, Equals, GT)
	c.Assert(y, Equals, "10")

	x, y = ExtractPrefixAndValue("lt10")
	c.Assert(x, Equals, LT)
	c.Assert(y, Equals, "10")

	x, y = ExtractPrefixAndValue("ge10")
	c.Assert(x, Equals, GE)
	c.Assert(y, Equals, "10")

	x, y = ExtractPrefixAndValue("le10")
	c.Assert(x, Equals, LE)
	c.Assert(y, Equals, "10")

	x, y = ExtractPrefixAndValue("ap10")
	c.Assert(x, Equals, AP)
	c.Assert(y, Equals, "10")
}

func (s *SearchPTSuite) TestPrefixDefault(c *C) {
	x, y := ExtractPrefixAndValue("10")
	c.Assert(x, Equals, EQ)
	c.Assert(y, Equals, "10")
}
