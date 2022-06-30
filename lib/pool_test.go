package lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	yes bool
}

func TestGenericPool(t *testing.T) {
	results := make(chan string, 10)
	p := NewPool[testStruct](10, func(par testStruct) {
		results <- fmt.Sprintf("yes? %t", par.yes)
	})
	p.Start()

	for i := 0; i < 10; i++ {
		p.Do(testStruct{yes: true})
	}

	<-p.Done()
	close(results)

	for res := range results {
		fmt.Println(res)
	}

	assert.Equal(t, p, p, "good")
}
