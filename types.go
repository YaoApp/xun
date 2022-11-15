package xun

import "net/textproto"

// R alias map[string]interface{}, R is the first letter of "Row"
type R map[string]interface{}

// N an numberic value,  R is the first letter of "Numberic"
type N struct {
	Number interface{}
}

// T an datetime value, T is the first letter of "Time"
type T struct {
	Time interface{}
}

// P an Paginator struct, P is the first letter of "Paginator"
type P struct {
	Items        []interface{}          `json:"items"`
	Total        int                    `json:"total"`
	TotalPages   int                    `json:"total_pages"`
	PageSize     int                    `json:"page_size"`
	CurrentPage  int                    `json:"current_page"`
	NextPage     int                    `json:"next_page"`
	PreviousPage int                    `json:"previous_page"`
	LastPage     int                    `json:"last_page"`
	Options      map[string]interface{} `json:"options,omtempty"`
}

// UploadFile deprecated -> gou.UploadFile upload file
type UploadFile struct {
	Name     string
	TempFile string
	Size     int64
	Header   textproto.MIMEHeader
}
