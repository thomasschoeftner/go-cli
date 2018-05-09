package tasks

type ProcessingConf struct {
	Type string //"sequential", "concurrent"
	Concurrent *ConcurrentProcessingConfig
}

type ConcurrentProcessingConfig struct {
	BufferSize int
}

