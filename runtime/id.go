package runtime

type Id uint32

var runtimeIds = map[*Runtime]Id{}
var runtimes []*Runtime

func (rt *Runtime) Id() Id {
	id, ok := runtimeIds[rt]
	if !ok {
		id = Id(len(runtimes))
		runtimeIds[rt] = id
		runtimes = append(runtimes, rt)
	}
	return id
}

func GetById(id Id) *Runtime {
	return runtimes[id]
}
