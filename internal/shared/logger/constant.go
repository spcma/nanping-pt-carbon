package logger

// 默认日志-V1
var loggersv1 = []*CategoryConfig{
	{
		Category: "runtime",
		Level:    "info",
	},
	{
		Category: "debug",
		Level:    "info",
	},
	{
		Category: "error",
		Level:    "info",
	},
	{
		Category: "traffic",
		Level:    "info",
	},
	{
		Category: "iam",
		Level:    "info",
	},
	{
		Category: "http",
		Level:    "info",
	},
	{
		Category: "initializer",
		Level:    "info",
	},
	{
		Category: "ipfs",
		Level:    "info",
	},
	{
		Category: "ipfsrate",
		Level:    "info",
	},
	{
		Category: "scheduler",
		Level:    "info",
	},
}
