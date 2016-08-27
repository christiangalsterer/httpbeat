package flag

import (
	goflag "flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/go-ucfg"
)

func TestFlagValuePrimitives(t *testing.T) {
	config, err := parseTestFlags("-D b=true -D b2 -D i=42 -D f=3.14 -D s=string")
	assert.NoError(t, err)

	//validate
	checkFields(t, config)
}

func TestFlagValueLast(t *testing.T) {
	config, err := parseTestFlags("-D b2=false -D i=23 -D s=test -D b=true -D b2 -D i=42 -D f=3.14 -D s=string")
	assert.NoError(t, err)

	//validate
	checkFields(t, config)

}

func TestFlagValueNested(t *testing.T) {
	config, err := parseTestFlags("-D c.b=true -D c.b2 -D c.i=42 -D c.f=3.14 -D c.s=string")
	assert.NoError(t, err)

	// validate
	sub, err := config.Child("c", -1)
	assert.NoError(t, err)
	assert.NotNil(t, sub)

	if sub != nil {
		checkFields(t, sub)
	}
}

func parseTestFlags(args string) (*ucfg.Config, error) {
	fs := goflag.NewFlagSet("test", goflag.ContinueOnError)
	config := Config(fs, "D", "overwrite", ucfg.PathSep("."))
	err := fs.Parse(strings.Split(args, " "))
	return config, err
}

func checkFields(t *testing.T, config *ucfg.Config) {
	b, err := config.Bool("b", -1, ucfg.PathSep("."))
	assert.NoError(t, err)
	assert.Equal(t, true, b)

	b, err = config.Bool("b2", -1, ucfg.PathSep("."))
	assert.NoError(t, err)
	assert.Equal(t, true, b)

	i, err := config.Int("i", -1, ucfg.PathSep("."))
	assert.NoError(t, err)
	assert.Equal(t, 42, int(i))

	f, err := config.Float("f", -1, ucfg.PathSep("."))
	assert.NoError(t, err)
	assert.Equal(t, 3.14, f)

	s, err := config.String("s", -1, ucfg.PathSep("."))
	assert.NoError(t, err)
	assert.Equal(t, "string", s)
}
