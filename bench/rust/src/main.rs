//! Hyperfine workload for the cross-language benchmark in issue #8.
//!
//! Iterates the full sequence of one `itertools` operation, sum-accumulating
//! (wrapping `u64`) every index it sees, and prints the accumulator once at the
//! end. The data-dependent sum defeats dead-code elimination and gives a cheap,
//! observable per-iteration cost matching the Go and Python drivers.
//!
//! Usage: rust-driver <op> <n> [k]
//!     op  in combinations | cwr | permutations | product | powerset
//!     n   size of the input (elements are the ints 0..n-1)
//!     k   required for every op except powerset
//!
//! Rust's itertools yields an owned `Vec` per item (no buffer reuse), so this
//! driver is most comparable to the Go driver's allocating `All()` path.

use itertools::Itertools;
use std::env;
use std::process::exit;

fn main() {
    let args: Vec<String> = env::args().skip(1).collect();
    if args.len() < 2 {
        eprintln!("usage: rust-driver <op> <n> [k]");
        exit(2);
    }
    let op = args[0].as_str();
    let n: usize = args[1].parse().expect("invalid n");
    let k: usize = if op == "powerset" {
        0
    } else {
        args.get(2).expect("op requires k").parse().expect("invalid k")
    };

    let data: Vec<u64> = (0..n as u64).collect();

    // Each branch folds the per-item index sum into a single wrapping accumulator.
    let fold = |acc: u64, item: Vec<u64>| {
        item.iter().fold(acc, |a, &x| a.wrapping_add(x))
    };

    let acc: u64 = match op {
        "combinations" => data.iter().copied().combinations(k).fold(0, fold),
        "cwr" => data
            .iter()
            .copied()
            .combinations_with_replacement(k)
            .fold(0, fold),
        "permutations" => data.iter().copied().permutations(k).fold(0, fold),
        "product" => itertools::repeat_n(data.clone(), k)
            .multi_cartesian_product()
            .fold(0, fold),
        "powerset" => data.iter().copied().powerset().fold(0, fold),
        other => {
            eprintln!("unknown op {other:?}");
            exit(2);
        }
    };

    println!("{acc}");
}
