package checksum

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Request struct {
	Path             string `form:"Path" json:"Path" query:"Path"`
	BlockSize        int32  `form:"BlockSize" json:"BlockSize" query:"BlockSize"`
	BeginChunkOffset int32  `form:"BeginChunkOffset" json:"BeginChunkOffset" query:"BeginChunkOffset"`
	EndChunkOffset   int32  `form:"EndChunkOffset" json:"EndChunkOffset" query:"EndChunkOffset"`
	RollingHashType  int32  `form:"RollingHashType" json:"RollingHashType" query:"RollingHashType"`
}

type Response struct {
	v20230530.Response `json:",inline"`
	Data               *Data `json:"Data,omitempty"`
}

type Data struct {
	Checksums []*ChunkChecksum `json:"Checksums"`
}

type ChunkChecksum struct {
	// an offset in terms of chunk count
	ChunkOffset uint `json:"ChunkOffset"`
	// the size of the block
	BlockSize      int64  `json:"Size"`
	WeakChecksum   []byte `json:"WeakChecksum"`
	StrongChecksum []byte `json:"StrongChecksum"`
}
