package geom_test

import (
	"strings"
	"testing"

	"calendar-puzzle/geom"
)

var strVariants = [4]string{
	"**\n* ",
	"**\n *",
	"* \n**",
	" *\n**",
}

func TestVariants(t *testing.T) {
	points := []geom.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 1}}
	shape := geom.NewShape(true, 3, points)

	variants := geom.Variants(shape)

	if len(variants) != 8 {
		t.Errorf("Expected 8 variants, got %d", len(variants))
	}

	strVariantToCount := make(map[string]int, len(variants)/2)

	for _, v := range strVariants {
		strVariantToCount[v] = 0
	}

	for _, variant := range variants {
		if len(variant) != len(points) {
			t.Errorf("Expected variant length %d, got %d", len(points), len(variant))
		}

		convert := func(p geom.Point) string {
			for _, point := range variant {
				if point == p {
					return string("*")
				}
			}
			return string(" ")
		}

		vstr := strings.Join(geom.Stringify(variant, convert), "\n")

		if _, ok := strVariantToCount[vstr]; ok {
			strVariantToCount[vstr]++
		} else {
			t.Errorf("Unexpected variant: \n%s", vstr)
		}
	}

	for _, v := range strVariants {
		if strVariantToCount[v] != 2 {
			t.Errorf("Expected 2 occurrences of \n%s, got %d", v, strVariantToCount[v])
		}
	}
}
