package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStatusPorcelain(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected workingTreeChanges
		display  string
	}{
		{
			name:     "clean",
			output:   "",
			expected: workingTreeChanges{},
			display:  "clean",
		},
		{
			name:     "staged only",
			output:   "M  a.go\nA  b.go\nD  c.go\nR  old.go -> new.go\n",
			expected: workingTreeChanges{Staged: 4},
			display:  "4 staged",
		},
		{
			name:     "unstaged only",
			output:   " M a.go\n D b.go\n",
			expected: workingTreeChanges{Unstaged: 2},
			display:  "2 unstaged",
		},
		{
			name:     "staged and unstaged in the same file",
			output:   "MM a.go\n",
			expected: workingTreeChanges{Staged: 1, Unstaged: 1},
			display:  "1 staged\n1 unstaged",
		},
		{
			name:     "untracked",
			output:   "?? new.go\n?? other/\n",
			expected: workingTreeChanges{Untracked: 2},
			display:  "2 untracked",
		},
		{
			name:     "conflicts",
			output:   "UU a.go\nAA b.go\nDD c.go\nAU d.go\nUD e.go\n",
			expected: workingTreeChanges{Conflicted: 5},
			display:  "5 conflicted",
		},
		{
			name:     "mixed",
			output:   "UU a.go\nM  b.go\n M c.go\n?? d.go\n",
			expected: workingTreeChanges{Staged: 1, Unstaged: 1, Untracked: 1, Conflicted: 1},
			display:  "1 conflicted\n1 staged\n1 unstaged\n1 untracked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := parseStatusPorcelain(tt.output)

			assert.Equal(t, tt.expected, changes)
			assert.Equal(t, tt.display, changes.String())
		})
	}
}
