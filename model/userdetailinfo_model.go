package model

type UserDetailInfo struct {
	User_Detail_Id int    `json:"id"`
	AuthId         int    `json:"auth_id"`
	Calling_name   string `json:"calling_name"`
	Sport          string `json:"sport"`
	Skill          string `json:"skill"`
	Avatar         string `json:"avatar"`
	Bio            string `json:"bio"`
	Created_at     string `json:"created_at"`
}

type Location struct {
	Id             int     `json:"id"`
	User_Detail_Id int     `json:"user_detail_id"`
	Country        string  `json:"country"`
	City           string  `json:"city"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
}

type Match struct {
	Id         int      `json:"id"`
	CreatorId  int      `json:"creator_id"`
	Sport      string   `json:"sport"`
	Title      string   `json:"title"`
	Location   Location `json:"location"`
	Time       string   `json:"time"`
	MaxPlayers int      `json:"max_players"`
	Status     string   `json:"status"`
	CreatedAt  string   `json:"created_at"`
}

type UserProfileResponse struct {
	AuthId      int     `json:"auth_id"`
	CallingName string  `json:"calling_name"`
	Sport       string  `json:"sport"`
	Skill       string  `json:"skill"`
	Avatar      string  `json:"avatar"`
	Bio         string  `json:"bio"`
	Country     string  `json:"country"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
