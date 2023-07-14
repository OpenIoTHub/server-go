// Package runtimeStorage 运行时存储数据接口
package runtimeStorage

// RuntimeStorage 保存运行时数据的基础接口，用于提供保存http proxy数据或者登录黑白名单等配置
// 可以基于此接口做redis或者内存等的实现
type RuntimeStorageIfce interface {
	GetAllKeys() ([]string, error)

	GetValueByKey(key string) (any, error)
	GetValueByKeyToBytes(key string) ([]byte, error)
	GetValueByKeyToString(key string) (string, error)
	GetValueByKeyToBool(key string) (bool, error)
	GetValueByKeyToInt(key string) (int, error)

	SetValueByKey(key string, value any) error
	DelValueByKey(key string) error
	CheckKeyExit(key string) (bool, error)
}
