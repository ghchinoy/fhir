package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	//"github.com/davecgh/go-spew/spew"
	"github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/dbtest"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MongoSearchSuite struct {
	DBServer      *dbtest.DBServer
	Session       *mgo.Session
	FhirModels    []interface{}
	MongoSearcher *MongoSearcher
}

var _ = Suite(&MongoSearchSuite{})

func (m *MongoSearchSuite) SetUpSuite(c *C) {
	// Set up the database
	m.DBServer = &dbtest.DBServer{}
	m.DBServer.SetPath(c.MkDir())

	m.Session = m.DBServer.Session()
	db := m.Session.DB("fhir-test")
	m.MongoSearcher = &MongoSearcher{db}

	// Read in the data in FHIR format
	data, err := ioutil.ReadFile("../fixtures/john_peters_direct.json")
	util.CheckErr(err)

	maps := make([]interface{}, 19)
	err = json.Unmarshal(data, &maps)
	util.CheckErr(err)

	for _, resourceMap := range maps {
		r := models.MapToResource(resourceMap, true)
		switch r.(type) {
		case *models.Condition:
			util.CheckErr(db.C("conditions").Insert(r))
		case *models.DiagnosticReport:
			util.CheckErr(db.C("diagnosticreports").Insert(r))
		case *models.Encounter:
			util.CheckErr(db.C("encounters").Insert(r))
		case *models.Immunization:
			util.CheckErr(db.C("immunizations").Insert(r))
		case *models.MedicationStatement:
			util.CheckErr(db.C("medicationstatements").Insert(r))
		case *models.Observation:
			util.CheckErr(db.C("observations").Insert(r))
		case *models.Patient:
			util.CheckErr(db.C("patients").Insert(r))
		case *models.Procedure:
			util.CheckErr(db.C("procedures").Insert(r))
		default:
			fmt.Printf("NOT FOUND: %T\n", r)
		}
	}
}

func (m *MongoSearchSuite) TearDownSuite(c *C) {
	m.Session.Close()
	m.DBServer.Wipe()
	m.DBServer.Stop()
}

// TODO: Test anything w/ multiple paths

func (m *MongoSearchSuite) TestConditionIdQueryObject(c *C) {
	q := Query{"Condition", "_id=123456789"}

	o := m.MongoSearcher.CreateQueryObject(q)
	c.Assert(o, DeepEquals, bson.M{"_id": bson.RegEx{Pattern: "^123456789$", Options: "i"}})
}

// TODO: Test other variations of token types

func (m *MongoSearchSuite) TestConditionCodeQueryObjectBySystemAndCode(c *C) {
	q := Query{"Condition", "code=http://snomed.info/sct|123641001"}
	o := m.MongoSearcher.CreateQueryObject(q)
	c.Assert(o, DeepEquals, bson.M{
		"code.coding": bson.M{
			"$elemMatch": bson.M{
				"system": ci("http://snomed.info/sct"),
				"code":   ci("123641001"),
			},
		},
	})
}

func (m *MongoSearchSuite) TestConditionCodeQueryObjectByCode(c *C) {
	q := Query{"Condition", "code=123641001"}

	o := m.MongoSearcher.CreateQueryObject(q)
	c.Assert(o, DeepEquals, bson.M{"code.coding.code": ci("123641001")})
}

func (m *MongoSearchSuite) TestConditionReferenceQueryObjectByPatientId(c *C) {
	q := Query{"Condition", "patient=123456789"}

	o := m.MongoSearcher.CreateQueryObject(q)
	c.Assert(o, DeepEquals, bson.M{"patient.referenceid": ci("123456789")})
}

func (m *MongoSearchSuite) TestConditionReferenceQueryObjectByPatientTypeAndId(c *C) {
	q := Query{"Condition", "patient=Patient/123456789"}

	o := m.MongoSearcher.CreateQueryObject(q)
	c.Assert(o, DeepEquals, bson.M{"patient.referenceid": ci("123456789"), "patient.type": "Patient"})
}

func (m *MongoSearchSuite) TestConditionReferenceQueryObjectByPatientURL(c *C) {
	q := Query{"Condition", "patient=http://acme.com/Patient/123456789"}

	o := m.MongoSearcher.CreateQueryObject(q)
	c.Assert(o, DeepEquals, bson.M{"patient.reference": ci("http://acme.com/Patient/123456789")})
}

func (m *MongoSearchSuite) TestConditionIdQuery(c *C) {
	q := Query{"Condition", "_id=8664777288161060797"}
	mq := m.MongoSearcher.CreateQuery(q)
	num, err := mq.Count()
	util.CheckErr(err)
	c.Assert(num, Equals, 1)

	cond := &models.Condition{}
	err = mq.One(cond)
	util.CheckErr(err)

	cond2 := &models.Condition{}
	err = m.Session.DB("fhir-test").C("conditions").FindId("8664777288161060797").One(cond2)

	c.Assert(cond, DeepEquals, cond2)
}

func (m *MongoSearchSuite) TestConditionCodeQueryBySystemAndCode(c *C) {
	var conditions []*models.Condition
	q := Query{"Condition", "code=http://snomed.info/sct|123641001"}
	mq := m.MongoSearcher.CreateQuery(q)
	err := mq.All(&conditions)
	util.CheckErr(err)
	c.Assert(conditions, HasLen, 2)
	foundIvd, foundCad := false, false
	for _, cond := range conditions {
		if strings.Contains(cond.Code.Text, "Ischemic Vascular Disease") {
			foundIvd = true
		} else if strings.Contains(cond.Code.Text, "Coronary Artery Disease No MI") {
			foundCad = true
		}
	}
	c.Assert(foundIvd && foundCad, Equals, true)
}

func (m *MongoSearchSuite) TestConditionCodeQueryByCode(c *C) {
	var conditions []*models.Condition
	q := Query{"Condition", "code=123641001"}
	mq := m.MongoSearcher.CreateQuery(q)
	err := mq.All(&conditions)
	util.CheckErr(err)
	c.Assert(conditions, HasLen, 2)
	foundIvd, foundCad := false, false
	for _, cond := range conditions {
		if strings.Contains(cond.Code.Text, "Ischemic Vascular Disease") {
			foundIvd = true
		} else if strings.Contains(cond.Code.Text, "Coronary Artery Disease No MI") {
			foundCad = true
		}
	}
	c.Assert(foundIvd && foundCad, Equals, true)
}

func (m *MongoSearchSuite) TestConditionCodeQueryByWrongCodeSystem(c *C) {
	var conditions []*models.Condition
	q := Query{"Condition", "code=http://hl7.org/fhir/sid/icd-9|123641001"}
	mq := m.MongoSearcher.CreateQuery(q)
	err := mq.All(&conditions)
	util.CheckErr(err)
	c.Assert(conditions, HasLen, 0)
}

func (m *MongoSearchSuite) TestConditionPatientQueryById(c *C) {
	var conditions []*models.Condition

	q := Query{"Condition", "patient=4954037118555241963"}
	mq := m.MongoSearcher.CreateQuery(q)
	err := mq.All(&conditions)
	util.CheckErr(err)
	c.Assert(conditions, HasLen, 5)
}

/*
func (m *MongoSearchSuite) TestConditionPatientQueryByTypeAndId(c *C) {
	var conditions []*models.Condition

	q := Query{"Condition", "patient=Patient/4954037118555241963"}
	mq := m.MongoSearcher.CreateQuery(q)
	err := mq.All(&conditions)
	util.CheckErr(err)
	c.Assert(conditions, HasLen, 5)
}
*/
