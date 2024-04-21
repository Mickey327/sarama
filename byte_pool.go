package sarama

import "github.com/IBM/sarama/mickey/bytepool"

type BytePoolOptions = bytepool.BytePoolOptions

var (
	DefaultBytePoolOptions      = bytepool.DefaultBytePoolOptions
	ConfigureDecompressBytePool = bytepool.ConfigureDecompressBytePool
)
