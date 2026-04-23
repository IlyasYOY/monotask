package version_test

import (
	"runtime/debug"
	"testing"

	"github.com/IlyasYOY/monotask/internal/pkg/version"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		info *debug.BuildInfo
		want string
	}{
		{
			name: "exact tag",
			info: buildInfo("v0.2.0",
				debug.BuildSetting{Key: "vcs.revision", Value: "1f4b76a64268bfcbbcadd37ca51224ed3326c0f6"},
			),
			want: "v0.2.0",
		},
		{
			name: "pseudo version from commit after stable tag",
			info: buildInfo("v0.2.1-0.20260423110655-1f4b76a64268",
				debug.BuildSetting{Key: "vcs.revision", Value: "1f4b76a64268bfcbbcadd37ca51224ed3326c0f6"},
			),
			want: "v0.2.0+1f4b76a",
		},
		{
			name: "pseudo version uses suffix commit without vcs metadata",
			info: buildInfo("v0.2.1-0.20260423110655-1f4b76a64268"),
			want: "v0.2.0+1f4b76a",
		},
		{
			name: "dirty build",
			info: buildInfo("v0.2.1-0.20260423110655-1f4b76a64268",
				debug.BuildSetting{Key: "vcs.revision", Value: "1f4b76a64268bfcbbcadd37ca51224ed3326c0f6"},
				debug.BuildSetting{Key: "vcs.modified", Value: "true"},
			),
			want: "v0.2.0+1f4b76a.dirty",
		},
		{
			name: "dirty module version suffix",
			info: buildInfo("v0.2.1-0.20260423110655-1f4b76a64268+dirty",
				debug.BuildSetting{Key: "vcs.revision", Value: "1f4b76a64268bfcbbcadd37ca51224ed3326c0f6"},
				debug.BuildSetting{Key: "vcs.modified", Value: "true"},
			),
			want: "v0.2.0+1f4b76a.dirty",
		},
		{
			name: "local devel build with revision",
			info: buildInfo("(devel)",
				debug.BuildSetting{Key: "vcs.revision", Value: "1f4b76a64268bfcbbcadd37ca51224ed3326c0f6"},
			),
			want: "devel+1f4b76a",
		},
		{
			name: "missing metadata",
			info: buildInfo("(devel)"),
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := version.Format(tt.info)
			if got != tt.want {
				t.Fatalf("Format() = %q, want %q", got, tt.want)
			}
		})
	}
}

func buildInfo(mainVersion string, settings ...debug.BuildSetting) *debug.BuildInfo {
	return &debug.BuildInfo{
		Main: debug.Module{
			Path:    "github.com/IlyasYOY/monotask",
			Version: mainVersion,
		},
		Settings: settings,
	}
}
