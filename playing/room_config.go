package playing

import (
	"time"
)

type RoomConfig struct {
	NeedPlayerNum			int32        `json:"need_player_num"`
	MaxPlayGameCnt			int	     `json:"max_play_game_cnt"`	//最大的游戏局数
	RandomDropNum			int          `json:"random_drop_num"`	//随机出牌张数
	SpringMultiple			int          `json:"spring_multiple"`	//春天倍数

	WaitPlayerEnterRoomTimeout	int        `json:"wait_player_enter_room_timeout"`
	WaitPlayerOperateTimeout	int        `json:"wait_player_operate_timeout"`
	WaitDropSec                	time.Duration      `json:"wait_drop_sec"`	//等待出牌时长
	WaitReadySec              	time.Duration      `json:"wait_ready_sec"`	//等待准备时长
	AfterDispatchSec              	time.Duration      `json:"after_dispatch_sec"`	//发牌之后等待时长
}

func NewRoomConfig() *RoomConfig {
	return &RoomConfig{}
}

func (config *RoomConfig) Init() {
	config.NeedPlayerNum = 4
	config.MaxPlayGameCnt = 4
	config.RandomDropNum = 1
	config.SpringMultiple = 4

	config.WaitPlayerEnterRoomTimeout = 300
	config.WaitPlayerOperateTimeout = 300
	config.WaitDropSec = 3
	config.WaitReadySec = 5
	config.AfterDispatchSec = 3
}