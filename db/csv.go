package db

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	warehouse "github.com/diegocmsantos/warehouse/core"
)

type LineValidator struct {
	Validators []func([]string) (bool, error)
}

// CSVSource represents a CSVSource of elements
type CSVSource struct {
	reader    io.Reader
	validator LineValidator
}

// New creates a new CSVSource based on the given filename
func New(file string, lineValidators ...func([]string) (bool, error)) (*CSVSource, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", file, err)
	}

	vals := LineValidator{Validators: make([]func([]string) (bool, error), 0)}
	for _, f := range lineValidators {
		vals.Validators = append(vals.Validators, f)
	}

	return &CSVSource{reader: f, validator: vals}, nil
}

// NewWithReader creates a new CSVSource based on the given reader
func NewWithReader(reader io.Reader, lineValidators ...func([]string) (bool, error)) (*CSVSource, error) {

	if reader == nil {
		return nil, errors.New("reader cannot be nil")
	}

	vals := LineValidator{Validators: make([]func([]string) (bool, error), 0)}
	for _, f := range lineValidators {
		vals.Validators = append(vals.Validators, f)
	}
	return &CSVSource{reader: reader, validator: vals}, nil
}

// GetElements build the tree based on the given CSV file
func (s *CSVSource) GetElements() (*warehouse.Element, error) {
	r := csv.NewReader(s.reader)

	record, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	return s.buildElementTree(record)
}

func (s *CSVSource) buildElementTree(ar [][]string) (*warehouse.Element, error) {

	header := ar[0]
	isRegularOrder := strings.Contains(header[len(header)-1], "item")

	e := warehouse.Element{}
	var currentElem *warehouse.Element
	for _, v := range ar[1:] {

		if !isRegularOrder {
			v = fixLineOrder(v)
		}

		for _, val := range s.validator.Validators {
			if isVal, err := val(v); !isVal {
				return nil, err
			}
		}

		i, ok := e.Children[v[0]]
		if ok {
			currentElem = i
		} else {
			currentElem = &warehouse.Element{}
			if e.Children == nil {
				e.Children = make(map[string]*warehouse.Element, 0)
			}
			e.Children[v[0]] = currentElem
		}
		for idx, k := range v[1:] {
			if k == "" {
				continue
			}
			li, ok := currentElem.Children[k]
			if ok {
				currentElem = li
			} else {
				if currentElem.Children == nil {
					currentElem.Children = make(map[string]*warehouse.Element, 0)
				}
				currentElem.Children[k] = &warehouse.Element{Item: len(v[1:])-1 == idx}
				currentElem = currentElem.Children[k]
			}
		}
	}
	return &e, nil
}

func fixLineOrder(line []string) []string {
	res := make([]string, 0, len(line))
	res = append(res, line[1:]...)
	res = append(res, line[0])
	return res
}
