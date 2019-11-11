package main

import (
	"fmt"
	"strings"
	"time"
)

/*import (
	"fmt"
	db "gochat/database"
)

func main() {
	/*list := make([]int, 8)
	list = append(list, 3)
	list2 := make([]int, 8, 8)
	list2 = append(list2, 3)
	list3 := make([]int, 0, 8)
	list3 = append(list3, 3)
	fmt.Println(list)
	fmt.Println(list2)
	fmt.Println(list3)*/
/*fmt.Println(db.VarifyUserByPasswd("c1", "1"))
}*/


type S struct {
	a, b int
}

// String implements the fmt.Stringer interface
func (s *S) String() string {
	return fmt.Sprintf("%s", s) // Sprintf will call s.String()
}

func main() {
	//s := &S{a: 1, b: 2}
	//fmt.Println(s)
	//runtime.Set
	start := time.Now()
	builder := strings.Builder{}
	for i := 0; i < 100000; i++ {
		builder.Write([]byte("hello"))
		builder.Write([]byte(" world!"))
	}
	_ = builder.String()
	fmt.Println(time.Since(start).Milliseconds())
	start = time.Now()
	str := ""
	for i := 0; i < 100000; i++ {
		str = str + "hello"
		str = str + " world!"
	}
	fmt.Println(time.Since(start).Milliseconds())
}
