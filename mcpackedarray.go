// CREATED BY RPH
//
// RELEASED TO THE PUBLIC DOMAIN. IF NOT POSSIBLE UNDER LOCAL LAW, LICENSED UNDER THE UNLICENSE.

package mcpackedarray

import (
	"math"
	"sync"
)

type PackedArray struct {
	BitsPerEntry byte
	EntryAmount	 int32
	Entries		 []uint32
	mut          sync.Mutex
}

func (p *PackedArray) lock() {
	p.mut.Lock()
}

func (p *PackedArray) unlock() {
	p.mut.Unlock()
}


func NewPackedArray(bitsPerEntry byte, entryAmount int32) *PackedArray {
	var array = new(PackedArray)
	array.lock()
	array.BitsPerEntry = bitsPerEntry
	array.Entries = make([]uint32, entryAmount)
	array.EntryAmount = entryAmount

	array.unlock()
	return array
}

func (p *PackedArray) Set(id int32, value uint32) {
	p.lock()
	p.Entries[id] = value
	p.unlock()
}

func (p *PackedArray) Get(id int32) uint32 {
	return p.Entries[id]
}

func (p *PackedArray) Serialise() []byte {
	amountOfEntries := int(math.Ceil((float64(p.EntryAmount) * float64(p.BitsPerEntry)) / float64(8)))
	required := amountOfEntries % 8
	if required != 0 {
		amountOfEntries += 8 - required
	}

	output := make([]byte, amountOfEntries)

	index := 0
	bitsLeft := 0
	indexOffset := 7

	for i := int32(0); i < p.EntryAmount; i++ {
		entry := p.Entries[i] % uint32(math.Pow(2, float64(p.BitsPerEntry)))

		for j := 0; j < int(p.BitsPerEntry); j++ {
			if entry & uint32(math.Pow(2, float64(j))) > 0 {
				output[index+indexOffset] = output[index+indexOffset] + byte(math.Pow(2, float64(bitsLeft)))
			}
			bitsLeft++
			if bitsLeft > 7 {
				indexOffset--
				bitsLeft = 0
			}

			if indexOffset < 0 {
				indexOffset = 7
				index += 8
			}
		}
	}

	return output
}

func PackedArrayFromData(data []byte, bitsPerEntry byte) *PackedArray {
	if len(data) % 8 != 0 {
		panic("Invalid packed array! [len(data) % 8] must be equal to 0!")
	}

	entryAmount := int32(
		math.Floor(
			float64(len(data) * 8) / float64(bitsPerEntry)))
	pa := NewPackedArray(bitsPerEntry, entryAmount)

	index := 0
	bitsLeft := 0
	indexOffset := 7


	for i := int32(0); i < entryAmount; i++ {
		value := uint32(0)
		for j := 0; j < int(bitsPerEntry); j++ {
			thisByte := data[index + indexOffset]
			if thisByte & byte(math.Pow(2, float64(bitsLeft))) > 0 {
				value = value + uint32(math.Pow(2, float64(j)))
			}

			bitsLeft++
			if bitsLeft > 7 {
				indexOffset--
				bitsLeft = 0
			}

			if indexOffset < 0 {
				indexOffset = 7
				index += 8
			}
		}
		pa.Set(i, value)
	}

	return pa
}
