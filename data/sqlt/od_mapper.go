package sqlt

import (
	"mngr/data"
	"strconv"
)

type OdMapper struct {
}

func (o *OdMapper) Map(source *OdEntity) *data.OdDto {
	ret := &data.OdDto{}
	ret.Id = strconv.FormatUint(uint64(source.ID), 10)
	ret.GroupId = source.GroupId
	ret.CreatedAt = source.CreatedAtStr
	ret.DetectedObject = &data.DetectedObjectDto{
		PredScore:   source.PredScore,
		PredClsIdx:  source.PredClsIdx,
		PredClsName: source.PredClsName,
	}
	ret.ImageFileName = source.ImageFileName

	ret.VideoFile = &data.VideoFileDto{}
	ret.VideoFile.Name = source.VideoFileName
	ret.VideoFile.CreatedDate = *source.VideoFileCreatedDate
	ret.VideoFile.Duration = source.VideoFileDuration
	ret.VideoFile.Merged = source.VideoFileMerged
	ret.VideoFile.ObjectAppearsAt = source.ObjectAppearsAt

	ret.AiClip = &data.AiClip{
		Enabled:        source.AiClipEnabled,
		FileName:       source.AiClipFileName,
		CreatedAt:      source.CreatedAtStr,
		LastModifiedAt: source.AiClipLastModifiedAtStr,
		Duration:       source.AiClipDuration,
	}

	return ret
}
