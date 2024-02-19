package driven

type IUserDriven interface{}

type userDriven struct {
}

func NewUserDriven() IUserDriven {
	return &userDriven{}
}

func (u *userDriven) SetTimeoutForUserRegister() {

}
