package main

//The Task: Given _N_ sorted slices of _K_ songs, implement a function that outputs the merged slice of sorted songs
//Hints:
// * The container package will be useful in this challenge.
// * Remember that the album slices are sorted.

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

const path = "songs.json"

// Playlist is a list of songs sorted by PlayCount that implements heap.Interface (based on https://pkg.go.dev/container/heap@go1.27.1#pkg-overview example)
type Playlist []*Song

func (pl Playlist) Len() int { return len(pl) }

func (pl Playlist) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, PlayCount so we use greater than here.
	return pl[i].PlayCount > pl[j].PlayCount
}

func (pl Playlist) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
	pl[i].index = i
	pl[j].index = j
}

func (pl *Playlist) Push(x any) {
	n := len(*pl)
	item := x.(*Song)
	item.index = n
	*pl = append(*pl, item)
}

func (pl *Playlist) Pop() any {
	old := *pl
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pl = old[0 : n-1]
	return item
}

// Song stores all the song related information
type Song struct {
	Name      string `json:"name"`
	Album     string `json:"album"`
	PlayCount int64  `json:"play_count"`
	index     int
}

// makePlaylist makes the merged sorted list of songs
func makePlaylist(albums [][]Song) []Song {
	/*The plan:
	* Implement a priority queue that contains the methods of the "heap" interface.  The song's play_count will represent its priority
	* For each album:
	* 	Push all songs into the queue
	* After queue has been populated, pop songs off into a slice and return it
	 */
	var pl Playlist
	var list []Song
	idx := 0

	//See note regarding the usage of heap.Push/.Pop (https://pkg.go.dev/container/heap@go1.27.1#Interface)
	heap.Init(&pl)
	for _, a := range albums {
		for _, s := range a {
			s.index = idx
			heap.Push(&pl, &s)
			idx++
		}
	}

	for range pl {
		s := heap.Pop(&pl).(*Song)
		list = append(list, *s)
	}
	return list
}

func main() {
	albums := importData()
	printTable(makePlaylist(albums))
}

// printTable prints merged playlist as a table
func printTable(songs []Song) {
	w := tabwriter.NewWriter(os.Stdout, 3, 3, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "####\tSong\tAlbum\tPlay count")
	for i, s := range songs {
		fmt.Fprintf(w, "[%d]:\t%s\t%s\t%d\n", i+1, s.Name, s.Album, s.PlayCount)
	}
	w.Flush()

}

// importData reads the input data from file and creates the friends map
func importData() [][]Song {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var data [][]Song
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
