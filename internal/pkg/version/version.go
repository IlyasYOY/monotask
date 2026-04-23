package version

import (
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

const unknown = "unknown"

var commitPseudoVersionPattern = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)-0\.[0-9]{14}-([0-9A-Za-z]+)$`)

func Current() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return unknown
	}

	return Format(info)
}

func Format(info *debug.BuildInfo) string {
	if info == nil {
		return unknown
	}

	revision, dirty := buildVCSInfo(info.Settings)
	mainVersion, moduleVersionDirty := normalizeModuleVersion(info.Main.Version)
	dirty = dirty || moduleVersionDirty
	if mainVersion == "" || mainVersion == "(devel)" {
		if revision == "" {
			return unknown
		}
		return addDirtySuffix("devel+"+shortRevision(revision), dirty)
	}

	if nearestTag, pseudoRevision, ok := nearestStableTag(mainVersion); ok {
		if revision == "" {
			revision = pseudoRevision
		}
		return addDirtySuffix(nearestTag+"+"+shortRevision(revision), dirty)
	}

	return addDirtySuffix(mainVersion, dirty)
}

func normalizeModuleVersion(moduleVersion string) (string, bool) {
	if strings.HasSuffix(moduleVersion, "+dirty") {
		return strings.TrimSuffix(moduleVersion, "+dirty"), true
	}
	return moduleVersion, false
}

func buildVCSInfo(settings []debug.BuildSetting) (string, bool) {
	var revision string
	var dirty bool
	for _, setting := range settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			dirty = setting.Value == "true"
		}
	}
	return revision, dirty
}

func nearestStableTag(moduleVersion string) (string, string, bool) {
	matches := commitPseudoVersionPattern.FindStringSubmatch(moduleVersion)
	if matches == nil {
		return "", "", false
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil || patch == 0 {
		return "", "", false
	}

	return "v" + matches[1] + "." + matches[2] + "." + strconv.Itoa(patch-1), matches[4], true
}

func shortRevision(revision string) string {
	if len(revision) <= 7 {
		return revision
	}
	return revision[:7]
}

func addDirtySuffix(version string, dirty bool) string {
	if dirty && !strings.HasSuffix(version, ".dirty") {
		return version + ".dirty"
	}
	return version
}
