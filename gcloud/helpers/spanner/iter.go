package spanner

import (
	"cloud.google.com/go/spanner"
)

// GetOneWithIterator is a helper function to get one row from a row iterator
// You can use it to simplify your code when you know there is only one row in the iterator.
func GetOneWithIterator(iter *spanner.RowIterator) (row *spanner.Row, exist bool, err error) {
	err = iter.Do(func(r *spanner.Row) error {
		exist = true
		row = r
		return nil
	})

	return
}

// IterToStructs is a helper function to process a row iterator and convert each row to a struct
// This function will return an array of structs and/or the spanner error if one occurs.
func IterToStructs[T any](iter *spanner.RowIterator) ([]T, error) {
	var results []T

	err := iter.Do(func(r *spanner.Row) error {
		var dest T
		if err := r.ToStruct(&dest); err != nil {
			return err
		}

		results = append(results, dest)
		return nil
	})

	return results, err
}

// IterToStruct is a helper function to process and convert one row from a row iterator
// You can use it to simplify your code when you know there is only one row in the iterator.
func IterToStruct[T any](iter *spanner.RowIterator) (result T, exist bool, err error) {
	results, errIter := IterToStructs[T](iter)

	if err != nil {
		return result, false, errIter
	}

	if exist = len(results) > 0; exist {
		result = results[0]
	}

	return
}
