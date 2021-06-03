// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2021 Datadog, Inc.

package store

import (
	"sort"

	"github.com/DataDog/sketches-go/ddsketch/pb/sketchpb"
)

type SparseStore struct {
	counts map[int]float64
}

func NewSparseStore() *SparseStore {
	return &SparseStore{counts: make(map[int]float64)}
}

func (s *SparseStore) Add(index int) {
	s.counts[index]++
}

func (s *SparseStore) AddBin(bin Bin) {
	s.AddWithCount(bin.index, bin.count)
}

func (s *SparseStore) AddWithCount(index int, count float64) {
	if count == 0 {
		return
	}
	s.counts[index] += count
}

func (s *SparseStore) Bins() <-chan Bin {
	orderedBins := s.orderedBins()
	ch := make(chan Bin)
	go func() {
		defer close(ch)
		for _, bin := range orderedBins {
			ch <- bin
		}
	}()
	return ch
}

func (s *SparseStore) orderedBins() []Bin {
	bins := make([]Bin, 0, len(s.counts))
	for index, count := range s.counts {
		bins = append(bins, Bin{index: index, count: count})
	}
	sort.Slice(bins, func(i, j int) bool { return bins[i].index < bins[j].index })
	return bins
}

func (s *SparseStore) Copy() Store {
	countsCopy := make(map[int]float64)
	for index, count := range s.counts {
		countsCopy[index] = count
	}
	return &SparseStore{counts: countsCopy}
}

func (s *SparseStore) IsEmpty() bool {
	return len(s.counts) == 0
}

func (s *SparseStore) MaxIndex() (int, error) {
	if s.IsEmpty() {
		return 0, errUndefinedMaxIndex
	}
	maxIndex := minInt
	for index := range s.counts {
		if index > maxIndex {
			maxIndex = index
		}
	}
	return maxIndex, nil
}

func (s *SparseStore) MinIndex() (int, error) {
	if s.IsEmpty() {
		return 0, errUndefinedMinIndex
	}
	minIndex := maxInt
	for index := range s.counts {
		if index < minIndex {
			minIndex = index
		}
	}
	return minIndex, nil
}

func (s *SparseStore) TotalCount() float64 {
	totalCount := float64(0)
	for _, count := range s.counts {
		totalCount += count
	}
	return totalCount
}

func (s *SparseStore) KeyAtRank(rank float64) int {
	orderedBins := s.orderedBins()
	cumulCount := float64(0)
	for _, bin := range orderedBins {
		cumulCount += bin.count
		if cumulCount > rank {
			return bin.index
		}
	}
	maxIndex, err := s.MaxIndex()
	if err == nil {
		return maxIndex
	} else {
		// FIXME: make Store's KeyAtRank consistent with MinIndex and MaxIndex
		return 0
	}
}

func (s *SparseStore) MergeWith(store Store) {
	for bin := range store.Bins() {
		s.AddBin(bin)
	}
}

func (s *SparseStore) ToProto() *sketchpb.Store {
	binCounts := make(map[int32]float64)
	for index, count := range binCounts {
		binCounts[index] = count
	}
	return &sketchpb.Store{BinCounts: binCounts}
}