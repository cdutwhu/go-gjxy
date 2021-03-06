package gjxy

// XMLTag : get the xml string's first tag
func XMLTag(xml string) string {
	XML := Str(xml).T(BLANK)
	pc(XML == "" || XML.C(0) != '<' || XML.C(LAST) != '>', fEf("Not a valid XML section"))
	return XML.S(XML.LIdx("</")+2, ALL-1).V()
}

// XMLTagEle : Looking for the first tag's xml string
func XMLTagEle(xml, tag string) (string, int, int) {
	XML := Str(xml).T(BLANK)
	pc(XML == "" || XML.C(0) != '<' || XML.C(LAST) != '>', fEf("Invalid XML"))

	XML = Str(xml) //                                                   *** we have to return original position, so use original xml ***
	s, sa := XML.Idx(fSf("<%s>", tag)), XML.Idx(fSf("<%s ", tag))
	if s < 0 && sa >= 0 {
		s = sa
	} else if s >= 0 && sa < 0 {
		// s = s
	} else if s >= 0 && sa >= 0 {
		min, _ := Min(I32s{s, sa}, "")
		s = min.(int)
	} else if s < 0 && sa < 0 {
		return "", -1, -1
	}

	if eR := XML.S(s, ALL).Idx(fSf("</%s>", tag)); eR > 0 {
		sNext := s + eR + Str(tag).L() + 3
		return XML.S(s, sNext).V(), s, sNext
	}

	pc(true, fEf("Invalid XML"))
	return "", -1, -1
}

// XMLTagEleEx : idx from 1
func XMLTagEleEx(xml, tag string, idx int) (string, int) {
	const LIMIT = 4096
	pc(idx > LIMIT || idx < 1, fEf("idx starts from 1, to %d", LIMIT))
	esum, rst := 0, ""
	for i := 1; i <= LIMIT; i++ { //              *** set a limit for searching ***
		XML := Str(xml).S(esum, ALL)
		if XML.V() == "" { //                     *** to the end, return ***
			return rst, i - 1
		}
		ele, _, e := XMLTagEle(XML.V(), tag)
		if e == -1 { //                           *** could not find, return ***
			return rst, i - 1
		}
		if i == idx { //                          *** find, return ***
			rst = ele
		}
		esum += e
	}
	pc(true, fEf("Should NOT Be Here!"))
	return "", -1
}

// XMLXPathEle :
func XMLXPathEle(xml, xpath, del string, indices ...int) (ele string, nArr int) {
	pc(xpath == "", fEf("At least one path must be provided"))

	segs := sSpl(xpath, del)
	pc(len(segs) != len(indices), fEf("path & seg's index count not match"))

	for i, seg := range segs {
		xml = IF(ele != "", ele, xml).(string)
		ele, nArr = XMLTagEleEx(xml, seg, indices[i])
	}
	return
}

// XMLAttributes is (ONLY LIKE  <SchoolInfo RefId="D3F5B90C-D85D-4728-8C6F-0D606070606C" Type="LGL">)
func XMLAttributes(xmlele string) (attributes, attriValues []string) { //       *** 'map' may cause mis-order, so use slice
	XMLELE := Str(xmlele).T(BLANK)
	pc(XMLELE == "" || XMLELE.C(0) != '<' || XMLELE.C(LAST) != '>', fEf("Not a valid XML section"))

	tag := Str(XMLTag(xmlele))
	if eol := XMLELE.Idx(`">`) + 1; XMLELE.C(tag.L()+1) == ' ' && eol > tag.L() { //    *** has attributes

		focus := XMLELE.S(tag.L()+2, eol)
		nAttri := sCnt(focus.V(), "=\"")
		// fPln(nAttri)

		for i := 1; i <= nAttri; i++ {
			av, _, _ := focus.QuotesPos(QDouble, i)
			attriValues = append(attriValues, av.RmQuotes(QDouble).V())
		}

	OUT:
		for i := 1; i <= nAttri; i++ {
			_, left, _ := focus.QuotesPos(QDouble, i)
			for p := left - 2; p >= 0; p-- {
				if p == 0 {
					attributes = append(attributes, focus.S(0, left-1).V())
					continue OUT
				}
				if focus.C(p) == ' ' {
					attributes = append(attributes, focus.S(p+1, left-1).V())
					continue OUT
				}
			}
		}
	}
	return attributes, attriValues
}

// XMLChildren : (NOT search grandchildren)
func XMLChildren(xmlele string, fNArr bool) (children []string) {
	XMLELE := Str(xmlele).T(BLANK)
	pc(XMLELE == "" || XMLELE.C(0) != '<' || XMLELE.C(LAST) != '>', fEf("Invalid XML section"))

	L := XMLELE.L()
	skip, childpos, level, inflag := false, []int{}, 0, false

	for i := 0; i < L; i++ {
		c := XMLELE.C(i)

		if c == '<' && XMLELE.S(i, i+4) == "<!--" {
			skip = true
		}
		if c == '>' && XMLELE.S(i-2, i+1) == "-->" {
			skip = false
		}
		if skip {
			continue
		}

		if c == '<' && XMLELE.C(i+1) != '/' {
			level++
		}
		if c == '<' && XMLELE.C(i+1) == '/' {
			level--
			if level == 1 {
				inflag = false
			}
		}

		if level == 2 {
			if !inflag {
				childpos, inflag = append(childpos, i+1), true
			}
		}
	}

	if len(childpos) == 0 {
		return
	}

	for _, p := range childpos {
		pe, peA := XMLELE.S(p, ALL).Idx(">"), XMLELE.S(p, ALL).Idx(" ")
		pe = IF(peA > 0 && peA < pe, peA, pe).(int)
		child := XMLELE.S(p, p+pe)
		children = append(children, child.V())
	}

	children = IArrFoldRep(Strs(children), IF(fNArr, "[n]", "[]").(string)).([]string)
	return
}

// XMLFamilyTree :
func XMLFamilyTree(xml, fName, del string, mapFT *map[string][]string) {
	pc(mapFT == nil, fEf("FamilyTree return map is not initialised !"))
	XML := Str(xml).T(BLANK)
	pc(XML == "" || XML.C(0) != '<' || XML.C(LAST) != '>', fEf("Invalid XML section"))

	fName = IF(fName == "", XMLTag(xml), fName).(string)
	if children := XMLChildren(xml, false); len(children) > 0 {
		// fPln(tag, children)

		(*mapFT)[fName] = children //                           *** record path ***

		for _, child := range children {
			if Str(child).HP("[") {
				child = Str(child).RmHeadToLast("]").V() //     *** remove array symbol ***
			}
			nextPath := Str(fName + del + child).T(del).V()
			subxml, _ := XMLTagEleEx(xml, child, 1)
			XMLFamilyTree(subxml, nextPath, del, mapFT)
		}
	}
}

// XMLCntByIPath : dump all leaves, non-array leaves can be removed later.
func XMLCntByIPath(xml, iPath, del string, mapFT *map[string][]string) (arrNames []string, arrCnts []int, nextIPaths []string) {
	path, indices := IPathToPathIndices(iPath, del) //         *** defined in <json.go> ***
	// fPln("indices:", indices)

	leaves := (*mapFT)[path]
	for _, leaf := range leaves {
		// fPln(leaf)
		LEAF := Str(leaf)

		arrName := LEAF.V()
		if LEAF.HP("[]") {
			arrName = LEAF.S(2, ALL).V()
		}
		_, nArr := XMLXPathEle(xml, path+del+arrName, del, append(indices, 1)...)
		// fPln(nArr)

		arrNames = append(arrNames, arrName)
		arrCnts = append(arrCnts, nArr)

		for i := 1; i <= nArr; i++ {
			nextIPath := iPath + del + arrName + fSf("#%d", i)
			// fPln(nextIPath)
			nextIPaths = append(nextIPaths, nextIPath)
		}
	}
	return
}

// XMLWholeCntByIPathByR :
func XMLWholeCntByIPathByR(xml, iPath, del, id string, mapFT *map[string][]string, mapIPathNID *map[string]struct {
	Count int
	ID    string
}) {
	pc(mapIPathNID == nil, fEf("result <mapIPathNID> is not initialized"))
	arrNames, arrCnts, subIPaths := XMLCntByIPath(xml, iPath, del, mapFT)
	pc(len(arrNames) != len(arrCnts), fEf("error in XMLCntByIPath"))

	nNames := len(arrNames)
	for i := 0; i < nNames; i++ {
		(*mapIPathNID)[iPath+del+arrNames[i]] = struct {
			Count int
			ID    string
		}{
			Count: arrCnts[i],
			ID:    id,
		}
	}
	for _, subIPath := range subIPaths {
		XMLWholeCntByIPathByR(xml, subIPath, del, id, mapFT, mapIPathNID)
	}
}

// XMLCntInfo :
func XMLCntInfo(xml, xpath, del, id string, mapFT *map[string][]string) (*map[string][]string, *map[string]struct {
	Count int
	ID    string
}) {
	if mapFT == nil {
		mapFT = &map[string][]string{}
		XMLFamilyTree(xml, xpath, del, mapFT)
	}

	// fPln(" ------------------------------------------ ")

	root := XMLTag(xml)
	// fPf("ROOT is <%s>\n", root)
	pc(root == "", fEf("Invalid path"))

	iRoot := root + "#1"
	mapIPathNID := &map[string]struct {
		Count int
		ID    string
	}{}
	XMLWholeCntByIPathByR(xml, iRoot, del, id, mapFT, mapIPathNID)
	return mapFT, mapIPathNID
}

/**********************************************************************************************************************************/

// XMLSegPos : level from 1, index from 1                                         &
func XMLSegPos(xml string, level, index int) (tag, str string, left, right int) {
	s := Str(xml)
	markS, markE1, markE2, markE3 := '<', '<', '/', '>'
	curLevel, curIndex, To := 0, 0, s.L()-1

	found := false
	i := 0
	for _, c := range s {
		if i < To {
			curLevel = IF(c == markS && s.C(i+1) != markE2, curLevel+1, curLevel).(int)
			curLevel = IF(c == markE1 && s.C(i+1) == markE2, curLevel-1, curLevel).(int)
			if curLevel == level && c == markS && s.C(i+1) != markE2 {
				left = i
			}
			if curLevel == level-1 && c == markE1 && s.C(i+1) == markE2 {
				right = i
				curIndex++
				if curIndex == index {
					found = true
					break
				}
			}
		}
		i++
	}

	if !found {
		return "", "", 0, 0
	}

	tagendRel := s.S(left+1, right).Idx(" ") // when tag has attribute(s)
	if tagendRel == -1 {
		tagendRel = s.S(left+1, right).Idx(string(markE3))
	}
	pc(tagendRel == -1, fEf("xml error"))

	tag = s.S(left+1, left+1+tagendRel).V()
	right += Str(tag).L() + 2
	return tag, s.S(left, right+1).V(), left, right
}

// XMLSegsCount : only count top level                                            &
func XMLSegsCount(xml string) (count int) {
	s := Str(xml)
	markS, markE1, markE2 := '<', '<', '/'

	level, inflag, To := 0, false, s.L()-1
	i := 0
	for _, c := range s {
		if i < To {
			if c == markS && s.C(i+1) != markE2 {
				level++
			}
			if c == markE1 && s.C(i+1) == markE2 {
				level--
				if level == 0 {
					inflag = false
				}
			}
			if level == 1 {
				if !inflag {
					count++
					inflag = true
				}
			}
		}
		i++
	}
	return count
}
