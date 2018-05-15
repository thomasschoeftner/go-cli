package task

type ProcessingConfig struct {
	Type string //"sequential", "concurrent"
	Sequential *SequentialProcessingConfig
	Concurrent *ConcurrentProcessingConfig
}

type ConcurrentProcessingConfig struct {
	BufferSize int
}

type SequentialProcessingConfig struct {
	StopAtError bool
}
