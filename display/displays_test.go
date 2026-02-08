package display

import "testing"

func TestGetShortProcessName(t *testing.T) {
	dtests := []struct {
		full  string
		short string
	}{
		{full: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			short: "Chrome"},
		{full: "/Applications/Visual Studio Code.app/Contents/MacOS/Electron",
			short: "Electron"},
		{full: "/System/Library/CoreServices/WindowServer",
			short: "WindowServer"},
		{full: "/usr/libexec/kafkactl-agent-helper (Renderer)",
			short: "kafkactl-agent"},
		{full: "very-long-process-name-that-needs-truncating",
			short: "very-long-pr..."},
	}

	for _, dt := range dtests {
		got := GetShortProcessName(dt.full)
		if got != dt.short {
			t.Errorf("Укоротили до %s, а должны были %s", got, dt.short)
		}
	}
}
