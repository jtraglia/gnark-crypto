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

package bn254

const MAX_BATCH_SIZE = 600

type batchOp struct {
	bucketID, pointID uint32
}

func (o batchOp) isNeg() bool {
	return o.pointID&1 == 1
}

// processChunkG1BatchAffine process a chunk of the scalars during the msm
// using affine coordinates for the buckets. To amortize the cost of the inverse in the affine addition
// we use a batch affine addition.
//
// this is derived from a PR by 0x0ece : https://github.com/ConsenSys/gnark-crypto/pull/249
// See Section 5.3: ia.cr/2022/1396
func processChunkG1BatchAffine[B ibG1Affine](chunk uint64,
	chRes chan<- g1JacExtended,
	c uint64,
	points []G1Affine,
	pscalars []uint32) {

	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	batch := newBatchG1Affine(&buckets, points)
	queue := make([]batchOp, 0, 4096) // TODO find right capacity here.
	nbBatches := 0
	for i := 0; i < len(pscalars); i++ {
		bits := pscalars[i]

		if bits == 0 {
			continue
		}

		op := batchOp{pointID: uint32(i) << 1}
		// if msbWindow bit is set, we need to substract
		if bits&1 == 0 {
			// add
			op.bucketID = uint32((bits >> 1) - 1)
			// buckets[bits-1].Add(&points[i], &buckets[bits-1])
		} else {
			// sub
			op.bucketID = (uint32((bits >> 1)))
			op.pointID += 1
			// op.isNeg = true
			// buckets[bits & ^msbWindow].Sub( &buckets[bits & ^msbWindow], &points[i])
		}
		if batch.CanAdd(op.bucketID) {
			batch.Add(op)
			if batch.IsFull() {
				batch.ExecuteAndReset()
				nbBatches++
				if len(queue) != 0 { // TODO @gbotrel this doesn't seem to help much? should minimize queue resizing
					batch.Add(queue[len(queue)-1])
					queue = queue[:len(queue)-1]
				}
			}
		} else {
			// put it in queue.
			queue = append(queue, op)
		}
	}
	// fmt.Printf("chunk %d\nlen(queue)=%d\nnbBatches=%d\nbatchSize=%d\nnbBuckets=%d\nnbPoints=%d\n",
	// 	chunk, len(queue), nbBatches, batch.batchSize, len(buckets), len(points))
	// batch.ExecuteAndReset()
	for len(queue) != 0 {
		queue = processQueueG1Affine(queue, &batch)
		batch.ExecuteAndReset() // execute batch even if not full.
	}

	// flush items in batch.
	batch.ExecuteAndReset()

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
type bucketG1AffineC4 [1 << (4 - 1)]G1Affine
type bucketG1AffineC5 [1 << (5 - 1)]G1Affine
type bucketG1AffineC6 [1 << (6 - 1)]G1Affine
type bucketG1AffineC7 [1 << (7 - 1)]G1Affine
type bucketG1AffineC8 [1 << (8 - 1)]G1Affine
type bucketG1AffineC9 [1 << (9 - 1)]G1Affine
type bucketG1AffineC10 [1 << (10 - 1)]G1Affine
type bucketG1AffineC11 [1 << (11 - 1)]G1Affine
type bucketG1AffineC12 [1 << (12 - 1)]G1Affine
type bucketG1AffineC13 [1 << (13 - 1)]G1Affine
type bucketG1AffineC14 [1 << (14 - 1)]G1Affine
type bucketG1AffineC15 [1 << (15 - 1)]G1Affine
type bucketG1AffineC16 [1 << (16 - 1)]G1Affine
type bucketG1AffineC20 [1 << (20 - 1)]G1Affine
type bucketG1AffineC21 [1 << (21 - 1)]G1Affine

type ibG1Affine interface {
	bucketG1AffineC4 |
		bucketG1AffineC5 |
		bucketG1AffineC6 |
		bucketG1AffineC7 |
		bucketG1AffineC8 |
		bucketG1AffineC9 |
		bucketG1AffineC10 |
		bucketG1AffineC11 |
		bucketG1AffineC12 |
		bucketG1AffineC13 |
		bucketG1AffineC14 |
		bucketG1AffineC15 |
		bucketG1AffineC16 |
		bucketG1AffineC20 |
		bucketG1AffineC21
}

type BatchG1Affine[B ibG1Affine] struct {
	P         [MAX_BATCH_SIZE]G1Affine
	R         [MAX_BATCH_SIZE]*G1Affine
	batchSize int
	cptP      int
	bucketIds map[uint32]struct{}
	points    []G1Affine
	buckets   *B
}

func newBatchG1Affine[B ibG1Affine](buckets *B, points []G1Affine) BatchG1Affine[B] {
	batchSize := len(*buckets) / 5
	if batchSize > MAX_BATCH_SIZE {
		batchSize = MAX_BATCH_SIZE
	}
	if batchSize <= 0 {
		batchSize = 1
	}
	return BatchG1Affine[B]{
		buckets:   buckets,
		points:    points,
		batchSize: batchSize,
		bucketIds: make(map[uint32]struct{}, len(*buckets)/2),
	}
}

func (b *BatchG1Affine[B]) IsFull() bool {
	return b.cptP == b.batchSize
}

func (b *BatchG1Affine[B]) ExecuteAndReset() {
	if b.cptP == 0 {
		return
	}
	// for i := 0; i < len(b.R); i++ {
	// 	b.R[i].Add(b.R[i], b.P[i])
	// }
	BatchAddG1Affine(b.R[:b.cptP], b.P[:b.cptP], b.cptP)
	for k := range b.bucketIds {
		delete(b.bucketIds, k)
	}
	// b.bucketIds = [MAX_BATCH_SIZE]uint32{}
	b.cptP = 0
}

func (b *BatchG1Affine[B]) CanAdd(bID uint32) bool {
	_, ok := b.bucketIds[bID]
	return !ok
}

func (b *BatchG1Affine[B]) Add(op batchOp) {
	// CanAdd must be called before --> ensures bucket is not "used" in current batch

	BK := &(*b.buckets)[op.bucketID]
	P := &b.points[op.pointID>>1]
	if P.IsInfinity() {
		return
	}
	// handle special cases with inf or -P / P
	if BK.IsInfinity() {
		if op.isNeg() {
			BK.Neg(P)
		} else {
			BK.Set(P)
		}
		return
	}
	if op.isNeg() {
		// if bucket == P --> -P == 0
		if BK.Equal(P) {
			BK.setInfinity()
			return
		}
	} else {
		// if bucket == -P, B == 0
		if BK.X.Equal(&P.X) && !BK.Y.Equal(&P.Y) {
			BK.setInfinity()
			return
		}
	}

	// b.bucketIds[b.cptP] = op.bucketID
	b.bucketIds[op.bucketID] = struct{}{}
	b.R[b.cptP] = BK
	if op.isNeg() {
		b.P[b.cptP].Neg(P)
	} else {
		b.P[b.cptP].Set(P)
	}
	b.cptP++
}

func processQueueG1Affine[B ibG1Affine](queue []batchOp, batch *BatchG1Affine[B]) []batchOp {
	for i := len(queue) - 1; i >= 0; i-- {
		if batch.CanAdd(queue[i].bucketID) {
			batch.Add(queue[i])
			if batch.IsFull() {
				batch.ExecuteAndReset()
			}
			queue[i] = queue[len(queue)-1]
			queue = queue[:len(queue)-1]
		}
	}
	return queue

}

// processChunkG2BatchAffine process a chunk of the scalars during the msm
// using affine coordinates for the buckets. To amortize the cost of the inverse in the affine addition
// we use a batch affine addition.
//
// this is derived from a PR by 0x0ece : https://github.com/ConsenSys/gnark-crypto/pull/249
// See Section 5.3: ia.cr/2022/1396
func processChunkG2BatchAffine[B ibG2Affine](chunk uint64,
	chRes chan<- g2JacExtended,
	c uint64,
	points []G2Affine,
	pscalars []uint32) {

	var buckets B
	for i := 0; i < len(buckets); i++ {
		buckets[i].setInfinity()
	}

	batch := newBatchG2Affine(&buckets, points)
	queue := make([]batchOp, 0, 4096) // TODO find right capacity here.
	nbBatches := 0
	for i := 0; i < len(pscalars); i++ {
		bits := pscalars[i]

		if bits == 0 {
			continue
		}

		op := batchOp{pointID: uint32(i) << 1}
		// if msbWindow bit is set, we need to substract
		if bits&1 == 0 {
			// add
			op.bucketID = uint32((bits >> 1) - 1)
			// buckets[bits-1].Add(&points[i], &buckets[bits-1])
		} else {
			// sub
			op.bucketID = (uint32((bits >> 1)))
			op.pointID += 1
			// op.isNeg = true
			// buckets[bits & ^msbWindow].Sub( &buckets[bits & ^msbWindow], &points[i])
		}
		if batch.CanAdd(op.bucketID) {
			batch.Add(op)
			if batch.IsFull() {
				batch.ExecuteAndReset()
				nbBatches++
				if len(queue) != 0 { // TODO @gbotrel this doesn't seem to help much? should minimize queue resizing
					batch.Add(queue[len(queue)-1])
					queue = queue[:len(queue)-1]
				}
			}
		} else {
			// put it in queue.
			queue = append(queue, op)
		}
	}
	// fmt.Printf("chunk %d\nlen(queue)=%d\nnbBatches=%d\nbatchSize=%d\nnbBuckets=%d\nnbPoints=%d\n",
	// 	chunk, len(queue), nbBatches, batch.batchSize, len(buckets), len(points))
	// batch.ExecuteAndReset()
	for len(queue) != 0 {
		queue = processQueueG2Affine(queue, &batch)
		batch.ExecuteAndReset() // execute batch even if not full.
	}

	// flush items in batch.
	batch.ExecuteAndReset()

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
type bucketG2AffineC4 [1 << (4 - 1)]G2Affine
type bucketG2AffineC5 [1 << (5 - 1)]G2Affine
type bucketG2AffineC6 [1 << (6 - 1)]G2Affine
type bucketG2AffineC7 [1 << (7 - 1)]G2Affine
type bucketG2AffineC8 [1 << (8 - 1)]G2Affine
type bucketG2AffineC9 [1 << (9 - 1)]G2Affine
type bucketG2AffineC10 [1 << (10 - 1)]G2Affine
type bucketG2AffineC11 [1 << (11 - 1)]G2Affine
type bucketG2AffineC12 [1 << (12 - 1)]G2Affine
type bucketG2AffineC13 [1 << (13 - 1)]G2Affine
type bucketG2AffineC14 [1 << (14 - 1)]G2Affine
type bucketG2AffineC15 [1 << (15 - 1)]G2Affine
type bucketG2AffineC16 [1 << (16 - 1)]G2Affine
type bucketG2AffineC20 [1 << (20 - 1)]G2Affine
type bucketG2AffineC21 [1 << (21 - 1)]G2Affine

type ibG2Affine interface {
	bucketG2AffineC4 |
		bucketG2AffineC5 |
		bucketG2AffineC6 |
		bucketG2AffineC7 |
		bucketG2AffineC8 |
		bucketG2AffineC9 |
		bucketG2AffineC10 |
		bucketG2AffineC11 |
		bucketG2AffineC12 |
		bucketG2AffineC13 |
		bucketG2AffineC14 |
		bucketG2AffineC15 |
		bucketG2AffineC16 |
		bucketG2AffineC20 |
		bucketG2AffineC21
}

type BatchG2Affine[B ibG2Affine] struct {
	P         [MAX_BATCH_SIZE]G2Affine
	R         [MAX_BATCH_SIZE]*G2Affine
	batchSize int
	cptP      int
	bucketIds map[uint32]struct{}
	points    []G2Affine
	buckets   *B
}

func newBatchG2Affine[B ibG2Affine](buckets *B, points []G2Affine) BatchG2Affine[B] {
	batchSize := len(*buckets) / 5
	if batchSize > MAX_BATCH_SIZE {
		batchSize = MAX_BATCH_SIZE
	}
	if batchSize <= 0 {
		batchSize = 1
	}
	return BatchG2Affine[B]{
		buckets:   buckets,
		points:    points,
		batchSize: batchSize,
		bucketIds: make(map[uint32]struct{}, len(*buckets)/2),
	}
}

func (b *BatchG2Affine[B]) IsFull() bool {
	return b.cptP == b.batchSize
}

func (b *BatchG2Affine[B]) ExecuteAndReset() {
	if b.cptP == 0 {
		return
	}
	// for i := 0; i < len(b.R); i++ {
	// 	b.R[i].Add(b.R[i], b.P[i])
	// }
	BatchAddG2Affine(b.R[:b.cptP], b.P[:b.cptP], b.cptP)
	for k := range b.bucketIds {
		delete(b.bucketIds, k)
	}
	// b.bucketIds = [MAX_BATCH_SIZE]uint32{}
	b.cptP = 0
}

func (b *BatchG2Affine[B]) CanAdd(bID uint32) bool {
	_, ok := b.bucketIds[bID]
	return !ok
}

func (b *BatchG2Affine[B]) Add(op batchOp) {
	// CanAdd must be called before --> ensures bucket is not "used" in current batch

	BK := &(*b.buckets)[op.bucketID]
	P := &b.points[op.pointID>>1]
	if P.IsInfinity() {
		return
	}
	// handle special cases with inf or -P / P
	if BK.IsInfinity() {
		if op.isNeg() {
			BK.Neg(P)
		} else {
			BK.Set(P)
		}
		return
	}
	if op.isNeg() {
		// if bucket == P --> -P == 0
		if BK.Equal(P) {
			BK.setInfinity()
			return
		}
	} else {
		// if bucket == -P, B == 0
		if BK.X.Equal(&P.X) && !BK.Y.Equal(&P.Y) {
			BK.setInfinity()
			return
		}
	}

	// b.bucketIds[b.cptP] = op.bucketID
	b.bucketIds[op.bucketID] = struct{}{}
	b.R[b.cptP] = BK
	if op.isNeg() {
		b.P[b.cptP].Neg(P)
	} else {
		b.P[b.cptP].Set(P)
	}
	b.cptP++
}

func processQueueG2Affine[B ibG2Affine](queue []batchOp, batch *BatchG2Affine[B]) []batchOp {
	for i := len(queue) - 1; i >= 0; i-- {
		if batch.CanAdd(queue[i].bucketID) {
			batch.Add(queue[i])
			if batch.IsFull() {
				batch.ExecuteAndReset()
			}
			queue[i] = queue[len(queue)-1]
			queue = queue[:len(queue)-1]
		}
	}
	return queue

}
