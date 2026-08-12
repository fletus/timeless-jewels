// Package embedded ships the game tables as compressed assets and hands them to the data
// package at init.
//
// It exists so that github.com/Vilsol/timeless-jewels/data — and therefore the calculator —
// can be linked WITHOUT the ~1.8 MB of assets. Consumers that want the shipped tables import
// this package for its side effect and everything behaves as it did before the split:
//
//	import _ "github.com/Vilsol/timeless-jewels/data/embedded"
//
// Consumers that hold their own copies of the tables call data.Initialize directly and never
// import this package, so the linker drops the assets entirely.
//
// ⚠️ THIS PACKAGE PANICS ON A BAD ASSET, DELIBERATELY. The assets are compiled in, so a
// failure here is a build-time mistake — a truncated file or a schema change — not a runtime
// condition a caller could handle. data.Initialize returns an error instead, because there the
// tables come from the caller.
package embedded

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"

	"github.com/Vilsol/timeless-jewels/data"

	_ "embed"
)

//go:embed alternate_passive_additions.json.gz
var alternatePassiveAdditionsGz []byte

//go:embed alternate_passive_skills.json.gz
var alternatePassiveSkillsGz []byte

//go:embed alternate_tree_versions.json.gz
var alternateTreeVersionsGz []byte

//go:embed passive_skills.json.gz
var passiveSkillsGz []byte

//go:embed stats.json.gz
var statsGz []byte

//go:embed SkillTree.json.gz
var skillTreeGz []byte

//go:embed stat_descriptions.json.gz
var statTranslationsGz []byte

//go:embed passive_skill_stat_descriptions.json.gz
var passiveSkillStatTranslationsGz []byte

//go:embed passive_skill_aura_stat_descriptions.json.gz
var passiveSkillAuraStatTranslationsGz []byte

//go:embed possible_stats.json.gz
var possibleStatsGz []byte

func init() {
	if err := data.Initialize(
		unzipJSONTo[[]*data.PassiveSkill](passiveSkillsGz),
		unzipJSONTo[[]*data.AlternatePassiveSkill](alternatePassiveSkillsGz),
		unzipJSONTo[[]*data.AlternatePassiveAddition](alternatePassiveAdditionsGz),
		unzipJSONTo[[]*data.AlternateTreeVersion](alternateTreeVersionsGz),
	); err != nil {
		panic(err)
	}

	data.SetStats(unzipJSONTo[[]*data.Stat](statsGz))

	// SkillTreeJSON is re-marshalled from the parsed tree rather than served as the raw
	// bytes, exactly as it was before the split — the frontend consumes the normalised shape.
	data.SkillTreeData = unzipJSONTo[data.SkillTree](skillTreeGz)

	skillTreeJSON, err := json.Marshal(data.SkillTreeData)
	if err != nil {
		panic(err)
	}
	data.SkillTreeJSON = skillTreeJSON

	data.StatTranslationsJSON = unzipTo(statTranslationsGz)
	data.PassiveSkillStatTranslationsJSON = unzipTo(passiveSkillStatTranslationsGz)
	data.PassiveSkillAuraStatTranslationsJSON = unzipTo(passiveSkillAuraStatTranslationsGz)
	data.PossibleStatsJSON = unzipTo(possibleStatsGz)
}

func unzipJSONTo[T any](in []byte) T {
	out := new(T)
	if err := json.Unmarshal(unzipTo(in), &out); err != nil {
		panic(err)
	}

	return *out
}

func unzipTo(in []byte) []byte {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		panic(err)
	}

	all, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	return all
}
