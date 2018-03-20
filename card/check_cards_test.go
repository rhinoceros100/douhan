package card

import (
	"testing"
	"github.com/bmizerany/assert"
)

func TestGetCardsType(t *testing.T) {
	card1 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card2 := &Card{CardType: CardType_Club,CardNo: 1, Weight:14,}		//梅花1
	card3 := &Card{CardType: CardType_Heart,CardNo: 1, Weight:14,}		//红桃1
	card4 := &Card{CardType: CardType_Spade,CardNo: 1, Weight:14,}		//黑桃1

	card5 := &Card{CardType: CardType_Diamond,CardNo: 2, Weight:15,}	//方片2
	card6 := &Card{CardType: CardType_Club,CardNo: 2, Weight:15,}		//梅花2
	card7 := &Card{CardType: CardType_Heart,CardNo: 2, Weight:15,}		//红桃2
	card8 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2

	card13 := &Card{CardType: CardType_Spade,CardNo: 3, Weight:3,}		//黑桃3
	//card14 := &Card{CardType: CardType_Spade,CardNo: 4, Weight:4,}		//黑桃4
	card15 := &Card{CardType: CardType_Spade,CardNo: 5, Weight:5,}		//黑桃5
	card16 := &Card{CardType: CardType_Spade,CardNo: 6, Weight:6,}		//黑桃6
	card17 := &Card{CardType: CardType_Spade,CardNo: 7, Weight:7,}		//黑桃7
	card18 := &Card{CardType: CardType_Spade,CardNo: 8, Weight:8,}		//黑桃8
	card19 := &Card{CardType: CardType_Spade,CardNo: 9, Weight:9,}		//黑桃9
	//card20 := &Card{CardType: CardType_Spade,CardNo: 10, Weight:10,}	//黑桃10
	//card21 := &Card{CardType: CardType_Spade,CardNo: 11, Weight:11,}	//黑桃J
	card22 := &Card{CardType: CardType_Spade,CardNo: 12, Weight:12,}	//黑桃Q
	//card23 := &Card{CardType: CardType_Spade,CardNo: 13, Weight:13,}	//黑桃K

	card31 := &Card{CardType: CardType_BlackJoker,CardNo: 14, Weight:16}	//小王
	card32 := &Card{CardType: CardType_RedJoker,CardNo: 14, Weight:17}	//大王

	cards1 := make([]*Card, 0)
	cards1 = append(cards1, card1)
	cards1 = append(cards1, card2)
	cards1 = append(cards1, card3)
	cards1 = append(cards1, card4)
	drop_cards1 := CreateNewCards(cards1)

	cards2 := make([]*Card, 0)
	cards2 = append(cards2, card5)
	cards2 = append(cards2, card6)
	cards2 = append(cards2, card7)
	cards2 = append(cards2, card8)
	cards2 = append(cards2, card31)
	cards2 = append(cards2, card32)
	drop_cards2 := CreateNewCards(cards2)

	cards3 := make([]*Card, 0)
	cards3 = append(cards3, card5)
	cards3 = append(cards3, card6)
	cards3 = append(cards3, card7)
	cards3 = append(cards3, card31)
	cards3 = append(cards3, card32)
	drop_cards3 := CreateNewCards(cards3)

	cards4 := make([]*Card, 0)
	cards4 = append(cards4, card5)
	cards4 = append(cards4, card6)
	cards4 = append(cards4, card32)
	drop_cards4 := CreateNewCards(cards4)

	cards5 := make([]*Card, 0)
	cards5 = append(cards5, card13)
	cards5 = append(cards5, card31)
	cards5 = append(cards5, card32)
	drop_cards5 := CreateNewCards(cards5)

	cards6 := make([]*Card, 0)
	cards6 = append(cards6, card13)
	cards6 = append(cards6, card32)
	drop_cards6 := CreateNewCards(cards6)

	cards7 := make([]*Card, 0)
	cards7 = append(cards7, card1)
	cards7 = append(cards7, card2)
	drop_cards7 := CreateNewCards(cards7)

	cards8 := make([]*Card, 0)
	cards8 = append(cards8, card31)
	cards8 = append(cards8, card32)
	drop_cards8 := CreateNewCards(cards8)

	cards9 := make([]*Card, 0)
	cards9 = append(cards9, card31)
	drop_cards9 := CreateNewCards(cards9)

	cards10 := make([]*Card, 0)
	cards10 = append(cards10, card6)
	drop_cards10 := CreateNewCards(cards10)

	cards11 := make([]*Card, 0)
	cards11 = append(cards11, card13)
	cards11 = append(cards11, card15)
	cards11 = append(cards11, card31)
	drop_cards11 := CreateNewCards(cards11)

	cards12 := make([]*Card, 0)
	cards12 = append(cards12, card22)
	cards12 = append(cards12, card1)
	cards12 = append(cards12, card31)
	cards12 = append(cards12, card32)
	drop_cards12 := CreateNewCards(cards12)

	cards13 := make([]*Card, 0)
	cards13 = append(cards13, card17)
	cards13 = append(cards13, card18)
	cards13 = append(cards13, card31)
	cards13 = append(cards13, card32)
	drop_cards13 := CreateNewCards(cards13)

	cards14 := make([]*Card, 0)
	cards14 = append(cards14, card16)
	cards14 = append(cards14, card19)
	cards14 = append(cards14, card31)
	cards14 = append(cards14, card32)
	drop_cards14 := CreateNewCards(cards14)

	t.Log(GetCardsType(drop_cards1, true, 0))		//四炸 14
	t.Log(GetCardsType(drop_cards2, false, 0))		//六炸 15
	t.Log(GetCardsType(drop_cards3, false, 0))		//五炸 15
	t.Log(GetCardsType(drop_cards4, false, 0))		//三炸 15
	t.Log(GetCardsType(drop_cards5, false, 0))		//三炸 3
	t.Log(GetCardsType(drop_cards6, false, 0))		//对子
	t.Log(GetCardsType(drop_cards7, false, 0))		//对子
	t.Log(GetCardsType(drop_cards8, true, 0))		//对王
	t.Log(GetCardsType(drop_cards8, false, 0))		//无牌型
	t.Log(GetCardsType(drop_cards9, false, 0))		//无牌型
	t.Log(GetCardsType(drop_cards10, false, 0))	//单牌
	t.Log(GetCardsType(drop_cards11, false, 0))	//顺子 3
	t.Log(GetCardsType(drop_cards12, false, 11))	//顺子 11
	t.Log(GetCardsType(drop_cards13, false, 6))	//顺子 7
	t.Log(GetCardsType(drop_cards14, false, 6))	//顺子 6
}

func TestGetSameCardsNum(t *testing.T) {
	card1 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card2 := &Card{CardType: CardType_Club,CardNo: 1, Weight:14,}		//梅花1
	card3 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card4 := &Card{CardType: CardType_Club,CardNo: 2, Weight:15,}		//梅花2
	card5 := &Card{CardType: CardType_Heart,CardNo: 2, Weight:15,}		//红桃2
	card6 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	card7 := &Card{CardType: CardType_Spade,CardNo: 5, Weight:5,}		//黑桃5
	card8 := &Card{CardType: CardType_Spade,CardNo: 8, Weight:8,}		//黑桃8
	card9 := &Card{CardType: CardType_BlackJoker,CardNo: 14, Weight:16,}	//小王
	card10 := &Card{CardType: CardType_RedJoker,CardNo: 14, Weight:17,}	//大王

	cards1 := make([]*Card, 0)
	cards1 = append(cards1, card1)
	cards1 = append(cards1, card2)
	cards1 = append(cards1, card3)
	cards1 = append(cards1, card4)

	cards2 := make([]*Card, 0)
	cards2 = append(cards2, card1)
	cards2 = append(cards2, card2)
	cards2 = append(cards2, card3)
	cards2 = append(cards2, card4)
	cards2 = append(cards2, card5)
	cards2 = append(cards2, card6)
	cards2 = append(cards2, card7)
	cards2 = append(cards2, card8)
	cards2 = append(cards2, card9)
	cards2 = append(cards2, card10)

	t.Log(GetSameCardsNum(cards1))
	t.Log(GetSameCardsNum(cards2))
}

func TestIsSameCardType(t *testing.T) {
	card1 := &Card{CardType: CardType_Diamond,CardNo: 1,}	//方片1
	card2 := &Card{CardType: CardType_Club,CardNo: 1,}	//梅花1
	card3 := &Card{CardType: CardType_Diamond,CardNo: 1,}	//方片1
	//card4 := &Card{CardType: CardType_Club,CardNo: 2,}	//梅花2
	card5 := &Card{CardType: CardType_Heart,CardNo: 2,}	//红桃2
	card6 := &Card{CardType: CardType_Spade,CardNo: 2,}	//黑桃2
	card7 := &Card{CardType: CardType_Spade,CardNo: 5,}	//黑桃5
	card8 := &Card{CardType: CardType_Spade,CardNo: 8,}	//黑桃8
	card9 := &Card{CardType: CardType_Spade,CardNo: 2,}	//黑桃2

	cards1 := make([]*Card, 0)
	cards1 = append(cards1, card1)
	cards1 = append(cards1, card2)
	cards1 = append(cards1, card3)
	assert.Equal(t, IsSameCardType(cards1), false)

	cards2 := make([]*Card, 0)
	cards2 = append(cards2, card6)
	cards2 = append(cards2, card7)
	cards2 = append(cards2, card8)
	cards2 = append(cards2, card9)
	assert.Equal(t, IsSameCardType(cards2), true)

	cards3 := make([]*Card, 0)
	cards3 = append(cards3, card5)
	assert.Equal(t, IsSameCardType(cards3), true)

	cards4 := make([]*Card, 0)
	assert.Equal(t, IsSameCardType(cards4), false)
}

func TestIsStraight(t *testing.T) {
	nums1 := []int{2,3,4}
	t.Log(IsStraight(nums1, 0))

	nums2 := []int{3,4}
	t.Log(IsStraight(nums2, 0))

	nums3 := []int{16,17,13,12}
	t.Log(IsStraight(nums3, 10))

	nums4 := []int{16,13,12}
	t.Log(IsStraight(nums4, 11))

	nums5 := []int{16,13,11}
	t.Log(IsStraight(nums5, 11))

	nums6 := []int{16,17,13,10}
	t.Log(IsStraight(nums6, 0))

	nums7 := []int{16,17,4,6}
	t.Log(IsStraight(nums7, 3))

	nums8 := []int{16,17,12,14}
	t.Log(IsStraight(nums8, 0))

	nums9 := []int{16,17,12,15}
	t.Log(IsStraight(nums9, 0))
}

func TestGetStraightWeight(t *testing.T) {
	nums1 := []int{2,3,4}
	t.Log(GetStraightWeight(nums1))

	nums2 := []int{3,4}
	t.Log(GetStraightWeight(nums2))

	nums3 := []int{1,2,13,12}
	t.Log(GetStraightWeight(nums3))
}

type roon struct {
	cardsType int
	planeNum int
	weight int
}

func canCover(cardsType, planeNum, weight int, roo *roon) (canCover bool) {
	canCover = false
	if roo.cardsType == CardsType_NO {
		return cardsType != CardsType_NO
	}
	//已经出的牌型非炸弹牌型
	if roo.cardsType < 20{
		if cardsType >= 20 {
			return true
		}
		//普通牌型打普通牌型必须为同一牌型，并且飞机数量必须相同
		if cardsType != roo.cardsType{
			return false
		}
		if cardsType == CardsType_STRAIGHT {
			if planeNum != roo.planeNum {
				return false
			}
			return weight == roo.weight + 1
		}
		//非炸弹中2最大
		if weight == 15 {
			return false
		}
		return weight == 15 || weight == roo.weight + 1
	}

	//更大的炸弹可以管住
	if cardsType > roo.cardsType{
		return true
	}
	return weight > roo.weight
}

func TestCanCover(t *testing.T) {
	roo1 := &roon{cardsType:0, planeNum:0, weight:0}
	assert.Equal(t, canCover(CardsType_STRAIGHT, 3, 4, roo1), true)
	assert.Equal(t, canCover(CardsType_PAIR, 1, 4, roo1), true)
	assert.Equal(t, canCover(CardsType_NO, 3, 2, roo1), false)

	roo2 := &roon{cardsType:CardsType_SINGLE, planeNum:0, weight:6}
	assert.Equal(t, canCover(CardsType_SINGLE, 0, 3, roo2), false)
	assert.Equal(t, canCover(CardsType_SINGLE, 0, 7, roo2), true)
	assert.Equal(t, canCover(CardsType_SINGLE, 0, 8, roo2), false)
	assert.Equal(t, canCover(CardsType_PAIR, 0, 7, roo2), false)
	assert.Equal(t, canCover(CardsType_BOMB3, 0, 4, roo2), true)

	roo3 := &roon{cardsType:CardsType_STRAIGHT, planeNum:3, weight:5}
	assert.Equal(t, canCover(CardsType_STRAIGHT, 3, 6, roo3), true)
	assert.Equal(t, canCover(CardsType_STRAIGHT, 3, 7, roo3), false)
	assert.Equal(t, canCover(CardsType_STRAIGHT, 4, 6, roo3), false)
	assert.Equal(t, canCover(CardsType_SINGLE, 1, 6, roo3), false)
	assert.Equal(t, canCover(CardsType_BOMB4, 1, 8, roo3), true)

	roo5 := &roon{cardsType:CardsType_BOMB3, planeNum:0, weight:5}
	assert.Equal(t, canCover(CardsType_BOMB3, 1, 3, roo5), false)
	assert.Equal(t, canCover(CardsType_BOMB3, 1, 8, roo5), true)
	assert.Equal(t, canCover(CardsType_BOMB4, 1, 3, roo5), true)
	assert.Equal(t, canCover(CardsType_BOMB5, 1, 3, roo5), true)
	assert.Equal(t, canCover(CardsType_BOMB6, 1, 4, roo5), true)
}

