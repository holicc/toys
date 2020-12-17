package main

import (
	"fmt"
	"gopkg.in/music-theory.v0/scale"
	"testing"
)

func TestScale(t *testing.T) {
	sc := scale.Of("C minor")
	fmt.Println(sc.ToYAML())
}
