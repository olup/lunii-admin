package lunii

type Job struct {
	IsComplete  bool `json:"isComplete"`
	TotalImages int  `json:"totalImages"`
	TotalAudios int  `json:"totalAudios"`
	ImagesDone  int  `json:"imagesDone"`
	AudiosDone  int  `json:"audiosDone"`

	HasError string `json:"hasError"`

	InitDone             bool `json:"initDone"`
	BinGenerationDone    bool `json:"binGenerationDone"`
	UnpackDone           bool `json:"unpackDone"`
	ImagesConversionDone bool `json:"imagesConversionDone"`
	AudiosConversionDone bool `json:"audiosConversionDone"`
	MetadataDone         bool `json:"metadataDone"`
	CopyingDone          bool `json:"copyingDone"`
	IndexDone            bool `json:"indexDone"`
}

var CurrentJob *Job
