package io_reverse

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

//给您道面试题：写入文件时，对数据做反转；读取文件时，再将数据反转回来。例如：数据为[]byte{5,4,3,2,1}，磁盘上的数据应为12345

// ReverseOpt 数据操作
type ReverseOpt struct{}

// ReverseBytes byte切片翻转
func (opt *ReverseOpt) ReverseBytes(data []byte) []byte {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

func (opt *ReverseOpt) Read(filename string) ([]byte, error) {
	fRead, err := os.Open(filename)
	defer fRead.Close()
	if err != nil {
		return nil, err
	}

	fInfo, err := fRead.Stat()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(fRead)
	data := make([]byte, fInfo.Size())
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}
	return opt.ReverseBytes(data), nil
}

func (opt *ReverseOpt) Write(filename string, data []byte) error {
	fWrite, err := os.Create(filename)
	defer fWrite.Close()

	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fWrite)
	_, err = writer.Write(opt.ReverseBytes(data))
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// CreateByteFile 生成包含制定个byte的测试文件
func (opt *ReverseOpt) CreateByteFile(filename string, num int) error {
	data := make([]byte, num)

	source := rand.NewSource(time.Now().UnixNano())
	rIns := rand.New(source)

	for i := 0; i < num; i++ {
		data[i] = byte(rIns.Intn(256))
	}
	err := opt.Write(filename, data)
	if err != nil {
		return err
	}
	return nil
}
