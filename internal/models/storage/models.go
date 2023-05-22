//go:generate reform
package storage

import (
	"time"
)

//reform:links
type Link struct {
	ID          string    `reform:"id,pk"`
	UserID      string    `reform:"user_id"`
	SSTEmail    string    `reform:"sst_email"`
	SSTPassword string    `reform:"sst_password"`
	CreatedAt   time.Time `reform:"created_at"`
	UpdatedAt   time.Time `reform:"updated_at"`
}

func (s *Link) BeforeUpdate() error {
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Link) Equal(o *Link) bool {
	return s.ID == o.ID &&
		s.SSTEmail == o.SSTEmail &&
		s.SSTPassword == o.SSTPassword
}

type LogLevel string

const (
	Error LogLevel = "Error"
	Info  LogLevel = "Info"
)

//reform:logs
type Log struct {
	ID      string    `reform:"id,pk"`
	LinkID  string    `reform:"link_id"`
	Time    time.Time `reform:"time"`
	Level   LogLevel  `reform:"level"`
	Message string    `reform:"message"`
}
