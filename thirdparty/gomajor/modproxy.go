package gomajor

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
)

// MaxVersion returns the latest version.
// If there are no versions, the empty string is returned.
// Prefix can be used to filter the versions based on a prefix.
// If pre is false, pre-release versions will are excluded.
func (m *Module) MaxVersion(prefix string, pre bool) string {
	var max string
	for _, v := range m.Versions {
		if !semver.IsValid(v) || !strings.HasPrefix(v, prefix) {
			continue
		}
		if !pre && semver.Prerelease(v) != "" {
			continue
		}
		max = MaxVersion(v, max)
	}
	return max
}

// IsNewerVersion returns true if newversion is greater than oldversion in terms of semver.
// If major is true, then newversion must be a major version ahead of oldversion to be considered newer.
func IsNewerVersion(oldversion, newversion string, major bool) bool {
	if major {
		return semver.Compare(semver.Major(oldversion), semver.Major(newversion)) < 0
	}
	return semver.Compare(oldversion, newversion) < 0
}

// MaxVersion returns the larger of two versions according to semantic version precedence.
// Incompatible versions are considered lower than non-incompatible ones.
// Invalid versions are considered lower than valid ones.
// If both versions are invalid, the empty string is returned.
func MaxVersion(v, w string) string {
	// sort by validity
	vValid := semver.IsValid(v)
	wValid := semver.IsValid(w)
	if !vValid && !wValid {
		return ""
	}
	if vValid != wValid {
		if vValid {
			return v
		}
		return w
	}
	// sort by compatibility
	vIncompatible := strings.HasSuffix(semver.Build(v), "+incompatible")
	wIncompatible := strings.HasSuffix(semver.Build(w), "+incompatible")
	if vIncompatible != wIncompatible {
		if wIncompatible {
			return v
		}
		return w
	}
	// sort by semver
	if semver.Compare(v, w) == 1 {
		return v
	}
	return w
}

// NextMajor returns the next major version after the provided version
func NextMajor(version string) (string, error) {
	major, err := strconv.Atoi(strings.TrimPrefix(semver.Major(version), "v"))
	if err != nil {
		return "", err
	}
	major++
	return fmt.Sprintf("v%d", major), nil
}

// WithMajorPath returns the module path for the provided version
func (m *Module) WithMajorPath(version string) string {
	prefix := ModPrefix(m.Path)
	return JoinPath(prefix, version, "")
}

// NextMajorPath returns the module path of the next major version
func (m *Module) NextMajorPath() (string, bool) {
	latest := m.MaxVersion("", true)
	if latest == "" {
		return "", false
	}
	if semver.Major(latest) == "v0" {
		return "", false
	}
	next, err := NextMajor(latest)
	if err != nil {
		return "", false
	}
	return m.WithMajorPath(next), true
}

// QueryCurrent finds the current major version of a module via go proxy.
// If the module does not exist, the second return parameter will be false
// cached sets the Disable-Module-Fetch: true header
func QueryCurrent(modpath string, cached bool) (*Module, bool, error) {
	escaped, err := module.EscapePath(modpath)
	if err != nil {
		return nil, false, err
	}
	url := fmt.Sprintf("https://proxy.golang.org/%s/@v/list", escaped)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", "snuuze/0.0")
	if cached {
		req.Header.Set("Disable-Module-Fetch", "true")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var body []byte
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, false, err
		}
		if res.StatusCode == http.StatusNotFound && bytes.HasPrefix(body, []byte("not found:")) {
			return nil, false, nil
		}
		msg := string(body)
		if msg == "" {
			msg = res.Status
		}
		return nil, false, fmt.Errorf("proxy: %s", msg)
	}
	var mod Module
	mod.Path = modpath
	sc := bufio.NewScanner(res.Body)
	for sc.Scan() {
		mod.Versions = append(mod.Versions, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, false, err
	}
	return &mod, true, nil
}

// Query finds the all versions of a module with major versions greater than or equal to current one.
// cached sets the Disable-Module-Fetch: true header
func Query(modpath string, cached bool) (*MultiModule, error) {
	multiModule := &MultiModule{}

	latest, ok, err := QueryCurrent(modpath, cached)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("module not found: %s", modpath)
	}
	multiModule.Modules = append(multiModule.Modules, latest)

	for i := 0; i < 100; i++ {
		nextpath, ok := latest.NextMajorPath()
		if !ok {
			return multiModule, nil
		}
		next, ok, err := QueryCurrent(nextpath, cached)
		if err != nil {
			return nil, err
		}
		if !ok {
			// handle the case where a project switched to modules
			// without incrementing the major version
			version := latest.MaxVersion("", true)
			if semver.Build(version) == "+incompatible" {
				nextpath = latest.WithMajorPath(semver.Major(version))
				if nextpath != latest.Path {
					next, ok, err = QueryCurrent(nextpath, cached)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if !ok {
			return multiModule, nil
		}
		multiModule.Modules = append(multiModule.Modules, next)
		latest = next
	}
	return nil, fmt.Errorf("request limit exceeded")
}

// Update reports a newer version of a module.
// The Err field will be set if an error occurred.
type Update struct {
	Module  module.Version
	Version string
	Err     error
}
