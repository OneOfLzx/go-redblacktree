## 简介
用GO语言实现的红黑树，可做工具库直接调用
## 类和接口介绍
`RedBlackTreeNodeValEntry:` 封装节点存的数据的接口  
Equal(): 判断两个数据是否相等  
Smaller()：判断数据是否小于  
`RedBlackTreeNode:` 红黑树的节点  
IsValidNode(): 返回当前节点是否有效  
Value(): 返回当前节点里的数据  
PrevNode(): 根据中序遍历，返回前一个节点的指针  
NextNode(): 根据中序遍历，返回后一个节点的指针  
`RedBlackTree:` 表示红黑树的数据结构  
FindNode(): 查找是否存在包含某个值的节点。如果存在返回该节点的指针；反之返回nil  
AddNode(): 添加一个节点，并返回添加的节点的指针。如果改值已经包含在这棵树中，则返回已存在的节点指针  
RemoveNodeByVal(): 通过数据来删除节点。如果没有节点包含该数据则什么也不做  
RemoveNode(): 删除一个指定的节点