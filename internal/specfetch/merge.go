package specfetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

var httpMethods = map[string]bool{
	"get": true, "put": true, "post": true, "delete": true,
	"options": true, "head": true, "patch": true, "trace": true,
}

// Merge combines the single-operation documents of the reference pages,
// keyed by page URL, into one OpenAPI document. Parts every page repeats
// (openapi, info, servers, components) must agree, so a page that drifted
// from the rest fails the run instead of being silently dropped.
func Merge(docs map[string][]byte) (map[string]any, error) {
	out := map[string]any{}
	for _, name := range slices.Sorted(maps.Keys(docs)) {
		if err := mergeDoc(out, docs[name]); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
	}
	return out, nil
}

func mergeDoc(out map[string]any, data []byte) error {
	// json.Number keeps bounds such as int64 maxima exact through the round trip.
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var doc map[string]any
	if err := dec.Decode(&doc); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	if doc == nil {
		return errors.New("not a JSON object")
	}
	ops := 0
	for key, val := range doc {
		switch key {
		case "paths":
			dst, err := child(out, "paths", "paths")
			if err != nil {
				return err
			}
			n, err := mergePaths(dst, val)
			if err != nil {
				return err
			}
			ops += n
		case "components":
			if err := mergeComponents(out, val); err != nil {
				return err
			}
		default:
			if err := setOrMatch(out, key, val, key); err != nil {
				return err
			}
		}
	}
	if ops == 0 {
		return errors.New("page defines no operation")
	}
	return nil
}

func mergePaths(dst map[string]any, val any) (int, error) {
	paths, err := object(val, "paths")
	if err != nil {
		return 0, err
	}
	ops := 0
	for path, rawItem := range paths {
		label := "paths." + path
		item, err := object(rawItem, label)
		if err != nil {
			return 0, err
		}
		dstItem, err := child(dst, path, label)
		if err != nil {
			return 0, err
		}
		for key, v := range item {
			if !httpMethods[key] {
				if err := setOrMatch(dstItem, key, v, label+"."+key); err != nil {
					return 0, err
				}
				continue
			}
			if _, dup := dstItem[key]; dup {
				return 0, fmt.Errorf("duplicate operation %s %s", strings.ToUpper(key), path)
			}
			dstItem[key] = v
			ops++
		}
	}
	return ops, nil
}

func mergeComponents(out map[string]any, val any) error {
	comps, err := object(val, "components")
	if err != nil {
		return err
	}
	dst, err := child(out, "components", "components")
	if err != nil {
		return err
	}
	for kind, rawEntries := range comps {
		label := "components." + kind
		entries, err := object(rawEntries, label)
		if err != nil {
			return err
		}
		dstKind, err := child(dst, kind, label)
		if err != nil {
			return err
		}
		for name, v := range entries {
			if err := setOrMatch(dstKind, name, v, label+"."+name); err != nil {
				return err
			}
		}
	}
	return nil
}

func object(v any, label string) (map[string]any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s is not an object", label)
	}
	return m, nil
}

// child returns dst[key] as an object, creating it when absent.
func child(dst map[string]any, key, label string) (map[string]any, error) {
	v, ok := dst[key]
	if !ok {
		m := map[string]any{}
		dst[key] = m
		return m, nil
	}
	return object(v, label)
}

// setOrMatch stores val, or checks that it equals what an earlier page stored.
func setOrMatch(dst map[string]any, key string, val any, label string) error {
	prev, ok := dst[key]
	if !ok {
		dst[key] = val
		return nil
	}
	if !reflect.DeepEqual(prev, val) {
		return fmt.Errorf("%s differs from earlier pages", label)
	}
	return nil
}
