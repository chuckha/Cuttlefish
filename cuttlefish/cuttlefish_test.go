package cuttlefish

import (
	"testing"
)

func TestGetUrlDiesOnError(t *testing.T) {
	t.Log("HELLO WORLD")
	csite := make(chan Site)
	death := make(chan struct{})

	t.Log("Before get")
	GetUrl([]byte("http://localhost"), csite, death)
	t.Log("After get")

	x := <-death
	if x != struct{}{} {
		t.Errorf("Did not die")
	}

}

//func TestGetUrl
//	site := <-csite
//	if site.Url != nil {
//		t.Errorf("Site did not return")
//	}
