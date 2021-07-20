package logic

import (
	"container/ring"
	"testing"
)

func TestRun(t *testing.T) {
	n:=5
	r:=ring.New(n)
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}
	/*r.Do(func(i interface{}) {
		t.Logf("%d\n", i.(int))
	})*/
	nr:=&ring.Ring{Value: 6}
	r.Link(nr)
	// n%r.Len()
	//t.Log(r.Link(nr).Value)
	/*r.Link(nr).Do(func(i interface{}) {
		//t.Logf("%d\n", i.(int))
	})*/
	//t.Log(r.Move(3).Value)
	//return
	cr:=r.Unlink(6)
	cr.Do(func(i interface{}) {
		t.Log(i)
	})
	//r.Link(r.Move(3))
}