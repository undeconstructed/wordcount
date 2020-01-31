
# wordcount

## What is

Uses no library structures or algorithms at all. Very compact, reusing buffers for reading file and building words, and combining an array and a tree in one.

I didn't want to reimplement a generic structure, so the tree is entirely specific to the task. It's not at all practical to rebalance it, but the input words are likely quite random re. alphabetical order, so I don't imagine it's too problematic.

The tree works by having each node reference left and right children, which are just indexes back into the containing array. This means that iteration in random order (used for finding topN list) is just iteration through a list.

There's no generic sort used, instead items are inserted in-order into the topN list, so the last in the list is always the first to be discarded.

## Concurrency ..

Could run parsing and tabulating concurrently, but that would mean more memory as couldn't immediately reuse the word buffer.

Could pre-divide the input, but there's no particularly natural way to do that.

So probably not enough benefit to bother.

## Bugs?

Not actually the same output as the given script in the case of words with same count. The script sorts these reverse-alphabetically, the app sorts randomly.

## TODO

* If the input finishes on a word character?
