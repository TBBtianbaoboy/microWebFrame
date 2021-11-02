package kua

import "strings"

//路由结点
type node struct{
	pattern string //待匹配路由
	part string //路由中的一部分
	chlidren []*node //当前结点的子结点
	isWild bool //是否精确匹配

}

//匹配当前的路由结点n是否有第一个part结点
func (n *node) matchChild(part string) *node {
	for _,child := range n.chlidren {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil 
}

//匹配当前路由结点n所有的part结点
func (n *node) matchChildren (part string) []*node {
	nodes := make([]*node,0)
	for _,child := range n.chlidren {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入路由结点
func (n *node) insert(pattern string,parts []string,height int){
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part:part,isWild:part[0]== '*' || part[0]== ':'}
		n.chlidren = append(n.chlidren, child)
	}
	child.insert(pattern,parts,height+1)
}

//递归搜寻路由 
func (n *node) search(parts []string,height int) *node {
	if len(parts)==height || strings.HasPrefix(n.part,"*"){
		if n.pattern ==""{
			return nil
		}
		return n
	}

	part := parts[height]
	chlidren := n.matchChildren(part)
	
	for _,child := range chlidren {
		result := child.search(parts,height+1)
		if result != nil {
			return result
		}
	}
	return nil 
}































