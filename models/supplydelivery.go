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

import (
	"encoding/json"
	"time"
)

type SupplyDelivery struct {
	Id           string           `json:"-" bson:"_id"`
	Identifier   *Identifier      `bson:"identifier,omitempty" json:"identifier,omitempty"`
	Status       string           `bson:"status,omitempty" json:"status,omitempty"`
	Patient      *Reference       `bson:"patient,omitempty" json:"patient,omitempty"`
	Type         *CodeableConcept `bson:"type,omitempty" json:"type,omitempty"`
	Quantity     *Quantity        `bson:"quantity,omitempty" json:"quantity,omitempty"`
	SuppliedItem *Reference       `bson:"suppliedItem,omitempty" json:"suppliedItem,omitempty"`
	Supplier     *Reference       `bson:"supplier,omitempty" json:"supplier,omitempty"`
	WhenPrepared *Period          `bson:"whenPrepared,omitempty" json:"whenPrepared,omitempty"`
	Time         *FHIRDateTime    `bson:"time,omitempty" json:"time,omitempty"`
	Destination  *Reference       `bson:"destination,omitempty" json:"destination,omitempty"`
	Receiver     []Reference      `bson:"receiver,omitempty" json:"receiver,omitempty"`
}

type SupplyDeliveryBundle struct {
	Type         string                      `json:"resourceType,omitempty"`
	Title        string                      `json:"title,omitempty"`
	Id           string                      `json:"id,omitempty"`
	Updated      time.Time                   `json:"updated,omitempty"`
	TotalResults int                         `json:"totalResults,omitempty"`
	Entry        []SupplyDeliveryBundleEntry `json:"entry,omitempty"`
	Category     SupplyDeliveryCategory      `json:"category,omitempty"`
}

type SupplyDeliveryBundleEntry struct {
	Title    string                 `json:"title,omitempty"`
	Id       string                 `json:"id,omitempty"`
	Content  SupplyDelivery         `json:"content,omitempty"`
	Category SupplyDeliveryCategory `json:"category,omitempty"`
}

type SupplyDeliveryCategory struct {
	Term   string `json:"term,omitempty"`
	Label  string `json:"label,omitempty"`
	Scheme string `json:"scheme,omitempty"`
}

func (resource *SupplyDelivery) MarshalJSON() ([]byte, error) {
	x := struct {
		ResourceType string `json:"resourceType"`
		SupplyDelivery
	}{
		ResourceType:   "SupplyDelivery",
		SupplyDelivery: *resource,
	}
	return json.Marshal(x)
}