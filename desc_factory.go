package gfsm

import (
	"github.com/kiexu/go-generic-collection"
	"github.com/kiexu/go-generic-collection/hashset"
)

type (
	// DefGFactory Default factory with basic config struct
	// As a regular FSM, {stateVal, eventVal} need to be unique
	DefGFactory[T, S comparable, U, V any] struct {
		DescList     []*DescCell[T, S, U, V] // Required. Describe FSM graph
		VertexValMap map[T]V                 // Optional. Store custom value in vertex
	}

	// DescCell Describe one eventE
	DescCell[T, S comparable, U, V any] struct {
		EventVal     S
		FromState    []T
		ToState      T
		EdgeStoreVal U // Every edge's EdgeStoreVal in this cell will be assigned this field
	}

	// stateEvent Deduplication helper
	stateEvent[T, S comparable] struct {
		stateVal T
		eventVal S
	}
)

// Ensure interface implement
var _ GraphFactory[struct{}, struct{}, struct{}, struct{}] = new(DefGFactory[struct{}, struct{}, struct{}, struct{}])

// NewG New a Graph
func (fac *DefGFactory[T, S, U, V]) NewG() (*Graph[T, S, U, V], error) {

	g := &Graph[T, S, U, V]{
		stoV: make(map[T]*Vertex[T, V]),
	}

	// Init itoV
	var stateValSet gcollection.Set[T] = hashset.NewHashSet[T]()
	for _, desc := range fac.DescList {
		if ok := stateValSet.Add(desc.ToState); ok {
			g.itoV = append(g.itoV, fac.newV(desc.ToState))
		}
		for _, fs := range desc.FromState {
			if ok := stateValSet.Add(fs); ok {
				g.itoV = append(g.itoV, fac.newV(fs))
			}
		}
	}

	// Init idx and stoV
	// Idx starts with 0
	for i, v := range g.itoV {
		v.idx = i
		g.stoV[v.stateVal] = v
	}

	// initial adj
	vl := len(g.itoV)
	var stateEventSet gcollection.Set[stateEvent[T, S]] = hashset.NewHashSet[stateEvent[T, S]]()
	g.adj = make([]*EdgeCollection[T, S, U, V], vl, vl)
	for _, d := range fac.DescList {
		toIdx := g.VertexByState(d.ToState).idx
		for _, s := range d.FromState {
			fromIdx := g.VertexByState(s).idx
			if g.adj[fromIdx] == nil {
				g.adj[fromIdx] = &EdgeCollection[T, S, U, V]{
					eList: make([]*Edge[T, S, U, V], 0),
					eFast: make(map[S][]*Edge[T, S, U, V]),
				}
			}
			uniqSE := stateEvent[T, S]{
				stateVal: s,
				eventVal: d.EventVal,
			}
			if ok := stateEventSet.Add(uniqSE); !ok {
				return nil, &DuplicateStateAndEventErr[T, S]{State: s, Event: d.EventVal}
			}
			e := &Edge[T, S, U, V]{
				fromV:    g.itoV[fromIdx],
				toV:      g.itoV[toIdx],
				eventVal: d.EventVal,
				storeVal: d.EdgeStoreVal,
			}
			g.adj[fromIdx].addE(e)
		}
	}

	return g, nil
}

// newV Without idx, autofill storeVal
func (fac *DefGFactory[T, S, U, V]) newV(state T) *Vertex[T, V] {
	genV := &Vertex[T, V]{
		stateVal: state,
	}
	if storeVal, ok := fac.VertexValMap[state]; ok {
		genV.storeVal = storeVal
	}
	return genV
}
