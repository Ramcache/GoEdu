package models

import "GoEdu/proto"

type Instructor struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (i *Instructor) ToProto() *proto.Instructor {
	return &proto.Instructor{
		Id:    i.ID,
		Name:  i.Name,
		Email: i.Email,
	}
}
