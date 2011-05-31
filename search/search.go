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
  "appengine"; "appengine/blobstore"; "appengine/taskqueue"; "appengine/user"
  "bytes"
  "fmt"
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
  http.HandleFunc("/process/gedcom", GedcomHandler)
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
  data["Base_person"] =
      "{\"name\": \"James Morrison\", \"Date of Birth\": [1981, 10, 2]}"
  upload_url, _ :=  blobstore.UploadURL(c, "/upload")
  data["Upload_Action"] = upload_url.String()

  template_err := searchTemplate.Execute(w, data)
  if template_err != nil {
    log.Print("Error rendering template ", template_err)
  }
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  u := user.Current(c)
  if u == nil {
    url, _ := user.LoginURL(c, r.URL.String())
    w.Header().Set("Location", url)
    w.WriteHeader(http.StatusFound)
    return
  }

  id := r.FormValue("id")
  if len(id) > 0 {
    w.Header().Set("Location", "/upload2?id=has_key:" + id)
    w.WriteHeader(http.StatusFound)
//    uploadTemplate.Execute(w, id)
    return
  }

  blobs, other_params, err := blobstore.ParseUpload(r)
  if len(blobs) == 0 {
//    w.WriteHeader(http.StatusBadRequest)
//    fmt.Fprintf(w, "No data '%v'", err)
    w.Header().Set("Location", "/upload2?id=Bad+upload:" + err.String())
    w.WriteHeader(http.StatusFound)
    return
  }
  file := blobs["file_data"]
  if len(file) == 0 {
//    w.WriteHeader(http.StatusBadRequest)
//    fmt.Fprintf(w, "No data")
    w.Header().Set("Location", "/upload2?id=No_file_data")
    w.WriteHeader(http.StatusFound)
    return
  }

  key := string(file[0].BlobKey)
  if other_params == nil {
    other_params = make(map[string] []string)
  }
  other_params["key"] = append(other_params["key"], key)
  task := taskqueue.NewPOSTTask("/process/gedcom", other_params)
  task.Name = key
  if err := taskqueue.Add(c, task, ""); err != nil {
//    http.Error(w, err.String(), http.StatusInternalServerError)
    w.Header().Set("Location", "/upload2?id=bad_task:" + err.String())
    w.WriteHeader(http.StatusFound)
    return
  }

  w.Header().Set("Location", "/upload?id=" + key)
  w.WriteHeader(http.StatusFound)
  return
}

func GedcomHandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  key := appengine.BlobKey(r.FormValue("key"))

  buffer := new(bytes.Buffer)
  buffer.ReadFrom(blobstore.NewReader(c, key))
  var raw_data RawGedCom
  if !raw_data.Parse(buffer) {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, "Bad data")
    return
  }
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "ok")
}
