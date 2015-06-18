// Copyright (c) 2011-2015, HL7, Inc & The MITRE Corporation
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
//     * Redistributions of source code must retain the above copyright notice, this
//       list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright notice,
//       this list of conditions and the following disclaimer in the documentation
//       and/or other materials provided with the distribution.
//     * Neither the name of HL7 nor the names of its contributors may be used to
//       endorse or promote products derived from this software without specific
//       prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
// INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package models

import "time"

type MedicationPrescription struct {
	Id                        string                                                                   `json:"-" bson:"_id"`
	Identifier                []Identifier                                                             `bson:"identifier,omitempty" json:"identifier,omitempty"`
	DateWritten               *FHIRDateTime                                                            `bson:"dateWritten,omitempty" json:"dateWritten,omitempty"`
	Status                    string                                                                   `bson:"status,omitempty" json:"status,omitempty"`
	Patient                   *Reference                                                               `bson:"patient,omitempty" json:"patient,omitempty"`
	Prescriber                *Reference                                                               `bson:"prescriber,omitempty" json:"prescriber,omitempty"`
	Encounter                 *Reference                                                               `bson:"encounter,omitempty" json:"encounter,omitempty"`
	ReasonCodeableConcept     *CodeableConcept                                                         `bson:"reasonCodeableConcept,omitempty" json:"reasonCodeableConcept,omitempty"`
	ReasonReference           *Reference                                                               `bson:"reasonReference,omitempty" json:"reasonReference,omitempty"`
	Note                      string                                                                   `bson:"note,omitempty" json:"note,omitempty"`
	MedicationCodeableConcept *CodeableConcept                                                         `bson:"medicationCodeableConcept,omitempty" json:"medicationCodeableConcept,omitempty"`
	MedicationReference       *Reference                                                               `bson:"medicationReference,omitempty" json:"medicationReference,omitempty"`
	DosageInstruction         []MedicationPrescriptionMedicationPrescriptionDosageInstructionComponent `bson:"dosageInstruction,omitempty" json:"dosageInstruction,omitempty"`
	Dispense                  *MedicationPrescriptionMedicationPrescriptionDispenseComponent           `bson:"dispense,omitempty" json:"dispense,omitempty"`
	Substitution              *MedicationPrescriptionMedicationPrescriptionSubstitutionComponent       `bson:"substitution,omitempty" json:"substitution,omitempty"`
}
type MedicationPrescriptionMedicationPrescriptionDosageInstructionComponent struct {
	Text                    string           `bson:"text,omitempty" json:"text,omitempty"`
	AdditionalInstructions  *CodeableConcept `bson:"additionalInstructions,omitempty" json:"additionalInstructions,omitempty"`
	ScheduledDateTime       *FHIRDateTime    `bson:"scheduledDateTime,omitempty" json:"scheduledDateTime,omitempty"`
	ScheduledPeriod         *Period          `bson:"scheduledPeriod,omitempty" json:"scheduledPeriod,omitempty"`
	ScheduledTiming         *Timing          `bson:"scheduledTiming,omitempty" json:"scheduledTiming,omitempty"`
	AsNeededBoolean         *bool            `bson:"asNeededBoolean,omitempty" json:"asNeededBoolean,omitempty"`
	AsNeededCodeableConcept *CodeableConcept `bson:"asNeededCodeableConcept,omitempty" json:"asNeededCodeableConcept,omitempty"`
	Site                    *CodeableConcept `bson:"site,omitempty" json:"site,omitempty"`
	Route                   *CodeableConcept `bson:"route,omitempty" json:"route,omitempty"`
	Method                  *CodeableConcept `bson:"method,omitempty" json:"method,omitempty"`
	DoseRange               *Range           `bson:"doseRange,omitempty" json:"doseRange,omitempty"`
	DoseQuantity            *Quantity        `bson:"doseQuantity,omitempty" json:"doseQuantity,omitempty"`
	Rate                    *Ratio           `bson:"rate,omitempty" json:"rate,omitempty"`
	MaxDosePerPeriod        *Ratio           `bson:"maxDosePerPeriod,omitempty" json:"maxDosePerPeriod,omitempty"`
}
type MedicationPrescriptionMedicationPrescriptionDispenseComponent struct {
	MedicationCodeableConcept *CodeableConcept `bson:"medicationCodeableConcept,omitempty" json:"medicationCodeableConcept,omitempty"`
	MedicationReference       *Reference       `bson:"medicationReference,omitempty" json:"medicationReference,omitempty"`
	ValidityPeriod            *Period          `bson:"validityPeriod,omitempty" json:"validityPeriod,omitempty"`
	NumberOfRepeatsAllowed    *uint32          `bson:"numberOfRepeatsAllowed,omitempty" json:"numberOfRepeatsAllowed,omitempty"`
	Quantity                  *Quantity        `bson:"quantity,omitempty" json:"quantity,omitempty"`
	ExpectedSupplyDuration    *Quantity        `bson:"expectedSupplyDuration,omitempty" json:"expectedSupplyDuration,omitempty"`
}
type MedicationPrescriptionMedicationPrescriptionSubstitutionComponent struct {
	Type   *CodeableConcept `bson:"type,omitempty" json:"type,omitempty"`
	Reason *CodeableConcept `bson:"reason,omitempty" json:"reason,omitempty"`
}

type MedicationPrescriptionBundle struct {
	Type         string                              `json:"resourceType,omitempty"`
	Title        string                              `json:"title,omitempty"`
	Id           string                              `json:"id,omitempty"`
	Updated      time.Time                           `json:"updated,omitempty"`
	TotalResults int                                 `json:"totalResults,omitempty"`
	Entry        []MedicationPrescriptionBundleEntry `json:"entry,omitempty"`
	Category     MedicationPrescriptionCategory      `json:"category,omitempty"`
}

type MedicationPrescriptionBundleEntry struct {
	Title    string                         `json:"title,omitempty"`
	Id       string                         `json:"id,omitempty"`
	Content  MedicationPrescription         `json:"content,omitempty"`
	Category MedicationPrescriptionCategory `json:"category,omitempty"`
}

type MedicationPrescriptionCategory struct {
	Term   string `json:"term,omitempty"`
	Label  string `json:"label,omitempty"`
	Scheme string `json:"scheme,omitempty"`
}
