# IP tree

store your ip ranges as a partial radix tree.

## Test

    make test

## Secret Sauce

Each significante bytes are a node, last node is a leaf.

class A has one node.

Class B has two nodes.

Class C has three nodes.

Nodes are sorted.

```
192.168.1.0/24 =>

-+- 0
 +- 1
 +- …
 +- 190    +- …
 +- 191    +- 167
 +- 192 ---+- 168 ----- 1 -> leaf
 +- 193    +- 169
 +- 194    +- …
 +- …
 +- 254
 +- 255
```
