package main

import "testing"

func TestDeduplicate(t *testing.T) {
	t.Parallel()

	cases := []struct{
		got path
		want path
		name string
	} {
		{
			path{},
			path{},
			"empty",
		},
		{
			path{step{1, true}},
			path{step{1, true}},
			"one element",
		},
		{
			path{step{1, true}, step{2, false}, step{2, true}},
			path{step{1, true}, step{2, true}},
			"same steps at the end with true last step",
		},
		{
			path{step{1, true}, step{2, false}, step{2, false}},
			path{step{1, true}, step{2, false}},
			"same steps at the end with false last step",
		},
		{
			path{step{1, true}, step{2, true}, step{2, false}},
			path{step{1, true}, step{2, true}},
			"same steps at the end with true step",
		},
		{
			path{step{1, true}, step{2, true}, step{2, false}, step{3, false}},
			path{step{1, true}, step{2, true}, step{3, false}},
			"same steps at the middle",
		},
		{
			path{step{1, true}, step{2, false}, step{2, true}, step{2, false}, step{3, false}},
			path{step{1, true}, step{2, true}, step{3, false}},
			"same steps at the middle, true in the middle",
		},
		{
			path{step{2, true}, step{2, false}, step{1, true}},
			path{step{2, true}, step{1, true}},
			"same steps at the start",
		},
		{
			path{step{1, true}, step{1, false}, step{0, true}, step{2, true}, step{2, false}, step{0, true}, step{3, false}, step{3, false}},
			path{step{1, true}, step{0, true}, step{2, true}, step{0, true}, step{3, false}},
			"same steps at the start, middle, end",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			c.got.deduplicate()
			if len(c.want) != len(c.got) {
				t.Fatalf("decuplication went wrong.\nGot:\t%v\nWant:\t%v\n", c.got, c.want)
			}
		})
	}
}