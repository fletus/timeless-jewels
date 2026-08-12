package data

// The tables the timeless-jewel calculation reads, and every index derived from them.
//
// ⚠️ NOTHING HERE IS POPULATED BY AN init(). The gzipped assets that used to live in this file
// now sit in the data/embedded package, whose init() calls Initialize below. Consumers that
// want the shipped tables import that package for its side effect:
//
//	import _ "github.com/Vilsol/timeless-jewels/data/embedded"
//
// Consumers that already hold these tables — go-pob loads all four from go-pob-data, keyed to
// the game version it targets — call Initialize directly and never link the assets in. That is
// the whole point of the split: the embedded copies are frozen at library-build time, so a
// consumer tracking a different game version would otherwise carry two copies that drift apart.
//
// ⚠️ ASSIGNING THE EXPORTED SLICES IS NOT ENOUGH AND NEVER WAS. Every lookup in manager.go
// except GetAlternatePassiveSkillKeyStone and GetApplicablePassives answers from one of the
// unexported maps below. Overwriting PassiveSkills without rebuilding idToPassiveSkill leaves
// the map pointing at the previous tables — a silent wrong answer rather than an error, which
// is exactly why this needs to be a function and not a documented convention.

import (
	"errors"
	"fmt"
)

var AlternatePassiveAdditions []*AlternatePassiveAddition

var (
	idToAlternatePassiveAddition     = make(map[uint32]*AlternatePassiveAddition)
	reverseAlternatePassiveAdditions = make(map[PassiveSkillType]map[uint32][]*AlternatePassiveAddition)
)

var AlternatePassiveSkills []*AlternatePassiveSkill

var (
	idToAlternatePassiveSkill     = make(map[uint32]*AlternatePassiveSkill)
	reverseAlternatePassiveSkills = make(map[PassiveSkillType]map[uint32][]*AlternatePassiveSkill)
)

var AlternateTreeVersions []*AlternateTreeVersion

var idToAlternateTreeVersion = make(map[uint32]*AlternateTreeVersion)

var PassiveSkills []*PassiveSkill

var idToPassiveSkill = make(map[uint32]*PassiveSkill)

var Stats []*Stat

var idToStat = make(map[uint32]*Stat)

var (
	SkillTreeJSON []byte
	SkillTreeData SkillTree
)

var (
	StatTranslationsJSON                 []byte
	PassiveSkillStatTranslationsJSON     []byte
	PassiveSkillAuraStatTranslationsJSON []byte
	PossibleStatsJSON                    []byte
)

var initialized bool

// ErrNotInitialized is returned by Initialize when a required table is empty.
var ErrNotInitialized = errors.New("timeless-jewels/data: table is empty")

// Initialized reports whether Initialize has completed successfully.
//
// ⚠️ CHECK THIS ONCE AT STARTUP. Every lookup returns a nil/zero value on an uninitialised
// package, and Calculate turns that into an empty AlternatePassiveSkillInformation rather than
// a panic — so a consumer that forgets to initialise gets plausible "this jewel changes
// nothing" results on every node. The check is deliberately NOT inside the lookups: they run
// once per passive per seed and a branch there is measurable.
func Initialized() bool {
	return initialized
}

// Initialize populates the four tables the calculation reads and rebuilds all six indexes
// derived from them. It is safe to call more than once; each call discards the previous
// indexes rather than merging into them.
//
// The tables correspond to these game data files:
//
//	passives      PassiveSkills
//	altSkills     AlternatePassiveSkills
//	altAdditions  AlternatePassiveAdditions
//	treeVersions  AlternateTreeVersions
//
// Stats, the skill tree and the stat translations are NOT required here — nothing the
// calculation does reads them. See SetStats for the one of those that has a derived index.
func Initialize(
	passives []*PassiveSkill,
	altSkills []*AlternatePassiveSkill,
	altAdditions []*AlternatePassiveAddition,
	treeVersions []*AlternateTreeVersion,
) error {
	// ⚠️ REFUSE RATHER THAN INITIALISE HALFWAY. An empty table produces no error at any later
	// point: the lookups just miss and the calculation reports that the jewel replaces nothing.
	for _, check := range []struct {
		name string
		n    int
	}{
		{"passives", len(passives)},
		{"altSkills", len(altSkills)},
		{"altAdditions", len(altAdditions)},
		{"treeVersions", len(treeVersions)},
	} {
		if check.n == 0 {
			return fmt.Errorf("%w: %s", ErrNotInitialized, check.name)
		}
	}

	newIDToAddition := make(map[uint32]*AlternatePassiveAddition, len(altAdditions))
	newReverseAdditions := make(map[PassiveSkillType]map[uint32][]*AlternatePassiveAddition)

	for _, alt := range altAdditions {
		newIDToAddition[alt.Index] = alt

		for _, skillType := range alt.PassiveType {
			if _, ok := newReverseAdditions[skillType]; !ok {
				newReverseAdditions[skillType] = make(map[uint32][]*AlternatePassiveAddition)
			}

			newReverseAdditions[skillType][alt.AlternateTreeVersionsKey] = append(newReverseAdditions[skillType][alt.AlternateTreeVersionsKey], alt)
		}
	}

	newIDToSkill := make(map[uint32]*AlternatePassiveSkill, len(altSkills))
	newReverseSkills := make(map[PassiveSkillType]map[uint32][]*AlternatePassiveSkill)

	for _, alt := range altSkills {
		newIDToSkill[alt.Index] = alt

		for _, skillType := range alt.PassiveType {
			if _, ok := newReverseSkills[skillType]; !ok {
				newReverseSkills[skillType] = make(map[uint32][]*AlternatePassiveSkill)
			}

			newReverseSkills[skillType][alt.AlternateTreeVersionsKey] = append(newReverseSkills[skillType][alt.AlternateTreeVersionsKey], alt)
		}
	}

	newIDToTreeVersion := make(map[uint32]*AlternateTreeVersion, len(treeVersions))
	for _, alt := range treeVersions {
		newIDToTreeVersion[alt.Index] = alt
	}

	newIDToPassive := make(map[uint32]*PassiveSkill, len(passives))
	for _, skill := range passives {
		newIDToPassive[skill.Index] = skill
	}

	AlternatePassiveAdditions = altAdditions
	idToAlternatePassiveAddition = newIDToAddition
	reverseAlternatePassiveAdditions = newReverseAdditions

	AlternatePassiveSkills = altSkills
	idToAlternatePassiveSkill = newIDToSkill
	reverseAlternatePassiveSkills = newReverseSkills

	AlternateTreeVersions = treeVersions
	idToAlternateTreeVersion = newIDToTreeVersion

	PassiveSkills = passives
	idToPassiveSkill = newIDToPassive

	initialized = true

	return nil
}

// SetStats populates the stat table and its index. The calculation never reads it — it is used
// for rendering stat text, which is why it is separate from Initialize.
func SetStats(stats []*Stat) {
	newIDToStat := make(map[uint32]*Stat, len(stats))
	for _, stat := range stats {
		newIDToStat[stat.Index] = stat
	}

	Stats = stats
	idToStat = newIDToStat
}
