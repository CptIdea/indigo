package radix

import (
	"context"
	"errors"
	routertypes "github.com/fakefloordiv/indigo/router/inbuilt/types"
	"github.com/fakefloordiv/indigo/valuectx"
)

var (
	ErrNotImplemented = errors.New(
		"different dynamic segment names are not allowed for common path prefix",
	)
)

type Tree interface {
	Insert(Template, routertypes.HandlerFunc) error
	MustInsert(Template, routertypes.HandlerFunc)
	Match(context.Context, string) (context.Context, routertypes.HandlerFunc)
}

type Node struct {
	staticSegments map[string]*Node
	isDynamic      bool
	dynamicName    string
	// Next is used only in case current node is dynamic
	next *Node

	handler routertypes.HandlerFunc
}

func NewTree() Tree {
	return newNode(nil, false, "")
}

func newNode(handler routertypes.HandlerFunc, isDyn bool, dynName string) *Node {
	return &Node{
		staticSegments: make(map[string]*Node),
		isDynamic:      isDyn,
		dynamicName:    dynName,
		handler:        handler,
	}
}

func (n *Node) Insert(template Template, handler routertypes.HandlerFunc) error {
	return n.insertRecursively(template.segments, handler)
}

func (n *Node) MustInsert(template Template, handler routertypes.HandlerFunc) {
	if err := n.Insert(template, handler); err != nil {
		panic(err.Error())
	}
}

func (n *Node) insertRecursively(segments []Segment, handler routertypes.HandlerFunc) error {
	if len(segments) == 0 {
		n.handler = handler

		return nil
	}

	segment := segments[0]

	if segment.IsDynamic {
		if n.isDynamic && segment.Payload != n.dynamicName {
			return ErrNotImplemented
		}

		n.isDynamic = true
		n.dynamicName = segment.Payload

		if n.next == nil {
			n.next = newNode(nil, false, "")
		}

		return n.next.insertRecursively(segments[1:], handler)
	}

	if node, found := n.staticSegments[segment.Payload]; found {
		return node.insertRecursively(segments[1:], handler)
	}

	node := newNode(nil, false, "")
	n.staticSegments[segment.Payload] = node

	return node.insertRecursively(segments[1:], handler)
}

func (n *Node) Match(ctx context.Context, path string) (context.Context, routertypes.HandlerFunc) {
	if path[0] != '/' {
		// all http request paths MUST have a leading slash
		return ctx, nil
	}

	path = path[1:]

	var (
		offset int
		node   = n
	)

	for i := range path {
		if path[i] == '/' {
			var ok bool

			ctx, node, ok = processSegment(ctx, path[offset:i], node)
			if !ok {
				return ctx, nil
			}

			offset = i + 1
		}
	}

	if offset < len(path) {
		var ok bool
		ctx, node, ok = processSegment(ctx, path[offset:], node)
		if !ok {
			return ctx, nil
		}
	}

	return ctx, node.handler
}

func processSegment(ctx context.Context, segment string, node *Node) (context.Context, *Node, bool) {
	if nextNode, found := node.staticSegments[segment]; found {
		return ctx, nextNode, true
	}

	if !node.isDynamic || len(segment) == 0 {
		return ctx, nil, false
	}

	if len(node.dynamicName) > 0 {
		ctx = valuectx.WithValue(ctx, node.dynamicName, segment)
	}

	return ctx, node.next, true
}
