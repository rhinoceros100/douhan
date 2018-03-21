package card

import "fmt"

const INIT_CARD_NUM int = 5      	//开局发牌数量
const INIT_CARD_NUM_MASTER int = 6      //上局赢家开局发牌数量
const TOTAL_CARD_NUM int = 54 		//牌的总数量

type PlayingCards struct {
	CardsInHand			*Cards		//手上的牌
}

func NewPlayingCards() *PlayingCards {
	return  &PlayingCards{
		CardsInHand: NewCards(),
	}
}

func (playingCards *PlayingCards) Reset() {
	playingCards.CardsInHand.Clear()
}

func (playingCards *PlayingCards) AddCards(cards *Cards) {
	playingCards.CardsInHand.AppendCards(cards)
	playingCards.CardsInHand.Sort()
}

//增加一张牌
func (playingCards *PlayingCards) AddCard(card *Card) {
	playingCards.CardsInHand.AddAndSort(card)
}

func (playingCards *PlayingCards) String() string{
	return fmt.Sprintf(
		"{%v}",
		playingCards.CardsInHand,
	)
}

func (playingCards *PlayingCards) Tail(num int) []*Card {
	return playingCards.CardsInHand.Tail(num)
}

func (playingCards *PlayingCards) RandomTail() []*Card {
	hand_cards := playingCards.CardsInHand
	cards_len := hand_cards.Len()
	tail_weight := hand_cards.At(cards_len - 1).Weight
	tail_num := 1
	for i:= 1; i < cards_len; i++ {
		if hand_cards.At(cards_len - 1 - i).Weight == tail_weight {
			tail_num++
		}
	}
	if cards_len == 2 && hand_cards.At(0).Weight > 15{
		tail_num = 2
	}
	if cards_len == 3 && hand_cards.At(0).Weight > 15 && (hand_cards.At(1).Weight == hand_cards.At(2).Weight){
		tail_num = 3
	}
	return playingCards.CardsInHand.Tail(tail_num)
}

func (playingCards *PlayingCards) Get2JokerNum() (num2, num_joker int32) {
	return playingCards.CardsInHand.Get2JokerNum()
}

func (playingCards *PlayingCards) GetCardNumByWeight(card_weight int) (int32) {
	return playingCards.CardsInHand.GetCardNumByWeight(card_weight)
}

//丢弃一张牌
func (playingCards *PlayingCards) DropCards(cards []*Card) bool {
	return playingCards.CardsInHand.TakeAwayGroup(cards)
}