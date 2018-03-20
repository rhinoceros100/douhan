package playing

import (
	"testing"
	"douhan/card"
	"douhan/util"
	"github.com/bmizerany/assert"
)

func TestGetCanDrop(t *testing.T) {
	conf := NewRoomConfig()
	conf.Init(1, 2, 1)
	room := NewRoom(util.UniqueId(), conf)

	player1 := NewPlayer(1)
	player1.room = room
	player2 := NewPlayer(2)
	player2.room = room
	player3 := NewPlayer(3)
	player3.room = room
	player4 := NewPlayer(4)
	player4.room = room

	card1 := &card.Card{CardType: card.CardType_Diamond,CardNo: 1, Weight:14,}		//方片1
	card2 := &card.Card{CardType: card.CardType_Club,CardNo: 1, Weight:14,}			//梅花1
	card3 := &card.Card{CardType: card.CardType_Heart,CardNo: 1, Weight:14,}		//红桃1
	card4 := &card.Card{CardType: card.CardType_Spade,CardNo: 1, Weight:14,}		//黑桃1
	card5 := &card.Card{CardType: card.CardType_Diamond,CardNo: 5, Weight:5,}		//方片1
	card6 := &card.Card{CardType: card.CardType_Club,CardNo: 5, Weight:5,}			//梅花1
	card7 := &card.Card{CardType: card.CardType_Heart,CardNo: 5, Weight:5,}			//黑桃
	card8 := &card.Card{CardType: card.CardType_BlackJoker,CardNo: 14, Weight:16,}		//小王
	card9 := &card.Card{CardType: card.CardType_Heart,CardNo: 2, Weight:15,}		//红桃2
	card10 := &card.Card{CardType: card.CardType_Heart,CardNo: 7, Weight:7,}		//红桃7
	card11 := &card.Card{CardType: card.CardType_Spade,CardNo: 7, Weight:7,}		//黑桃7
	card12 := &card.Card{CardType: card.CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	card13 := &card.Card{CardType: card.CardType_RedJoker,CardNo: 14, Weight:17,}		//大王
	card14 := &card.Card{CardType: card.CardType_Spade,CardNo: 8, Weight:8,}		//黑桃8

	//单牌
	room.SetCardsType(card.CardsType_SINGLE)
	room.SetWeight(6)
	room.SetPlaneNum(1)
	player1.AddCard(card1)
	assert.Equal(t, player1.GetCanDrop(), false)
	player2.AddCard(card8)
	assert.Equal(t, player2.GetCanDrop(), false)
	player3.AddCard(card9)
	assert.Equal(t, player3.GetCanDrop(), true)
	player4.AddCard(card10)
	assert.Equal(t, player4.GetCanDrop(), true)

	//对子
	room.SetCardsType(card.CardsType_PAIR)
	room.SetWeight(6)
	room.SetPlaneNum(1)
	player1.AddCard(card2)
	assert.Equal(t, player1.GetCanDrop(), false)
	player2.AddCard(card7)
	assert.Equal(t, player2.GetCanDrop(), false)
	player3.AddCard(card12)
	assert.Equal(t, player3.GetCanDrop(), true)
	player4.AddCard(card11)
	assert.Equal(t, player4.GetCanDrop(), true)
	player4.AddCard(card8)
	assert.Equal(t, player4.GetCanDrop(), true)

	//炸弹
	room.SetCardsType(card.CardsType_BOMB3)
	room.SetWeight(6)
	room.SetPlaneNum(1)
	player1.AddCard(card3)
	assert.Equal(t, player1.GetCanDrop(), true)
	player1.AddCard(card4)
	assert.Equal(t, player1.GetCanDrop(), true)
	player2.AddCard(card6)
	assert.Equal(t, player2.GetCanDrop(), false)
	player2.AddCard(card5)
	assert.Equal(t, player2.GetCanDrop(), true)
	assert.Equal(t, player4.GetCanDrop(), true)
	player4.AddCard(card13)
	assert.Equal(t, player4.GetCanDrop(), true)

	//顺子
	room.SetCardsType(card.CardsType_STRAIGHT)
	room.SetWeight(4)
	room.SetPlaneNum(3)//456
	assert.Equal(t, player4.GetCanDrop(), true)
	player5 := NewPlayer(5)
	player5.room = room
	player5.AddCard(card7)
	player5.AddCard(card8)
	player5.AddCard(card10)
	assert.Equal(t, player5.GetCanDrop(), true)
	player6 := NewPlayer(6)
	player6.room = room
	player6.AddCard(card9)
	player6.AddCard(card10)
	player6.AddCard(card13)
	assert.Equal(t, player6.GetCanDrop(), false)

	room.SetPlaneNum(4)//4567
	assert.Equal(t, player5.GetCanDrop(), false)
	player5.AddCard(card13)
	assert.Equal(t, player5.GetCanDrop(), true)
	player6.AddCard(card14)
	assert.Equal(t, player6.GetCanDrop(), false)

}