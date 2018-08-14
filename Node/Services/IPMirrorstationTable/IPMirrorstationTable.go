package IPMirrorstationtable

import (
	"github.com/gogo/protobuf/proto"
	"kaleido/Common/Services/IPMirrorstationTableMessage"
)

func Unmarshal(data []byte) map[uint32]map[uint32][]string {
	var table IPMirrorStationtableMessages.Table
	result := make(map[uint32]map[uint32][]string)
	proto.Unmarshal(data, &table)
	for _, relation := range table.Relations {
		_, exist := result[relation.MaskBitLength]
		if !exist {
			result[relation.MaskBitLength] = make(map[uint32][]string)
		}
		result[relation.MaskBitLength][relation.Ip] = relation.Url
	}
	return result
}
