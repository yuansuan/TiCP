package checksumsFindChunks

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Request struct {
	Path            string           `json:"Path"`
	Offset          int64            `json:"Offset"`
	Length          int64            `json:"Length"`
	BlockSize       int32            `json:"BlockSize"`
	RollingHashType int32            `json:"RollingHashType"`
	Checksums       []*ChunkChecksum `json:"Checksums,omitempty"`
}

type Response struct {
	v20230530.Response `json:",inline"`
	Data               *Data `json:"Data,omitempty"`
}

type ChunkChecksum struct {
	// an offset in terms of chunk count
	ChunkOffset uint `json:"ChunkOffset"`
	// the size of the block
	BlockSize      int64  `json:"Size"`
	WeakChecksum   []byte `json:"WeakChecksum"`
	StrongChecksum []byte `json:"StrongChecksum"`
	//
	Id      string `json:"Id"`
	RoundId int64  `json:"RoundId"`
}

type Data struct {
	MatchChunks []*Chunk `json:"MatchChunks"`
}

type Chunk struct {
	// to find source chunk of dir file
	Id       string `json:"Id"`
	RoundId  int64  `json:"RoundId"` // 客户端 chunk 的时间标签
	Priority int    `json:"Priority"`
	// mark location in source file
	Offset         int64  `json:"Offset"`
	Length         int64  `json:"Length"`
	WeakChecksum   []byte `json:"WeakChecksum"`
	StrongChecksum []byte `json:"StrongChecksum"`
}
