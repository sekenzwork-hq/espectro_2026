package unit

import (
	"espectro/pkg"
	"testing"
)

func TestParseLimitAndPage(t *testing.T) {

	limitStr := "10"
	pageStr := "2"

	limit, page, err := pkg.ParsePageAndLimit(limitStr, pageStr)

	if err != nil {
		t.Error("Limit parsing failed (correct)", err)
		return
	} else if limit != 10 {
		t.Error("Limit parsing failed (correct)")
		return
	} else if page != 2 {
		t.Error("Limit parsing failed (correct)")
		return
	}

	limitStr = ""
	pageStr = ""

	limit, page, err = pkg.ParsePageAndLimit(limitStr, pageStr)

	if limit != 50 {
		t.Error("Limit parsing failed (incorrect)")
		return
	} else if page != 1 {
		t.Error("Limit parsing failed (incorrect)")
		return
	}

	limitStr = "1000"
	pageStr = "5"

	limit, page, err = pkg.ParsePageAndLimit(limitStr, pageStr)

	if err == nil {
		t.Error("Limit restriction failed", err)
		return
	}
}
