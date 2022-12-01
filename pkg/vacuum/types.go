package vacuum

type OperationStats struct {
	CopiedFiles uint
	CopiedBytes uint
	Errors      []*error
}

func (o *OperationStats) Add(co *OperationStats) {
	o.CopiedFiles += co.CopiedFiles
	o.CopiedBytes += co.CopiedBytes
	o.Errors += co.Errors
}
