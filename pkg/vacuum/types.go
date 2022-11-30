package vacuum

type OperationStats struct {
	CopiedFiles uint
	CopiedBytes uint
}

func (o *OperationStats) Add(co *OperationStats) {
	o.CopiedFiles += co.CopiedFiles
	o.CopiedBytes += co.CopiedBytes
}
