package pipeline

import (
	"encoding/binary"
	"io"
	"math/rand"
	"sort"
)

func ArraySource(a ...int) <-chan int {
	out := make(chan int) // 无缓冲
	go func() {
		for _, val := range a {
			out <- val
		}
		close(out)
	}()
	return out
}

func InMemSort(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		a := []int{}
		for v := range in {
			a = append(a, v)
		}

		sort.Ints(a)

		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

// Merge 将两个排序好的chan 进行merge合并
func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
	}()
	return out
}

// MergeN N条chan数据归并
func MergeN(pips ...<-chan int) <-chan int {
	if len(pips) == 1 {
		return pips[0]
	}
	m := len(pips) / 2
	return Merge(MergeN(pips[:m]...), MergeN(pips[m:]...))

}

// ReaderSource 文件io 获取数据
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int)
	go func() {
		buffer := make([]byte, 8)
		bytesSize := 0
		for {
			n, err := reader.Read(buffer)
			bytesSize += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && bytesSize > chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		_, err := writer.Write(buffer)
		if err != nil {
			return
		}
	}
}

func RandomSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}
