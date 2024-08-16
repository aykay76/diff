package diff

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io"
)

type hashBlock struct {
	hash             string
	marked           bool
	line             string
	sourceLineNumber int
	numberOfLines    int
}

type Diff struct {
	sourceHashList []*hashBlock
	targetHashList []*hashBlock
	different      bool
}

func NewDiff(a io.Reader, b io.Reader) *Diff {
	diff := &Diff{}

	// create hash blocks
	diff.sourceHashList = diff.createHashList(a)
	diff.targetHashList = diff.createHashList(b)

	i1 := 0
	i2 := 0
	for i1 < len(diff.sourceHashList) {
		max := 0
		length := 0
		max_start := -1
		hb1 := diff.sourceHashList[i1]

		i2 = 0
		for i2 < len(diff.targetHashList) {
			hb2 := diff.targetHashList[i2]

			if !hb2.marked {
				if hb1.hash == hb2.hash {
					length = 1
					i12 := i1 + 1
					i22 := i2 + 1
					for i12 < len(diff.sourceHashList) && i22 < len(diff.targetHashList) && !diff.targetHashList[i22].marked {
						if diff.sourceHashList[i12].hash == diff.targetHashList[i22].hash {
							length++
						} else {
							break
						}

						i12++
						i22++
					}

					if length > max {
						max_start = i2
						max = length
					}
				}
			}

			i2++
		}

		if max_start == 0 && max == len(diff.sourceHashList) && max == len(diff.targetHashList) {
			diff.different = false
			return diff
		}

		if max_start == -1 {
			i1++
		} else {
			i2 = max_start
			hb2 := diff.targetHashList[i2]
			hb2.sourceLineNumber = hb1.sourceLineNumber
			hb2.numberOfLines = max

			for i := 0; i < max; i++ {
				diff.sourceHashList[i1].marked = true
				diff.targetHashList[i2].marked = true
				i1++
				i2++
			}
		}
	}

	return diff
}

func (diff *Diff) createHashList(input io.Reader) (hashes []*hashBlock) {
	hashes = make([]*hashBlock, 0)

	i := 0
	reader := bufio.NewReader(input)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		checksum := md5.Sum([]byte(line))
		hash := hex.EncodeToString(checksum[:])

		block := &hashBlock{
			hash:             hash,
			line:             line,
			sourceLineNumber: i,
		}
		i++

		hashes = append(hashes, block)
	}
}
