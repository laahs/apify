//package storage provides all necessary bits for abstraction of the storage process and query handling.
package storage

import (
	"github.com/laahs/gopkg/apify/query"
)

// ----- Storage Handler = Recorder -----
// ------------------

//interface Recorder is a generic type representing a storage handler handling where the Items are stored, and from where they are
//retrieved or deleted. Also acts as a search engine intermediary.
type Recorder interface {
	Insert(Item) error
	Delete(string) error
	Update(string) error
	Find(string) Item
	FindList(string) []Item
	Search(query.QueryParams) ([]Items, []string)
	Type() string //returns the storage type
}
