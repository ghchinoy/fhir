package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/intervention-engine/fhir/models"
	"gopkg.in/mgo.v2/bson"
)

func HealthcareServiceIndexHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result []models.HealthcareService
	c := Database.C("healthcareservices")
	iter := c.Find(nil).Limit(100).Iter()
	err := iter.All(&result)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	var healthcareserviceEntryList []models.HealthcareServiceBundleEntry
	for _, healthcareservice := range result {
		var entry models.HealthcareServiceBundleEntry
		entry.Title = "HealthcareService " + healthcareservice.Id
		entry.Id = healthcareservice.Id
		entry.Content = healthcareservice
		healthcareserviceEntryList = append(healthcareserviceEntryList, entry)
	}

	var bundle models.HealthcareServiceBundle
	bundle.Type = "Bundle"
	bundle.Title = "HealthcareService Index"
	bundle.Id = bson.NewObjectId().Hex()
	bundle.Updated = time.Now()
	bundle.TotalResults = len(result)
	bundle.Entry = healthcareserviceEntryList

	log.Println("Setting healthcareservice search context")
	context.Set(r, "HealthcareService", result)
	context.Set(r, "Resource", "HealthcareService")
	context.Set(r, "Action", "search")

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(bundle)
}

func LoadHealthcareService(r *http.Request) (*models.HealthcareService, error) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		return nil, errors.New("Invalid id")
	}

	c := Database.C("healthcareservices")
	result := models.HealthcareService{}
	err := c.Find(bson.M{"_id": id.Hex()}).One(&result)
	if err != nil {
		return nil, err
	}

	log.Println("Setting healthcareservice read context")
	context.Set(r, "HealthcareService", result)
	context.Set(r, "Resource", "HealthcareService")
	return &result, nil
}

func HealthcareServiceShowHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, "Action", "read")
	_, err := LoadHealthcareService(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(rw).Encode(context.Get(r, "HealthcareService"))
}

func HealthcareServiceCreateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	decoder := json.NewDecoder(r.Body)
	healthcareservice := &models.HealthcareService{}
	err := decoder.Decode(healthcareservice)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("healthcareservices")
	i := bson.NewObjectId()
	healthcareservice.Id = i.Hex()
	err = c.Insert(healthcareservice)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting healthcareservice create context")
	context.Set(r, "HealthcareService", healthcareservice)
	context.Set(r, "Resource", "HealthcareService")
	context.Set(r, "Action", "create")

	host, err := os.Hostname()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Add("Location", "http://"+host+":3001/HealthcareService/"+i.Hex())
}

func HealthcareServiceUpdateHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	decoder := json.NewDecoder(r.Body)
	healthcareservice := &models.HealthcareService{}
	err := decoder.Decode(healthcareservice)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	c := Database.C("healthcareservices")
	healthcareservice.Id = id.Hex()
	err = c.Update(bson.M{"_id": id.Hex()}, healthcareservice)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Setting healthcareservice update context")
	context.Set(r, "HealthcareService", healthcareservice)
	context.Set(r, "Resource", "HealthcareService")
	context.Set(r, "Action", "update")
}

func HealthcareServiceDeleteHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var id bson.ObjectId

	idString := mux.Vars(r)["id"]
	if bson.IsObjectIdHex(idString) {
		id = bson.ObjectIdHex(idString)
	} else {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
	}

	c := Database.C("healthcareservices")

	err := c.Remove(bson.M{"_id": id.Hex()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Setting healthcareservice delete context")
	context.Set(r, "HealthcareService", id.Hex())
	context.Set(r, "Resource", "HealthcareService")
	context.Set(r, "Action", "delete")
}
