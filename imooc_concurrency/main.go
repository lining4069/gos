package main

import (
	"bufio"
	"github.com/lining4069/gos/imooc_concurrency/pipeline"
	"os"
)

func OutSort(filename, sortedFilename string, randomCount, pipeNum int) {

	fileWrite, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fileWrite.Close()

	p := pipeline.RandomSource(randomCount)
	writer := bufio.NewWriter(fileWrite)
	pipeline.WriterSink(writer, p)
	writer.Flush() // 刷新缓冲区以确保所有数据都被写入文件

	// 第二步：从文件读取数据，进行排序，并将结果保存到新文件
	fileRead, errRead := os.Open(filename)
	if errRead != nil {
		panic(errRead)
	}
	defer fileRead.Close()

	sortedFileWrite, errSorted := os.Create(sortedFilename)
	if errSorted != nil {
		panic(errSorted)
	}
	defer sortedFileWrite.Close()

	pipes := []<-chan int{}
	for i := 0; i < pipeNum; i++ {
		pipes = append(pipes, pipeline.InMemSort(pipeline.ReaderSource(bufio.NewReader(fileRead), randomCount/pipeNum*8)))
	}
	p = pipeline.MergeN(pipes...)
	pipeline.WriterSink(fileWrite, p)
}
