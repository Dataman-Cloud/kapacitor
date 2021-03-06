package pipeline

import (
	"errors"
	"fmt"

	"github.com/influxdata/kapacitor/tick/ast"
)

// A GroupByNode will group the incoming data.
// Each group is then processed independently for the rest of the pipeline.
// Only tags that are dimensions in the grouping will be preserved;
// all other tags are dropped.
//
// Example:
//    stream
//        |groupBy('service', 'datacenter')
//        ...
//
// The above example groups the data along two dimensions `service` and `datacenter`.
// Groups are dynamically created as new data arrives and each group is processed
// independently.
type GroupByNode struct {
	chainnode
	//The dimensions by which to group to the data.
	// tick:ignore
	Dimensions []interface{}
}

func newGroupByNode(wants EdgeType, dims []interface{}) *GroupByNode {
	return &GroupByNode{
		chainnode:  newBasicChainNode("groupby", wants, wants),
		Dimensions: dims,
	}
}
func (n *GroupByNode) validate() error {
	return validateDimensions(n.Dimensions)
}

func validateDimensions(dimensions []interface{}) error {
	hasStar := false
	for _, d := range dimensions {
		switch dim := d.(type) {
		case string:
			if len(dim) == 0 {
				return errors.New("dimensions cannot not be the empty string")
			}
		case *ast.StarNode:
			hasStar = true
		default:
			return fmt.Errorf("invalid dimension object of type %T", d)
		}
	}
	if hasStar && len(dimensions) > 1 {
		return errors.New("cannot group by both '*' and named dimensions.")
	}
	return nil
}
