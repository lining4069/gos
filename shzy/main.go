package main

import (
	"bufio"
	"log"
	"os"
)

//给您道面试题：写入文件时，对数据做反转；读取文件时，再将数据反转回来。例如：数据为[]byte{5,4,3,2,1}，磁盘上的数据应为12345

type ReverseOpt struct{}

// ReverseBytes byte切片翻转
func (opt *ReverseOpt) ReverseBytes(data []byte) []byte {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}

// ByteRead  复用读物byte数据
func (opt *ReverseOpt) ByteRead(filename string) ([]byte, error) {
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
	return data, nil
}

// Read 读取
func (opt *ReverseOpt) Read(filename string) ([]byte, error) {
	data, err := opt.ByteRead(filename)
	if err != nil {
		return nil, err
	}
	return opt.ReverseBytes(data), nil
}

// Write 存储
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

func main() {
	opt := ReverseOpt{}
	data := []byte{1, 3, 5, 7, 9}
	log.Printf("源数据：%v", data)
	filename := "./reverse_test.bin"
	wErr := opt.Write(filename, data)
	if wErr != nil {
		log.Fatalf("写数据时发生异常 %v", wErr)
	}
	ValidReadData, vErr := opt.ByteRead(filename)
	if vErr != nil {
		log.Fatalf("读数据时发生异常 %v", vErr)
	}
	log.Printf("文件中存储的byte数据正序读取为：%v", ValidReadData)
	rReadData, rRrr := opt.Read(filename)
	if rRrr != nil {
		log.Fatalf("读数据时发生异常 %v", wErr)
	}
	log.Printf("读取出的数据为:%v", rReadData)

}
