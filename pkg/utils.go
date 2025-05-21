package pkg

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomFloat 生成 [min, max) 范围内的随机浮点数
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// 工具函数可在此文件实现
