package search

// TODO: Generate this from FHIR

var SearchParameterDictionary = map[string]map[string]SearchParamInfo{
	"Condition": map[string]SearchParamInfo{
		"_id":     SearchParamInfo{Name: "_id", Type: "string", Paths: map[string]string{"_id": "string"}},
		"code":    SearchParamInfo{Name: "code", Type: "token", Paths: map[string]string{"code": "CodeableConcept"}},
		"patient": SearchParamInfo{Name: "patient", Type: "reference", Paths: map[string]string{"patient": "Reference"}},
	},
}
