package getx

import (
	"bytes"
	"testing"

	"github.com/desal/richtext"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestSingleRepo(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s2", "gh/u1/p1"))
	repos.AddRepo("gh/u2/p1",
		Pkg("gh/u2/p1", "gh/u1/p1/s1"))
	repos.AddRepo("gh/u2/p2",
		Pkg("gh/u2/p2", "gh/u1/p1/s1", "gh/u1/p1", "gh/u1/p1/s2"))

	buf := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		ctx := New(richtext.Debug(buf), goPath, ruleSet, "", Verbose, MustPanic, Install, RecurseTopLevel)
		ctx.Get(".", "gh/u1/p1/s2", false, false)
	})

	expected := stringSet{
		"./pkg/gh/u1/p1.a":         empty{},
		"./pkg/gh/u1/p1/s1.a":      empty{},
		"./pkg/gh/u1/p1/s2.a":      empty{},
		"./src/gh/u1/p1/gen.go":    empty{},
		"./src/gh/u1/p1/s1/gen.go": empty{},
		"./src/gh/u1/p1/s2/gen.go": empty{},
	}
	assert.Equal(t, expected, fileList)
	assert.Equal(t, "gh/u1/p1\n", buf.String())

}

func TestDependentRepo(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s2", "gh/u1/p1"))
	repos.AddRepo("gh/u2/p1",
		Pkg("gh/u2/p1", "gh/u1/p1/s1"))
	repos.AddRepo("gh/u2/p2",
		Pkg("gh/u2/p2", "gh/u1/p1/s1", "gh/u1/p1", "gh/u1/p1/s2"))

	buf := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		ctx := New(richtext.Debug(buf), goPath, ruleSet, "", Verbose, MustPanic, Install, DeepScan)
		ctx.Get(".", "gh/u2/p2", false, false)
	})

	expected := stringSet{
		"./pkg/gh/u1/p1.a":         empty{},
		"./pkg/gh/u1/p1/s1.a":      empty{},
		"./pkg/gh/u1/p1/s2.a":      empty{},
		"./pkg/gh/u2/p2.a":         empty{},
		"./src/gh/u1/p1/gen.go":    empty{},
		"./src/gh/u1/p1/s1/gen.go": empty{},
		"./src/gh/u1/p1/s2/gen.go": empty{},
		"./src/gh/u2/p2/gen.go":    empty{},
	}

	assert.Equal(t, expected, fileList)
	assert.Equal(t, `gh/u1/p1/s1
gh/u1/p1
gh/u1/p1/s2
gh/u2/p2
`, buf.String())
}

func TestMultiDepOk(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s2", "gh/u1/p1"))
	repos.AddRepo("gh/u2/p1",
		Pkg("gh/u2/p1", "gh/u1/p1/s1"))
	repos.AddRepo("gh/u2/p2",
		Pkg("gh/u2/p2", "gh/u1/p1/s1", "gh/u1/p1", "gh/u2/p1"))

	buf := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		ctx := New(richtext.Debug(buf), goPath, ruleSet, "", Verbose, MustPanic, Install, RecurseTopLevel)
		ctx.Get(".", "gh/u2/p2", false, false)
	})

	expected := stringSet{
		"./pkg/gh/u1/p1.a":         empty{},
		"./pkg/gh/u1/p1/s1.a":      empty{},
		"./pkg/gh/u1/p1/s2.a":      empty{},
		"./pkg/gh/u2/p1.a":         empty{},
		"./pkg/gh/u2/p2.a":         empty{},
		"./src/gh/u1/p1/gen.go":    empty{},
		"./src/gh/u1/p1/s1/gen.go": empty{},
		"./src/gh/u1/p1/s2/gen.go": empty{},
		"./src/gh/u2/p1/gen.go":    empty{},
		"./src/gh/u2/p2/gen.go":    empty{},
	}
	assert.Equal(t, expected, fileList)
	assert.Equal(t, "gh/u1/p1\ngh/u2/p1\ngh/u2/p2\n", buf.String())
}

func TestMultiDepPartialFail(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s2", "gh/u1/p1", "gh/u1/p1/missing1"))
	repos.AddRepo("gh/u2/p1",
		Pkg("gh/u2/p1", "gh/u1/p1/s1"))
	repos.AddRepo("gh/u2/p2",
		Pkg("gh/u2/p2", "gh/u1/p1/s1", "gh/u1/p1", "gh/u2/p1"))

	buf := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		ctx := New(richtext.Debug(buf), goPath, ruleSet, "", Verbose, MustPanic, Install, RecurseTopLevel)
		ctx.Get(".", "gh/u2/p2", false, false)
	})

	expected := stringSet{
		//"./pkg/gh/u1/p1/s2.a":      empty{},
		"./pkg/gh/u1/p1.a":         empty{},
		"./pkg/gh/u1/p1/s1.a":      empty{},
		"./pkg/gh/u2/p1.a":         empty{},
		"./pkg/gh/u2/p2.a":         empty{},
		"./src/gh/u1/p1/gen.go":    empty{},
		"./src/gh/u1/p1/s1/gen.go": empty{},
		"./src/gh/u1/p1/s2/gen.go": empty{},
		"./src/gh/u2/p1/gen.go":    empty{},
		"./src/gh/u2/p2/gen.go":    empty{},
	}
	assert.Equal(t, expected, fileList)
	assert.Equal(t, "[WARN]gh/u1/p1/... [Failed: .../s2][]\ngh/u2/p1\ngh/u2/p2\n", buf.String())
}

func TestDepNoRootUpgradeRecurseTopLevel(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u1/p2/s1"))
	repos.AddRepo("gh/u1/p2",
		Pkg("gh/u1/p2/s1"),
		Pkg("gh/u1/p2/s2", "gh/u1/p3/s1"))
	repos.AddRepo("gh/u1/p3",
		Pkg("gh/u1/p3/s1"))

	//Test Get
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		{
			ctx := New(richtext.Debug(buf1), goPath, ruleSet, "", Verbose, RecurseTopLevel)
			ctx.Get(".", "gh/u1/p1", false, false)
		}
		{
			ctx := New(richtext.Debug(buf2), goPath, ruleSet, "", Verbose, Update, MustPanic, Install, RecurseTopLevel)
			ctx.Get(".", "gh/u1/p1", false, false)
		}
	})

	expected := stringSet{
		"./pkg/gh/u1/p1.a":         empty{},
		"./pkg/gh/u1/p2/s1.a":      empty{},
		"./pkg/gh/u1/p2/s2.a":      empty{},
		"./pkg/gh/u1/p3/s1.a":      empty{},
		"./src/gh/u1/p1/gen.go":    empty{},
		"./src/gh/u1/p2/s1/gen.go": empty{},
		"./src/gh/u1/p2/s2/gen.go": empty{},
		"./src/gh/u1/p3/s1/gen.go": empty{},
	}
	assert.Equal(t, expected, fileList)
	assert.Equal(t, `gh/u1/p3
gh/u1/p2
gh/u1/p1
`, buf1.String())
	assert.Equal(t, `gh/u1/p3
gh/u1/p2
gh/u1/p1
`, buf2.String())

}

func TestDepNoRootUpgradeNoRecurseTopLevel(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u1/p2/s1"))
	repos.AddRepo("gh/u1/p2",
		Pkg("gh/u1/p2/s1"),
		Pkg("gh/u1/p2/s2", "gh/u1/p3/s1"))
	repos.AddRepo("gh/u1/p3",
		Pkg("gh/u1/p3/s1"))

	//Test Get
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		{
			ctx := New(richtext.Debug(buf1), goPath, ruleSet, "", Verbose, DeepScan)
			ctx.Get(".", "gh/u1/p1", false, false)
		}
		{
			ctx := New(richtext.Debug(buf2), goPath, ruleSet, "", Verbose, Update, MustPanic, Install, DeepScan)
			ctx.Get(".", "gh/u1/p1", false, false)
		}
	})

	expected := stringSet{
		"./pkg/gh/u1/p1.a":         empty{},
		"./pkg/gh/u1/p2/s1.a":      empty{},
		"./src/gh/u1/p1/gen.go":    empty{},
		"./src/gh/u1/p2/s1/gen.go": empty{},
		"./src/gh/u1/p2/s2/gen.go": empty{},
	}
	assert.Equal(t, expected, fileList)
	assert.Equal(t, `gh/u1/p2/s1
gh/u1/p1
`, buf1.String())
	assert.Equal(t, `gh/u1/p2/s1
gh/u1/p1
`, buf2.String())

}

func TestNoRecurseNoSubPkgs(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u1/p2"))
	repos.AddRepo("gh/u1/p2",
		Pkg("gh/u1/p2"),
		Pkg("gh/u1/p2/sub", "gh/u1/p3"))
	repos.AddRepo("gh/u1/p3",
		Pkg("gh/u1/p3"))

	//Test Get
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		{
			ctx := New(richtext.Debug(buf1), goPath, ruleSet, "", Verbose, DeepScan)
			ctx.Get(".", "gh/u1/p1", false, false)
		}
		{
			ctx := New(richtext.Debug(buf2), goPath, ruleSet, "", Verbose, Update, MustPanic, Install, DeepScan)
			ctx.Get(".", "gh/u1/p1", false, false)
		}
	})

	expected := stringSet{
		"./pkg/gh/u1/p1.a":          empty{},
		"./pkg/gh/u1/p2.a":          empty{},
		"./src/gh/u1/p1/gen.go":     empty{},
		"./src/gh/u1/p2/gen.go":     empty{},
		"./src/gh/u1/p2/sub/gen.go": empty{},
	}
	assert.Equal(t, expected, fileList)
	assert.Equal(t, `gh/u1/p2
gh/u1/p1
`, buf1.String())
	assert.Equal(t, `gh/u1/p2
gh/u1/p1
`, buf2.String())
}

func TestNoRecurseMultiPass(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u1/p2/s1", "gh/u1/p2/s2"))
	repos.AddRepo("gh/u1/p2",
		Pkg("gh/u1/p2/s1", "gh/u1/p3"),
		Pkg("gh/u1/p2/s2", "gh/u1/p4"))
	repos.AddRepo("gh/u1/p3",
		Pkg("gh/u1/p3"))
	repos.AddRepo("gh/u1/p4",
		Pkg("gh/u1/p4"))

	//Test Get
	buf1 := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		{
			ctx := New(richtext.Debug(buf1), goPath, ruleSet, "", Verbose, Install, MustPanic, DeepScan)
			ctx.Get(".", "gh/u1/p1", false, false)
		}
	})

	expected := stringSet{
		"./pkg/gh/u1/p1.a":         empty{},
		"./pkg/gh/u1/p2/s1.a":      empty{},
		"./pkg/gh/u1/p2/s2.a":      empty{},
		"./pkg/gh/u1/p3.a":         empty{},
		"./pkg/gh/u1/p4.a":         empty{},
		"./src/gh/u1/p1/gen.go":    empty{},
		"./src/gh/u1/p2/s1/gen.go": empty{},
		"./src/gh/u1/p2/s2/gen.go": empty{},
		"./src/gh/u1/p3/gen.go":    empty{},
		"./src/gh/u1/p4/gen.go":    empty{},
	}
	assert.Equal(t, expected, fileList)
	assert.Equal(t, `gh/u1/p3
gh/u1/p2/s1
gh/u1/p4
gh/u1/p2/s2
gh/u1/p1
`, buf1.String())
}

func TestHooks(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	pkg := Pkg("hookrepo")
	repo := repos.AddRepo("hookrepo", pkg)

	repo.hookBeforeUpdate = `#!/usr/bin/env sh
if [ -e ".hook1" ]; then
	touch .hookfail
fi
touch .hook1
`
	repo.hookBeforeInstall = `#!/usr/bin/env sh
if [ -e ".hook2a" ]; then
	touch .hook2b
fi
touch .hook2a`

	repo.hookAfterInstall = `#!/usr/bin/env sh
if [ -e ".hook3a" ]; then
	touch .hook3b
fi
touch .hook3a`

	repo.gitIgnore = `
.hook*
`

	//Install
	buf1 := &bytes.Buffer{}
	fileList1 := repos.Test(func(goPath []string, ruleSet RuleSet) {
		ctx := New(richtext.Debug(buf1), goPath, ruleSet, "", Verbose, MustPanic, Install, ApplyHooks, DeepScan)
		ctx.Get(".", "hookrepo", false, false)
	})
	assert.Equal(t, "hookrepo\n", buf1.String())

	expected1 := stringSet{
		"./pkg/hookrepo.a":                     empty{},
		"./src/hookrepo/.gitignore":            empty{},
		"./src/hookrepo/.hook2a":               empty{},
		"./src/hookrepo/.hook3a":               empty{},
		"./src/hookrepo/get-before-update.sh":  empty{},
		"./src/hookrepo/get-before-install.sh": empty{},
		"./src/hookrepo/get-after-install.sh":  empty{},
		"./src/hookrepo/gen.go":                empty{},
	}
	assert.Equal(t, expected1, fileList1)

	//Install AND Upgrade
	buf2 := &bytes.Buffer{}
	buf3 := &bytes.Buffer{}
	fileList2 := repos.Test(func(goPath []string, ruleSet RuleSet) {
		{
			ctx := New(richtext.Debug(buf2), goPath, ruleSet, "", Verbose, MustPanic, Install, ApplyHooks, DeepScan)
			ctx.Get(".", "hookrepo", false, false)
		}

		{
			ctx := New(richtext.Debug(buf3), goPath, ruleSet, "", Verbose, MustPanic, Install, ApplyHooks, Update, DeepScan)
			ctx.Get(".", "hookrepo", false, false)
		}
	})
	assert.Equal(t, "hookrepo\n", buf2.String())
	assert.Equal(t, "hookrepo\n", buf3.String())

	expected2 := stringSet{
		"./pkg/hookrepo.a":                     empty{},
		"./src/hookrepo/.gitignore":            empty{},
		"./src/hookrepo/.hook1":                empty{}, //MISSING
		"./src/hookrepo/.hook2a":               empty{},
		"./src/hookrepo/.hook2b":               empty{}, //MISSING
		"./src/hookrepo/.hook3a":               empty{},
		"./src/hookrepo/.hook3b":               empty{}, //MISSING
		"./src/hookrepo/get-before-update.sh":  empty{},
		"./src/hookrepo/get-before-install.sh": empty{},
		"./src/hookrepo/get-after-install.sh":  empty{},
		"./src/hookrepo/gen.go":                empty{},
	}
	assert.Equal(t, expected2, fileList2)

}

func TestAccessing2SeparatePkgsFromSameRepoWithDistinctTransDeps(t *testing.T) {
	format := richtext.Test(t)

	repos := NewRepos(format)

	repos.AddRepo("gh/u1/p1",
		Pkg("gh/u1/p1", "gh/u2/p1"),
		Pkg("gh/u1/p1/s1"),
		Pkg("gh/u1/p1/s1/a1", "gh/u2/p2"))
	repos.AddRepo("gh/u2/p1",
		Pkg("gh/u2/p1"))
	repos.AddRepo("gh/u2/p2",
		Pkg("gh/u2/p2"))
	repos.AddRepo("gh/u3/p1",
		Pkg("gh/u3/p1", "gh/u1/p1/s1", "gh/u1/p1/s1/a1"))

	buf := &bytes.Buffer{}
	fileList := repos.Test(func(goPath []string, ruleSet RuleSet) {
		ctx := New(richtext.Debug(buf), goPath, ruleSet, "", Verbose, MustPanic, Install, RecurseTopLevel)
		ctx.Get(".", "gh/u3/p1", false, false)
	})

	expected := stringSet{
		"./pkg/gh/u1/p1.a":      empty{},
		"./pkg/gh/u1/p1/s1.a":      empty{},
		"./pkg/gh/u1/p1/s1/a1.a":      empty{},
		"./pkg/gh/u2/p1.a":         empty{},
		"./pkg/gh/u2/p2.a":         empty{},
		"./pkg/gh/u3/p1.a":         empty{},
		"./src/gh/u1/p1/gen.go": empty{},
		"./src/gh/u1/p1/s1/gen.go": empty{},
		"./src/gh/u1/p1/s1/a1/gen.go": empty{},
		"./src/gh/u2/p1/gen.go":    empty{},
		"./src/gh/u2/p2/gen.go":    empty{},
		"./src/gh/u3/p1/gen.go":    empty{},
	}
	assert.Equal(t, expected, fileList)
	assertOutputEquivTo(t,
		buf.String(),
		stringSet{
			"gh/u2/p1": empty{},
			"gh/u2/p2": empty{},
			"gh/u1/p1": empty{},
			"gh/u3/p1": empty{},
		},
	)
}

func assertOutputEquivTo(t *testing.T, output string, expectedLines stringSet) {
	actualLines := stringSet{}
	for _, line := range strings.Split(output, "\n") {
		line := strings.TrimSpace(line)
		if len(line) > 0 {
			actualLines[line] = empty{}
		}
	}
	assert.Equal(t, expectedLines, actualLines)
}


//Poor test coverage:
// non-git repos
// tags
