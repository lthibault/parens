package parens

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"math"
	"unsafe"
)

const shiftSize = 5
const nodeSize = 32
const shiftBitMask = 0x1F

type commonNode interface{}

var emptyCommonNode commonNode = []commonNode{}

func uintMin(a, b uint) uint {
	if a < b {
		return a
	}

	return b
}

func newPath(shift uint, node commonNode) commonNode {
	if shift == 0 {
		return node
	}

	return newPath(shift-shiftSize, commonNode([]commonNode{node}))
}

func assertSliceOk(start, stop, len int) {
	if start < 0 {
		panic(fmt.Sprintf("Invalid slice index %d (index must be non-negative)", start))
	}

	if start > stop {
		panic(fmt.Sprintf("Invalid slice index: %d > %d", start, stop))
	}

	if stop > len {
		panic(fmt.Sprintf("Slice bounds out of range, start=%d, stop=%d, len=%d", start, stop, len))
	}
}

const upperMapLoadFactor float64 = 8.0
const lowerMapLoadFactor float64 = 2.0
const initialMapLoadFactor float64 = (upperMapLoadFactor + lowerMapLoadFactor) / 2

//////////////////////////
//// Hash functions //////
//////////////////////////

func hash(x []byte) uint32 {
	return crc32.ChecksumIEEE(x)
}

//go:noescape
//go:linkname nilinterhash runtime.nilinterhash
func nilinterhash(p unsafe.Pointer, h uintptr) uintptr

func interfaceHash(x interface{}) uint32 {
	return uint32(nilinterhash(unsafe.Pointer(&x), 0))
}

func byteHash(x byte) uint32 {
	return hash([]byte{x})
}

func uint8Hash(x uint8) uint32 {
	return byteHash(byte(x))
}

func int8Hash(x int8) uint32 {
	return uint8Hash(uint8(x))
}

func uint16Hash(x uint16) uint32 {
	bX := make([]byte, 2)
	binary.LittleEndian.PutUint16(bX, x)
	return hash(bX)
}

func int16Hash(x int16) uint32 {
	return uint16Hash(uint16(x))
}

func uint32Hash(x uint32) uint32 {
	bX := make([]byte, 4)
	binary.LittleEndian.PutUint32(bX, x)
	return hash(bX)
}

func int32Hash(x int32) uint32 {
	return uint32Hash(uint32(x))
}

func uint64Hash(x uint64) uint32 {
	bX := make([]byte, 8)
	binary.LittleEndian.PutUint64(bX, x)
	return hash(bX)
}

func int64Hash(x int64) uint32 {
	return uint64Hash(uint64(x))
}

func intHash(x int) uint32 {
	return int64Hash(int64(x))
}

func uintHash(x uint) uint32 {
	return uint64Hash(uint64(x))
}

func boolHash(x bool) uint32 {
	if x {
		return 1
	}

	return 0
}

func runeHash(x rune) uint32 {
	return int32Hash(int32(x))
}

func stringHash(x string) uint32 {
	return hash([]byte(x))
}

func float64Hash(x float64) uint32 {
	return uint64Hash(math.Float64bits(x))
}

func float32Hash(x float32) uint32 {
	return uint32Hash(math.Float32bits(x))
}

///////////
/// Map ///
///////////

//////////////////////
/// Backing vector ///
//////////////////////

type privatePersistentScopeItemBucketVector struct {
	tail  []privatePersistentScopeItemBucket
	root  commonNode
	len   uint
	shift uint
}

// PersistentScopeItem .
type PersistentScopeItem struct {
	Key   string
	Value scopeEntry
}

type privatePersistentScopeItemBucket []PersistentScopeItem

var emptyPersistentScopeItemBucketVectorTail = make([]privatePersistentScopeItemBucket, 0)
var emptyPersistentScopeItemBucketVector *privatePersistentScopeItemBucketVector = &privatePersistentScopeItemBucketVector{root: emptyCommonNode, shift: shiftSize, tail: emptyPersistentScopeItemBucketVectorTail}

func (v *privatePersistentScopeItemBucketVector) Get(i int) privatePersistentScopeItemBucket {
	if i < 0 || uint(i) >= v.len {
		panic("Index out of bounds")
	}

	return v.sliceFor(uint(i))[i&shiftBitMask]
}

func (v *privatePersistentScopeItemBucketVector) sliceFor(i uint) []privatePersistentScopeItemBucket {
	if i >= v.tailOffset() {
		return v.tail
	}

	node := v.root
	for level := v.shift; level > 0; level -= shiftSize {
		node = node.([]commonNode)[(i>>level)&shiftBitMask]
	}

	return node.([]privatePersistentScopeItemBucket)
}

func (v *privatePersistentScopeItemBucketVector) tailOffset() uint {
	if v.len < nodeSize {
		return 0
	}

	return ((v.len - 1) >> shiftSize) << shiftSize
}

func (v *privatePersistentScopeItemBucketVector) Set(i int, item privatePersistentScopeItemBucket) *privatePersistentScopeItemBucketVector {
	if i < 0 || uint(i) >= v.len {
		panic("Index out of bounds")
	}

	if uint(i) >= v.tailOffset() {
		newTail := make([]privatePersistentScopeItemBucket, len(v.tail))
		copy(newTail, v.tail)
		newTail[i&shiftBitMask] = item
		return &privatePersistentScopeItemBucketVector{root: v.root, tail: newTail, len: v.len, shift: v.shift}
	}

	return &privatePersistentScopeItemBucketVector{root: v.doAssoc(v.shift, v.root, uint(i), item), tail: v.tail, len: v.len, shift: v.shift}
}

func (v *privatePersistentScopeItemBucketVector) doAssoc(level uint, node commonNode, i uint, item privatePersistentScopeItemBucket) commonNode {
	if level == 0 {
		ret := make([]privatePersistentScopeItemBucket, nodeSize)
		copy(ret, node.([]privatePersistentScopeItemBucket))
		ret[i&shiftBitMask] = item
		return ret
	}

	ret := make([]commonNode, nodeSize)
	copy(ret, node.([]commonNode))
	subidx := (i >> level) & shiftBitMask
	ret[subidx] = v.doAssoc(level-shiftSize, ret[subidx], i, item)
	return ret
}

func (v *privatePersistentScopeItemBucketVector) pushTail(level uint, parent commonNode, tailNode []privatePersistentScopeItemBucket) commonNode {
	subIdx := ((v.len - 1) >> level) & shiftBitMask
	parentNode := parent.([]commonNode)
	ret := make([]commonNode, subIdx+1)
	copy(ret, parentNode)
	var nodeToInsert commonNode

	if level == shiftSize {
		nodeToInsert = tailNode
	} else if subIdx < uint(len(parentNode)) {
		nodeToInsert = v.pushTail(level-shiftSize, parentNode[subIdx], tailNode)
	} else {
		nodeToInsert = newPath(level-shiftSize, tailNode)
	}

	ret[subIdx] = nodeToInsert
	return ret
}

func (v *privatePersistentScopeItemBucketVector) Append(item ...privatePersistentScopeItemBucket) *privatePersistentScopeItemBucketVector {
	result := v
	itemLen := uint(len(item))
	for insertOffset := uint(0); insertOffset < itemLen; {
		tailLen := result.len - result.tailOffset()
		tailFree := nodeSize - tailLen
		if tailFree == 0 {
			result = result.pushLeafNode(result.tail)
			result.tail = emptyPersistentScopeItemBucketVector.tail
			tailFree = nodeSize
			tailLen = 0
		}

		batchLen := uintMin(itemLen-insertOffset, tailFree)
		newTail := make([]privatePersistentScopeItemBucket, 0, tailLen+batchLen)
		newTail = append(newTail, result.tail...)
		newTail = append(newTail, item[insertOffset:insertOffset+batchLen]...)
		result = &privatePersistentScopeItemBucketVector{root: result.root, tail: newTail, len: result.len + batchLen, shift: result.shift}
		insertOffset += batchLen
	}

	return result
}

func (v *privatePersistentScopeItemBucketVector) pushLeafNode(node []privatePersistentScopeItemBucket) *privatePersistentScopeItemBucketVector {
	var newRoot commonNode
	newShift := v.shift

	// Root overflow?
	if (v.len >> shiftSize) > (1 << v.shift) {
		newNode := newPath(v.shift, node)
		newRoot = commonNode([]commonNode{v.root, newNode})
		newShift = v.shift + shiftSize
	} else {
		newRoot = v.pushTail(v.shift, v.root, node)
	}

	return &privatePersistentScopeItemBucketVector{root: newRoot, tail: v.tail, len: v.len, shift: newShift}
}

func (v *privatePersistentScopeItemBucketVector) Len() int {
	return int(v.len)
}

func (v *privatePersistentScopeItemBucketVector) Range(f func(privatePersistentScopeItemBucket) bool) {
	var currentNode []privatePersistentScopeItemBucket
	for i := uint(0); i < v.len; i++ {
		if i&shiftBitMask == 0 {
			currentNode = v.sliceFor(uint(i))
		}

		if !f(currentNode[i&shiftBitMask]) {
			return
		}
	}
}

// persistentScope is a persistent key - value map
type persistentScope struct {
	backingVector *privatePersistentScopeItemBucketVector
	len           int
}

func (m *persistentScope) pos(key string) int {
	return int(uint64(stringHash(key)) % uint64(m.backingVector.Len()))
}

// Helper type used during map creation and reallocation
type privatePersistentScopeItemBuckets struct {
	buckets []privatePersistentScopeItemBucket
	length  int
}

func newPrivatePersistentScopeItemBuckets(itemCount int) *privatePersistentScopeItemBuckets {
	size := int(float64(itemCount)/initialMapLoadFactor) + 1
	buckets := make([]privatePersistentScopeItemBucket, size)
	return &privatePersistentScopeItemBuckets{buckets: buckets}
}

func (b *privatePersistentScopeItemBuckets) AddItem(item PersistentScopeItem) {
	ix := int(uint64(stringHash(item.Key)) % uint64(len(b.buckets)))
	bucket := b.buckets[ix]
	if bucket != nil {
		// Hash collision, merge with existing bucket
		for keyIx, bItem := range bucket {
			if item.Key == bItem.Key {
				bucket[keyIx] = item
				return
			}
		}

		b.buckets[ix] = append(bucket, PersistentScopeItem{Key: item.Key, Value: item.Value})
		b.length++
	} else {
		bucket := make(privatePersistentScopeItemBucket, 0, int(math.Max(initialMapLoadFactor, 1.0)))
		b.buckets[ix] = append(bucket, item)
		b.length++
	}
}

func (b *privatePersistentScopeItemBuckets) AddItemsFromMap(m *persistentScope) {
	m.backingVector.Range(func(bucket privatePersistentScopeItemBucket) bool {
		for _, item := range bucket {
			b.AddItem(item)
		}
		return true
	})
}

func newPersistentScope(items []PersistentScopeItem) *persistentScope {
	buckets := newPrivatePersistentScopeItemBuckets(len(items))
	for _, item := range items {
		buckets.AddItem(item)
	}

	return &persistentScope{backingVector: emptyPersistentScopeItemBucketVector.Append(buckets.buckets...), len: buckets.length}
}

// Len returns the number of items in m.
func (m *persistentScope) Len() int {
	return int(m.len)
}

// Load returns value identified by key. ok is set to true if key exists in the map, false otherwise.
func (m *persistentScope) Load(key string) (value scopeEntry, ok bool) {
	bucket := m.backingVector.Get(m.pos(key))
	if bucket != nil {
		for _, item := range bucket {
			if item.Key == key {
				return item.Value, true
			}
		}
	}

	var zeroValue scopeEntry
	return zeroValue, false
}

// Store returns a new persistentScope containing value identified by key.
func (m *persistentScope) Store(key string, value scopeEntry) *persistentScope {
	// Grow backing vector if load factor is too high
	if m.Len() >= m.backingVector.Len()*int(upperMapLoadFactor) {
		buckets := newPrivatePersistentScopeItemBuckets(m.Len() + 1)
		buckets.AddItemsFromMap(m)
		buckets.AddItem(PersistentScopeItem{Key: key, Value: value})
		return &persistentScope{backingVector: emptyPersistentScopeItemBucketVector.Append(buckets.buckets...), len: buckets.length}
	}

	pos := m.pos(key)
	bucket := m.backingVector.Get(pos)
	if bucket != nil {
		for ix, item := range bucket {
			if item.Key == key {
				// Overwrite existing item
				newBucket := make(privatePersistentScopeItemBucket, len(bucket))
				copy(newBucket, bucket)
				newBucket[ix] = PersistentScopeItem{Key: key, Value: value}
				return &persistentScope{backingVector: m.backingVector.Set(pos, newBucket), len: m.len}
			}
		}

		// Add new item to bucket
		newBucket := make(privatePersistentScopeItemBucket, len(bucket), len(bucket)+1)
		copy(newBucket, bucket)
		newBucket = append(newBucket, PersistentScopeItem{Key: key, Value: value})
		return &persistentScope{backingVector: m.backingVector.Set(pos, newBucket), len: m.len + 1}
	}

	item := PersistentScopeItem{Key: key, Value: value}
	newBucket := privatePersistentScopeItemBucket{item}
	return &persistentScope{backingVector: m.backingVector.Set(pos, newBucket), len: m.len + 1}
}

// Delete returns a new persistentScope without the element identified by key.
func (m *persistentScope) Delete(key string) *persistentScope {
	pos := m.pos(key)
	bucket := m.backingVector.Get(pos)
	if bucket != nil {
		newBucket := make(privatePersistentScopeItemBucket, 0)
		for _, item := range bucket {
			if item.Key != key {
				newBucket = append(newBucket, item)
			}
		}

		removedItemCount := len(bucket) - len(newBucket)
		if removedItemCount == 0 {
			return m
		}

		if len(newBucket) == 0 {
			newBucket = nil
		}

		newMap := &persistentScope{backingVector: m.backingVector.Set(pos, newBucket), len: m.len - removedItemCount}
		if newMap.backingVector.Len() > 1 && newMap.Len() < newMap.backingVector.Len()*int(lowerMapLoadFactor) {
			// Shrink backing vector if needed to avoid occupying excessive space
			buckets := newPrivatePersistentScopeItemBuckets(newMap.Len())
			buckets.AddItemsFromMap(newMap)
			return &persistentScope{backingVector: emptyPersistentScopeItemBucketVector.Append(buckets.buckets...), len: buckets.length}
		}

		return newMap
	}

	return m
}

// Range calls f repeatedly passing it each key and value as argument until either
// all elements have been visited or f returns false.
func (m *persistentScope) Range(f func(string, scopeEntry) bool) {
	m.backingVector.Range(func(bucket privatePersistentScopeItemBucket) bool {
		for _, item := range bucket {
			if !f(item.Key, item.Value) {
				return false
			}
		}
		return true
	})
}

// ToNativeMap returns a native Go map containing all elements of m.
func (m *persistentScope) ToNativeMap() map[string]scopeEntry {
	result := make(map[string]scopeEntry)
	m.Range(func(key string, value scopeEntry) bool {
		result[key] = value
		return true
	})

	return result
}

////////////////////
/// Constructors ///
////////////////////

// // NewPersistentScope returns a new persistentScope containing all items in items.
// func NewPersistentScope(items ...PersistentScopeItem) *persistentScope {
// 	return newPersistentScope(items)
// }

// // NewPersistentScopeFromNativeMap returns a new persistentScope containing all items in m.
// func NewPersistentScopeFromNativeMap(m map[string]scopeEntry) *persistentScope {
// 	buckets := newPrivatePersistentScopeItemBuckets(len(m))
// 	for key, value := range m {
// 		buckets.AddItem(PersistentScopeItem{Key: key, Value: value})
// 	}

// 	return &persistentScope{backingVector: emptyPersistentScopeItemBucketVector.Append(buckets.buckets...), len: buckets.length}
// }
