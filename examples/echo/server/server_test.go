package main

import "testing"

func TestServerID(t *testing.T) {

	rid := MakeRid(1, 2, 3, 78)
	t.Errorf("rid=%d", rid)

	t.Errorf("pkey=%d,aid=%d,skey=%d,rsn=%d", GetPkeySn(rid), GetAidSn(rid), GetSkeySn(rid), GetRoleSn(rid))

	islandId := MakeSceneObjID(99, 1, 100)
	if IsObjType(islandId, 1) {
		t.Errorf("scene=%d,islandId=%d", GetSceneSN(islandId), GetSceneObjSN(islandId))
	} else {
		t.Errorf("failed islandID = %d", islandId)
	}

}
