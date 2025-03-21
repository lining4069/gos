package io_reverse

import (
	"bufio"
	"fmt"
	"io"
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

// ReverseFile 翻转文件内容
func (opt *ReverseOpt) ReverseFile(inputPath, outputPath string, chunkSize int) error {
	// 打开输入文件
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开输入文件: %w", err)
	}
	defer inputFile.Close()

	// 获取输入文件的大小
	fileInfo, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("无法获取文件信息: %w", err)
	}
	fileSize := fileInfo.Size()

	// 打开输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %w", err)
	}
	defer outputFile.Close()

	// 初始化输出文件的大小
	err = outputFile.Truncate(fileSize)
	if err != nil {
		return fmt.Errorf("无法调整输出文件大小: %w", err)
	}

	// 分块处理文件
	for offset := fileSize; offset > 0; offset -= int64(chunkSize) {
		// 计算当前块的大小
		if offset < int64(chunkSize) {
			chunkSize = int(offset)
		}
		// 定位到当前块的起始位置
		_, err = inputFile.Seek(offset-int64(chunkSize), io.SeekStart)
		if err != nil {
			return fmt.Errorf("无法定位到文件位置: %w", err)
		}
		// 读取当前块的数据
		chunk := make([]byte, chunkSize)
		_, err = bufio.NewReader(inputFile).Read(chunk)
		if err != nil {
			return fmt.Errorf("无法读取文件块: %w", err)
		}
		// 反转当前块的数据
		reversedChunk := opt.ReverseBytes(chunk)
		// 定位到输出文件的相应位置
		_, err = outputFile.Seek(fileSize-offset, io.SeekStart)
		if err != nil {
			return fmt.Errorf("无法定位到输出文件位置: %w", err)
		}
		// 将反转后的数据写入输出文件
		_, err = outputFile.Write(reversedChunk)
		if err != nil {
			return fmt.Errorf("无法写入文件块: %w", err)
		}
	}
	return nil
}
