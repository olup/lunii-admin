package lunii

type Job struct {
	IsComplete  bool
	TotalImages int
	TotalAudios int
	ImagesDone  int
	AudiosDone  int

	Err string

	InitDone             bool
	BinGenerationDone    bool
	UnpackDone           bool
	ImagesConversionDone bool
	AudiosConversionDone bool
	MetadataDone         bool
	CopyingDone          bool
	IndexDone            bool
}

var CurrentJob *Job
