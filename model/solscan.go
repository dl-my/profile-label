package model

type ApiResponse struct {
	Data []NoteItem `json:"data"`
}

type NoteItem struct {
	ID        string `json:"_id"`
	Hash      string `json:"hash"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	Label     string `json:"label"`
	Note      string `json:"note"`
	Status    int    `json:"status"`
	Type      string `json:"type"`
	UpdatedAt string `json:"updatedAt"`
}

type Solscan struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Address   string `json:"address"`
	Tag       string `json:"tag"`
	Note      string `json:"note"`
	Type      string `json:"type"`
	CreatedAt int64  `json:"created_at"`
}

func (Solscan) TableName() string {
	return "solscan"
}
