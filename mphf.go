package mphf

import "sort"

// Hash computes the minimal perfect hash value for the given string and offset value d
func Hash(d int, str string) int {
	if d == 0 {
		d = 0x811c9dc5
	}

	for i := 0; i < len(str); i++ {
		// http://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function
		// http://isthe.com/chongo/src/fnv/hash_32.c
		// multiply by the 32 bit FNV magic prime mod 2^32
		d += (d << 1) + (d << 4) + (d << 7) + (d << 8) + (d << 24)
		// xor the bottom with the current octet
		d ^= int(str[i])
	}

	return d & 0x7fffffff
}

type bucket struct {
	keys []string
}

type bucketslice []*bucket

func (a bucketslice) Len() int           { return len(a) }
func (a bucketslice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bucketslice) Less(i, j int) bool { return len(a[i].keys) > len(a[j].keys) }

// MinHashTable is a wrapper for the minimal perfect hash tables
type MinHashTable struct {
	G []int
	V []*int
}

// Lookup the given key in the minimal perfect hash table
func (table *MinHashTable) Lookup(key string) *int {
	d := table.G[Hash(0, key)%len(table.G)]
	var value *int
	if d < 0 {
		value = table.V[0-d-1]
	} else {
		value = table.V[Hash(d, key)%len(table.V)]
	}

	if value != nil {
		return value
	}

	return new(int)
}

// Create the minimal perfect hash table based on the given dictionary
func Create(dict map[string]int) *MinHashTable {
	size := len(dict)
	buckets := make(map[int]*bucket)
	g := make([]int, size)
	v := make([]*int, size)

	// Place all of the keys into buckets
	for k := range dict {
		bkey := Hash(0, k) % size
		if b, ok := buckets[bkey]; ok {
			b.keys = append(b.keys, k)
		} else {
			buckets[bkey] = &bucket{keys: []string{k}}
		}
	}

	// Sort the buckets and process the ones with the most items first
	var sortedBuckets []*bucket
	for _, bucket := range buckets {
		sortedBuckets = append(sortedBuckets, bucket)
	}
	sort.Sort(bucketslice(sortedBuckets))
	buckets = nil

	b := 0
	for ; b < int(size); b++ {
		if len(sortedBuckets[b].keys) <= 1 {
			break
		}

		bucket := sortedBuckets[b]
		d := 1
		item := 0
		slots := []int{}
		used := make(map[int]struct{})
		var empty struct{}

		for item < len(bucket.keys) {
			slot := Hash(d, bucket.keys[item]) % size
			_, ok := used[slot]
			if v[slot] != nil || ok {
				d++
				item = 0
				slots = []int{}
				used = make(map[int]struct{})
			} else {
				used[slot] = empty
				slots = append(slots, slot)
				item++
			}
		}

		g[Hash(0, bucket.keys[0])%size] = d
		for i := 0; i < len(bucket.keys); i++ {
			value := new(int)
			*value = dict[bucket.keys[i]]
			v[slots[i]] = value
		}
	}

	// Only buckets with 1 item remain. Process them more quickly by directly
	// placing them into a free slot. Use a negative value of d to indicate this.

	freelist := []int{}
	for i := 0; i < int(size); i++ {
		if v[i] == nil {
			freelist = append(freelist, i)
		}
	}

	count := len(freelist) - 1
	for ; b < int(size); b++ {
		if len(freelist) == 0 || b >= len(sortedBuckets) || sortedBuckets[b] == nil || len(sortedBuckets[b].keys) == 0 {
			break
		}

		bucket := sortedBuckets[b]
		slot := freelist[count]
		count--
		if count > 0 {
			freelist = freelist[:count]
		} else {
			freelist = []int{}
		}

		// We subtract one to ensure it's negative even if the zeroith slot was used
		g[Hash(0, bucket.keys[0])%size] = 0 - slot - 1
		value := new(int)
		*value = dict[bucket.keys[0]]
		v[slot] = value
	}

	return &MinHashTable{G: g, V: v}
}
