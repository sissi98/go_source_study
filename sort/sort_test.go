package sort

import (
	"fmt"
	"sort"
	"testing"
)

// 自定义类型
type Person struct {
	Name string
	Age  int
}

// 实现 sort.Interface 接口的 Len、Less 和 Swap 方法
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func Test_SortReverse(t *testing.T) {
	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"Dave", 35},
	}

	fmt.Println("Before sorting:")
	for _, person := range people {
		fmt.Println(person)
	}

	// 使用 Reverse 函数进行排序
	sort.Sort(sort.Reverse(ByAge(people)))

	fmt.Println("\nAfter sorting in reverse order:")
	for _, person := range people {
		fmt.Println(person)
	}
}
func Test_Sort(t *testing.T) {
	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"Barlie", 20},
		{"harlie", 20},
		{"Dave0", 35},
		{"Dave1", 35},
		{"Dave2", 35},
		{"Dave3", 35},
		{"Dave4", 60},
		{"Dave5", 35},
		{"Dave6", 35},
		{"Dave7", 35},
		{"Dave8", 80},
		{"Dave9", 35},
		{"Dave10", 35},
	}

	fmt.Println("Before sorting:")
	for _, person := range people {
		fmt.Println(person)
	}

	// 使用 Reverse 函数进行排序
	sort.Sort(ByAge(people))

	fmt.Println("\nAfter sorting :")
	for _, person := range people {
		fmt.Println(person)
	}
}
func Test_Stable(t *testing.T) {
	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 20},
		{"Barlie", 20},
		{"harlie", 20},
		{"Dave0", 35},
		{"Dave1", 35},
		{"Dave2", 35},
		{"Dave3", 35},
		{"Dave4", 60},
		{"Dave5", 35},
		{"Dave6", 35},
		{"Dave7", 35},
		{"Dave8", 80},
		{"Dave9", 35},
		{"Dave10", 35},
	}

	fmt.Println("Before sorting:")
	for _, person := range people {
		fmt.Println(person)
	}

	// 使用 Reverse 函数进行排序
	sort.Stable(ByAge(people))

	fmt.Println("\nAfter sorting :")
	for _, person := range people {
		fmt.Println(person)
	}
}
