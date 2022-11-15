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

import (
	"github.com/consensys/gnark-crypto/ecc/bw6-756/fp"
)

type batchOpG1Affine struct {
	bucketID uint16
	point    G1Affine
}

func (o batchOpG1Affine) isNeg() bool {
	return o.bucketID&1 == 1
}

// processChunkG1BatchAffine process a chunk of the scalars during the msm
// using affine coordinates for the buckets. To amortize the cost of the inverse in the affine addition
// we use a batch affine addition.
//
// this is derived from a PR by 0x0ece : https://github.com/ConsenSys/gnark-crypto/pull/249
// See Section 5.3: ia.cr/2022/1396
func processChunkG1BatchAffine[B ibG1Affine, BS bitSet, TP pG1Affine, TPP ppG1Affine, TQ qOpsG1Affine, TC cG1Affine](
	chunk uint64,
	chRes chan<- g1JacExtended,
	c uint64,
	points []G1Affine,
	digits []uint16) {

	// init the buckets
	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	// setup for the batch affine;
	var (
		bucketIds BS  // bitSet to signify presence of a bucket in current batch
		cptAdd    int // count the number of bucket + point added to current batch
		R         TPP // bucket references
		P         TP  // points to be added to R (buckets); it is beneficial to store them on the stack (ie copy)
		queue     TQ  // queue of points that conflict the current batch
		qID       int // current position in queue
	)

	batchSize := len(P)

	isFull := func() bool {
		return cptAdd == batchSize
	}

	executeAndReset := func() {
		batchAddG1Affine[TP, TPP, TC](&R, &P, cptAdd)
		var tmp BS
		bucketIds = tmp
		cptAdd = 0
	}

	add := func(bucketID uint16, PP *G1Affine, isAdd bool) {
		// @precondition: ensures bucket is not "used" in current batch
		BK := &buckets[bucketID]
		// handle special cases with inf or -P / P
		if BK.IsInfinity() {
			if isAdd {
				BK.Set(PP)
			} else {
				BK.Neg(PP)
			}
			return
		}
		if BK.X.Equal(&PP.X) {
			if BK.Y.Equal(&PP.Y) {
				// P + P: doubling, which should be quite rare --
				// TODO FIXME @gbotrel / @yelhousni this path is not taken by our tests.
				// need doubling in affine implemented ?
				if isAdd {
					BK.Add(BK, BK)
				} else {
					BK.setInfinity()
				}

				return
			}
			if isAdd {
				BK.setInfinity()
			} else {
				BK.Add(BK, BK)
			}
			return
		}

		bucketIds[bucketID] = true
		R[cptAdd] = BK
		if isAdd {
			P[cptAdd].Set(PP)
		} else {
			P[cptAdd].Neg(PP)
		}
		cptAdd++
	}

	processQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if bucketIds[queue[i].bucketID] {
				continue
			}
			add(queue[i].bucketID, &queue[i].point, true)
			if isFull() {
				executeAndReset()
			}
			queue[i] = queue[qID-1]
			qID--
		}
	}

	for i, digit := range digits {

		if digit == 0 || points[i].IsInfinity() {
			continue
		}

		bucketID := uint16((digit >> 1))
		isAdd := digit&1 == 0
		if isAdd {
			// add
			bucketID -= 1
		}

		if bucketIds[bucketID] {
			// put it in queue
			queue[qID].bucketID = bucketID
			if isAdd {
				queue[qID].point = points[i]
			} else {
				queue[qID].point.Neg(&points[i])
			}
			qID++

			// queue is full, flush it.
			if qID == len(queue)-1 {
				executeAndReset()
				processQueue()
			}
			continue
		}

		// we add the point to the batch.
		add(bucketID, &points[i], isAdd)
		if isFull() {
			executeAndReset()
			processQueue()
		}
	}

	// empty the queue
	for qID != 0 {
		processQueue()
		executeAndReset()
	}

	// flush items in batch.
	executeAndReset()

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, total g1JacExtended
	runningSum.setInfinity()
	total.setInfinity()
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].IsInfinity() {
			runningSum.addMixed(&buckets[k])
		}
		total.add(&runningSum)
	}

	chRes <- total

}

// we declare the buckets as fixed-size array types
// this allow us to allocate the buckets on the stack
type bucketG1AffineC16 [1 << (16 - 1)]G1Affine

// buckets: array of G1Affine points of size 1 << (c-1)
type ibG1Affine interface {
	bucketG1AffineC16
}

// array of coordinates fp.Element
type cG1Affine interface {
	cG1AffineC16
}

// buckets: array of G1Affine points (for the batch addition)
type pG1Affine interface {
	pG1AffineC16
}

// buckets: array of *G1Affine points (for the batch addition)
type ppG1Affine interface {
	ppG1AffineC16
}

// buckets: array of G1Affine queue operations (for the batch addition)
type qOpsG1Affine interface {
	qG1AffineC16
}

// batch size 640 when c = 16
type cG1AffineC16 [640]fp.Element
type pG1AffineC16 [640]G1Affine
type ppG1AffineC16 [640]*G1Affine
type qG1AffineC16 [640]batchOpG1Affine

type batchOpG2Affine struct {
	bucketID uint16
	point    G2Affine
}

func (o batchOpG2Affine) isNeg() bool {
	return o.bucketID&1 == 1
}

// processChunkG2BatchAffine process a chunk of the scalars during the msm
// using affine coordinates for the buckets. To amortize the cost of the inverse in the affine addition
// we use a batch affine addition.
//
// this is derived from a PR by 0x0ece : https://github.com/ConsenSys/gnark-crypto/pull/249
// See Section 5.3: ia.cr/2022/1396
func processChunkG2BatchAffine[B ibG2Affine, BS bitSet, TP pG2Affine, TPP ppG2Affine, TQ qOpsG2Affine, TC cG2Affine](
	chunk uint64,
	chRes chan<- g2JacExtended,
	c uint64,
	points []G2Affine,
	digits []uint16) {

	// init the buckets
	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	// setup for the batch affine;
	var (
		bucketIds BS  // bitSet to signify presence of a bucket in current batch
		cptAdd    int // count the number of bucket + point added to current batch
		R         TPP // bucket references
		P         TP  // points to be added to R (buckets); it is beneficial to store them on the stack (ie copy)
		queue     TQ  // queue of points that conflict the current batch
		qID       int // current position in queue
	)

	batchSize := len(P)

	isFull := func() bool {
		return cptAdd == batchSize
	}

	executeAndReset := func() {
		batchAddG2Affine[TP, TPP, TC](&R, &P, cptAdd)
		var tmp BS
		bucketIds = tmp
		cptAdd = 0
	}

	add := func(bucketID uint16, PP *G2Affine, isAdd bool) {
		// @precondition: ensures bucket is not "used" in current batch
		BK := &buckets[bucketID]
		// handle special cases with inf or -P / P
		if BK.IsInfinity() {
			if isAdd {
				BK.Set(PP)
			} else {
				BK.Neg(PP)
			}
			return
		}
		if BK.X.Equal(&PP.X) {
			if BK.Y.Equal(&PP.Y) {
				// P + P: doubling, which should be quite rare --
				// TODO FIXME @gbotrel / @yelhousni this path is not taken by our tests.
				// need doubling in affine implemented ?
				if isAdd {
					BK.Add(BK, BK)
				} else {
					BK.setInfinity()
				}

				return
			}
			if isAdd {
				BK.setInfinity()
			} else {
				BK.Add(BK, BK)
			}
			return
		}

		bucketIds[bucketID] = true
		R[cptAdd] = BK
		if isAdd {
			P[cptAdd].Set(PP)
		} else {
			P[cptAdd].Neg(PP)
		}
		cptAdd++
	}

	processQueue := func() {
		for i := qID - 1; i >= 0; i-- {
			if bucketIds[queue[i].bucketID] {
				continue
			}
			add(queue[i].bucketID, &queue[i].point, true)
			if isFull() {
				executeAndReset()
			}
			queue[i] = queue[qID-1]
			qID--
		}
	}

	for i, digit := range digits {

		if digit == 0 || points[i].IsInfinity() {
			continue
		}

		bucketID := uint16((digit >> 1))
		isAdd := digit&1 == 0
		if isAdd {
			// add
			bucketID -= 1
		}

		if bucketIds[bucketID] {
			// put it in queue
			queue[qID].bucketID = bucketID
			if isAdd {
				queue[qID].point = points[i]
			} else {
				queue[qID].point.Neg(&points[i])
			}
			qID++

			// queue is full, flush it.
			if qID == len(queue)-1 {
				executeAndReset()
				processQueue()
			}
			continue
		}

		// we add the point to the batch.
		add(bucketID, &points[i], isAdd)
		if isFull() {
			executeAndReset()
			processQueue()
		}
	}

	// empty the queue
	for qID != 0 {
		processQueue()
		executeAndReset()
	}

	// flush items in batch.
	executeAndReset()

	// reduce buckets into total
	// total =  bucket[0] + 2*bucket[1] + 3*bucket[2] ... + n*bucket[n-1]

	var runningSum, total g2JacExtended
	runningSum.setInfinity()
	total.setInfinity()
	for k := len(buckets) - 1; k >= 0; k-- {
		if !buckets[k].IsInfinity() {
			runningSum.addMixed(&buckets[k])
		}
		total.add(&runningSum)
	}

	chRes <- total

}

// we declare the buckets as fixed-size array types
// this allow us to allocate the buckets on the stack
type bucketG2AffineC16 [1 << (16 - 1)]G2Affine

// buckets: array of G2Affine points of size 1 << (c-1)
type ibG2Affine interface {
	bucketG2AffineC16
}

// array of coordinates fp.Element
type cG2Affine interface {
	cG2AffineC16
}

// buckets: array of G2Affine points (for the batch addition)
type pG2Affine interface {
	pG2AffineC16
}

// buckets: array of *G2Affine points (for the batch addition)
type ppG2Affine interface {
	ppG2AffineC16
}

// buckets: array of G2Affine queue operations (for the batch addition)
type qOpsG2Affine interface {
	qG2AffineC16
}

// batch size 640 when c = 16
type cG2AffineC16 [640]fp.Element
type pG2AffineC16 [640]G2Affine
type ppG2AffineC16 [640]*G2Affine
type qG2AffineC16 [640]batchOpG2Affine

type bitSetC4 [1 << (4 - 1)]bool
type bitSetC5 [1 << (5 - 1)]bool
type bitSetC8 [1 << (8 - 1)]bool
type bitSetC16 [1 << (16 - 1)]bool

type bitSet interface {
	bitSetC4 |
		bitSetC5 |
		bitSetC8 |
		bitSetC16
}
