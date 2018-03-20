package card

import (
	"testing"
	"time"
)

func TestNewPlayingCardsIsOk(t *testing.T) {
	start := time.Now()
	hu15 := &Cards{
		Data: []*Card{
			&Card{CardType:CardType_Diamond, CardNo:1},

			&Card{CardType:CardType_Diamond, CardNo:2},
			&Card{CardType:CardType_Diamond, CardNo:2},

			&Card{CardType:CardType_Club, CardNo:2},
			&Card{CardType:CardType_Club, CardNo:2},

			&Card{CardType:CardType_Diamond, CardNo:3},

			&Card{CardType:CardType_Diamond, CardNo:3},
			&Card{CardType:CardType_Club, CardNo:3},
			&Card{CardType:CardType_Club, CardNo:3},

			&Card{CardType:CardType_Diamond, CardNo:6},
			&Card{CardType:CardType_Diamond, CardNo:7},
			&Card{CardType:CardType_Diamond, CardNo:8},

			&Card{CardType:CardType_Diamond, CardNo:4},
			&Card{CardType:CardType_Diamond, CardNo:4},
			&Card{CardType:CardType_Club, CardNo:4},

		},
	}

	playingCards := NewPlayingCards()
	playingCards.AddCards(hu15)
	t.Log(time.Now().Sub(start))
}