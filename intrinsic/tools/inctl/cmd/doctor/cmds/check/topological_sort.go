// Copyright 2023 Intrinsic Innovation LLC

package check

import (
	"fmt"
	"slices"
)

func topologicalSortVisit[T any](
	entities *map[string]*T,
	dependencies func(entity *T) *[]string,
	entityName string,
	entity *T,
	unvisited *[]string,
	sorted *[]*T,
	current *[]string,
) error {
	if !slices.Contains(*unvisited, entityName) {
		// Already visited.
		return nil
	}
	if slices.Contains(*current, entityName) {
		return fmt.Errorf("circular dependency in topological sort of entities: %v", *current)
	}

	*current = append(*current, entityName)

	for _, dep := range *dependencies(entity) {
		depEntity, ok := (*entities)[dep]
		if !ok {
			return fmt.Errorf("entity '%q' has dependency '%q' which is not a known entity", entityName, dep)
		}
		if err := topologicalSortVisit(
			entities,
			dependencies,
			dep,
			depEntity,
			unvisited,
			sorted,
			current,
		); err != nil {
			return err
		}
	}

	for idx, checkName := range *current {
		if checkName == entityName {
			*current = append((*current)[:idx], (*current)[idx+1:]...)
			break
		}
	}
	for idx, checkName := range *unvisited {
		if checkName == entityName {
			*unvisited = append((*unvisited)[:idx], (*unvisited)[idx+1:]...)
			break
		}
	}

	*sorted = append(*sorted, entity)

	return nil
}

func topologicalSortWithUnvisited[T any](
	entities *map[string]*T,
	unvisited []string,
	dependencies func(entity *T) *[]string,
) ([]*T, error) {
	var sorted []*T

	for len(unvisited) > 0 {
		first, ok := (*entities)[unvisited[0]]
		if !ok {
			return nil, fmt.Errorf("entity '%q' unexpectedly not in the entities map", unvisited[0])
		}
		if err := topologicalSortVisit(
			entities,
			dependencies,
			unvisited[0],
			first,
			&unvisited,
			&sorted,
			&[]string{},
		); err != nil {
			return nil, err
		}
	}

	return sorted, nil
}

func topologicalSort[T any](
	entities *map[string]*T,
	dependencies func(entity *T) *[]string,
) ([]*T, error) {
	// The iteration order of the keys of an ordered map is not guaranteed, so we arbitrarily sort
	// the keys here.
	unvisited := make([]string, 0, len(*entities))
	for k := range *entities {
		unvisited = append(unvisited, k)
	}
	slices.Sort(unvisited)
	return topologicalSortWithUnvisited(
		entities,
		unvisited,
		dependencies,
	)
}
