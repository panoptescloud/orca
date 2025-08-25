package slices_test

import (
	"testing"

	"github.com/panoptescloud/orca/pkg/slices"
	"github.com/stretchr/testify/assert"
)

type Element struct {
	name string
	prop string
}

func (self Element) GetName() string {
	return self.name
}

func Test_GetNamedElementIndex(t *testing.T) {
	tests := []struct {
		name   string
		in     []Element
		search string
		expect int
	}{
		{
			name:   "only value exists",
			search: "one",
			expect: 0,
			in: []Element{
				{
					name: "one",
				},
			},
		},

		{
			name:   "second value",
			search: "two",
			expect: 1,
			in: []Element{
				{
					name: "one",
				},
				{
					name: "two",
				},
				{
					name: "three",
				},
			},
		},

		{
			name:   "not found",
			search: "blah",
			expect: -1,
			in: []Element{
				{
					name: "one",
				},
				{
					name: "two",
				},
				{
					name: "three",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert.Equal(tt, test.expect, slices.GetNamedElementIndex(test.in, test.search))
		})
	}
}

func Test_NamedElementExists(t *testing.T) {
	tests := []struct {
		name   string
		in     []Element
		search string
		expect bool
	}{
		{
			name:   "only value exists",
			search: "one",
			expect: true,
			in: []Element{
				{
					name: "one",
				},
			},
		},

		{
			name:   "second value",
			search: "two",
			expect: true,
			in: []Element{
				{
					name: "one",
				},
				{
					name: "two",
				},
				{
					name: "three",
				},
			},
		},

		{
			name:   "not found",
			search: "blah",
			expect: false,
			in: []Element{
				{
					name: "one",
				},
				{
					name: "two",
				},
				{
					name: "three",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert.Equal(tt, test.expect, slices.NamedElementExists(test.in, test.search))
		})
	}
}

func Test_GetNamedElement(t *testing.T) {
	tests := []struct {
		name   string
		in     []Element
		search string
		expect *Element
	}{
		{
			name:   "only value exists",
			search: "one",
			expect: &Element{
				name: "one",
			},
			in: []Element{
				{
					name: "one",
				},
			},
		},

		{
			name:   "second value",
			search: "two",
			expect: &Element{
				name: "two",
			},
			in: []Element{
				{
					name: "one",
				},
				{
					name: "two",
				},
				{
					name: "three",
				},
			},
		},

		{
			name:   "not found",
			search: "blah",
			expect: nil,
			in: []Element{
				{
					name: "one",
				},
				{
					name: "two",
				},
				{
					name: "three",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			v := slices.GetNamedElement(test.in, test.search)

			if test.expect == nil {
				assert.Nil(tt, v)
				return
			}

			if v == nil {
				tt.Error("Expected non-nil result, got nil")
				return
			}

			assert.Equal(tt, *test.expect, *v)
		})
	}
}

func Test_UpsertNamedElement(t *testing.T) {
	tests := []struct {
		name    string
		initial []Element
		upsert  Element
		expect  []Element
	}{
		{
			name: "adds to empty list",
			expect: []Element{
				{
					name: "one",
					prop: "blah",
				},
			},
			upsert: Element{
				name: "one",
				prop: "blah",
			},
			initial: []Element{},
		},

		{
			name: "appends to list",
			expect: []Element{
				{
					name: "one",
					prop: "blah",
				},
				{
					name: "two",
					prop: "blah",
				},
			},
			upsert: Element{
				name: "two",
				prop: "blah",
			},
			initial: []Element{
				{
					name: "one",
					prop: "blah",
				},
			},
		},

		{
			name: "replace existing element",
			expect: []Element{
				{
					name: "one",
					prop: "blah",
				},
				{
					name: "two",
					prop: "meh",
				},
				{
					name: "three",
					prop: "blah",
				},
			},
			upsert: Element{
				name: "two",
				prop: "meh",
			},
			initial: []Element{
				{
					name: "one",
					prop: "blah",
				},
				{
					name: "two",
					prop: "blah",
				},
				{
					name: "three",
					prop: "blah",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			result := slices.UpsertNamedElement(test.initial, test.upsert)

			assert.Equal(tt, test.expect, result)
		})
	}
}
