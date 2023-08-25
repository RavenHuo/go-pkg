package collections

// TrieNode 前缀树节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
			isEnd:    false,
		},
	}
}

func (t *Trie) Insert(word string) {
	current := t.root
	for _, ch := range word {
		if current.children[ch] == nil {
			current.children[ch] = &TrieNode{
				children: make(map[rune]*TrieNode),
				isEnd:    false,
			}
		}
		current = current.children[ch]
	}
	current.isEnd = true
}

func (t *Trie) Search(word string) bool {
	current := t.root
	for _, ch := range word {
		if current.children[ch] == nil {
			return false
		}
		current = current.children[ch]
	}
	return current.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
	current := t.root
	for _, ch := range prefix {
		if current.children[ch] == nil {
			return false
		}
		current = current.children[ch]
	}
	return true
}
