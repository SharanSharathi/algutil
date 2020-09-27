package quicksort

type partitionList struct {
	start int
	end   int
	next  *partitionList
}

// Non-recursive implementation.
func Ints(a []int) {
	p := &partitionList{
		start: 0,
		end:   len(a), // not inclusive
	}

	// used for appending next partition
	lastPartition := p

	for ; p != nil; p = p.next {
		if p.end-p.start <= 1 {
			continue
		}

		// taking last element as pivot
		// since end is non-inclusive, end = end - 1
		pivotIndex := p.end - 1

		// used to swap elements less than pivot's value
		swapIndex := p.start

		for i := p.start; i < p.end-1; i++ {
			if a[i] < a[pivotIndex] {
				// swapping smaller elements
				a[i], a[swapIndex] = a[swapIndex], a[i]
				swapIndex++
			}
		}

		// swapping pivot
		a[pivotIndex], a[swapIndex] = a[swapIndex], a[pivotIndex]
		pivotIndex = swapIndex

		leftPartition := &partitionList{
			start: p.start,
			end:   pivotIndex,
		}

		rightPartition := &partitionList{
			start: pivotIndex + 1,
			end:   p.end,
		}

		lastPartition.next = leftPartition
		leftPartition.next = rightPartition
		lastPartition = rightPartition
	}
}
