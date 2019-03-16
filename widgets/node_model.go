package widgets

import (
	"encoding/json"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"

	"schannel-qt5/geoip"
	"schannel-qt5/parser"
)

var CountryFlags = make([]map[string]string, 0)

func init() {
	// 获取countryCode对应的flag emoji
	flagData := core.NewQFile2(":/flags/data.json")
	flagData.Open(core.QIODevice__ReadOnly)
	err := json.Unmarshal([]byte(flagData.ReadAll().Data()), &CountryFlags)
	if err != nil {
		CountryFlags = nil
	}
	flagData.Close()
}

// NodeTreeItem 将节点保存成地区/名字的树形结构
type NodeTreeItem struct {
	core.QObject
	// 节点或地区的名字
	name string
	node *parser.SSRNode

	parent   *NodeTreeItem
	children []*NodeTreeItem
}

// NewNodeTreeItem2 用name创建NodeTreeItem
func NewNodeTreeItem2(name string) *NodeTreeItem {
	item := NewNodeTreeItem(nil)
	item.name = name
	item.children = make([]*NodeTreeItem, 0)
	// 清理资源，goqt无法自动清理
	item.ConnectDestroyNodeTreeItem(item.destroyTreeItem)

	return item
}

// 释放所有children
func (n *NodeTreeItem) destroyTreeItem() {
	for _, child := range n.children {
		child.DestroyNodeTreeItem()
	}

	n.DestroyNodeTreeItemDefault()
}

func (n *NodeTreeItem) ChildCount() int {
	return len(n.children)
}

func (n *NodeTreeItem) ColumnCount() int {
	return 1
}

// 返回自己的直接子节点
func (n *NodeTreeItem) Child(row int) *NodeTreeItem {
	return n.children[row]
}

func (n *NodeTreeItem) ParentItem() *NodeTreeItem {
	return n.parent
}

// 添加节点至当前节点的children
func (n *NodeTreeItem) AppendChild(child *NodeTreeItem) {
	child.parent = n
	n.children = append(n.children, child)
}

func (n *NodeTreeItem) Data() string {
	return n.name
}

// Row 返回当前节点在父节点的children中的索引
func (n *NodeTreeItem) Row() int {
	if n.parent == nil || len(n.parent.children) < 1 {
		return 0
	}

	for i, child := range n.parent.children {
		if child == n {
			return i
		}
	}

	return 0
}

// 根据name查找当前tree的子节点（不进行递归查找）
func (n *NodeTreeItem) FindChild(name string) *NodeTreeItem {
	for _, item := range n.children {
		if item.name == name {
			return item
		}
	}

	return nil
}

func (n *NodeTreeItem) SetNode(node *parser.SSRNode) {
	n.node = node
}

func (n *NodeTreeItem) Node() *parser.SSRNode {
	return n.node
}

// LatestChild 返回末端节点
// 如果已经是末端节点则返回自己
func (n *NodeTreeItem) LatestChild(row int) *NodeTreeItem {
	if len(n.children) == 0 {
		return n
	}

	child := n.Child(row)
	for len(child.children) != 0 {
		child = child.Child(row)
	}

	return child
}

// NodeTreeModel 保存按地理信息分类的nodes
type NodeTreeModel struct {
	core.QAbstractItemModel

	// 根节点
	rootItem *NodeTreeItem
}

// NewNodeTreeModel2 用nodes初始化model
func NewNodeTreeModel2(nodes []*parser.SSRNode) *NodeTreeModel {
	model := NewNodeTreeModel(nil)
	model.rootItem = NewNodeTreeItem2("")
	for _, node := range nodes {
		// 通过geo信息逐层查找，找到或新建对应的最底层地理区域节点
		geo := strings.Split(getGeoName(node.IP), "-")
		baseItem := model.rootItem
		for _, area := range geo {
			areaItem := baseItem.FindChild(area)
			if areaItem == nil {
				areaItem = NewNodeTreeItem2(area)
				baseItem.AppendChild(areaItem)
			}

			baseItem = areaItem
		}

		nodeItem := NewNodeTreeItem2(node.NameNumber())
		nodeItem.SetNode(node)
		baseItem.AppendChild(nodeItem)
	}

	model.ConnectIndex(model.index)
	model.ConnectParent(model.parent)
	model.ConnectHeaderData(model.headerData)
	model.ConnectRowCount(model.rowCount)
	model.ConnectColumnCount(model.columnCount)
	model.ConnectData(model.data)

	return model
}

// 返回子节点的index
func (n *NodeTreeModel) index(row int, column int, parent *core.QModelIndex) *core.QModelIndex {
	if !n.HasIndex(row, column, parent) {
		return core.NewQModelIndex()
	}

	var parentItem *NodeTreeItem
	if !parent.IsValid() {
		parentItem = n.rootItem
	} else {
		parentItem = NewNodeTreeItemFromPointer(parent.InternalPointer())
	}

	childItem := parentItem.Child(row)
	if childItem != nil {
		return n.CreateIndex(row, column, childItem.Pointer())
	}

	return core.NewQModelIndex()
}

// 标题信息
func (n *NodeTreeModel) headerData(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if role == int(core.Qt__DisplayRole) && orientation == core.Qt__Horizontal {
		return core.NewQVariant17("节点")
	}

	return n.HeaderDataDefault(section, orientation, role)
}

func (n *NodeTreeModel) parent(index *core.QModelIndex) *core.QModelIndex {
	if !index.IsValid() {
		return core.NewQModelIndex()
	}

	item := NewNodeTreeItemFromPointer(index.InternalPointer())
	parentItem := item.ParentItem()

	if parentItem.Pointer() == n.rootItem.Pointer() {
		return core.NewQModelIndex()
	}

	return n.CreateIndex(parentItem.Row(), 0, parentItem.Pointer())
}

func (n *NodeTreeModel) columnCount(_ *core.QModelIndex) int {
	return n.rootItem.ColumnCount()
}

// 每个节点的子节点数目
func (n *NodeTreeModel) rowCount(parent *core.QModelIndex) int {
	if !parent.IsValid() {
		return n.rootItem.ChildCount()
	}

	return NewNodeTreeItemFromPointer(parent.InternalPointer()).ChildCount()
}

// view取得节点name
func (n *NodeTreeModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	item := NewNodeTreeItemFromPointer(index.InternalPointer())
	// 处理顶层节点
	if item.parent.Data() == "" {
		switch role {
		case int(core.Qt__FontRole):
			font := gui.NewQFont2("noto color emoji", -1, -1, false)
			return core.NewQVariant3(int(core.QVariant__Font), font.Pointer())
		case int(core.Qt__DisplayRole):
			if CountryFlags == nil {
				break
			}

			// 在顶层节点显示国旗
			if country, err := geoip.GetCountryISOCode(item.LatestChild(0).node.IP); err == nil {
				for i := range CountryFlags {
					if country == CountryFlags[i]["code"] {
						return core.NewQVariant17(CountryFlags[i]["emoji"] + " " + item.Data())
					}
				}
			}
		}
	}

	switch role {
	case int(core.Qt__DisplayRole):
		return core.NewQVariant17(item.Data())
	}

	return core.NewQVariant()
}

// FindNodeIndex 返回node所在的item的index
func (n *NodeTreeModel) FindNodeIndex(node *parser.SSRNode) *core.QModelIndex {
	names := strings.Split(getGeoName(node.IP), "-")
	names = append(names, node.NameNumber())

	parentItem := n.rootItem
	// 递归查找，找不到就返回无效index
	for _, name := range names {
		childItem := parentItem.FindChild(name)
		if childItem == nil {
			return core.NewQModelIndex()
		}

		if childItem.name == name {
			parentItem = childItem
		}
	}

	// 处理无children的情况
	if parentItem.Pointer() == n.rootItem.Pointer() {
		return core.NewQModelIndex()
	}

	return n.CreateIndex(parentItem.Row(), 0, parentItem.Pointer())
}
