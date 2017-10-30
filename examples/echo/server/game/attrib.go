package game

const (
	Attr_hp      = 1 // 生命值
	Attr_mp      = 2 // 能量值
	Attr_attack  = 3 // 攻击力
	Attr_defence = 4 // 防御力
)

const (
	// MaxRoleAttr 最多角色属性数量
	MaxRoleAttr = uint32(4)
)

// RoleAttr 角色属性变动机制
type RoleAttr struct {
	Origins []int64 // 源
	Numbers []int64 // 变动值
	Scales  []int64 // 缩放比例
	Lasts   []int64 // 最终值
	Changes []bool  // 变动过
}

// 初始化
func (r *RoleAttr) Init() {
	r.Origins = make([]int64, MaxRoleAttr)
	r.Numbers = make([]int64, MaxRoleAttr)
	r.Scales = make([]int64, MaxRoleAttr)
	r.Lasts = make([]int64, MaxRoleAttr)
	r.Changes = make([]bool, MaxRoleAttr)
}

// CheckAttrIdx 检查属性ID
func (r *RoleAttr) CheckAttrIdx(attrID uint32) bool {
	return attrID > 0 && attrID < MaxRoleAttr
}

// Apply 对attr_id应用val值
func (r *RoleAttr) Apply(attrID uint32, val int64, number bool) {
	if r.CheckAttrIdx(attrID) {
		if number {
			r.Numbers[attrID-1] += val
		} else {
			r.Scales[attrID-1] += val
		}
		r.calcAttrLast(attrID)
	}
}

// Cancel 取消对attr_id应用的val值
func (r *RoleAttr) Cancel(attrID uint32, val int64, number bool) {
	if r.CheckAttrIdx(attrID) {
		if number {
			r.Numbers[attrID-1] -= val
		} else {
			r.Scales[attrID-1] -= val
		}
		r.calcAttrLast(attrID)
	}
}

// calcAttrLast 计算属性最终值
func (r *RoleAttr) calcAttrLast(attrID uint32) {
	attrID = attrID - 1
	// 计算 Lasts 最终值
	fixScale := r.Scales[attrID]
	if fixScale < 0 {
		fixScale = 0
	}

	r.Lasts[attrID] = (r.Origins[attrID] + r.Numbers[attrID]) * (100 + fixScale) / 100

	// 根据属性不同, 会有多种计算方式
	// 1. HP  : MaxHP 是上限, 下限是0, 增加百分比时用的是MaxHP作为基数, 只有 Apply 生效, Cancel 无效
	// 2. EXP : 无上限, 无下限, 只能增加值, Apply生效, Cancel无效
	// HP
}
