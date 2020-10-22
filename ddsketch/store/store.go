// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020 Datadog, Inc.

package store

type Store interface {
	Add(index int)
	AddBin(bin Bin)
	Bins() <-chan Bin
	IsEmpty() bool
	MaxIndex() int
	MinIndex() int
	TotalCount() float64
	KeyAtRank(rank float64) int
	MergeWith(store Store)
}
