#!/usr/bin/env python3
"""Hyperfine workload for the cross-language benchmark in issue #8.

Iterates the full sequence of one operation, sum-accumulating a data-dependent
value per item (masked to 64 bits to match the Go and Rust drivers' wrapping
``u64``), and prints the accumulator once at the end. The data-dependent sum
defeats dead-code elimination and gives a cheap, observable per-iteration cost.

Usage: driver.py <op> <args...>
    combinations | cwr | permutations | product  <n> <k>
    powerset                                     <n>
    productof                                    <size>...   (one axis per size)
    intpartitions                                <n>
    setpartitions                                <n> [k]

The first six ops use stdlib ``itertools`` and sum every index. intpartitions
uses ``sympy`` and sums every part; setpartitions uses ``more-itertools`` and
sums assignment[i] * i over the restricted-growth-string block assignment
(blocks renumbered by first appearance to match the Go and Rust drivers).
sympy/more_itertools are imported lazily so the itertools ops don't pay for
them.
"""
import sys
from itertools import (
    chain,
    combinations,
    combinations_with_replacement,
    permutations,
    product,
)

MASK = 0xFFFFFFFFFFFFFFFF


def index_iter(op, nums):
    if op == "combinations":
        return combinations(range(nums[0]), nums[1])
    if op == "cwr":
        return combinations_with_replacement(range(nums[0]), nums[1])
    if op == "permutations":
        return permutations(range(nums[0]), nums[1])
    if op == "product":
        return product(range(nums[0]), repeat=nums[1])
    if op == "powerset":
        # itertools powerset recipe
        n = nums[0]
        return chain.from_iterable(combinations(range(n), r) for r in range(n + 1))
    if op == "productof":
        return product(*(range(size) for size in nums))
    raise SystemExit(f"unknown op {op!r}")


def accumulate(op, nums):
    if op == "intpartitions":
        from sympy.utilities.iterables import partitions

        # Each partition arrives as a {part: multiplicity} dict; sum the parts.
        return sum(
            part * mult for p in partitions(nums[0]) for part, mult in p.items()
        )

    if op == "setpartitions":
        from more_itertools import set_partitions

        n = nums[0]
        k = nums[1] if len(nums) > 1 else None
        acc = 0
        elem_block = [0] * n
        for blocks in set_partitions(range(n), k):
            for j, block in enumerate(blocks):
                for x in block:
                    elem_block[x] = j
            # Renumber blocks by first appearance (restricted growth string).
            labels = {}
            for i in range(n):
                b = elem_block[i]
                if b not in labels:
                    labels[b] = len(labels)
                acc += labels[b] * i
        return acc

    return sum(idx for tup in index_iter(op, nums) for idx in tup)


def main():
    args = sys.argv[1:]
    if len(args) < 2:
        raise SystemExit("usage: driver.py <op> <args...>")
    op = args[0]
    nums = [int(a) for a in args[1:]]
    print(accumulate(op, nums) & MASK)


if __name__ == "__main__":
    main()
