package dxtypes

import (
	"sort"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func createPackage(name, version, file, analyzer string) Package {
	return Package{
		Name:    name,
		Version: version,
		FromFile: []string{
			file,
		},
		FromAnalyzer: []string{
			analyzer,
		},
	}
}

type pkgtestdata struct {
	FromFile           []string
	FromAnalyzer       []string
	UpStreamPackages   []string
	DownStreamPackages []string
}

func check(t *testing.T, p *Package, want *pkgtestdata) {
	if len(p.FromAnalyzer) != len(want.FromAnalyzer) {
		t.Fatalf("fromAnalyzer len error: %v", p.FromAnalyzer)
	}
	if slices.CompareFunc(p.FromAnalyzer, want.FromAnalyzer, strings.Compare) != 0 {
		t.Fatalf("fromAnalyzer error: %v", p.FromAnalyzer)
	}
	if len(p.FromFile) != len(want.FromFile) {
		t.Fatal("fromFile len error")
	}
	if slices.CompareFunc(p.FromFile, want.FromFile, strings.Compare) != 0 {
		t.Fatalf("fromFile error: %v", p.FromFile)
	}

	downPkgNames := []string{}
	for _, pkg := range p.DownStreamPackages {
		if _, ok := pkg.UpStreamPackages[p.Identifier()]; !ok {
			t.Fatalf("downSteram not link to pa %s", pkg)
		}
		downPkgNames = append(downPkgNames, pkg.Name)
	}

	sort.Strings(downPkgNames)
	if len(downPkgNames) != len(want.DownStreamPackages) {
		t.Fatal("down package len error")
	}
	if slices.CompareFunc(downPkgNames, want.DownStreamPackages, strings.Compare) != 0 {
		t.Fatalf("down pkgname error: %v", downPkgNames)
	}
	upPkgNames := []string{}
	for _, pkg := range p.UpStreamPackages {
		if _, ok := pkg.DownStreamPackages[p.Identifier()]; !ok {
			t.Fatalf("UpStream not link to pa %s", pkg)
		}
		upPkgNames = append(upPkgNames, pkg.Name)
	}

	sort.Strings(upPkgNames)
	if len(upPkgNames) != len(want.UpStreamPackages) {
		t.Fatal("up package len error")
	}
	if slices.CompareFunc(upPkgNames, want.UpStreamPackages, strings.Compare) != 0 {
		t.Fatalf("up pkgname error: %v", upPkgNames)
	}
}

func TestPackageMergeNormal(t *testing.T) {
	pa := createPackage("pa", "0.0.1", "/path/pa", "pa-analyzer")
	pa_down := createPackage("pa-down", "0.0.2", "/path/pa", "pa-analyzer")
	pa_down.LinkDepend(&pa)
	pb := createPackage("pb", "0.0.1", "/path/pb", "pb-analyzer")
	pb_down := createPackage("pb-down", "0.0.3", "/path/pd", "pb-analyzer")
	// pb.LinkDepend(&pb_down)
	pb_down.LinkDepend(&pb)

	// Merge(&pa, &pb)
	pa.Merge(&pb)
	// fmt.Printf("%s", pa)
	want := &pkgtestdata{
		FromFile:         []string{"/path/pa", "/path/pb"},
		FromAnalyzer:     []string{"pa-analyzer", "pb-analyzer"},
		UpStreamPackages: []string{},
		DownStreamPackages: []string{
			"pa-down",
			"pb-down",
		},
	}
	check(t, &pa, want)
}

func TestPackageMergeRepeat(t *testing.T) {
	// same file and analyzer
	pa := createPackage("pa", "0.0.1", "/path/pa", "pa-analyzer")
	pb := createPackage("pb", "0.0.1", "/path/pa", "pa-analyzer")

	pa_down := createPackage("pa-down", "0.0.2", "/path/pa", "pa-analyzer")
	pb_down := createPackage("pb-down", "0.0.3", "/path/pd", "pb-analyzer")

	pa_down.LinkDepend(&pa)
	pb_down.LinkDepend(&pb)

	// Merge(&pa, &pb)
	pa.Merge(&pb)
	// fmt.Printf("%s", pa)
	want := &pkgtestdata{
		FromFile:         []string{"/path/pa"},
		FromAnalyzer:     []string{"pa-analyzer"},
		UpStreamPackages: []string{},
		DownStreamPackages: []string{
			"pa-down",
			"pb-down",
		},
	}
	check(t, &pa, want)
}

func TestPackageCanMerge(t *testing.T) {
	// same
	// merge (pa to pb) or (pb to pa)
	t.Run("same-name-and-version", func(t *testing.T) {
		pa := createPackage("pa", "0.0.1", "", "")
		pb := createPackage("pa", "0.0.1", "", "")
		if CanMerge(&pa, &pb) != 1 {
			t.Fatal("same name and version shoud merge pa(pb)")
		}
		if CanMerge(&pb, &pa) != 1 {
			t.Fatal("same name and version shoud merge pb(pa)")
		}
	})
	// not same
	// not merge
	t.Run("same-version-wrong-name", func(t *testing.T) {
		pa := createPackage("pa", "0.0.1", "", "")
		pb := createPackage("pb", "0.0.1", "", "")
		if CanMerge(&pa, &pb) != 0 {
			t.Fatal("same version wrong name shoud not merge pa(pb)")
		}
		if CanMerge(&pb, &pa) != 0 {
			t.Fatal("same version wrong name shoud not merge pb(pa)")
		}
	})

	// pa is range : "*"
	// merge pa to pb
	t.Run("merge-same-version-range-*", func(t *testing.T) {
		pa := createPackage("pa", "*", "", "")
		pa.IsVersionRange = true
		pb := createPackage("pa", "0.0.1", "", "")
		if CanMerge(&pa, &pb) != -1 {
			t.Fatal("same name and version is * shoud merge pa(pb)")
		}
		if CanMerge(&pb, &pa) != 1 {
			t.Fatal("same name and version is * shoud merge pb(pa)")
		}
	})
	// pa range: "<0.0.3"
	// merge pa to pb
	t.Run("merge-same-version-range-match", func(t *testing.T) {
		pa := createPackage("pa", "<0.0.3", "", "")
		pa.IsVersionRange = true
		pb := createPackage("pa", "0.0.1", "", "")
		if CanMerge(&pa, &pb) != -1 {
			t.Fatal("same name and version match range shoud merge pa(pb)")
		}
		if CanMerge(&pb, &pa) != 1 {
			t.Fatal("same name and version match range shoud merge pb(pa)")
		}
	})

	// pa range: "=0.0.3"
	// merge pa to pb
	t.Run("merge-same-version-range-exact", func(t *testing.T) {
		pa := createPackage("pa", "=0.0.3", "", "")
		pa.IsVersionRange = true
		pb := createPackage("pa", "0.0.3", "", "")
		if CanMerge(&pa, &pb) != -1 {
			t.Fatal("same name and version match range shoud merge pa(pb)")
		}
		if CanMerge(&pb, &pa) != 1 {
			t.Fatal("same name and version match range shoud merge pb(pa)")
		}
	})

	// pa range ">0.0.3"
	// not merge
	t.Run("not-merge-version-range-not-match", func(t *testing.T) {
		pa := createPackage("pa", ">0.0.3", "", "")
		pa.IsVersionRange = true
		pb := createPackage("pa", "0.0.1", "", "")
		if CanMerge(&pa, &pb) != 0 {
			t.Fatal("same name and version not match range shoud not merge pa(pb)")
		}
		if CanMerge(&pb, &pa) != 0 {
			t.Fatal("same name and version not match range shoud not merge pb(pa)")
		}
	})

	// ignore pa.IsVersionRange
	// pa range "<0.0.3"
	// merge pa to pb
	t.Run("merge-version-range-match-ignore-IsVersionRange", func(t *testing.T) {
		pa := createPackage("pa", "<0.0.3", "", "")
		pb := createPackage("pa", "0.0.1", "", "")
		if CanMerge(&pa, &pb) != -1 {
			t.Fatal("same name and version match range shoud merge even not set IsVersionRange pa(pb)")
		}
		if CanMerge(&pb, &pa) != 1 {
			t.Fatal("same name and version match range shoud merge even not set IsVersionRange pb(pa)")
		}
	})

	// pa version no number
	// not merge
	t.Run("not-merge-same-name-no-version", func(t *testing.T) {
		pa := createPackage("pa", "a", "", "")
		pb := createPackage("pa", "0.0.1", "", "")
		if CanMerge(&pa, &pb) != 0 {
			t.Fatal("same name but no version number, don't merge: pa-pb")
		}
		if CanMerge(&pb, &pa) != 0 {
			t.Fatal("same name but no version number, don't merge: pb-pa")
		}
	})

	// pa version no number with range
	// not merge
	t.Run("not-merge-same-name-no-version-with-range", func(t *testing.T) {
		pa := createPackage("pa", ">a", "", "")
		pb := createPackage("pa", "0.0.1", "", "")
		if CanMerge(&pa, &pb) != 0 {
			t.Fatal("same name but no version number, don't merge: pa-pb")
		}
		if CanMerge(&pb, &pa) != 0 {
			t.Fatal("same name but no version number, don't merge: pb-pa")
		}
	})

	// pa range ">0.0.1 && <0.0.3"
	// merge pa to pb
	t.Run("merge-same-name-version-range-match", func(t *testing.T) {
		pa := createPackage("pa", ">0.0.1 && <0.0.3", "", "")
		pb := createPackage("pa", "0.0.2", "", "")
		if CanMerge(&pa, &pb) != -1 {
			t.Fatal("same name and version range match, merge")
		}
		if CanMerge(&pb, &pa) != 1 {
			t.Fatal("same name and version range match, merge")
		}
	})

	// pa range  and pb range
	// not merge
	t.Run("not-merge-two-packages-with-version-range", func(t *testing.T) {
		pa := createPackage("pa", ">0.0.1", "", "")
		pb := createPackage("pa", ">0.0.2", "", "")
		if CanMerge(&pa, &pb) != 0 {
			t.Fatal("two package with version range, don't merge")
		}
		if CanMerge(&pb, &pa) != 0 {
			t.Fatal("two package with version range, don't merge")
		}
	})
}
