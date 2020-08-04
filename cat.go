package ktail

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/service/kinesis"
)

func (app *App) Cat(ctx context.Context, partitionKey string, src io.Reader) error {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		b := scanner.Bytes()
		if app.AppendLF {
			b = append(b, LF...)
		}
		in := &kinesis.PutRecordInput{
			Data:       b,
			StreamName: &app.StreamName,
		}
		if partitionKey == "" {
			pk := fmt.Sprintf("%x", sha256.Sum256(b))
			in.PartitionKey = &pk
		} else {
			in.PartitionKey = &partitionKey
		}
		if _, err := app.kinesis.PutRecord(in); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
