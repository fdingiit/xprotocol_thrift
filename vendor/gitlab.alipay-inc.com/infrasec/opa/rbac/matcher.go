package rbac

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	rbactypes "gitlab.alipay-inc.com/infrasec/api/types"
)

// StringMatcher
type StringMatcher interface {
	isStringMatcher()
	Equal(string) bool
}

func (matcher *StringMatcherExactMatch) isStringMatcher()  {}
func (matcher *StringMatcherPrefixMatch) isStringMatcher() {}
func (matcher *StringMatcherSuffixMatch) isStringMatcher() {}
func (matcher *StringMatcherRegexMatch) isStringMatcher()  {}

// StringMatcherConf_Exact
type StringMatcherExactMatch struct {
	ExactMatch  string
	InvertMatch bool
}

func (matcher *StringMatcherExactMatch) Equal(targetValue string) bool {
	isMatch := matcher.ExactMatch == targetValue
	// InvertMatch xor isMatch
	return isMatch != matcher.InvertMatch
}

// StringMatcherConf_Prefix
type StringMatcherPrefixMatch struct {
	PrefixMatch string
	InvertMatch bool
}

func (matcher *StringMatcherPrefixMatch) Equal(targetValue string) bool {
	isMatch := strings.HasPrefix(targetValue, matcher.PrefixMatch)
	// InvertMatch xor isMatch
	return isMatch != matcher.InvertMatch
}

// StringMatcherConf_Suffix
type StringMatcherSuffixMatch struct {
	SuffixMatch string
	InvertMatch bool
}

func (matcher *StringMatcherSuffixMatch) Equal(targetValue string) bool {
	isMatch := strings.HasSuffix(targetValue, matcher.SuffixMatch)
	// InvertMatch xor isMatch
	return isMatch != matcher.InvertMatch
}

// StringMatcherConf_Regex
type StringMatcherRegexMatch struct {
	RegexMatch  *regexp.Regexp
	InvertMatch bool
}

func (matcher *StringMatcherRegexMatch) Equal(targetValue string) bool {
	isMatch := matcher.RegexMatch.MatchString(targetValue)
	// InvertMatch xor isMatch
	return isMatch != matcher.InvertMatch
}

func NewStringMatcher(matcher *rbactypes.StringMatcherConf) (StringMatcher, error) {
	switch matcher.GetMatchPattern().(type) {
	case *rbactypes.StringMatcherConf_ExactMatch:
		return &StringMatcherExactMatch{
			ExactMatch:  matcher.GetMatchPattern().(*rbactypes.StringMatcherConf_ExactMatch).ExactMatch,
			InvertMatch: matcher.GetInvertMatch(),
		}, nil
	case *rbactypes.StringMatcherConf_PrefixMatch:
		return &StringMatcherPrefixMatch{
			PrefixMatch: matcher.GetMatchPattern().(*rbactypes.StringMatcherConf_PrefixMatch).PrefixMatch,
			InvertMatch: matcher.GetInvertMatch(),
		}, nil
	case *rbactypes.StringMatcherConf_SuffixMatch:
		return &StringMatcherSuffixMatch{
			SuffixMatch: matcher.GetMatchPattern().(*rbactypes.StringMatcherConf_SuffixMatch).SuffixMatch,
			InvertMatch: matcher.GetInvertMatch(),
		}, nil
	case *rbactypes.StringMatcherConf_RegexMatch:
		if rePattern, err := regexp.Compile(
			matcher.GetMatchPattern().(*rbactypes.StringMatcherConf_RegexMatch).RegexMatch); err != nil {
			return nil, fmt.Errorf("[NewStringMatcher] failed to build regex, error: %v", err)
		} else {
			return &StringMatcherRegexMatch{
				RegexMatch:  rePattern,
				InvertMatch: matcher.GetInvertMatch(),
			}, nil
		}
	default:
		return nil, fmt.Errorf(
			"[NewStringMatcher] not support MatchPattern type found, detail: %v",
			reflect.TypeOf(matcher.GetMatchPattern()))
	}
}

// HeaderMatcher
type HeaderMatcher interface {
	isHeaderMatcher()
	Equal(string) bool
}

func (matcher *StringMatcherExactMatch) isHeaderMatcher()   {}
func (matcher *StringMatcherPrefixMatch) isHeaderMatcher()  {}
func (matcher *StringMatcherSuffixMatch) isHeaderMatcher()  {}
func (matcher *StringMatcherRegexMatch) isHeaderMatcher()   {}
func (matcher *HeaderMatcherPresentMatch) isHeaderMatcher() {}
func (matcher *HeaderMatcherRangeMatch) isHeaderMatcher()   {}

// HeaderMatcherConf_PresentMatch
type HeaderMatcherPresentMatch struct {
	PresentMatch bool
}

func (matcher *HeaderMatcherPresentMatch) Equal(targetValue string) bool {
	return matcher.PresentMatch
}

// HeaderMatcherConf_RangeMatch
type HeaderMatcherRangeMatch struct {
	Start       int64 // inclusive
	End         int64 // exclusive
	InvertMatch bool
}

func (matcher *HeaderMatcherRangeMatch) Equal(targetValue string) bool {
	intValue, err := strconv.ParseInt(targetValue, 10, 64)
	if err != nil {
		return false
	}
	isMatch := intValue >= matcher.Start && intValue < matcher.End
	// InvertMatch xor isMatch
	return isMatch != matcher.InvertMatch
}

func NewHeaderMatcher(matcher *rbactypes.HeaderMatcherConf) (HeaderMatcher, error) {
	switch matcher.GetHeaderMatchSpecifier().(type) {
	case *rbactypes.HeaderMatcherConf_ExactMatch:
		return &StringMatcherExactMatch{
			ExactMatch: matcher.GetHeaderMatchSpecifier().(*rbactypes.HeaderMatcherConf_ExactMatch).ExactMatch,
		}, nil
	case *rbactypes.HeaderMatcherConf_PrefixMatch:
		return &StringMatcherPrefixMatch{
			PrefixMatch: matcher.GetHeaderMatchSpecifier().(*rbactypes.HeaderMatcherConf_PrefixMatch).PrefixMatch,
		}, nil
	case *rbactypes.HeaderMatcherConf_SuffixMatch:
		return &StringMatcherSuffixMatch{
			SuffixMatch: matcher.GetHeaderMatchSpecifier().(*rbactypes.HeaderMatcherConf_SuffixMatch).SuffixMatch,
		}, nil
	case *rbactypes.HeaderMatcherConf_RegexMatch:
		if rePattern, err := regexp.Compile(
			matcher.GetHeaderMatchSpecifier().(*rbactypes.HeaderMatcherConf_RegexMatch).RegexMatch); err != nil {
			return nil, fmt.Errorf("[NewHeaderMatcher] failed to build regex, error: %v", err)
		} else {
			return &StringMatcherRegexMatch{
				RegexMatch: rePattern,
			}, nil
		}
	case *rbactypes.HeaderMatcherConf_PresentMatch:
		return &HeaderMatcherPresentMatch{
			PresentMatch: matcher.GetHeaderMatchSpecifier().(*rbactypes.HeaderMatcherConf_PresentMatch).PresentMatch,
		}, nil
	case *rbactypes.HeaderMatcherConf_RangeMatch:
		return &HeaderMatcherRangeMatch{
			Start:       matcher.GetHeaderMatchSpecifier().(*rbactypes.HeaderMatcherConf_RangeMatch).RangeMatch.Start,
			End:         matcher.GetHeaderMatchSpecifier().(*rbactypes.HeaderMatcherConf_RangeMatch).RangeMatch.End,
			InvertMatch: matcher.GetHeaderMatchSpecifier().(*rbactypes.HeaderMatcherConf_RangeMatch).RangeMatch.InvertMatch,
		}, nil
	default:
		return nil, fmt.Errorf(
			"[NewHeaderMatcher] not support HeaderMatchSpecifier type found, detail: %v",
			reflect.TypeOf(matcher.GetHeaderMatchSpecifier()))
	}
}
