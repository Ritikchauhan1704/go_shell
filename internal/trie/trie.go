package trie

// Node represents a node in the trie data structure
type Node struct {
	links [26]*Node
	flag  bool
}

// containsKey checks if the node contains a link for the given character
func (n *Node) containsKey(ch byte) bool {
	return n.links[ch-'a'] != nil
}

// put adds a new link from this node
func (n *Node) put(ch byte, node *Node) {
	n.links[ch-'a'] = node
}

// get retrieves the next node for the given character
func (n *Node) get(ch byte) *Node {
	return n.links[ch-'a']
}

// setEnd marks this node as the end of a word
func (n *Node) setEnd() {
	n.flag = true
}

// isEnd checks if this node is the end of a word
func (n *Node) isEnd() bool {
	return n.flag
}

// Trie is a data structure for efficient string prefix operations
type Trie struct {
	root *Node
}

// New creates a new Trie instance
func New() *Trie {
	return &Trie{root: &Node{}}
}

// Insert adds a word to the trie
func (t *Trie) Insert(word string) {
	node := t.root
	for i := 0; i < len(word); i++ {
		ch := word[i]
		if !node.containsKey(ch) {
			node.put(ch, &Node{})
		}
		node = node.get(ch)
	}
	node.setEnd()
}

// getNode finds the node that represents the given prefix
func (t *Trie) getNode(prefix string) *Node {
	node := t.root
	for i := 0; i < len(prefix); i++ {
		ch := prefix[i]
		if !node.containsKey(ch) {
			return nil
		}
		node = node.get(ch)
	}
	return node
}

// AutoComplete returns all words in the trie that start with the given prefix
func (t *Trie) AutoComplete(prefix string) []string {
	node := t.getNode(prefix)
	var result []string
	if node == nil {
		return result
	}
	t.dfs(node, prefix, &result)
	return result
}

// dfs performs depth-first search to find all words from a node
func (t *Trie) dfs(node *Node, path string, result *[]string) {
	if node.isEnd() {
		*result = append(*result, path)
	}
	for ch := byte('a'); ch <= byte('z'); ch++ {
		if next := node.get(ch); next != nil {
			t.dfs(next, path+string(ch), result)
		}
	}
}