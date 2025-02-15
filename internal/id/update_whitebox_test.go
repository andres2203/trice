// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.
package id

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

type idCheck struct {
	nbTRICE string
	nbID    string
	id      TriceID
	ok      bool
	tf      TriceFmt
}

// Only the first occurrence of a double used ID gets in the list.
// Only the list generation out of sources ist tested here.
func TestRefreshIDListSingle0(t *testing.T) {
	text := `
	TRICE16_3( Id(12345), "hi %2d, %13u, %64b\n",1,2,3 );	
	TRICE16_3( Id(12345), "DIFFERENT! %2d, %13u, %64b\n",1,2,3 );	
	Trice16_1( Id(12344), "hi %2d\n",1 );	
	Trice16_1i( Id(12344), "hi %2d\n",1 );	
	Trice16_1( Id(12344), "DIFFERENT! %2d\n", 2 );	
`
	expJSON := `{
	"12344": {
		"Type": "Trice16_1",
		"Strg": "hi %2d\\n"
	},
	"12345": {
		"Type": "TRICE16_3",
		"Strg": "hi %2d, %13u, %64b\\n"
	}
}`
	check(t, text, expJSON)
}

func TestRefreshIDListSingle1(t *testing.T) {
	text := `
	TRICE16_3( Id(12345), "hi %2d, %13u, %64b\n",1,2,3 );	
	TRICE16_3( Id(  123), "hi %2d, %13u, %64b\n",1,2,3 );	
	Trice16_1( Id(   13), "hi %13u\n", 3 );	
	Trice16_1i( Id(   13), "hi %13u\n", 3 );	
`
	expJSON := `{
	"123": {
		"Type": "TRICE16_3",
		"Strg": "hi %2d, %13u, %64b\\n"
	},
	"12345": {
		"Type": "TRICE16_3",
		"Strg": "hi %2d, %13u, %64b\\n"
	},
	"13": {
		"Type": "Trice16_1",
		"Strg": "hi %13u\\n"
	}
}`
	check(t, text, expJSON)
}

func TestRefreshIDListSingle2(t *testing.T) {
	text := `
	TRICE16_3( Id(12345), "hi %2d, %13u, %64b\n",1,2,3 );
	trice16_3( Id(12345), "hi %2d, %13u, %64b\n",1,2,3 );
	Trice16_1i( Id(12344), "ho %64b\n",1 );
	Trice16_1( Id(12344), "ho %64b\n",1 );
`
	expJSON := `{
	"12344": {
		"Type": "Trice16_1i",
		"Strg": "ho %64b\\n"
	},
	"12345": {
		"Type": "trice16_3",
		"Strg": "hi %2d, %13u, %64b\\n"
	}
}`
	check(t, text, expJSON)
}

func check(t *testing.T, text, expJSON string) {
	lu := make(TriceIDLookUp)
	tflu := lu.reverse()
	refreshIDs(text, lu, tflu)
	b, err := lu.toJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, expJSON, string(b))
}

// Because of the parallel test execution the global variables must be equal for all tests
func TestInsertSharedIDs0ZeroParam(t *testing.T) {
	SearchMethod = "downward"
	tt := testTable{
		{`... Trice0 ( "hi"); ...`, `... Trice0 ( id(32767), "hi"); ...`, true, true},
		{`... TRICE0 ( "hi"); ...`, `... TRICE0 ( Id(    99), "hi"); ...`, true, true},
		{`... trice0 ( "hi"); ...`, `... trice0 ( Id(    99), "hi"); ...`, true, false},
		{`... Trice0i( "hi"); ...`, `... Trice0i( id(32766), "hi"); ...`, true, true},
		{`... TRICE0i( "hi"); ...`, `... TRICE0i( Id(    98), "hi"); ...`, true, true},
		{`... trice0i( "hi"); ...`, `... trice0i( Id(    98), "hi"); ...`, true, false},
	}
	eList := `map[98:{trice0i hi} 99:{trice0 hi} 32766:{Trice0i hi} 32767:{Trice0 hi}]
`
	checkList(t, true, 1, 32767, 10, 99, tt, eList, true)
	checkList(t, true, 1, 32767, 10, 99, tt, eList, false)
}

// Because of the parallel test execution the global variables must be equal for all tests
func TestInsertSharedIDs1WithExtendN(t *testing.T) {
	SearchMethod = "upward"
	tt := testTable{
		{`...  Trice8 ( "hi %03u", 5); ...`, `...  Trice8_1 ( id(10000), "hi %03u", 5); ...`, true, true},
		{`...  TRICE8 ( "hi %03u", 5); ...`, `...  TRICE8_1 ( Id( 40000), "hi %03u", 5); ...`, true, true},
		{`...  trice8 ( "hi %03u", 5); ...`, `...  trice8_1 ( Id( 40000), "hi %03u", 5); ...`, true, false},

		{`... Trice16 ( "hi %03u", 5); ...`, `... Trice16_1 ( id(10001), "hi %03u", 5); ...`, true, true},
		{`... TRICE16 ( "hi %03u", 5); ...`, `... TRICE16_1 ( Id( 40001), "hi %03u", 5); ...`, true, true},
		{`... trice16 ( "hi %03u", 5); ...`, `... trice16_1 ( Id( 40001), "hi %03u", 5); ...`, true, false},

		{`... Trice32 ( "hi %03u", 5); ...`, `... Trice32_1 ( id(10002), "hi %03u", 5); ...`, true, true}, // does not exist but allowed to match
		{`... TRICE32 ( "hi %03u", 5); ...`, `... TRICE32_1 ( Id( 40002), "hi %03u", 5); ...`, true, true},
		{`... trice32 ( "hi %03u", 5); ...`, `... trice32_1 ( Id( 40002), "hi %03u", 5); ...`, true, false},

		{`... Trice64 ( "hi %03u", 5); ...`, `... Trice64_1 ( id(10003), "hi %03u", 5); ...`, true, true},
		{`... TRICE64 ( "hi %03u", 5); ...`, `... TRICE64_1 ( Id( 40003), "hi %03u", 5); ...`, true, true},
		{`... trice64 ( "hi %03u", 5); ...`, `... trice64_1 ( Id( 40003), "hi %03u", 5); ...`, true, false},

		{`...  Trice8i( "hi %03u", 5); ...`, `...  Trice8_1i( id(10004), "hi %03u", 5); ...`, true, true},
		{`...  TRICE8i( "hi %03u", 5); ...`, `...  TRICE8_1i( Id( 40004), "hi %03u", 5); ...`, true, true},
		{`...  trice8i( "hi %03u", 5); ...`, `...  trice8_1i( Id( 40004), "hi %03u", 5); ...`, true, false},

		{`... Trice16i( "hi %03u", 5); ...`, `... Trice16_1i( id(10005), "hi %03u", 5); ...`, true, true},
		{`... TRICE16i( "hi %03u", 5); ...`, `... TRICE16_1i( Id( 40005), "hi %03u", 5); ...`, true, true},
		{`... trice16i( "hi %03u", 5); ...`, `... trice16_1i( Id( 40005), "hi %03u", 5); ...`, true, false},

		{`... Trice32i( "hi %03u", 5); ...`, `... Trice32_1i( id(10006), "hi %03u", 5); ...`, true, true},
		{`... TRICE32i( "hi %03u", 5); ...`, `... TRICE32_1i( Id( 40006), "hi %03u", 5); ...`, true, true},
		{`... trice32i( "hi %03u", 5); ...`, `... trice32_1i( Id( 40006), "hi %03u", 5); ...`, true, false},

		{`... Trice64i( "hi %03u", 5); ...`, `... Trice64_1i( id(10007), "hi %03u", 5); ...`, true, true},
		{`... TRICE64i( "hi %03u", 5); ...`, `... TRICE64_1i( Id( 40007), "hi %03u", 5); ...`, true, true},
		{`... trice64i( "hi %03u", 5); ...`, `... trice64_1i( Id( 40007), "hi %03u", 5); ...`, true, false},
	}
	eList := `map[10000:{Trice8_1 hi %03u} 10001:{Trice16_1 hi %03u} 10002:{Trice32_1 hi %03u} 10003:{Trice64_1 hi %03u} 10004:{Trice8_1i hi %03u} 10005:{Trice16_1i hi %03u} 10006:{Trice32_1i hi %03u} 10007:{Trice64_1i hi %03u} 40000:{trice8_1 hi %03u} 40001:{trice16_1 hi %03u} 40002:{trice32_1 hi %03u} 40003:{trice64_1 hi %03u} 40004:{trice8_1i hi %03u} 40005:{trice16_1i hi %03u} 40006:{trice32_1i hi %03u} 40007:{trice64_1i hi %03u}]
`
	checkList(t, true, 10000, 20000, 40000, 50000, tt, eList, true)
}

// Because of the parallel test execution the global variables must be equal for all tests.
// The trice map does distinguish between TRICE8 and TRICE8_2 for example, so even "sameID" is selected,
// there are 2 different IDs used when the format string is identical.
func TestInsertSharedIDs2NoExtendN(t *testing.T) {
	SearchMethod = "upward"
	tt := []struct {
		text, exp        string
		fileMod, listMod bool
	}{
		{`...  Trice8_2 ( "hi %03u, %5x", 5, 7); ...`, `...  Trice8_2 ( id(10000), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`...  Trice8_2 ( "hi %03u, %5x", 5, 7); ...`, `...  Trice8_2 ( id(10000), "hi %03u, %5x", 5, 7); ...`, true, false},
		{`...  Trice8_2 ( id(10000), "hi %03u, %5x", 5, 7); ...`, `...  Trice8_2 ( id(10000), "hi %03u, %5x", 5, 7); ...`, false, false},

		{`...  TRICE8 ( "hi %03u, %5x", 5, 7); ...`, `...  TRICE8 ( Id( 40000), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`...  trice8 ( "hi %03u, %5x", 5, 7); ...`, `...  trice8 ( Id( 40000), "hi %03u, %5x", 5, 7); ...`, true, false},

		{`... Trice16 ( "hi %03u, %5x", 5, 7); ...`, `... Trice16 ( id(10001), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`... TRICE16 ( "hi %03u, %5x", 5, 7); ...`, `... TRICE16 ( Id( 40001), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`... trice16 ( "hi %03u, %5x", 5, 7); ...`, `... trice16 ( Id( 40001), "hi %03u, %5x", 5, 7); ...`, true, false},

		{`... Trice32 ( "hi %03u, %5x", 5, 7); ...`, `... Trice32 ( id(10002), "hi %03u, %5x", 5, 7); ...`, true, true}, // does not exist but allowed to match
		{`... TRICE32 ( "hi %03u, %5x", 5, 7); ...`, `... TRICE32 ( Id( 40002), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`... trice32 ( "hi %03u, %5x", 5, 7); ...`, `... trice32 ( Id( 40002), "hi %03u, %5x", 5, 7); ...`, true, false},

		{`... Trice64 ( "hi %03u, %5x", 5, 7); ...`, `... Trice64 ( id(10003), "hi %03u, %5x", 5, 7); ...`, true, true}, // does not exist but allowed to match
		{`... TRICE64 ( "hi %03u, %5x", 5, 7); ...`, `... TRICE64 ( Id( 40003), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`... trice64 ( "hi %03u, %5x", 5, 7); ...`, `... trice64 ( Id( 40003), "hi %03u, %5x", 5, 7); ...`, true, false},

		{`...  Trice8i( "hi %03u, %5x", 5, 7); ...`, `...  Trice8i( id(10004), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`...  TRICE8i( "hi %03u, %5x", 5, 7); ...`, `...  TRICE8i( Id( 40004), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`...  trice8i( "hi %03u, %5x", 5, 7); ...`, `...  trice8i( Id( 40004), "hi %03u, %5x", 5, 7); ...`, true, false},

		{`... Trice16i( "hi %03u, %5x", 5, 7); ...`, `... Trice16i( id(10005), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`... TRICE16i( "hi %03u, %5x", 5, 7); ...`, `... TRICE16i( Id( 40005), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`... trice16i( "hi %03u, %5x", 5, 7); ...`, `... trice16i( Id( 40005), "hi %03u, %5x", 5, 7); ...`, true, false},

		{`... Trice32i( "hi %03u, %5x", 5, 7); ...`, `... Trice32i( id(10006), "hi %03u, %5x", 5, 7); ...`, true, true}, // does not exist but allowed to match
		{`... TRICE32i( "hi %03u, %5x", 5, 7); ...`, `... TRICE32i( Id( 40006), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`... trice32i( "hi %03u, %5x", 5, 7); ...`, `... trice32i( Id( 40006), "hi %03u, %5x", 5, 7); ...`, true, false},

		{`... Trice64i( "hi %03u, %5x", 5, 7); ...`, `... Trice64i( id(10007), "hi %03u, %5x", 5, 7); ...`, true, true}, // does not exist but allowed to match
		{`... TRICE64i( "hi %03u, %5x", 5, 7); ...`, `... TRICE64i( Id( 40007), "hi %03u, %5x", 5, 7); ...`, true, true},
		{`... trice64i( "hi %03u, %5x", 5, 7); ...`, `... trice64i( Id( 40007), "hi %03u, %5x", 5, 7); ...`, true, false},

		{`... trice64_2i( "hi %03u, %5x", 5, 7); ...`, `... trice64_2i( Id( 40008), "hi %03u, %5x", 5, 7); ...`, true, true},
	}
	eList := `map[10000:{Trice8_2 hi %03u, %5x} 10001:{Trice16_2 hi %03u, %5x} 10002:{Trice32_2 hi %03u, %5x} 10003:{Trice64_2 hi %03u, %5x} 10004:{Trice8_2i hi %03u, %5x} 10005:{Trice16_2i hi %03u, %5x} 10006:{Trice32_2i hi %03u, %5x} 10007:{Trice64_2i hi %03u, %5x} 40000:{trice8_2 hi %03u, %5x} 40001:{trice16_2 hi %03u, %5x} 40002:{trice32_2 hi %03u, %5x} 40003:{trice64_2 hi %03u, %5x} 40004:{trice8_2i hi %03u, %5x} 40005:{trice16_2i hi %03u, %5x} 40006:{trice32_2i hi %03u, %5x} 40007:{trice64_2i hi %03u, %5x} 40008:{trice64_2i hi %03u, %5x}]
`
	checkList(t, true, 10000, 20000, 40000, 50000, tt, eList, false)
}

// Because of the parallel test execution the global variables must be equal for all tests
func TestInsertSharedIDs0WithParamCount(t *testing.T) {
	SearchMethod = "downward"

	tt := testTable{
		{
			`Trice0(id(0), "Trallala");`,
			`Trice0(id(20000), "Trallala");`, true, true},
		{
			`trice8_1( Id(0), "hi %d", 5); // first id`,
			`trice8_1( Id(    99), "hi %d", 5); // first id`, true, true},
		{
			`trice8_1( Id(0), "Hi %d", 5); // different format string needs a new id`,
			`trice8_1( Id(    98), "Hi %d", 5); // different format string needs a new id`, true, true},
		{
			`TRICE8_1( Id(       0   ), "Hi %d", 5); // different type case gets same id`,
			`TRICE8_1( Id(    98), "Hi %d", 5); // different type case gets same id`, true, false},
		{
			`TRICE8_1( Id(0), "hi %d", 5); // same format string gets same id`,
			`TRICE8_1( Id(    99), "hi %d", 5); // same format string gets same id`, true, false},
		{
			`trice8_1(  Id( 0       ),  "hi %d", 5); // same format string gets same id`,
			`trice8_1(  Id(    99),  "hi %d", 5); // same format string gets same id`, true, false},
		{
			`Trice8_1( Id(0), "hi %d", 5); // first id`,
			`Trice8_1( id(19999), "hi %d", 5); // first id`, true, true},
		{
			`Trice8_2( id(0), "hi %d %u", 5, 7); // first id`,
			`Trice8_2( id(19998), "hi %d %u", 5, 7); // first id`, true, true},
		{
			`Trice16_1( Id(0), "hi %d", 5); // first id`,
			`Trice16_1( id(19997), "hi %d", 5); // first id`, true, true},
		{
			`Trice16_1i( Id(0), "hi %d", 5); // first id`,
			`Trice16_1i( id(19996), "hi %d", 5); // first id`, true, true},
	}
	eList := `map[98:{TRICE8_1 Hi %d} 99:{trice8_1 hi %d} 19996:{Trice16_1i hi %d} 19997:{Trice16_1 hi %d} 19998:{Trice8_2 hi %d %u} 19999:{Trice8_1 hi %d} 20000:{Trice0 Trallala}]
`
	checkList(t, true, 10000, 20000, 10, 99, tt, eList, true)
	checkList(t, true, 10000, 20000, 10, 99, tt, eList, false)
}

func TestInsertSharedIDsInvalid0(t *testing.T) {
	SearchMethod = "downward"
	tt := testTable{
		{
			`trice8_1( Id(0x0), "Hi %d", 5); // this gets not changed because no decimal number`,
			`trice8_1( Id(0x0), "Hi %d", 5); // this gets not changed because no decimal number`, false, false},
		{
			`trice8_1( Id(-0), "hi %d", 5); // minus id's are not touched`,
			`trice8_1( Id(-0), "hi %d", 5); // minus id's are not touched`, false, false},
	}
	eList := `map[]
`
	checkList(t, SharedIDs, 10000, 20000, 40000, 50000, tt, eList, true)
	checkList(t, SharedIDs, 10000, 20000, 40000, 50000, tt, eList, false)
}

var tryOkSet = []idCheck{
	{`Trice0i(id(   59), "tt" )`, "id(   59)", 59, true, TriceFmt{"Trice0i", "tt"}},
	{`Trice0(Id(   59), "tt" )`, "Id(   59)", 59, true, TriceFmt{"Trice0", "tt"}}, // should be false
	{`Trice0(id(   59), "tt" )`, "id(   59)", 59, true, TriceFmt{"Trice0", "tt"}},
	{`TRICE0(Id(   59), "tt" )`, "Id(   59)", 59, true, TriceFmt{"TRICE0", "tt"}},
	{`TRICE0(id(   59), "tt" )`, "id(   59)", 59, true, TriceFmt{"TRICE0", "tt"}}, // should be false
	{`TRICE0(Id(59   ), "tt" )`, "Id(59   )", 59, true, TriceFmt{"TRICE0", "tt"}},
	{`TRICE0(Id(59), "tt" )`, "Id(59)", 59, true, TriceFmt{"TRICE0", "tt"}},
	{`TRICE0(Id( 59 ), "tt" )`, "Id( 59 )", 59, true, TriceFmt{"TRICE0", "tt"}},
	{`trice0(Id(59), "tt" )`, "Id(59)", 59, true, TriceFmt{"trice0", "tt"}},
	{`trice64_2(Id(59), "%d,%x", -3, -4 )`, "Id(59)", 59, true, TriceFmt{"trice64_2", "%d,%x"}},
}

func TestTriceParseOK(t *testing.T) {
	set := tryOkSet
	for i := range set {
		nbID, id, tf, ok := triceParse(set[i].nbTRICE)
		assert.True(t, ok == set[i].ok)
		if ok {
			checkID(t, set, i, id)
			checkNbID(t, set, i, nbID)
			assert.True(t, tf == tryOkSet[i].tf, fmt.Sprint(tf, tryOkSet[i].tf))
		}
	}
}

func TestTriceIDParseOK(t *testing.T) {
	set := tryOkSet
	for i := range set {
		nbID, id, ok := triceIDParse(set[i].nbTRICE)
		assert.True(t, ok == set[i].ok)
		if ok {
			checkID(t, set, i, id)
			checkNbID(t, set, i, nbID)
		}
	}
}

func TestTriceFmtParse(t *testing.T) {
	for i := range tryOkSet {
		tf, ok := triceFmtParse(tryOkSet[i].nbTRICE)
		assert.True(t, ok == tryOkSet[i].ok)
		assert.True(t, tf == tryOkSet[i].tf, fmt.Sprint(tf, tryOkSet[i].tf))
	}
}

func TestTriceIDParseNotOK(t *testing.T) {
	var tryNotOkSetID = []idCheck{
		{`TRICE0(Id(0x5), "tt" )`, "Id(0x5)", 5, false, TriceFmt{"TRICE0", "tt"}},
		{`TRICE0(id(0x5 ), "tt" )`, "id(0x59)", 0x59, false, TriceFmt{"TRICE0", "tt"}},
	}
	set := tryNotOkSetID
	for i := range set {
		nbID, id, ok := triceIDParse(set[i].nbTRICE)
		assert.True(t, ok == set[i].ok, fmt.Sprint(i))
		if ok {
			checkID(t, set, i, id)
			checkNbID(t, set, i, nbID)
		}
	}
}

func checkID(t *testing.T, set []idCheck, i int, id TriceID) {
	assert.True(t, id == set[i].id, fmt.Sprint(i, id))
}

func checkNbID(t *testing.T, set []idCheck, i int, nbID string) {
	assert.True(t, nbID == set[i].nbID, fmt.Sprint(i, nbID))
}

type testTable []struct {
	text    string // file content as input
	exp     string // file content as expected output
	fileMod bool   // expected file modification flag
	listMod bool   // expected list modification flag
}

func checkList(t *testing.T, sharedIDs bool, mins, maxs, min, max TriceID, tt testTable, eList string, extend bool) {
	lu := make(TriceIDLookUp)
	tflu := lu.reverse()
	Verbose = true
	for _, x := range tt {
		act0, _ := updateParamCountAndID0(x.text, extend)
		listModified := false
		act, fileModified := updateIDsUniqOrShared(sharedIDs, mins, maxs, min, max, act0, lu, tflu, &listModified)
		assert.Equal(t, x.fileMod, fileModified)
		assert.Equal(t, x.listMod, listModified)
		assert.Equal(t, x.exp, act)
	}
	lu.AddFmtCount()
	aListN := fmt.Sprintln(lu)
	assert.Equal(t, eList, aListN)
}
