// Exercise 7.9: Use the html/template package (§4.6) to replace printTracks with a function
// that displays the tracks as an HTML table. Use the solution to the previous exercise to arrange
// that each click on a column head makes an HTTP request to sort the table.

package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Arcterix", "From the Roots Up", 2010, length("3m38s")},
	{"Go", "Arcterix", "From the Roots Up", 1990, length("3m38s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

var tracklist = template.Must(template.New("tracklist").Parse(`
<h1>{{len .}} tracks</h1>
<table>
<tr style='text-align: left'>
  <th><a href='/?sort=title'>Title</a></th>
  <th><a href='/?sort=artist'>Artist</a></th>
  <th><a href='/?sort=album'>Album</a></th>
  <th><a href='/?sort=year'>Year<a/></th>
  <th><a href='/?sort=length'>Length<a/></th>
</tr>
{{range .}}
<tr>
  <td>{{.Title}}</td>
  <td>{{.Artist}}</td>
  <td>{{.Album}}</td>
  <td>{{.Year}}</td>
  <td>{{.Length}}</td>
</tr>
{{end}}
</table>
`))

type term func(x, y *Track) bool
type statefullByTerms struct {
	tracks []*Track
	terms  []term
}

func (x *statefullByTerms) AddTerms(t ...term) { x.terms = append(x.terms, t...) }
func (x statefullByTerms) Len() int            { return len(x.tracks) }
func (x statefullByTerms) Less(i, j int) bool {
	a, b := x.tracks[i], x.tracks[j]

	for _, t := range x.terms {
		if t(a, b) {
			return true
		} else if t(b, a) {
			return false
		}
	}
	return false
}
func (x statefullByTerms) Swap(i, j int) { x.tracks[i], x.tracks[j] = x.tracks[j], x.tracks[i] }

func sortTracks(terms []string) {
	byColumns := statefullByTerms{tracks, make([]term, 0)}

	for _, term := range terms {
		switch term {
		case "title":
			byColumns.AddTerms(func(x, y *Track) bool {
				return x.Title < y.Title
			})
		case "artist":
			byColumns.AddTerms(func(x, y *Track) bool {
				return x.Artist < y.Artist
			})
		case "album":
			byColumns.AddTerms(func(x, y *Track) bool {
				return x.Album < y.Album
			})
		case "year":
			byColumns.AddTerms(func(x, y *Track) bool {
				return x.Year < y.Year
			})
		case "length":
			byColumns.AddTerms(func(x, y *Track) bool {
				return x.Length.Nanoseconds() < y.Length.Nanoseconds()
			})
		default:
			log.Printf("unknown term provided for sort: %s", term)
		}
	}
	sort.Sort(byColumns)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if u, err := url.ParseQuery(r.URL.RawQuery); err != nil {
			log.Printf("failed parsing query %s: %s", r.URL.String(), err.Error())
		} else {
			if u.Has("sort") {
				sortTracks(u["sort"])
			}
		}

		if err := tracklist.Execute(w, tracks); err != nil {
			log.Fatal(err)
		}
	})
	http.ListenAndServe(":8080", nil)
}
