package entity

type User struct {
	id             uint
	userName       string
	hashedPassword string
}

func (u *User) SetID(id uint) { u.id = id }
func (u *User) ID() uint      { return u.id }

func (u *User) SetUsername(s string) { u.userName = s }
func (u *User) Username() string     { return u.userName }

func (u *User) SetPassword(l string) { u.hashedPassword = l }
func (u *User) Password() string     { return u.hashedPassword }
