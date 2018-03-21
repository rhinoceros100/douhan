package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"time"
	"douhan/playing"
	"douhan/util"
	"douhan/log"
	"douhan/card"
)

func help() {
	log.Debug("-----------------help---------------------")
	log.Debug("h")
	log.Debug("exit")
	log.Debug("mycards")
	log.Debug(playing.OperateEnterRoom, int(playing.OperateEnterRoom))
	log.Debug(playing.OperateReadyRoom, int(playing.OperateReadyRoom))
	log.Debug(playing.OperateLeaveRoom, int(playing.OperateLeaveRoom))
	log.Debug(playing.OperateDrop, int(playing.OperateDrop), "1(type) 7(cardno)")
	log.Debug(playing.OperatePass, int(playing.OperatePass))
	log.Debug("-----------------help---------------------")
}

type PlayerObserver struct {}
func (ob *PlayerObserver) OnMsg(player *playing.Player, msg *playing.Message) {
	log_time := time.Now().Unix()
	switch msg.Type {
	case playing.MsgEnterRoom:
		if enter_data, ok := msg.Data.(*playing.EnterRoomMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgEnterRoom, EnterPlayer:", enter_data.EnterPlayer)
		}
	case playing.MsgReadyRoom:
		if enter_data, ok := msg.Data.(*playing.ReadyRoomMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgReadyRoom, ReadyPlayer:", enter_data.ReadyPlayer)
		}
	case playing.MsgLeaveRoom:
		if enter_data, ok := msg.Data.(*playing.LeaveRoomMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgLeaveRoom, LeavePlayer:", enter_data.LeavePlayer)
		}
	case playing.MsgGameEnd:
		if _, ok := msg.Data.(*playing.GameEndMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgGameEnd")
		}
	case playing.MsgRoomClosed:
		if close_data, ok := msg.Data.(*playing.RoomClosedMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgRoomClosed")
			for _, data := range close_data.Summaries	{
				log.Debug(data.P, "Win:", data.WinNum, "SpringNum:", data.SpringNum,
					"TotalCoin:", data.TotalCoin, "IsWin:", data.IsWinner, "IsMostWin:", data.IsMostWinner, "IsMostLos:", data.IsMostLoser)
			}
		}
	case playing.MsgGetInitCards:
		if init_data, ok := msg.Data.(*playing.GetInitCardsMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgGetInitCards, PlayingCards:", init_data.PlayingCards)
		}
	case playing.MsgStartPlay:
		if sp_data, ok := msg.Data.(*playing.StartPlayMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgStartPlay, master:", sp_data.Master)
		}
	case playing.MsgSwitchOperator:
		if sp_data, ok := msg.Data.(*playing.SwitchOperatorMsgData); ok {
			log.Debug(log_time, player, "******OnMsg MsgSwitchOperator", msg.Owner, "NeedDropCard", sp_data.NeedDropCard)
		}
	case playing.MsgDrop:
		if drop_data, ok := msg.Data.(*playing.DropMsgData); ok {
			log.Debug(log_time, player, "MsgDrop", msg.Owner, "CardsType", drop_data.CardsType, "Weight", drop_data.Weight, "cards", drop_data.WhatGroup)
		}
	case playing.MsgDispatchCard:
		if _, ok := msg.Data.(*playing.DispatchCardMsgData); ok {
			log.Debug(log_time, player, "MsgDispatchCard", msg.Owner)
		}
	case playing.MsgPass:
		if _, ok := msg.Data.(*playing.PassMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgPass", msg.Owner)
		}
	case playing.MsgSummary:
		if summary_data, ok := msg.Data.(*playing.SummaryMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgSummary, summary_data:", summary_data.InfoType, "Multiple:", summary_data.Multiple)
			for _, score_data := range summary_data.Scores	{
				log.Debug(score_data.P, "Rank:", score_data.Rank, "IsWin:", score_data.IsWin, "LeftCardNum:", score_data.LeftCardNum,
					"IsSpring:", score_data.IsSpring, "Coin:", score_data.Coin, "TotalCoin:", score_data.TotalCoin)
			}
		}
	}
}

func main() {
	running := true

	//init room
	conf := playing.NewRoomConfig()
	conf.Init()
	room := playing.NewRoom(util.UniqueId(), conf)
	room.Start()

	robots := []*playing.Player{
		playing.NewPlayer(0),
		playing.NewPlayer(1),
		playing.NewPlayer(2),

		playing.NewPlayer(3),
	}

	for _, robot := range robots {
		robot.OperateEnterRoom(room)
		robot.AddObserver(&PlayerObserver{})
	}

	curPlayer := playing.NewPlayer(4)
	curPlayer.AddObserver(&PlayerObserver{})

	go func() {
		time.Sleep(time.Second * 1)
		robots[0].OperateDoReady()
		time.Sleep(time.Second * 1)
		robots[1].OperateDoReady()
		time.Sleep(time.Second * 1)
		robots[2].OperateDoReady()

		time.Sleep(time.Second * 1)
		robots[3].OperateDoReady()
	}()

	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		cmd := string(data)
		if cmd == "h" {
			help()
		} else if cmd == "exit" {
			return
		} else if cmd == "mycards" {
			log.Debug(curPlayer.GetPlayingCards())
		}
		splits := strings.Split(cmd, " ")
		c, _ := strconv.Atoi(splits[0])
		switch playing.OperateType(c) {
		case playing.OperateEnterRoom:
			curPlayer.OperateEnterRoom(room)
		case playing.OperateReadyRoom:
			curPlayer.OperateDoReady()
		case playing.OperateLeaveRoom:
			curPlayer.OperateLeaveRoom()
		case playing.OperateDrop:
			if len(splits) > 2 {
				card1 := &card.Card{}
				card1.CardType, _ = strconv.Atoi(splits[1])
				card1.CardNo, _ = strconv.Atoi(splits[2])
				card1.MakeIDWeight(1)
				cards := make([]*card.Card, 0)
				cards = append(cards, card1)
				curPlayer.OperateDropCard(cards)
			}else {
				help()
			}
		case playing.OperatePass:
			curPlayer.OperatePass()
		}
	}
}
