// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package id List is responsible for id List managing
package id

// List management

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/rokath/trice/pkg/msg"
)

// NewLut returns a look-up map generated from JSON map file named fn.
func NewLut(fn string) TriceIDLookUp {
	lu := make(TriceIDLookUp)
	msg.FatalOnErr(lu.fromFile(fn))
	if true == Verbose {
		fmt.Println("Read ID List file", fn, "with", len(lu), "items.")
	}
	return lu
}

// newID() gets a random ID not used so far.
// The delivered id is usable as key for lu, but not added. So calling fn twice without adding to lu could give the same value back.
// It is important that lu was refreshed before with all sources to avoid finding as a new ID an ID which is already used in the source tree.
func (lu TriceIDLookUp) newID(min, max TriceID) TriceID {
	if Verbose {
		fmt.Println("IDMin=", min, "IDMax=", max, "IDMethod=", SearchMethod)
	}
	switch SearchMethod {
	case "random":
		return lu.newRandomID(min, max)
	case "upward":
		return lu.newUpwardID(min, max)
	case "downward":
		return lu.newDownwardID(min, max)
	}
	msg.Info(fmt.Sprint("ERROR:", SearchMethod, "is unknown ID search method."))
	return 0
}

// newRandomID provides a random free ID inside interval [min,max].
// The delivered id is usable as key for lu, but not added. So calling fn twice without adding to lu could give the same value back.
func (lu TriceIDLookUp) newRandomID(min, max TriceID) (id TriceID) {
	interval := int(max - min + 1)
	freeIDs := interval - len(lu)
	msg.FatalInfoOnFalse(freeIDs > 0, "no new ID possible, "+fmt.Sprint(min, max, len(lu)))
	wrnLimit := interval >> 2 // 25%
	msg.InfoOnTrue(freeIDs < wrnLimit, "WARNING: Less than 25% IDs free!")
	id = min + TriceID(rand.Intn(interval))
	if 0 == len(lu) {
		return
	}
	for {
	nextTry:
		for k := range lu {
			if id == k { // id used
				fmt.Println("ID", id, "used, next try...")
				id = min + TriceID(rand.Intn(interval))
				goto nextTry
			}
		}
		return
	}
}

// newUpwardID provides the smallest free ID inside interval [min,max].
// The delivered id is usable as key for lut, but not added. So calling fn twice without adding to lu gives the same value back.
func (lu TriceIDLookUp) newUpwardID(min, max TriceID) (id TriceID) {
	interval := int(max - min + 1)
	freeIDs := interval - len(lu)
	msg.FatalInfoOnFalse(freeIDs > 0, "no new ID possible: "+fmt.Sprint("min=", min, ", max=", max, ", used=", len(lu)))
	id = min
	if 0 == len(lu) {
		return
	}
	for {
	nextTry:
		for k := range lu {
			if id == k { // id used
				id++
				goto nextTry
			}
		}
		return
	}
}

// newDownwardID provides the biggest free ID inside interval [min,max].
// The delivered id is usable as key for lut, but not added. So calling fn twice without adding to lu gives the same value back.
func (lu TriceIDLookUp) newDownwardID(min, max TriceID) (id TriceID) {
	interval := int(max - min + 1)
	freeIDs := interval - len(lu)
	msg.FatalInfoOnFalse(freeIDs > 0, "no new ID possible: "+fmt.Sprint("min=", min, ", max=", max, ", used=", len(lu)))
	id = max
	if 0 == len(lu) {
		return
	}
	for {
	nextTry:
		for k := range lu {
			if id == k { // id used
				id--
				goto nextTry
			}
		}
		return
	}
}

// FromJSON converts JSON byte slice to lu.
func (lu TriceIDLookUp) FromJSON(b []byte) (err error) {
	if 0 < len(b) {
		err = json.Unmarshal(b, &lu)
	}
	return
}

// fromFile reads file fn into lut. Existing keys are overwritten, lut is extended with new keys.
func (lu TriceIDLookUp) fromFile(fn string) error {
	b, err := ioutil.ReadFile(fn)
	msg.FatalInfoOnErr(err, "May be need to create an empty file first? (Safety feature)")
	return lu.FromJSON(b)
}

// AddFmtCount adds inside lu to all trice type names without format specifier count the appropriate count.
// example change:
// `map[10000:{Trice8_2 hi %03u, %5x} 10001:{TRICE16 hi %03u, %5x}]
// `map[10000:{Trice8_2 hi %03u, %5x} 10001:{TRICE16_2 hi %03u, %5x}]
func (lu TriceIDLookUp) AddFmtCount() {
	for i, x := range lu {
		if strings.ContainsAny(x.Type, "0_") {
			continue
		}
		n := FormatSpecifierCount(x.Strg)
		x.Type = addFormatSpecifierCount(x.Type, n)
		lu[i] = x
	}
}

// toJSON converts lut into JSON byte slice in human readable form.
func (lu TriceIDLookUp) toJSON() ([]byte, error) {
	return json.MarshalIndent(lu, "", "\t")
}

// toFile writes lut into file fn as indented JSON.
func (lu TriceIDLookUp) toFile(fn string) (err error) {
	var b []byte
	b, err = lu.toJSON()
	msg.FatalOnErr(err)
	var f *os.File
	f, err = os.Create(fn)
	msg.FatalOnErr(err)
	defer func() {
		err = f.Close()
	}()
	_, err = f.Write(b)
	return
}

// reverse returns a reversed map. If different triceID's assigned to several equal TriceFmt only one of the TriceID gets it into tflu.
func (lu TriceIDLookUp) reverse() (tflu TriceFmtLookUp) {
	tflu = make(TriceFmtLookUp)
	for id, tF := range lu {
		tF.Type = strings.ToUpper(tF.Type) // no distinction for lower and upper case Type
		tflu[tF] = id
	}
	return
}

//  // reverse returns a reversed map.  If different TriceFmt's assigned to several equal TriceFmt, this is an unexpected and unhandled error and only one of the TriceFmt's gets it into lu.
//  func (tflu TriceFmtLookUp) reverse() (lu TriceIDLookUp) {
//  	lu = make(TriceIDLookUp)
//  	for fm, id := range tflu {
//  		lu[id] = fm
//  	}
//  	return
//  }
