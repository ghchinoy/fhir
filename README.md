Intervention Engine FHIR Server [![Build Status](https://travis-ci.org/intervention-engine/fhir.svg?branch=master)](https://travis-ci.org/intervention-engine/fhir)[![GoDoc](https://godoc.org/github.com/intervention-engine/fhir?status.svg)](https://godoc.org/github.com/intervention-engine/fhir)
===================================================================================================================================================================

This project provides [HL7 FHIR DSTU2](http://hl7.org/fhir/DSTU2/index.html) models and server components implemented in Go and using MongoDB as storage. Other versions of FHIR are also supported on other branches within this repository:

* [`stu3_jan2017`](https://github.com/intervention-engine/fhir/tree/stu3_jan2017) - FHIR 1.8 (San Antonio)
* [`stu3_aug2016`](https://github.com/intervention-engine/fhir/tree/stu3_aug2016) - FHIR 1.6 (Baltimore)
* [`stu3_mar2016`](https://github.com/intervention-engine/fhir/tree/stu3_mar2016) - FHIR 1.4 (Montreal)

This is a library that can be embedded into other server applications. The library is not a complete implementation of FHIR, as features that are selected are driven by the
[Intervention Engine](https://github.com/intervention-engine/ie), [eCQM Engine](https://github.com/mitre/ecqm), [Patient Matching Test Harness](https://github.com/mitre/ptmatch)
and [Synthetic Mass](https://github.com/synthetichealth/syntheticmass) projects.

Currently, this server library supports:

-	JSON representations of all resources
-	Create/Read/Update/Delete (CRUD) operations
-	Conditional update and delete
-	Some but not all search features
	-	All defined resource-specific search parameters except composite types and contact (email/phone) searches
	-	Chained searches
	-	\_include and \_revinclude searches (*without* \_recurse)
-	Batch bundle uploads (POST, PUT, and DELETE entries)

Currently, this server does *not* support the following major features:

-	XML representations of resources
-	History (versions, etc.)
-	Extension of primitive types and resource sub-components

As mentioned above, this repository contains a library. It does not build into a stand alone application. If you would like to build a server application, please see [fhir-server project](https://github.com/mitre/fhir-server).

*NOTE: Most of the fhir source code is generated by the [fhir-golang-generator](https://github.com/intervention-engine/fhir-golang-generator). In most cases, updates to source code in the fhir repository need to be accompanied by corresponding updates in the fhir-golang-generator.*

Development
-----------

This project uses Go 1.7 and MongoDB. To test the library, first, install all of the dependencies:

```
$ go get -t ./...
```

Once the dependencies are installed, you should make sure that MongoDB is also running. The test suite
will create a `fhir-test` database in a local instance of MongoDB to execute some tests. You can then
run the test suite with the following:

```
$ go test ./...
```

Usage
-----

Users of this library should work with the [FHIRServer](https://godoc.org/github.com/intervention-engine/fhir/server#FHIRServer) struct. Web request
handlers in this library are implemented using [Gin](https://gin-gonic.github.io/gin/).

Examples of usage can be found in the [server set up of the eCQM Engine](https://github.com/mitre/ecqm/blob/master/server.go) or the
[server set up of Intervention Engine](https://github.com/intervention-engine/ie/blob/master/server.go).

License
-------

Copyright 2017 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
