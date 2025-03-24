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

func (opt *ReverseOpt) ReverseFile(inputPath string) error {
	// 打开输入文件
	inputFile, err := os.OpenFile(inputPath, os.O_RDWR, 0666)
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

	left := int64(0)
	right := fileSize - 1
	for left < right {
		// 读取左侧字节
		_, err = inputFile.Seek(left, io.SeekStart)
		if err != nil {
			return fmt.Errorf("无法定位到文件左侧位置: %w", err)
		}
		leftByte := make([]byte, 1)
		_, err = bufio.NewReader(inputFile).Read(leftByte)
		if err != nil {
			return fmt.Errorf("无法读取文件左侧字节: %w", err)
		}

		// 读取右侧字节
		_, err = inputFile.Seek(right, io.SeekStart)
		if err != nil {
			return fmt.Errorf("无法定位到文件右侧位置: %w", err)
		}
		rightByte := make([]byte, 1)
		_, err = bufio.NewReader(inputFile).Read(rightByte)
		if err != nil {
			return fmt.Errorf("无法读取文件右侧字节: %w", err)
		}

		// 交换左右字节
		_, err = inputFile.Seek(left, io.SeekStart)
		if err != nil {
			return fmt.Errorf("无法定位到文件左侧位置进行写入: %w", err)
		}
		_, err = inputFile.Write(rightByte)
		if err != nil {
			return fmt.Errorf("无法写入文件左侧字节: %w", err)
		}

		_, err = inputFile.Seek(right, io.SeekStart)
		if err != nil {
			return fmt.Errorf("无法定位到文件右侧位置进行写入: %w", err)
		}
		_, err = inputFile.Write(leftByte)
		if err != nil {
			return fmt.Errorf("无法写入文件右侧字节: %w", err)
		}

		left++
		right--
	}
	return nil
}
