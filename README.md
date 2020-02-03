
# wordcount

```
$ go run . test.txt
```

## What?

Counts words in an input file. Uses no library structures or algorithms at all, doesn't use string type. Very compact, reusing buffers for reading file and building words, and combining an array and a tree in one.

The tree works by storing each node in array, and each node having references to left and right children, which are just indexes back into the containing array. This means that iteration in random order (used for finding topN list) is just iteration through a list.

I didn't want to reimplement generic structures, so the tree is entirely specific to the task. It's not at all practical to rebalance it, but as long as the input words are fairly random re. alphabetical order, it won't be a major problem.

There's no generic sort used, instead items are inserted in-order into the topN list, so the last in the list is always the first to be discarded. Items are inserted by scanning from the start, but with only ~20 elements in the list that's probably not much slower than a proper search.

## Concurrency?

Could run parsing and tabulating concurrently, but that would mean more memory as couldn't immediately reuse the word buffer.

Could pre-divide the input, but there's no particularly natural way to do that, and the results would have to merged.

So probably not enough benefit to bother.

## Bugs?

Not actually the same output as the given script in the case of words with same count. The script sorts these reverse-alphabetically, the app sorts arbitrarily.

ASCII only.
