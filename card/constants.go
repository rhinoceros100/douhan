package card

const (
	CardType_Diamond	int = iota + 1 		//方片
	CardType_Club	 			        //梅花
	CardType_Heart		 			//红桃
	CardType_Spade					//黑桃
	CardType_BlackJoker				//小王
	CardType_RedJoker				//大王
)

const (
	CardsType_NO		int = iota		//没有任何牌型
	CardsType_SINGLE				//单牌
	CardsType_STRAIGHT				//顺子
	CardsType_PAIR			           	//对子
	CardsType_BOMB3		= 20      		//三炸
	CardsType_BOMB4		= 23           		//四炸
	CardsType_BOMB5		= 24	           	//五炸
	CardsType_BOMB6		= 25        		//六炸
)

const (
	InfoType_Normal		int32 = iota		//普通结果
	InfoType_Shuangji				//双基
	InfoType_PlayAloneSucc				//打独成功
	InfoType_PlayAloneFail				//打独失败
)
