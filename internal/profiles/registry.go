package profiles

import (
	"fmt"
	"sync"

	"github.com/supperdoggy/vr_api/internal/apperr"
)

type PageType string

const (
	PageSearch     PageType = "search"
	PageLawProject PageType = "law_project"
)

type SearchSelectorSet interface {
	ResultsContainer() string
	ResultRows() string
	ID() string
	Title() string
	TitleLink() string
	Status() string
	RegistrationDate() string
	Subject() string
	LinkURLAttr() string
}

type LawProjectSelectorSet interface {
	RootContainer() string
	Title() string
	SummaryRows() string
	SummaryLabel() string
	SummaryValue() string
	IDNode() string
	TimelineHeaderStatus() string
	TimelineRows() string
	TimelineDate() string
	TimelineStatus() string
	TimelineNote() string
	RelatedRows() string
	RelatedNumberLink() string
	RelatedRegistrationDate() string
	RelatedTitle() string
	LinkURLAttr() string
	ChronologyEmbed() string
	VotingResultsEmbed() string
}

type SearchProfile interface {
	ID() string
	Version() string
	Selectors() SearchSelectorSet
}

type LawProjectProfile interface {
	ID() string
	Version() string
	Selectors() LawProjectSelectorSet
}

type Registry struct {
	mu sync.RWMutex

	searchProfiles map[string]SearchProfile
	lawProfiles    map[string]LawProjectProfile

	defaultSearchVersion string
	defaultLawVersion    string
}

func NewRegistry() *Registry {
	r := &Registry{
		searchProfiles: make(map[string]SearchProfile),
		lawProfiles:    make(map[string]LawProjectProfile),
	}
	r.RegisterSearch(newSearchV1Profile())
	r.RegisterLawProject(newLawProjectV1Profile())
	return r
}

func (r *Registry) RegisterSearch(profile SearchProfile) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.searchProfiles[profile.Version()] = profile
	if r.defaultSearchVersion == "" {
		r.defaultSearchVersion = profile.Version()
	}
}

func (r *Registry) RegisterLawProject(profile LawProjectProfile) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lawProfiles[profile.Version()] = profile
	if r.defaultLawVersion == "" {
		r.defaultLawVersion = profile.Version()
	}
}

func (r *Registry) Search(version string) (SearchProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if version == "" {
		version = r.defaultSearchVersion
	}

	profile, ok := r.searchProfiles[version]
	if !ok {
		return nil, apperr.NewUnsupportedProfileError(
			"profiles.search",
			fmt.Sprintf("page=%s version=%q", PageSearch, version),
			nil,
		)
	}

	return profile, nil
}

func (r *Registry) LawProject(version string) (LawProjectProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if version == "" {
		version = r.defaultLawVersion
	}

	profile, ok := r.lawProfiles[version]
	if !ok {
		return nil, apperr.NewUnsupportedProfileError(
			"profiles.law_project",
			fmt.Sprintf("page=%s version=%q", PageLawProject, version),
			nil,
		)
	}

	return profile, nil
}

func (r *Registry) DefaultSearchVersion() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.defaultSearchVersion
}

func (r *Registry) DefaultLawProjectVersion() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.defaultLawVersion
}
