package main


import (
	"time"
	"github.com/patrickmn/go-cache"
	"fmt"
)

func main() {
	c := cache.New(5*time.Minute, 10 *time.Minute)
	c.Set("foo", "bar", cache.DefaultExpiration)
	c.Set("bar", 42, cache.NoExpiration)
	foo, found := c.Get("foo")
	if found {
		fmt.Println(foo)
	}


}
