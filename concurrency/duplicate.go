package concurrency

import (
	"context"
)

// Consume from the 'in' channel and produce the same value to all of the output channels
func DuplicateOutputs[T any](ctx context.Context, in <-chan T, outs ...chan T) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case data, ok := <-in:
			if !ok {
				break loop
			}
			for _, out := range outs {
				out <- data
			}
		}
	}

	for _, out := range outs {
		close(out)
	}
}

// Consume from the 'in' channel and produce the a transformed value to the output channels
// Meaning consume T and produce V
func DuplicateTransformedOutputs[T any, V any](ctx context.Context,
	transformer func(in T) V,
	in <-chan T, outs ...chan V) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case data, ok := <-in:
			if !ok {
				break loop
			}
			for _, out := range outs {
				out <- transformer(data)
			}
		}
	}

	for _, out := range outs {
		close(out)
	}
}
