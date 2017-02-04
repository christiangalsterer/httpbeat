package ucfg

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEnv map[string]string

func TestStructMergeUnpackTyped(t *testing.T) {
	type test struct {
		t   interface{}
		cfg interface{}
		env testEnv
	}

	tests := []test{
		{
			t: &struct {
				Strings []string
			}{
				Strings: []string{"string1", "abc"},
			},
			cfg: map[string]interface{}{
				"strings": "${env_strings}",
			},
			env: testEnv{"env_strings": "string1,abc"},
		},
		{
			t: &struct {
				Strings []string
			}{
				Strings: []string{"string1", "abc"},
			},
			cfg: map[string]interface{}{
				"strings": "${env_strings:string1,abc}",
			},
		},
		{
			t: &struct {
				Strings []string
			}{
				Strings: []string{"one string"},
			},
			cfg: map[string]interface{}{
				"strings": "${env_strings} string",
			},
			env: testEnv{"env_strings": "one"},
		},

		{
			t: &struct {
				Hosts []string
			}{
				Hosts: []string{"host1:1234", "host2:4567"},
			},
			cfg: map[string]interface{}{
				"hosts": "${hosts_from_env}",
			},
			env: testEnv{"hosts_from_env": "host1:1234,host2:4567"},
		},
		{
			t: &struct {
				Hosts []string
			}{
				Hosts: []string{"host1:1234", "host2:4567"},
			},
			cfg: map[string]interface{}{
				"hosts": "${missing_env:host1:1234,host2:4567}",
			},
		},
	}

	for i, test := range tests {
		t.Logf("run test (%v): %v, %v", i, test.t, test.cfg)

		opts := []Option{
			VarExp,
			resolveTestEnv(test.env),
		}

		// unpack input
		c, err := NewFrom(test.t, opts...)
		if err != nil {
			t.Fatal(err)
		}

		// compute expected outcome
		var expected map[string]interface{}
		if err := c.Unpack(&expected, opts...); err != nil {
			t.Fatal(err)
		}

		// reset test.t to zero value
		v := chaseValue(reflect.ValueOf(test.t))
		v.Set(reflect.Zero(v.Type()))

		// create new config from test
		t.Logf("new from: %v", test.cfg)
		if c, err = NewFrom(test.cfg, opts...); err != nil {
			t.Fatal(err)
		}

		// unpack config into zeroed out config
		if err := c.Unpack(test.t, opts...); err != nil {
			t.Fatal(err)
		}

		// parse restored input config
		if c, err = NewFrom(test.t, opts...); err != nil {
			t.Fatal(err)
		}

		// validate
		var actual map[string]interface{}
		if err := c.Unpack(&actual, opts...); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, actual)
	}
}

func resolveTestEnv(e testEnv) Option {
	fail := func(name string) (string, error) {
		return "", fmt.Errorf("empty environment variable %v", name)
	}

	if e == nil {
		return Resolve(fail)
	}
	return Resolve(func(name string) (string, error) {
		if v := e[name]; v != "" {
			return v, nil
		}
		return fail(name)
	})
}
