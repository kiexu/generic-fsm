package gfsm

type (
	// EdgeCollection fast query supported
	EdgeCollection[T, S comparable, U, V any] struct {
		eList []*Edge[T, S, U, V]       // Regular edge list
		eFast map[S][]*Edge[T, S, U, V] // Redundancy edge map for O(1) query
	}

	// Edge Event value included
	Edge[T, S comparable, U, V any] struct {
		fromV    *Vertex[T, V] // From vertex
		toV      *Vertex[T, V] // To vertex
		eventVal S             // Event value. Not unique. {stateVal, eventVal} can be not unique but not recommended
		storeVal U             // Anything you want. e.g. Real callback function(use CallBacks to invoke)
	}
)

// EdgeCollection

// addE add an edge to EdgeCollection
func (c *EdgeCollection[T, S, U, V]) addE(e *Edge[T, S, U, V]) {
	if e == nil {
		return
	}
	c.eList = append(c.eList, e)
	c.eFast[e.eventVal] = append(c.eFast[e.eventVal], e)
}

// EdgeByEventVal get eventE value by eventE value
func (c *EdgeCollection[T, S, U, V]) EdgeByEventVal(eventVal S) []*Edge[T, S, U, V] {
	return c.eFast[eventVal]
}

// Edge

func (e *Edge[T, S, U, V]) FromV() *Vertex[T, V] {
	return e.fromV
}

func (e *Edge[T, S, U, V]) SetFromV(fromV *Vertex[T, V]) {
	e.fromV = fromV
}

func (e *Edge[T, S, U, V]) ToV() *Vertex[T, V] {
	return e.toV
}

func (e *Edge[T, S, U, V]) SetToV(toV *Vertex[T, V]) {
	e.toV = toV
}

func (e *Edge[T, S, U, V]) EventVal() S {
	return e.eventVal
}

func (e *Edge[T, S, U, V]) SetEventVal(eventVal S) {
	e.eventVal = eventVal
}

func (e *Edge[T, S, U, V]) StoreVal() U {
	return e.storeVal
}

func (e *Edge[T, S, U, V]) SetStoreVal(storeVal U) {
	e.storeVal = storeVal
}
