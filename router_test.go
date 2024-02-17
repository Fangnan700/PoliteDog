package PoliteDog

import (
	"fmt"
	"testing"
)

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.GET("/user/info", nil)

	g := NewRouterGroup("admin")
	g.POST("/login/:id", nil)

	fmt.Printf("%+v\n", g.Routers.RouterTrie.next.children[0].children[0])
	fmt.Printf("%+v\n", g.Routers.RouterTrie.next.Search("/admin/login/1"))
}
