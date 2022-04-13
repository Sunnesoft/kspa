package kspa

func EntitiesToDexes(entities []*Entity) (dexes []*DexBase) {
	calcTokenIndeces(entities)

	dexes = make([]*DexBase, len(entities))

	for i, e := range entities {
		dexes[i] = EntityToDex(i, e)
	}
	return
}

func EntityToDex(id int, e *Entity) *DexBase {
	return NewDexBase(id, e.EntityId, e.Id1i, e.Id2i)
}

func calcTokenIndeces(entities []*Entity) {
	index := make(map[int]int)

	for _, v := range entities {
		index[v.Id1] = -1
		index[v.Id2] = -1
	}

	j := 0
	for i, v := range entities {
		if index[v.Id1] == -1 {
			index[v.Id1] = j
			j++
		}
		entities[i].Id1i = index[v.Id1]

		if index[v.Id2] == -1 {
			index[v.Id2] = j
			j++
		}
		entities[i].Id2i = index[v.Id2]
	}
}
