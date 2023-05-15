package orderedmap

import (
	"fmt"
	"testing"

	"github.com/elliotchance/orderedmap"
	"github.com/gogf/gf/v2/util/gconv"
	orderedmap1 "github.com/iancoleman/orderedmap"
)

func Test_orderedMap(t *testing.T) {
	var mp []*orderedmap1.OrderedMap
	for i := 0; i < 10; i++ {
		m := orderedmap1.New()
		m.Set("name", "apple"+gconv.String(i))
	}
	fmt.Println(mp)
}

func Test_orderedMap1(t *testing.T) {
	// 创建一个有序Map
	m := orderedmap.NewOrderedMap()

	// 添加键值对
	m.Set("key3", 3)
	m.Set("key1", 1)
	m.Set("key2", 2)

	// 遍历有序Map


	// 获取指定键的值
	value, exists := m.Get("key2")
	if exists {
		fmt.Println("Value of key2:", value)
	} else {
		fmt.Println("Key2 not found")
	}

	// 删除指定键的键值对
	m.Delete("key1")

	// 遍历更新后的有序Map

}
