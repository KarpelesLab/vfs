package zipfs

type Index int

const (
	NoIndex   Index = iota // no index, for small files with low access count
	IndexMap               // map index, quick for opening specific files but not adapted for listing/etc
	IndexFull              // full tree index
)
