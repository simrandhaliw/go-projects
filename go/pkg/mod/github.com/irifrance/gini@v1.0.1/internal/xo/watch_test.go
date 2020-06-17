// Copyright 2016 The Gini Authors. All rights reserved.  Use of this source
// code is governed by a license that can be found in the License file.

package xo

import (
	"fmt"
	"testing"

	"github.com/irifrance/gini/z"
)

func TestLocOverflow(t *testing.T) {
	loc := z.C(3)
	w := MakeWatch(loc, 7, true)
	loc2 := w.C()
	if w.C() != loc {
		t.Errorf("error isbin overflow?: %s != %s", loc, loc2)
	}
}

func TestWatch(t *testing.T) {
	loc := z.C(77)
	m := z.Lit(1024)
	isBin := true
	w := MakeWatch(loc, m, isBin)
	fmt.Printf("%s\n", w)
	if w.Other() != m {
		t.Errorf("other decode: %s != %s", w.Other(), m)
	}
	if w.IsBinary() != isBin {
		t.Errorf("isBind decode: %t != %t", w.IsBinary(), isBin)
	}
	if w.C() != loc {
		t.Errorf("loc en/decode: %s != %s", w.C(), loc)
	}

	newLoc := z.C(22)
	w0 := w.Relocate(newLoc)
	if w0.Other() != m {
		t.Errorf("relocate other: %s != %s", w0.Other(), m)
	}
	if w0.IsBinary() != isBin {
		t.Errorf("isBin decode %t != %t", w0.IsBinary(), isBin)
	}
	if w0.C() != newLoc {
		t.Errorf("relocate  newloc %s != %s", w0.C(), newLoc)
	}
}
