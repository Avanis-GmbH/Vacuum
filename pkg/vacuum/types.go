package vacuum

type OperationStats struct {
	CopiedFiles  uint64
	CopiedBytes  uint64
	DeletedFiles uint64
	Errors       []*error
}
