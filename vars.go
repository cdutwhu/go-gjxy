package gjxy

import (
	"fmt"
	"sort"
	"strings"

	u "github.com/cdutwhu/go-util"
	w "github.com/cdutwhu/go-wrappers"
)

type (
	Str  = w.Str
	Strs = w.Strs
	C32s = w.C32s
	I32s = w.I32s
)

const (
	BRound  = w.BRound
	BCurly  = w.BCurly
	BBox    = w.BBox
	QDouble = w.QDouble
	LAST    = w.LAST
	ALL     = w.ALL
)

var (
	IF          = u.IF
	pc          = u.PanicOnCondition
	pe          = u.PanicOnError
	pe1         = u.PanicOnError1
	ph          = u.PanicHandle
	must        = u.Must
	matchAssign = u.MatchAssign
	trueAssign  = u.TrueAssign
	MapKeys     = u.MapKeys
	XIn         = u.XIn

	IArrSearchOne = w.IArrSearchOne
	IArrIsSameEle = w.IArrIsSameEle
	IArrFoldRep   = w.IArrFoldRep
	Min           = w.Min

	sortByLess = sort.Sort

	fEf  = fmt.Errorf
	fSf  = fmt.Sprintf
	fPln = fmt.Println
	fPf  = fmt.Printf

	sCnt = strings.Count
	sSpl = strings.Split
	sJ   = strings.Join
	sRep = strings.Replace
	sFF  = strings.FieldsFunc
)
