package ftp

import (
	"context"
	"fmt"
	"strconv"

	"go.beyondstorage.io/v5/types"
)

func (s *Storage) listDirNext(ctx context.Context, page *types.ObjectPage) (err error) {
	input := page.Status.(*listDirInput)
	if input.objList == nil {
		input.objList, err = s.connection.List(input.rp)
		if err != nil {
			return err
		}
	}
	if !input.started {
		input.counter, err = strconv.Atoi(input.continuationToken)
		if err != nil {
			input.counter = 0
		}
		input.started = true
	}
	n := len(input.objList)
	input.continuationToken = fmt.Sprintf("%x", input.counter)
	if input.counter >= n {
		return types.IterateDone
	}

	v := input.objList[input.counter]

	obj, err := s.formatFileObject(v, input.rp)
	if err != nil {
		return err
	}

	page.Data = append(page.Data, obj)

	input.counter++

	return
}
