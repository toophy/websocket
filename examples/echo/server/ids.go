package main

/*
ID 编号规则
最大 1000 00000 00000000
	 平台 SKEY  ACCID
最多支持1000平台, 10万区服, 1亿帐号, 100角色
0平台以内的所有数值供系统使用
也就是第一个帐号ID是 1平台1区的1帐号

平台 帐号      区    角色
001  00000001 00001 01
999  99999999 99999 99
*/

// 帐号ID
func MakeAid(pkey int64, acc int64) int64 {
	if pkey < 0 || pkey > 999 {
		return 0
	}
	if acc < 0 || acc > 99999999 {
		return 0
	}

	return pkey*1000000000000000 + acc*10000000
}

// 从ID获取平台序号
func GetPkeySn(id int64) int64 {
	return id / 1000000000000000
}

// 从ID获取帐号序号
func GetAidSn(id int64) int64 {
	return id/10000000 - (id/1000000000000000)*100000000
}

// 从ID获取区服序号
func GetSkeySn(id int64) int64 {
	return id/100 - (id/10000000)*100000
}

// 从ID获取角色序号
func GetRoleSn(id int64) int64 {
	return id - (id/100)*100
}

// 帐号的角色ID
func MakeRid(pkey int64, acc int64, skey int64, role int64) int64 {
	if pkey < 0 || pkey > 999 {
		return 0
	}
	if acc < 0 || acc > 99999999 {
		return 0
	}
	if skey < 0 || skey > 99999 {
		return 0
	}
	if role < 0 || role > 99 {
		return 0
	}

	return pkey*1000000000000000 + acc*10000000 + skey*100 + role
}

// 房间ID
func MakeRoomID(pkey int64, skey int64, room int64) int64 {
	if pkey < 0 || pkey > 999 {
		return 0
	}
	if room < 0 || room > 99999999 {
		return 0
	}
	if skey < 0 || skey > 99999 {
		return 0
	}
	return pkey*1000000000000000 + room*10000000 + skey*100
}

// 场景ID
func MakeSceneID(pkey int64, skey int64, room int64, scene int64) int64 {
	if pkey < 0 || pkey > 999 {
		return 0
	}
	if room < 0 || room > 99999999 {
		return 0
	}
	if skey < 0 || skey > 99999 {
		return 0
	}
	if scene < 0 || scene > 99 {
		return 0
	}
	return pkey*1000000000000000 + room*10000000 + skey*100 + scene
}

// 场景对象ID
// 场景   类型  浮岛编号
// 01    01    00001
func MakeSceneObjID(scene int64, objType int64, id int64) int64 {
	if scene < 0 || scene > 99 {
		return 0
	}
	if objType < 0 || objType > 99 {
		return 0
	}
	if id < 0 || id > 99999 {
		return 0
	}
	return scene*10000000 + objType*100000 + id
}

// 场景对象序号
func GetSceneObjSN(id int64) int64 {
	return id - (id/100000)*100000
}

// 场景序号
func GetSceneSN(id int64) int64 {
	return id/10000000 - (id/1000000000)*100
}

// 场景ID类型
func GetSceneType(id int64) int64 {
	return id/100000 - (id/10000000)*100
}

// 是指定场景对象类型
func IsObjType(id int64, objType int64) bool {
	if id > 0 && id < 1000000000 {
		scene := GetSceneSN(id)
		return scene > 0 && scene <= 99 && GetSceneType(id) == objType
	}
	return false
}
