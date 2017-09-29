// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"log"
)

func Broadcast_Scene_Skill(a *AccountData, data *EchoProto) {

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		// 技能目标
		// 技能xx
		retData := struct {
			Name    string `json:"name"` // 技能name
			Account string `json:"account"`
			IsDown  bool   `json:"isDown"`
			Ret     string `json:"ret"`
		}{
			data.Data["name"].(string),
			a.Account,
			data.Data["isDown"].(bool),
			"ok"}

		ret, _ := json.Marshal(struct {
			C    string `json:"c"`
			M    string `json:"m"`
			Data struct {
				Name    string `json:"name"` // 技能name
				Account string `json:"account"`
				IsDown  bool   `json:"isDown"`
				Ret     string `json:"ret"`
			} `json:"data"`
		}{
			C:    "Scene",
			M:    "Skill",
			Data: retData})

		err := gAccounts[k].C.WriteMessage(a.Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Scene_Skill(a *AccountData, mt int, data *EchoProto) bool {

	// 技能目标
	// 技能xx
	retData := struct {
		Name    string `json:"name"` // 技能name
		Account string `json:"account"`
		Ret     string `json:"ret"`
	}{
		data.Data["name"].(string),
		a.Account,
		"ok"}

	ret, _ := json.Marshal(struct {
		C    string `json:"c"`
		M    string `json:"m"`
		Data struct {
			Name    string `json:"name"` // 技能name
			Account string `json:"account"`
			Ret     string `json:"ret"`
		} `json:"data"`
	}{
		C:    data.C,
		M:    data.M,
		Data: retData})

	err := a.C.WriteMessage(mt, ret)
	if err != nil {
		log.Println("write:", err)
		return false
	}

	Broadcast_Scene_Skill(a, data)

	return true
}

func Broadcast_Scene_PlayerEnter(a *AccountData, data *EchoProto) {

	data.C = "Scene"
	data.M = "PlayerEnter"

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		pos_x, _ := data.Data["pos_x"].(int)
		pos_y, _ := data.Data["pos_y"].(int)

		retData := struct {
			Account string `json:"account"`
			Pwd     string `json:"pwd"`
			PosX    int    `json:"pos_x"`
			PosY    int    `json:"pos_y"`
			Ret     string `json:"ret"`
			Msg     string `json:"msg"`
		}{
			data.Data["account"].(string),
			data.Data["pwd"].(string),
			pos_x,
			pos_y,
			"ok",
			""}

		ret, _ := json.Marshal(struct {
			C    string `json:"c"`
			M    string `json:"m"`
			Data struct {
				Account string `json:"account"`
				Pwd     string `json:"pwd"`
				PosX    int    `json:"pos_x"`
				PosY    int    `json:"pos_y"`
				Ret     string `json:"ret"`
				Msg     string `json:"msg"`
			} `json:"data"`
		}{
			C:    data.C,
			M:    data.M,
			Data: retData})

		err := gAccounts[k].C.WriteMessage(gAccounts[k].Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Broadcast_Scene_PlayerLeave(a *AccountData) {

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		retData := struct {
			Account string `json:"account"`
			Ret     string `json:"ret"`
			Msg     string `json:"msg"`
		}{
			a.Account,
			"ok",
			""}

		ret, _ := json.Marshal(struct {
			C    string `json:"c"`
			M    string `json:"m"`
			Data struct {
				Account string `json:"account"`
				Ret     string `json:"ret"`
				Msg     string `json:"msg"`
			} `json:"data"`
		}{
			C:    "Scene",
			M:    "PlayerLeave",
			Data: retData})

		err := gAccounts[k].C.WriteMessage(gAccounts[k].Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Index_Login(a *AccountData, mt int, data *EchoProto) bool {

	pos_x, _ := data.Data["pos_x"].(int)
	pos_y, _ := data.Data["pos_y"].(int)

	retData := struct {
		Account string `json:"account"`
		Pwd     string `json:"pwd"`
		PosX    int    `json:"pos_x"`
		PosY    int    `json:"pos_y"`
		Ret     string `json:"ret"`
		Msg     string `json:"msg"`
	}{
		data.Data["account"].(string),
		data.Data["pwd"].(string),
		pos_x,
		pos_y,
		"ok",
		""}

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

	ret, _ := json.Marshal(struct {
		C    string `json:"c"`
		M    string `json:"m"`
		Data struct {
			Account string `json:"account"`
			Pwd     string `json:"pwd"`
			PosX    int    `json:"pos_x"`
			PosY    int    `json:"pos_y"`
			Ret     string `json:"ret"`
			Msg     string `json:"msg"`
		} `json:"data"`
	}{
		C:    data.C,
		M:    data.M,
		Data: retData})

	err := a.C.WriteMessage(mt, ret)
	if err != nil {
		log.Println("write:", err)
		return false
	}

	data.C = "Scene"
	data.M = "PlayerEnter"
	Broadcast_Scene_PlayerEnter(a, data)

	return true
}
