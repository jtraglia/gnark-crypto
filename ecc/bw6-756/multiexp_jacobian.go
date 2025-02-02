// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bw6756

func processChunkG1Jacobian[B ibg1JacExtended](chunk uint64,
	chRes chan<- g1JacExtended,
	c uint64,
	points []G1Affine,
	digits []uint16) {

	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	// for each scalars, get the digit corresponding to the chunk we're processing.
	for i, digit := range digits {
		if digit == 0 {
			continue
		}

		// if msbWindow bit is set, we need to subtract
		if digit&1 == 0 {
			// add
			buckets[(digit>>1)-1].addMixed(&points[i])
		} else {
			// sub
			buckets[(digit >> 1)].subMixed(&points[i])
		}
	}

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, total g1JacExtended
	runningSum.setInfinity()
	total.setInfinity()
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].ZZ.IsZero() {
			runningSum.add(&buckets[k])
		}
		total.add(&runningSum)
	}

	chRes <- total
}

// we declare the buckets as fixed-size array types
// this allow us to allocate the buckets on the stack
type bucketg1JacExtendedC3 [4]g1JacExtended
type bucketg1JacExtendedC4 [8]g1JacExtended
type bucketg1JacExtendedC5 [16]g1JacExtended
type bucketg1JacExtendedC8 [128]g1JacExtended
type bucketg1JacExtendedC11 [1024]g1JacExtended
type bucketg1JacExtendedC16 [32768]g1JacExtended

type ibg1JacExtended interface {
	bucketg1JacExtendedC3 |
		bucketg1JacExtendedC4 |
		bucketg1JacExtendedC5 |
		bucketg1JacExtendedC8 |
		bucketg1JacExtendedC11 |
		bucketg1JacExtendedC16
}

func processChunkG2Jacobian[B ibg2JacExtended](chunk uint64,
	chRes chan<- g2JacExtended,
	c uint64,
	points []G2Affine,
	digits []uint16) {

	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	// for each scalars, get the digit corresponding to the chunk we're processing.
	for i, digit := range digits {
		if digit == 0 {
			continue
		}

		// if msbWindow bit is set, we need to subtract
		if digit&1 == 0 {
			// add
			buckets[(digit>>1)-1].addMixed(&points[i])
		} else {
			// sub
			buckets[(digit >> 1)].subMixed(&points[i])
		}
	}

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, total g2JacExtended
	runningSum.setInfinity()
	total.setInfinity()
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].ZZ.IsZero() {
			runningSum.add(&buckets[k])
		}
		total.add(&runningSum)
	}

	chRes <- total
}

// we declare the buckets as fixed-size array types
// this allow us to allocate the buckets on the stack
type bucketg2JacExtendedC3 [4]g2JacExtended
type bucketg2JacExtendedC4 [8]g2JacExtended
type bucketg2JacExtendedC5 [16]g2JacExtended
type bucketg2JacExtendedC8 [128]g2JacExtended
type bucketg2JacExtendedC11 [1024]g2JacExtended
type bucketg2JacExtendedC16 [32768]g2JacExtended

type ibg2JacExtended interface {
	bucketg2JacExtendedC3 |
		bucketg2JacExtendedC4 |
		bucketg2JacExtendedC5 |
		bucketg2JacExtendedC8 |
		bucketg2JacExtendedC11 |
		bucketg2JacExtendedC16
}
