package ast

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/maddiesch/marble/pkg/core/token"
)

// Node defines the base interface for a AST node
type Node interface {
	fmt.Stringer

	json.Marshaler

	Name() string

	SourceToken() token.Token
}

type Statement interface {
	Node

	_statementNode()
}

type Expression interface {
	Node

	_expressionNode()
}

func marshalNode(n Node) ([]byte, error) {
	data := map[string]any{
		"$self": map[string]any{
			"type": n.Name(),
			"token": map[string]any{
				"type":            n.SourceToken().Kind,
				"source_location": n.SourceToken().Location.String(),
			},
		},
	}

	t := reflect.TypeOf(n)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tag := f.Tag.Get("json")

		if strings.HasPrefix(tag, "-") {
			continue
		}

		parts := strings.SplitN(tag, ",", 2)

		key := f.Name
		if len(parts[0]) > 0 {
			key = parts[0]
		}

		if key == "Token" {
			continue
		}

		data[key] = reflect.Indirect(reflect.ValueOf(n)).FieldByName(f.Name).Interface()
	}

	return json.Marshal(data)
}
