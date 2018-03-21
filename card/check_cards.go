package card

func GetCardsType(the_cards *Cards, is_last_cards bool, straight_weight int) (cards_type int, plane_num int, weight int) {
	the_cards.Sort()
	drop_cards := the_cards.GetData()
	cards_type = CardsType_NO
	plane_num = 0
	weight = 0
	cards_len := len(drop_cards)

	if cards_len == 0 {
		return
	}
	if cards_len == 1 {
		if drop_cards[0].CardNo != 14 {
			cards_type = CardsType_SINGLE
		}
		weight = drop_cards[0].Weight
		return
	}

	most, sames, joker_num := GetSameCardsNum(drop_cards)
	if cards_len == 2 {
		if most == 2 {
			//大王、小王不能成对
			if drop_cards[0].Weight != drop_cards[1].Weight {
				return CardsType_NO, plane_num, weight
			}
			return CardsType_PAIR, plane_num, drop_cards[0].Weight
		}else{
			if joker_num == 2 {
				if is_last_cards {
					return CardsType_PAIR, plane_num, drop_cards[0].Weight
				}
			}else if joker_num == 1 {
				if drop_cards[0].Weight > 15 {
					weight = drop_cards[1].Weight
				}else{
					weight = drop_cards[0].Weight
				}
				return CardsType_PAIR, plane_num, weight
			}
		}
		return
	}
	//单牌和双牌需要区分大小王的权重，其他不需要
	weight = GetStraightWeight(sames)
	if cards_len == 3 {
		if most == 3 {
			return CardsType_BOMB3, plane_num, weight
		}else if most == 2 {
			if joker_num == 1 {
				return CardsType_BOMB3, plane_num, weight
			}
		}else if most == 1 {
			if joker_num == 2 {
				return CardsType_BOMB3, plane_num, weight
			}
		}
		is_straight, weight := IsStraight(sames, straight_weight)
		if is_straight {
			cards_type = CardsType_STRAIGHT
			plane_num = 3
		}
		return cards_type, plane_num, weight
	}

	//四张及以上，只有炸弹和顺子牌型
	if most + joker_num == cards_len {
		switch cards_len {
		case 4:
			cards_type = CardsType_BOMB4
		case 5:
			cards_type = CardsType_BOMB5
		case 6:
			cards_type = CardsType_BOMB6
		}
		return cards_type, plane_num, weight
	}
	if most == 1 {
		is_straight, weight := IsStraight(sames, straight_weight)
		if is_straight {
			plane_num = cards_len
			return CardsType_STRAIGHT, plane_num, weight
		}
	}

	return CardsType_NO, plane_num, weight
}

//获取一组牌中数量最多的数字相同的牌的数量
func GetSameCardsNum(drop_cards []*Card) (most int, same_card_nums []int, joker_num int) {
	arr := [18]int{}
	for _, drop_card := range drop_cards {
		arr[drop_card.Weight] ++
	}

	same_card_nums = make([]int, 0)
	most = 0
	joker_num = arr[16] + arr[17]
	for i, num := range arr {
		if num > most {
			same_card_nums = same_card_nums[0:0]
			most = num
			same_card_nums = append(same_card_nums, i)
		}else if num == most {
			same_card_nums = append(same_card_nums, i)
		}
	}
	return
}

func IsSameCardType(drop_cards []*Card) (bool) {
	if len(drop_cards) == 0 {
		return false
	}

	card_type := drop_cards[0].CardType
	for _, drop_card := range drop_cards {
		if drop_card.CardType != card_type {
			return false
		}
	}
	return true
}

/*是否为顺子，算入赖子
* 返回是否为顺子，及顺子权重
* param weights:要校验的牌的权重
* param check_weight:要组成的顺子的权重
*/
func IsStraight(weights []int, check_weight int) (is_straight bool, straight_weight int) {
	weight_len := len(weights)
	is_straight = false
	joker_num, straight_weight := 0, 0
	smallest, biggest := 15, 1

	if weight_len <= 2 {
		return
	}

	//获取牌的权重3-14
	weight_arr := [15]bool{}
	for _, weight := range weights {
		if weight == 15 {
			return
		}else if weight > 15 {
			joker_num ++
		}else{
			if weight_arr[weight] {
				return
			}
			weight_arr[weight] = true
			if weight < smallest{
				smallest = weight
			}
			if weight > biggest{
				biggest = weight
			}
		}
	}

	//检查是否可以组成指定的顺子
	if check_weight >= 3 {
		//组成的顺子最大不能大于A
		straight_weight = check_weight
		if check_weight + weight_len - 1 > 14 {
			is_straight = false
			return
		}
		empty_weight := 0
		for i:= 0; i < weight_len; i++ {
			if !weight_arr[check_weight + i] {
				empty_weight++
			}
		}
		is_straight = empty_weight == joker_num
		return
	}

	if biggest - smallest + 1 == weight_len {
		//4xx56
		straight_weight = smallest
		return true, straight_weight
	}else if biggest - smallest + 1 + joker_num == weight_len {
		//45xx
		straight_weight = smallest
		if biggest + joker_num > 14 {
			straight_weight = 14 + 1 - weight_len
		}
		return true, straight_weight
	}else if biggest - smallest + joker_num == weight_len {
		//4x6x
		straight_weight = smallest
		if biggest + joker_num - 1 > 14 {
			straight_weight = 14 + 1 - weight_len
		}
		return true, straight_weight
	}
	return
}

func GetStraightWeight(nums []int) (weight int) {
	weight = 18
	for _, num := range nums {
		w := num
		if w < weight {
			weight = w
		}
	}
	return
}
