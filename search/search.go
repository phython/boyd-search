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

// The package with the appengine http handlers.

package search

import (
  "appengine"; "appengine/user"
//  "./gedcom"
  "http"
  "log"
  "template"
)

var (
  searchTemplate *template.Template
  uploadTemplate *template.Template
)

func create_template(file string) *template.Template {
  new_template := template.New(nil)
  new_template.SetDelims("{{", "}}")
  err := new_template.ParseFile("canvas.html")
  if err != nil {
    panic(file + " is not parseable" + err.String())
  }

  return new_template
}

func init() {
  http.HandleFunc("/", SearchHandler)
  http.HandleFunc("/upload", UploadHandler)
  searchTemplate = create_template("canvas.html")
  uploadTemplate = create_template("upload.html")
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
  data := make(map[string] string)
  c := appengine.NewContext(r)
  u := user.Current(c)
  if u != nil {
    data["Email"] = u.String()
  } else {
    url, err := user.LoginURL(c, r.URL.String())
    if err != nil {
      http.Error(w, err.String(), http.StatusInternalServerError)
      return
    }
    data["Login_url"] = url
  }
  data["Base_person"] = "{\"name\": \"James Morrison\", \"Date of Birth\": [1981, 10, 2]}"

  template_err := searchTemplate.Execute(w, data)
  if template_err != nil {
    log.Print("Error rendering template ", template_err)
  }
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
  id := r.FormValue("id")
  if len(id) > 0 {
    uploadTemplate.Execute(w, id)
    return
  }

  w.Header().Set("Location", r.URL.Path + "?id=")
  w.WriteHeader(http.StatusFound)
}
