// Copyright 2011 James A. Morrison

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gedcom

import (
  "bytes"
  "container/vector"
  "log"
  "strconv"
  "strings"
  "os"
  "unicode"
)

type Person struct {
  name string
  dob []int  // Assume all dates can be converted to and stored in the
             // gregorian calendar.
}

type Family struct {
  id int
}

type Header struct {
  source string
}

type Trailer struct {
  bs string
}

type GedCom struct {
  persons []Person
  families []Family
  header Header
  trailer Trailer
}

// section_name is always non-empty.
type RawGedCom struct {
  section_name string
  value string
  reference string
  data vector.Vector  // Contains RawGedCom objects
}

func WhiteSpaceOrBom(rune int) bool {
  bom := []unicode.Range{unicode.Range{0xFEFF, 0xFEFF, 1}}
  return unicode.IsSpace(rune) || unicode.Is(bom, rune)
}

func (gedcom *RawGedCom) parse_data(line string, previous_level *int) bool {
  trimed_contents := strings.TrimFunc(line, WhiteSpaceOrBom)
  if len(trimed_contents) == 0 {
    return true
  }

  contents := strings.Split(trimed_contents, " ", 3)
  if len(contents) < 2 {
    log.Print("Not enough columns for: ", line)
    return false
  }
  level, err := strconv.Atoi(contents[0])
  if err != nil {
    log.Print("Can't convert ", contents[0], ":", err)
    return false
  }

  if level > *previous_level + 1 {
    log.Print("Level (", level, ") is too much larger than ", *previous_level)
    return false;
  }
  *previous_level = level

  var current_gedcom *RawGedCom = gedcom
  for i := 0; i < level; i++ {
    if current_gedcom.data.Len() == 0 {
      current_gedcom.data.Push(new(RawGedCom))
    }
    current_gedcom = current_gedcom.data.Last().(*RawGedCom)
  }
  current_gedcom.data.Push(new(RawGedCom))
  current_gedcom = current_gedcom.data.Last().(*RawGedCom)

  if contents[1][0] == '@' {
    if len(contents) < 3 {
      return false
    }
    current_gedcom.reference = contents[1]
    current_gedcom.section_name = contents[2]
  } else {
    current_gedcom.section_name = contents[1]
    if len(contents) > 2 {
      current_gedcom.value = contents[2]
    }
  }
  return true
}

func (gedcom *RawGedCom) Parse(contents *bytes.Buffer) bool {
  var line string
  var err os.Error = nil
  previous_level := 0
  for err == nil {
    line, err = contents.ReadString('\n')
    good := gedcom.parse_data(line, &previous_level)
    if !good {
      log.Print("Trouble parsing ", line)
    }
  }

  return true 
}
