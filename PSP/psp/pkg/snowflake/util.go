package snowflake

// BatchParseString 批量解析字符串ID，返回值中第二个数组是解析失败的
func BatchParseString(ids []string) ([]ID, []string) {
	idSize := len(ids)
	if idSize > 0 {
		sids := make([]ID, 0, idSize)
		// failed parse ids
		fids := make([]string, 0, idSize)
		for _, id := range ids {
			sid, err := ParseString(id)
			if err != nil {
				fids = append(fids, id)
				continue
			}

			sids = append(sids, sid)
		}

		return sids, fids
	}

	return []ID{}, []string{}
}

func BatchParseStringToID(ids []string) []ID {
	idList, _ := BatchParseString(ids)
	return idList
}
