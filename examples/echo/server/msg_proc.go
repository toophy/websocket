// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket/examples/echo/server/game"
)

func Broadcast_Scene_Skill(a *game.AccountData, data *EchoProto) {

	ret, _ := json.Marshal(struct {
		C    string                 `json:"c"`
		M    string                 `json:"m"`
		Data map[string]interface{} `json:"data"`
		Ret  string                 `json:"ret"`
		Msg  string                 `json:"msg"`
	}{
		C:    data.C,
		M:    data.M,
		Data: data.Data,
		Ret:  "ok",
		Msg:  ""})

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		err := gAccounts[k].C.WriteMessage(a.Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Scene_Skill(a *game.AccountData, mt int, data *EchoProto) bool {

	data.Data["account"] = a.Account

	ret, _ := json.Marshal(struct {
		C    string                 `json:"c"`
		M    string                 `json:"m"`
		Data map[string]interface{} `json:"data"`
		Ret  string                 `json:"ret"`
		Msg  string                 `json:"msg"`
	}{
		C:    data.C,
		M:    data.M,
		Data: data.Data,
		Ret:  "ok",
		Msg:  ""})

	err := a.C.WriteMessage(mt, ret)
	if err != nil {
		log.Println("write:", err)
		return false
	}

	Broadcast_Scene_Skill(a, data)

	return true
}

func Scene_PlayerPoint(a *game.AccountData, mt int, data *EchoProto) bool {

	data.Data["account"] = a.Account

	ret, _ := json.Marshal(struct {
		C    string                 `json:"c"`
		M    string                 `json:"m"`
		Data map[string]interface{} `json:"data"`
		Ret  string                 `json:"ret"`
		Msg  string                 `json:"msg"`
	}{
		C:    data.C,
		M:    data.M,
		Data: data.Data,
		Ret:  "ok",
		Msg:  ""})

	err := a.C.WriteMessage(mt, ret)
	if err != nil {
		log.Println("write:", err)
		return false
	}

	a.PosX = data.Data["x"].(float64)
	a.PosY = data.Data["y"].(float64)

	Broadcast_Scene_Skill(a, data)

	return true
}

func Broadcast_Scene_PlayerEnter(a *game.AccountData, data *EchoProto) {

	ret, _ := json.Marshal(struct {
		C    string                 `json:"c"`
		M    string                 `json:"m"`
		Data map[string]interface{} `json:"data"`
		Ret  string                 `json:"ret"`
		Msg  string                 `json:"msg"`
	}{
		C:    "Scene",
		M:    "PlayerEnter",
		Data: data.Data,
		Ret:  "ok",
		Msg:  ""})

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		err := gAccounts[k].C.WriteMessage(gAccounts[k].Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Broadcast_Scene_PlayerEnter_ToMe(a *game.AccountData) {

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		data := map[string]interface{}{"account": gAccounts[k].Account, "pwd": gAccounts[k].Pwd, "pos_x": gAccounts[k].PosX, "pos_y": gAccounts[k].PosY}

		ret, _ := json.Marshal(struct {
			C    string                 `json:"c"`
			M    string                 `json:"m"`
			Data map[string]interface{} `json:"data"`
			Ret  string                 `json:"ret"`
			Msg  string                 `json:"msg"`
		}{
			C:    "Scene",
			M:    "PlayerEnter",
			Data: data,
			Ret:  "ok",
			Msg:  ""})

		err := a.C.WriteMessage(a.Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Broadcast_Scene_PlayerLeave(a *game.AccountData) {

	data := map[string]interface{}{"account": a.Account}

	ret, _ := json.Marshal(struct {
		C    string                 `json:"c"`
		M    string                 `json:"m"`
		Data map[string]interface{} `json:"data"`
		Ret  string                 `json:"ret"`
		Msg  string                 `json:"msg"`
	}{
		C:    "Scene",
		M:    "PlayerLeave",
		Data: data,
		Ret:  "ok",
		Msg:  ""})

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		err := gAccounts[k].C.WriteMessage(gAccounts[k].Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Index_Login(a *game.AccountData, mt int, data *EchoProto) bool {

	retData := struct {
		C    string                 `json:"c"`
		M    string                 `json:"m"`
		Data map[string]interface{} `json:"data"`
		Ret  string                 `json:"ret"`
		Msg  string                 `json:"msg"`
	}{
		C:    data.C,
		M:    data.M,
		Data: data.Data,
		Ret:  "ok",
		Msg:  ""}

	log.Printf("account:%s,pwd:%s", data.Data["account"], data.Data["pwd"])
	if _, ok := gAccounts[data.Data["account"].(string)]; ok {
		retData.Ret = "fail"
		retData.Msg = "重复登录"
	} else {
		a.Account = data.Data["account"].(string)
		a.Pwd = data.Data["pwd"].(string)
		a.Mt = mt
		gAccounts[data.Data["account"].(string)] = a
	}

	ret, _ := json.Marshal(retData)

	err := a.C.WriteMessage(mt, ret)
	if err != nil {
		log.Println("write:", err)
		return false
	}

	data.C = "Scene"
	data.M = "PlayerEnter"
	Broadcast_Scene_PlayerEnter(a, data)
	Broadcast_Scene_PlayerEnter_ToMe(a)

	return true
}
