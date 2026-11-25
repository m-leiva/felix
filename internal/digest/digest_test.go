package digest

import (
	"strings"
	"testing"
	"time"
)

func TestGetSummaryFormattedEmptyAndPopulated(t *testing.T) {
	w := &Watcher{}
	if got := w.GetSummaryFormatted(); !strings.Contains(got, "no disponible") {
		t.Errorf("mensaje inesperado para snapshot vacio: %q", got)
	}

	w.current = SystemSnapshot{
		CPUTemp:    45.2,
		RAMUsedMB:  4096,
		RAMTotalMB: 8192,
		DiskRoot:   "50G disp (uso 40%)",
		PM2Online:  10,
		PM2Total:   10,
		UpdatedAt:  time.Now(),
	}

	got := w.GetSummaryFormatted()
	for _, want := range []string{"Sistema [", "45.2C", "4096/8192 MB", "50G disp (uso 40%)", "10/10 online (todos OK)"} {
		if !strings.Contains(got, want) {
			t.Errorf("resumen no contiene %q: %q", want, got)
		}
	}
}
