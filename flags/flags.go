package flags

import (
	"fmt"
	"os"
	"strings"

	"github.com/launchdarkly/ld-find-code-refs/coderefs"
	"github.com/launchdarkly/ld-find-code-refs/element"
	"github.com/launchdarkly/ld-find-code-refs/internal/helpers"
	"github.com/launchdarkly/ld-find-code-refs/internal/ld"
	"github.com/launchdarkly/ld-find-code-refs/internal/log"
	"github.com/launchdarkly/ld-find-code-refs/internal/version"
	"github.com/launchdarkly/ld-find-code-refs/options"
)

const (
	minFlagKeyLen = 3 // Minimum flag key length helps reduce the number of false positives
)

func GenerateSearchElements(opts options.Options, repoParams ld.RepoParams, delims string) element.ElementMatcher {
	matcher := element.ElementMatcher{
		Directory: opts.Dir,
	}

	projKey := opts.ProjKey

	ldApi := ld.InitApiClient(ld.ApiOptions{ApiKey: opts.AccessToken, BaseUri: opts.BaseUri, ProjKey: projKey, UserAgent: "LDFindCodeRefs/" + version.Version})
	isDryRun := opts.DryRun

	ignoreServiceErrors := opts.IgnoreServiceErrors
	if !isDryRun {
		err := ldApi.MaybeUpsertCodeReferenceRepository(repoParams)
		if err != nil {
			helpers.FatalServiceError(err, ignoreServiceErrors)
		}
	}

	flags, err := getFlags(ldApi)
	if err != nil {
		helpers.FatalServiceError(fmt.Errorf("could not retrieve flag keys from LaunchDarkly: %w", err), ignoreServiceErrors)
	}

	filteredFlags, omittedFlags := filterShortFlagKeys(flags)
	if len(filteredFlags) == 0 {
		log.Info.Printf("no flag keys longer than the minimum flag key length (%v) were found for project: %s, exiting early",
			minFlagKeyLen, projKey)
		os.Exit(0)
	} else if len(omittedFlags) > 0 {
		log.Warning.Printf("omitting %d flags with keys less than minimum (%d)", len(omittedFlags), minFlagKeyLen)
	}
	matcher.Elements = filteredFlags

	matcher.Aliases, err = coderefs.GenerateAliases(filteredFlags, opts.Aliases, opts.Dir)
	matcher.DelimitedFlags = buildDelimiterList(flags, delims)
	if err != nil {
		log.Error.Fatalf("failed to create flag key aliases: %v", err)
	}

	return matcher
}

func buildDelimiterList(flags []string, delimiters string) map[string][]string {
	delimiterMap := make(map[string][]string)
	if delimiters == "" {
		return delimiterMap
	}
	for _, flag := range flags {
		//flagsDelimited := []string{}
		tempFlags := []string{}
		for _, left := range delimiters {
			for _, right := range delimiters {
				var sb strings.Builder
				sb.Grow(len(flag) + 2)
				sb.WriteRune(left)
				sb.WriteString(flag)
				sb.WriteRune(right)
				tempFlags = append(tempFlags, sb.String())
			}
		}
		delimiterMap[flag] = tempFlags
	}
	return delimiterMap
}

// Very short flag keys lead to many false positives when searching in code,
// so we filter them out.
func filterShortFlagKeys(flags []string) (filtered []string, omitted []string) {
	filteredFlags := []string{}
	omittedFlags := []string{}
	for _, flag := range flags {
		if len(flag) >= minFlagKeyLen {
			filteredFlags = append(filteredFlags, flag)
		} else {
			omittedFlags = append(omittedFlags, flag)
		}
	}
	return filteredFlags, omittedFlags
}

func getFlags(ldApi ld.ApiClient) ([]string, error) {
	flags, err := ldApi.GetFlagKeyList()
	if err != nil {
		return nil, err
	}
	return flags, nil
}
