package vcgen

import (
	"errors"
	"math"
	"math/rand"
	"strings"
	"time"
)

var (
	// charset types
	numbers      = "0123456789"
	alphabetic   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanumeric = numbers + alphabetic

	// defaults
	minCount  uint16 = 1
	minLength uint16 = 6
)

const patternChar = "#"

var (
	ErrNotFeasible       = errors.New("Not feasible to generate requested number of codes")
	ErrPatternIsNotMatch = errors.New("Pattern is not match with the length value")
)

// initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

type Generator struct {
	// Length of the code
	Length uint16 `json:"length"`

	// Count of the codes
	Count uint16 `json:"count"`

	// Charset to use
	Charset string `json:"charset"`

	// Prefix of the code
	Prefix string `json:"prefix"`

	// Suffix of the code
	Suffix string `json:"suffix"`

	// Pattern of the code
	Pattern string `json:"pattern"`

	// Suffix of the code
	//
	// Deprecated: use Suffix instead
	Postfix string `json:"postfix"`
}

// Creates a new generator with options
func NewWithOptions(opts ...Option) (*Generator, error) {
	g := Default()
	if err := setOptions(opts...)(g); err != nil {
		return nil, err
	}

	return g, nil
}

// Creates a new generator with default values
func Default() *Generator {
	return &Generator{
		Length:  minLength,
		Count:   minCount,
		Charset: alphanumeric,
		Pattern: repeatStr(minLength, patternChar),
	}
}

// New generator config
//
// Deprecated: use NewWithOptions or Default instead
func New(i *Generator) *Generator {

	// check functione entry args
	if i == nil {
		rps := repeatStr(minLength, "#")

		// --
		i = &Generator{
			Length:  minLength,
			Count:   minCount,
			Charset: alphanumeric,
			Pattern: rps,
		}
		goto AfterChecker
	}

	// check pattern
	if i.Pattern == "" {
		rps := repeatStr(i.Length, "#")
		i.Pattern = rps
	}

	// check length
	if i.Length == 0 {
		i.Length = uint16(strings.Count(i.Pattern, "#"))
	}

	// check count
	if i.Count == 0 {
		i.Count = minCount
	}

	// check charset
	if i.Charset == "" {
		i.Charset = alphanumeric
	}

AfterChecker:

	// return
	return i
}

// Generates a list of codes
func (g *Generator) Run() ([]string, error) {
	if !isFeasible(g.Charset, g.Pattern, patternChar, g.Count) {
		return nil, ErrNotFeasible
	}

	result := make([]string, g.Count)

	var i uint16
	for i = 0; i < g.Count; i++ {
		result[i] = g.one()
	}

	return result, nil
}

// one generates one code
func (g *Generator) one() string {
	pts := strings.Split(g.Pattern, "")
	for i, v := range pts {
		if v == patternChar {
			pts[i] = randomChar([]byte(g.Charset))
		}
	}

	suffix := g.Suffix

	// TODO: Remove it when we deprecate Postfix
	if g.Postfix != "" {
		suffix = g.Postfix
	}

	return g.Prefix + strings.Join(pts, "") + suffix
}

var (
	ErrInvalidCount   = errors.New("invalid count, It should be greater than 0")
	ErrInvalidCharset = errors.New("invalid charset, charset length should be greater than 0")
	ErrInvalidPattern = errors.New("invalid pattern, pattern cannot be empty")
)

type Option func(*Generator) error

// SetLength sets the length of the code
func SetLength(length uint16) Option {
	return func(g *Generator) error {
		if length == 0 {
			length = numberOfChar(g.Pattern, patternChar)
		}
		g.Length = length
		return nil
	}
}

// SetCount sets the count of the code
func SetCount(count uint16) Option {
	return func(g *Generator) error {
		if count == 0 {
			return ErrInvalidCount
		}
		g.Count = count
		return nil
	}
}

// SetCharset sets the charset of the code
func SetCharset(charset string) Option {
	return func(g *Generator) error {
		if len(charset) == 0 {
			return ErrInvalidCharset
		}
		g.Charset = charset
		return nil
	}
}

// SetPrefix sets the prefix of the code
func SetPrefix(prefix string) Option {
	return func(g *Generator) error {
		g.Prefix = prefix
		return nil
	}
}

// SetSuffix sets the suffix of the code
func SetSuffix(suffix string) Option {
	return func(g *Generator) error {
		g.Suffix = suffix
		return nil
	}
}

// SetPattern sets the pattern of the code
func SetPattern(pattern string) Option {
	return func(g *Generator) error {
		if pattern == "" {
			return ErrInvalidPattern
		}

		numPatternChar := numberOfChar(pattern, patternChar)
		if g.Length == 0 || g.Length != numPatternChar {
			g.Length = numPatternChar
		}

		g.Pattern = pattern
		return nil
	}
}

func setOptions(opts ...Option) Option {
	return func(g *Generator) error {
		for _, opt := range opts {
			if err := opt(g); err != nil {
				return err
			}
		}
		return nil
	}
}

// HELPER

// return random int in the range min...max
func randomInt(min, max int) int {
	return min + rand.Intn(1+max-min)
}

// return random char string from charset
func randomChar(cs []byte) string {
	return string(cs[randomInt(0, len(cs)-1)])
}

// repeat string with one str (#)
func repeatStr(count uint16, str string) string {
	return strings.Repeat(str, int(count))
}

func numberOfChar(str, char string) uint16 {
	return uint16(strings.Count(str, char))
}

func isFeasible(charset, pattern, char string, count uint16) bool {
	ls := numberOfChar(pattern, char)
	return math.Pow(float64(len(charset)), float64(ls)) >= float64(count)
}
