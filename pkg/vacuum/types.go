package vacuum

type OperationStats struct {
	CopiedFiles uint
	CopiedBytes uint
	Errors      []*error
}
