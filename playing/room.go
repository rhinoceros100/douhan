package playing

import (
	"douhan/card"
	"douhan/log"
	"douhan/util"
	"time"
	"fmt"
)

type RoomStatusType int
const (
	RoomStatusWaitAllPlayerEnter	RoomStatusType = iota	// 等待玩家进入房间
	RoomStatusWaitAllPlayerReady				// 等待玩家准备
	RoomStatusGameStart					// 发牌开始打
	RoomStatusPlayGame					// 正在进行游戏，结束后会进入RoomStatusEndPlayGame
	RoomStatusEndPlayGame					// 游戏结束后会回到等待游戏开始状态，或者进入结束房间状态
	RoomStatusRoomEnd					// 房间结束状态
)

func (status RoomStatusType) String() string {
	switch status {
	case RoomStatusWaitAllPlayerEnter :
		return "RoomStatusWaitAllPlayerEnter"
	case RoomStatusWaitAllPlayerReady:
		return "RoomStatusWaitAllPlayerReady"
	case RoomStatusGameStart:
		return "RoomStatusGameStart"
	case RoomStatusPlayGame:
		return "RoomStatusPlayGame"
	case RoomStatusEndPlayGame:
		return "RoomStatusEndPlayGame"
	case RoomStatusRoomEnd:
		return "RoomStatusRoomEnd"
	}
	return "unknow RoomStatus"
}

type RoomObserver interface {
	OnRoomClosed(room *Room)
}

type Room struct {
	id			uint64					//房间id
	config 			*RoomConfig				//房间配置
	players 		[]*Player				//当前房间的玩家列表

	observers		[]RoomObserver				//房间观察者，需要实现OnRoomClose，房间close的时候会通知它
	roomStatus		RoomStatusType				//房间当前的状态
	playedGameCnt		int					//已经玩了的游戏的次数

	//begin playingGameData, reset when start playing game
	cardPool		*card.Pool				//洗牌池
	creatorUid		uint64					//创建房间玩家的uid
	opMaster		*Player					//当前出过牌的玩家
	waitOperator		*Player					//等待出牌的玩家
	masterPlayer 		*Player					//庄
	lastMaster		*Player					//上一把的庄家

	cardsType 		int					//上一次出的牌型
	planeNum 		int					//飞机数量
	weight	 		int					//权重
	bombNum	 		int					//炸弹数量
	//end playingGameData, reset when start playing game

	roomOperateCh	chan *Operate
	dropCardCh	[]chan *Operate				//出牌
	passCh		[]chan *Operate				//过牌
	roomReadyCh	[]chan *Operate

	isFirstRound 		bool
	isHuangju 		bool
	stop 			bool
}

func NewRoom(id uint64, config *RoomConfig) *Room {
	room := &Room{
		id:			id,
		config:			config,
		players:		make([]*Player, 0),
		cardPool:		card.NewPool(),
		observers:		make([]RoomObserver, 0),
		roomStatus:		RoomStatusWaitAllPlayerEnter,
		creatorUid:		0,
		playedGameCnt:		0,
		cardsType:		0,
		planeNum:		0,
		weight:			0,
		bombNum:		0,
		opMaster:		nil,
		waitOperator:		nil,
		masterPlayer:		nil,
		lastMaster:		nil,
		isFirstRound:		true,
		isHuangju:		false,

		roomOperateCh: make(chan *Operate, 1024),
		dropCardCh: make([]chan *Operate, config.NeedPlayerNum),
		passCh: make([]chan *Operate, config.NeedPlayerNum),
		roomReadyCh: make([]chan *Operate, config.NeedPlayerNum),
	}
	for idx := 0; idx < int(config.NeedPlayerNum); idx ++ {
		room.dropCardCh[idx] = make(chan *Operate, 1)
		room.passCh[idx] = make(chan *Operate, 1)
		room.roomReadyCh[idx] = make(chan *Operate, 1)
	}
	return room
}

func (room *Room) GetId() uint64 {
	return room.id
}

func (room *Room) PlayerOperate(op *Operate) {
	pos := op.Operator.position
	log.Debug(time.Now().Unix(), room, op.Operator, "PlayerOperate", op.Op, " pos:", pos)

	switch op.Op {
	case OperateEnterRoom, OperateLeaveRoom:
		room.roomOperateCh <- op
	case OperateReadyRoom:
		if room.roomStatus == RoomStatusWaitAllPlayerEnter {
			room.roomOperateCh <- op
		}else {
			room.roomReadyCh[pos] <- op
		}
	case OperateDrop:
		room.dropCardCh[pos] <- op
	case OperatePass:
		room.passCh[pos] <- op
	}
}

func (room *Room) addObserver(observer RoomObserver) {
	room.observers = append(room.observers, observer)
}

func (room *Room) Start() {
	go func() {
		start_time := time.Now().Unix()
		for  {
			if !room.stop {
				room.checkStatus()
				time.Sleep(time.Microsecond * 10)
			}else{
				break
			}
		}
		end_time := time.Now().Unix()
		log.Debug(end_time - start_time, "over^^")
	}()
}

func (room *Room) checkStatus() {
	switch room.roomStatus {
	case RoomStatusWaitAllPlayerEnter:
		room.waitAllPlayerEnter()
	case RoomStatusWaitAllPlayerReady:
		room.waitAllPlayerReady()
	case RoomStatusGameStart:
		room.gameStart()
	case RoomStatusPlayGame:
		room.playGame()
	case RoomStatusEndPlayGame:
		room.endPlayGame()
	case RoomStatusRoomEnd:
		room.close()
	}
}

func (room *Room) GetPlayerNum() int32 {
	return int32(len(room.players))
}

func (room *Room) isRoomEnd() bool {
	return room.playedGameCnt >= room.config.MaxPlayGameCnt
}

func (room *Room) GetCardsType() int {
	return room.cardsType
}

func (room *Room) SetCardsType(cardsType int) {
	room.cardsType = cardsType
}

func (room *Room) GetPlaneNum() int {
	return room.planeNum
}

func (room *Room) SetPlaneNum(planeNum int) {
	room.planeNum = planeNum
}

func (room *Room) GetWeight() int {
	return room.weight
}

func (room *Room) SetWeight(weight int) {
	room.weight = weight
}

func (room *Room) GetBombNum() int {
	return room.bombNum
}

func (room *Room) IncBombNum() int {
	room.bombNum++
	return room.bombNum
}

func (room *Room) close() {
	log.Debug(time.Now().Unix(), room, "Room.close")
	room.stop = true
	for _, observer := range room.observers {
		observer.OnRoomClosed(room)
	}

	msg := room.totalSummary()
	for _, player := range room.players {
		player.OnRoomClosed(msg)
	}
}

func (room *Room) isEnterPlayerEnough() bool {
	length := room.GetPlayerNum()
	log.Debug(time.Now().Unix(), room, "Room.isEnterPlayerEnough, player num :", length, ", need :", room.config.NeedPlayerNum)
	return length >= room.config.NeedPlayerNum
}

func (room *Room) switchStatus(status RoomStatusType) {
	log.Debug(time.Now().Unix(), room, "room status switch,", room.roomStatus, " =>", status)
	room.roomStatus = status
	log.Debug("---------------------------------------")
}

func (room *Room) canCover(cardsType, planeNum, weight int) (canCover bool) {
	canCover = false
	if room.GetCardsType() == card.CardsType_NO {
		return cardsType != card.CardsType_NO
	}
	//已经出的牌型非炸弹牌型
	if room.GetCardsType() < 20{
		if cardsType >= 20 {
			return true
		}
		//普通牌型打普通牌型必须为同一牌型，并且飞机数量必须相同
		if cardsType != room.GetCardsType(){
			return false
		}
		if cardsType == card.CardsType_STRAIGHT {
			if planeNum != room.GetPlaneNum() {
				return false
			}
			return weight == room.GetWeight() + 1
		}
		//非炸弹中2最大
		if weight == 15 {
			return false
		}
		return weight == 15 || weight == room.GetWeight() + 1
	}

	//更大的炸弹可以管住
	if cardsType > room.GetCardsType(){
		return true
	}
	return weight > room.GetWeight()
}

//等待游戏开局
func (room *Room) waitAllPlayerEnter() {
	log.Debug(time.Now().Unix(), room, "waitAllPlayerEnter......")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerEnterRoomTimeout) * time.Second
	for {
		timer := timeout - breakTimerTime
		select {
		case <-time.After(timer):
			log.Debug(time.Now().Unix(), room, "waitAllPlayerEnter timeout", timeout)
			room.switchStatus(RoomStatusRoomEnd) //超时发现没有足够的玩家都进入房间了，则结束
			return
		case op := <-room.roomOperateCh:
			if op.Op == OperateEnterRoom || op.Op == OperateLeaveRoom || op.Op == OperateReadyRoom {
				log.Debug(time.Now().Unix(), room, "waitAllPlayerEnter catch operate:", op)
				room.dealPlayerOperate(op)
				if room.isAllPlayerEnter() {
					room.switchStatus(RoomStatusWaitAllPlayerReady)
					return
				}
			}
		}
	}
}

func (room *Room) isAllPlayerEnter() bool {
	length := len(room.players)
	log.Debug(room, "Room.isAllPlayerEnter, num:", length, "need:", room.config.NeedPlayerNum)
	if length < int(room.config.NeedPlayerNum) {
		return false
	}
	for _, player := range room.players{
		if !player.GetIsReady() {
			return false
		}
	}

	return true
}

func (room *Room) waitDropCard(player *Player, mustDrop bool, canDrop bool) bool{
	wait_time := room.config.WaitDropSec
	if !canDrop{
		wait_time = time.Duration(3)
	}
	for{
		select {
		case <- time.After(time.Second * wait_time):
			random := util.RandomN(4)
			log.Debug(time.Now().Unix(), player, "waitDropCard do PlayerOperate, random:", random)

			if mustDrop {
				tailCards := player.GetTailCard()
				dropCards := card.CopyCards(tailCards)
				drop_cards := card.CreateNewCards(dropCards)

				cards_num := player.playingCards.CardsInHand.Len()
				is_last_cards := false
				if cards_num == len(dropCards) && room.opMaster == player{
					is_last_cards = true
				}

				to_cover_weight := 0
				if room.GetCardsType() == card.CardsType_STRAIGHT {
					to_cover_weight = room.GetWeight() + 1
				}
				data := &OperateDropData{whatGroup:dropCards}
				data.cardsType, data.planeNum, data.weight = card.GetCardsType(drop_cards, is_last_cards, to_cover_weight)
				can_cover := room.canCover(data.cardsType, data.planeNum, data.weight)
				log.Debug("******can_cover:", can_cover)

				op := NewOperateDrop(player, data)
				room.dealPlayerOperate(op)
				return true
			}else{
				//TODO 超时出牌只检查单牌和对子，顺子和炸弹先不做
				//查找手牌是否可以压住前面的牌
				can_cover := false
				cards_type := room.GetCardsType()
				cards_weight := room.GetWeight()
				num2, num_joker := player.Get2JokerNum()
				num_card := player.GetCardNumByWeight(cards_weight + 1)

				to_cover_cards := make([]*card.Card, 0)
				cards_data := player.GetPlayingCards().CardsInHand.GetData()
				switch cards_type {
				case card.CardsType_SINGLE:
					if cards_weight < 15 {
						check_weight := 0
						have_append := false
						if num_card >= 1 {
							check_weight = cards_weight + 1
						}else if num2 >= 1 {
							check_weight = 15
						}
						for _, hand_card := range cards_data {
							if hand_card.Weight == check_weight && !have_append{
								to_cover_cards = append(to_cover_cards, hand_card)
								can_cover = true
								have_append = true
							}
						}
					}
				case card.CardsType_PAIR:
					if cards_weight < 15 {
						check_weight := 0
						need_check_joker := false
						if num_card >= 2 || (num_card == 1 && num_joker >= 1) {
							check_weight = cards_weight + 1
							if num_card == 1{
								need_check_joker = true
							}
						}else if num2 >= 2 || (num2 == 1 && num_joker >= 1) {
							check_weight = 15
							if num2 == 1{
								need_check_joker = true
							}
						}
						for _, hand_card := range cards_data {
							if hand_card.Weight == check_weight {
								to_cover_cards = append(to_cover_cards, hand_card)
								can_cover = true
							}
							if need_check_joker && hand_card.Weight > 15{
								to_cover_cards = append(to_cover_cards, hand_card)
							}
						}
					}
				}

				if can_cover {
					drop_cards := card.CreateNewCards(to_cover_cards)
					data := &OperateDropData{whatGroup:to_cover_cards}
					data.cardsType, data.planeNum, data.weight = card.GetCardsType(drop_cards, false, 0)
					can_cover2 := room.canCover(data.cardsType, data.planeNum, data.weight)
					log.Debug("******can_cover2:", can_cover2)

					op := NewOperateDrop(player, data)
					room.dealPlayerOperate(op)
					return true
				}else{
					data := &OperatePassData{}
					op := NewOperatePass(player, data)
					room.dealPlayerOperate(op)
					return false
				}
			}
		case op := <-room.dropCardCh[player.position]:
			log.Debug(time.Now().Unix(), player, "Player.waitDropCard:", op.Data)
			room.dealPlayerOperate(op)
			return true
		case op := <-room.passCh[player.position] :
			log.Debug(room, "Room.waitDropCard operate :", op)
			room.dealPlayerOperate(op)
			return false
		}
	}

	log.Debug(time.Now().Unix(), player, "Player.waitBet fasle")
	return false
}

func (room *Room) getRandomDropNum(player *Player) int{
	num := room.config.RandomDropNum
	hand_cards_len := player.GetPlayingCards().CardsInHand.Len()
	if hand_cards_len < num{
		num = hand_cards_len
	}
	return num
}

func (room *Room) waitInitPlayerReady(player *Player) {
	time.Sleep(time.Second * room.config.WaitReadySec)
	if (room.roomStatus == RoomStatusWaitAllPlayerEnter || room.roomStatus == RoomStatusWaitAllPlayerReady) && !player.GetIsReady() {
		data := &OperateReadyRoomData{}
		op := NewOperateReadyRoom(player, data)
		log.Debug(player, "waitInitPlayerReady do PlayerOperate")
		room.PlayerOperate(op)
	}
}

func (room *Room) waitPlayerReady(player *Player) bool {
	log.Debug(time.Now().Unix(), player, "waitPlayerReady")
	for{
		select {
		case <- time.After(time.Second * room.config.WaitReadySec):
			data := &OperateReadyRoomData{}
			op := NewOperateReadyRoom(player, data)
			log.Debug("******")
			log.Debug(time.Now().Unix(), player, "waitPlayerReady do PlayerOperate")
			room.PlayerOperate(op)
			continue
		case op := <-room.roomReadyCh[player.GetPosition()]:
			log.Debug(time.Now().Unix(), player, "Player.waitPlayerReady")
			room.dealPlayerOperate(op)
			return true
		}
	}

	log.Debug(time.Now().Unix(), player, "Player.waitPlayerReady fasle")
	return false
}

func (room *Room) waitAllPlayerReady() {
	log.Debug(time.Now().Unix(), room, room.playedGameCnt, "Room.waitAllPlayerReady......")
	if room.playedGameCnt == 0 {
		room.switchStatus(RoomStatusGameStart)
		return
	}

	//等待所有玩家准备
	for _, player := range room.players {
		go room.waitPlayerReady(player)
	}
	timeout := int64(room.config.WaitPlayerOperateTimeout)
	start_time := time.Now().Unix()
	for  {
		//房间结束
		if room.roomStatus == RoomStatusRoomEnd {
			//如果此时房间已经结束，则直接返回，房间结束
			log.Debug(time.Now().Unix(), "waitAllPlayerReady room.roomStatus == RoomStatusRoomEnd")
			return
		}

		//所有人都已准备
		if room.isAllPlayerReady() {
			room.switchStatus(RoomStatusGameStart)
			return
		}

		//超时结束
		time_now := time.Now().Unix()
		if start_time + timeout < time_now {
			log.Debug(time.Now().Unix(), room, "waitAllPlayerReady timeout", timeout)
			room.switchStatus(RoomStatusRoomEnd) //超时发现没有足够的玩家都进入房间了，则结束
			return
		}

		time.Sleep(time.Millisecond * 2)
	}
}

func (room *Room) gameStart() {
	log.Debug(time.Now().Unix(), room, "gameStart", room.playedGameCnt)

	// 重置牌池, 洗牌
	room.Reset()
	room.cardPool.ReGenerate()

	//发牌
	master_pos := int32(util.RandomN(int(room.config.NeedPlayerNum)))
	if room.lastMaster != nil {
		master_pos = room.lastMaster.GetPosition()
	}
	room.masterPlayer = room.putCardsToPlayers(card.INIT_CARD_NUM, master_pos)
	room.switchOpMaster(room.masterPlayer, true, true, false)
	log.Debug(time.Now().Unix(), "master", room.masterPlayer)

	//通知所有玩家手上的牌
	for _, player := range room.players {
		player.OnGetInitCards()
	}

	//切换状态，开始打牌
	room.switchStatus(RoomStatusPlayGame)
	//log.Debug(time.Now().Unix(), room, "Room.playGame", room.playedGameCnt)

	//通知开始出牌
	sp_data := &StartPlayMsgData{
		Master:room.masterPlayer,
	}
	msg := NewStartPlayMsg(nil, sp_data)
	for _, player := range room.players {
		player.OnStartPlay(msg)
	}
}

func (room *Room) playGame() {
	log.Debug(time.Now().Unix(), room, "Room.playGame", room.playedGameCnt)

	is_round_end := false
	have_bomb := false
	curPlayer := room.opMaster
	if curPlayer.GetNeedDrop() {
		if !room.isFirstRound {
			//发牌
			is_dispatch_nil := room.dispatchCard(curPlayer)
			if is_dispatch_nil {
				//荒牌
				room.isHuangju = true
				room.dealRoundEnd()
				return
			}
			log.Debug(room, curPlayer, "dispatch now===", curPlayer.GetPlayingCards().CardsInHand.GetData())
			time.Sleep(room.config.AfterDispatchSec)
		}

		room.switchWaitPlayer(curPlayer, true, true, true)
		room.waitDropCard(curPlayer, true, true)
		if room.GetCardsType() >= card.CardsType_BOMB3 {
			have_bomb = true
		}

		//查看玩家是否出完手牌
		if room.isAllCardsDropped(curPlayer) {
			curPlayer.SetIsEndPlaying(true)
			curPlayer.SetRank(1)
			is_round_end = true
		}else {
			curPlayer.SetNeedDrop(false)
		}
	}

	if is_round_end {
		if have_bomb {
			room.IncBombNum()
		}
		room.dealRoundEnd()
		return
	}

	room.isFirstRound = false
	tmpPlayer := room.opMaster
	for {
		tmpPlayer = room.nextPlayer(tmpPlayer)
		if tmpPlayer == room.opMaster{
			//重置已出牌型
			room.SetCardsType(card.CardsType_NO)
			room.SetPlaneNum(0)
			room.SetWeight(0)

			room.opMaster.SetNeedDrop(true)
			break
		}

		canDrop := tmpPlayer.GetCanDrop()
		room.switchWaitPlayer(tmpPlayer, false, canDrop, true)
		is_drop := room.waitDropCard(tmpPlayer, false, canDrop)
		if is_drop {
			//查看玩家是否出完手牌
			if room.isAllCardsDropped(tmpPlayer) {
				tmpPlayer.SetIsEndPlaying(true)
				tmpPlayer.SetRank(1)
				is_round_end = true
			}
			if room.GetCardsType() >= card.CardsType_BOMB3 {
				have_bomb = true
			}

			room.switchOpMaster(tmpPlayer, false, true, false)
			break
		}
	}

	//每一小轮重复出现的炸弹只记一次
	if have_bomb {
		room.IncBombNum()
	}
	if is_round_end {
		room.dealRoundEnd()
	}
}

func (room *Room) dealRoundEnd() {
	room.switchStatus(RoomStatusEndPlayGame)
	//通知单局结算信息
	msg := room.summary()
	for _, player := range room.players {
		player.OnSummary(msg)
	}
}

func (room *Room) dispatchCard(curPlayer *Player) (is_dispatch_nil bool){
	dispatched_card := room.putCardToPlayer(curPlayer)
	log.Debug("#######dispatchCard:", curPlayer, dispatched_card)
	log.Debug("CardPool left card num:", room.cardPool.GetCardNum())
	if nil == dispatched_card {
		return true
	}

	//通知发牌
	dc_data := &DispatchCardMsgData{
		DispatchedPlayer:curPlayer,
		DispatchedCard:dispatched_card,
	}
	msg := NewDispatchCardMsg(curPlayer, dc_data)
	for _, player := range room.players {
		player.OnDispatchCard(msg)
	}
	return false
}

func (room *Room) isAllCardsDropped(player * Player) bool{
	return player.GetLeftCardNum() == 0
}

//在一个玩家出完手牌时判断此局是否已经结束
func (room *Room) isRoundEnd(endPlayingPlayer * Player) bool{
	return true
}

func (room *Room) switchOpMaster(player *Player, mustDrop bool, canDrop bool, needNotify bool) {
	log.Debug(time.Now().Unix(), room, "switchOperator", room.opMaster, "=>", player)
	room.opMaster = player
	player.SetNeedDrop(mustDrop)

	if needNotify {
		op := room.makeSwitchOperatorOperate(player, mustDrop, canDrop)
		for _, player := range room.players {
			player.OnPlayerSuccessOperated(op)
		}
	}
}

func (room *Room) switchWaitPlayer(player *Player, mustDrop bool, canDrop bool, needNotify bool) {
	log.Debug(time.Now().Unix(), room, "switchWaitPlayer", room.waitOperator, "=>", player)
	room.waitOperator = player

	if needNotify {
		op := room.makeSwitchOperatorOperate(player, mustDrop, canDrop)
		for _, player := range room.players {
			player.OnPlayerSuccessOperated(op)
		}
	}
}

func (room *Room) makeSwitchOperatorOperate(operator *Player, mustDrop bool, canDrop bool) *Operate {
	return NewSwitchOperator(operator, &OperateSwitchOperatorData{
		MustDrop:mustDrop,
		CanDrop:canDrop,
	})
}

/*func (room *Room) switchOperator(player *Player, mustDrop bool) {
	log.Debug(time.Now().Unix(), room, "switchOperator", room.opMaster, "=>", player)
	room.opMaster = player
	player.SetNeedDrop(mustDrop)

	op := room.makeSwitchOperatorOperate(player, mustDrop)
	for _, player := range room.players {
		player.OnPlayerSuccessOperated(op)
	}
}

func (room *Room) makeSwitchOperatorOperate(operator *Player, mustDrop bool) *Operate {
	return NewSwitchOperator(operator, &OperateSwitchOperatorData{MustDrop:mustDrop})
}*/

func (room *Room) endPlayGame() {
	room.playedGameCnt++
	log.Debug(time.Now().Unix(), room, "Room.endPlayGame cnt :", room.playedGameCnt)
	if room.isRoomEnd() {
		//log.Debug(time.Now().Unix(), room, "Room.endPlayGame room end")
		room.switchStatus(RoomStatusRoomEnd)
	} else {
		for _, player := range room.players {
			player.OnEndPlayGame()
		}
		//log.Debug(time.Now().Unix(), room, "Room.endPlayGame restart play game")
		room.switchStatus(RoomStatusWaitAllPlayerReady)
		log.Debug("============================================================================")
	}
}

func (room *Room) summary() *Message {
	info_type := card.InfoType_Normal
	bomb_num := room.GetBombNum()
	coin_multiple := 1
	for i := 0; i < bomb_num; i++ {
		coin_multiple *= 2
	}

	data := &SummaryMsgData{}
	data.Scores = make([]*PlayerSummaryData, 0)
	winner_win_coin := int32(0)
	for _, player := range room.players {
		left_card_num := player.GetPlayingCards().CardsInHand.Len()
		is_spring := false
		lose_coin := int32(0)
		is_win := false

		if left_card_num == 0 {
			is_win = true
		}else if left_card_num == 1{
		}else if left_card_num == 5{
			is_spring = true
			lose_coin = 0 - int32(coin_multiple * left_card_num*room.config.SpringMultiple)
		}else{
			lose_coin = 0 - int32(coin_multiple * left_card_num)
		}
		winner_win_coin -= lose_coin
		player.SetIsWin(is_win)
		player.AddCoin(lose_coin)
		player.AddTotalCoin(lose_coin)

		player_summary_data := &PlayerSummaryData{
			P:player,
			Rank:player.GetRank(),
			Coin:player.GetCoin(),
			TotalCoin:player.GetTotalCoin(),
			IsWin:player.GetIsWin(),
			LeftCardNum:int32(left_card_num),
			IsSpring:is_spring,
		}
		data.Scores = append(data.Scores, player_summary_data)

		if player.GetIsWin() {
			player.IncWinNum()
		}
		if is_spring {
			room.opMaster.IncSpringNum()
		}
	}

	//计算胜家赢多少
	room.opMaster.AddCoin(winner_win_coin)
	room.opMaster.AddTotalCoin(winner_win_coin)
	for _, score := range data.Scores {
		if score.P == room.opMaster {
			score.Coin = room.opMaster.GetCoin()
			score.TotalCoin = room.opMaster.GetTotalCoin()
		}
	}

	data.InfoType = info_type
	data.Multiple = int32(coin_multiple)
	return NewSummaryMsg(nil, data)
}

func (room *Room) totalSummary() *Message {
	var max_win, max_lose int32 = 0, 0
	for _, player := range room.players {
		total_coin := player.GetTotalCoin()
		if total_coin > max_win {
			max_win = total_coin
		}
		if total_coin < max_lose {
			max_lose = total_coin
		}
	}

	data := &RoomClosedMsgData{}
	data.Summaries = make([]*TotalSummaryData, 0)
	for _, player := range room.players {
		summary_data := &TotalSummaryData{
			P:player,
			WinNum:player.GetWinNum(),
			SpringNum:player.GetSpringNum(),
			TotalCoin:player.GetTotalCoin(),
			IsWinner:false,
			IsMostWinner:false,
			IsMostLoser:false,
			IsCreator:false,
		}
		if summary_data.TotalCoin > 0 {
			summary_data.IsWinner = true
		}
		if player.GetId() == room.creatorUid {
			summary_data.IsCreator = true
		}
		if summary_data.TotalCoin == max_lose {
			summary_data.IsMostLoser = true
		}
		if summary_data.TotalCoin == max_win {
			summary_data.IsMostWinner = true
		}
		data.Summaries = append(data.Summaries, summary_data)
	}
	return NewRoomClosedMsg(nil, data)
}

//取指定玩家的下一个玩家
func (room *Room) getPlayerByPos(position int32) *Player {
	for _, room_player := range room.players {
		if room_player.GetPosition() == position {
			return room_player
		}
	}
	if room.GetPlayerNum() > 0 {
		return room.players[0]
	}
	return nil
}

//取指定玩家的下一个玩家
func (room *Room) nextPlayer(player *Player) *Player {
	pos := player.GetPosition()

	need_player_num := int32(room.config.NeedPlayerNum)
	for i := int32(1); i <= need_player_num; i++ {
		next_pos := (pos + i) % need_player_num
		for _, room_player := range room.players {
			if room_player.GetPosition() == next_pos{
				if room_player == room.opMaster {
					//log.Debug(time.Now().Unix(), "nextPlayer", "pos:", pos, "next_pos:", next_pos)
					return room_player
				}
				if !room_player.GetIsEndPlaying() {
					//log.Debug(time.Now().Unix(), "nextPlayer", "pos:", pos, "next_pos:", next_pos)
					return room_player
				}
			}
		}
	}

	return room.players[0]
}

//取指定玩家的对家
func (room *Room) oppositePlayer(player *Player) *Player {
	if nil == player{
		return nil
	}

	pos := player.GetPosition()
	opp_pos := (pos + 2) % int32(room.config.NeedPlayerNum)
	for _, room_player := range room.players {
		if room_player.GetPosition() == opp_pos && !room_player.GetIsEndPlaying(){
			return room_player
		}
	}
	return nil
}

func (room *Room) isAllPlayerReady() bool{
	for _, player := range room.players {
		if !player.isReady {
			return false
		}
	}
	return true
}

//处理玩家操作
func (room *Room) dealPlayerOperate(op *Operate) bool{
	//log_time := time.Now().Unix()
	//log.Debug(log_time, room, "Room.dealPlayerOperate :", op)
	switch op.Op {
	case OperateEnterRoom:
		if _, ok := op.Data.(*OperateEnterRoomData); ok {
			if room.addPlayer(op.Operator) {
				//玩家进入成功
				player_pos := room.getMinUsablePosition()
				op.Operator.EnterRoom(room, player_pos)
				//log.Debug(log_time, room, "Room.dealPlayerOperate player enter :", op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperateReadyRoom:
		if _, ok := op.Data.(*OperateReadyRoomData); ok {
			if room.readyPlayer(op.Operator) { //	玩家确认开始游戏
				op.Operator.ReadyRoom(room)
				//log.Debug(log_time, room, "Room.dealPlayerOperate player ready :", op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperateLeaveRoom:
		if _, ok := op.Data.(*OperateLeaveRoomData); ok {
			//log.Debug(log_time, room, "Room.dealPlayerOperate player leave :", op.Operator)
			room.delPlayer(op.Operator)
			op.Operator.LeaveRoom()
			op.ResultCh <- true
			room.broadcastPlayerSuccessOperated(op)
			return true
		}

	case OperateDrop:
		if drop_data, ok := op.Data.(*OperateDropData); ok {
			if op.Operator.Drop(drop_data.whatGroup) {
				log.Debug(room, op.Operator, "drop left===", op.Operator.GetPlayingCards().CardsInHand.GetData())
				//出牌，计算分数和奖
				prize := room.getPrizeByCardsType(drop_data.cardsType)
				op.Operator.AddPrize(prize)
				if drop_data.cardsType != card.CardsType_NO{
					room.SetCardsType(drop_data.cardsType)
					room.SetPlaneNum(drop_data.planeNum)
					room.SetWeight(drop_data.weight)
				}

				log.Debug(time.Now().Unix(), room, "Room.dealPlayerOperate player drop :", op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperatePass:
		if _, ok := op.Data.(*OperatePassData); ok {
			log.Debug("-------", time.Now().Unix(), room, "Room.dealPlayerOperate player pass :", op.Operator)
			op.ResultCh <- true
			room.broadcastPlayerSuccessOperated(op)
			return true
		}

	}
	op.ResultCh <- false
	return false
}

func (room *Room) getCardsScores(drop_cards []*card.Card) (int32)  {
	score := int32(0)
	for _, drop_card := range drop_cards {
		if drop_card.CardNo == 5 {
			score += 5
		}else if drop_card.CardNo == 10 {
			score += 10
		}else if drop_card.CardNo == 13 {
			score += 10
		}
	}

	return score
}

func (room *Room) getPrizeByCardsType(cards_type int) (int32) {
	prize := int32(0)
	switch cards_type {
	case 24:
		prize = 1
	case 25:
		prize = 2
	case 26:
		prize = 3
	case 27, 28:
		prize = 5
	}
	return prize
}

//查找房间中未被占用的最新的position
func (room *Room) getMinUsablePosition() (int32)  {
	//log.Debug(time.Now().Unix(), room, "getMinUsablePosition")
	//获取所有已经被占用的position
	player_positions := make([]int32, 0)
	for _, room_player := range room.players {
		player_positions = append(player_positions, room_player.GetPosition())
	}

	//查找未被占用的position中最小的
	room_max_position := int32(room.config.NeedPlayerNum - 1)
	for i := int32(0); i <= room_max_position; i++ {
		is_occupied := false
		for _, occupied_pos := range player_positions{
			if occupied_pos == i {
				is_occupied = true
				break
			}
		}
		if !is_occupied {
			return i
		}
	}
	return room_max_position
}

//给所有玩家发牌
func (room *Room) putCardsToPlayers(init_num int, master_pos int32) (master *Player) {
	//log.Debug(time.Now().Unix(), room, "Room.putCardsToPlayers, init_type:", init_type)
	//确定master
	master = nil
	for _, room_player := range room.players {
		if room_player.GetPosition() == int32(master_pos) {
			master = room_player
		}
	}
	if nil == master {
		log.Error("nil == master")
	}

	//发牌
	for num := 0; num < init_num; num++ {
		for _, player := range room.players {
			room.putCardToPlayer(player)
			//log.Debug("put_card:", put_card)
		}
	}
	room.putCardToPlayer(master)
	return
}

//添加玩家
func (room *Room) addPlayer(player *Player) bool {
	/*if room.roomStatus != RoomStatusWaitAllPlayerEnter {
		return false
	}*/
	if room.GetPlayerNum() >= room.config.NeedPlayerNum {
		return false
	}
	room.players = append(room.players, player)
	return true
}

func (room *Room) readyPlayer(player *Player) bool {
	if room.roomStatus != RoomStatusWaitAllPlayerEnter && room.roomStatus != RoomStatusWaitAllPlayerReady{
		return false
	}
	player.SetIsReady(true)
	return true
}

func (room *Room) delPlayer(player *Player)  {
	for idx, p := range room.players {
		if p == player {
			room.players = append(room.players[0:idx], room.players[idx+1:]...)
			return
		}
	}
}

func (room *Room) broadcastPlayerSuccessOperated(op *Operate) {
	//log.Debug(time.Now().Unix(), room, "Room.broadcastPlayerSucOp :", op)
	for _, player := range room.players {
		player.OnPlayerSuccessOperated(op)
	}
}

//发牌给指定玩家
func (room *Room) putCardToPlayer(player *Player) *card.Card {
	dis_card := room.cardPool.PopFront()
	if dis_card == nil {
		return nil
	}
	player.AddCard(dis_card)
	return dis_card
}

func (room *Room) Reset() {
	room.lastMaster = room.opMaster
	room.opMaster = nil
	room.waitOperator = nil
	room.masterPlayer = nil
	room.cardsType = 0
	room.planeNum = 0
	room.weight = 0
	room.bombNum = 0
	room.isFirstRound = true
	room.isHuangju = false
}

func (room *Room) String() string {
	if room == nil {
		return "{room=nil}"
	}
	return fmt.Sprintf("{room=%v}", room.GetId())
}

func (room *Room) clearChannel() {
	for idx := 0 ; idx < int(room.config.NeedPlayerNum); idx ++ {
		select {
		case op := <-room.dropCardCh[idx]:
			op.ResultCh <- false
		default:
		}

		select {
		case op := <-room.passCh[idx]:
			op.ResultCh <- false
		default:
		}

		select {
		case op := <-room.roomReadyCh[idx]:
			op.ResultCh <- false
		default:
		}
	}
}
