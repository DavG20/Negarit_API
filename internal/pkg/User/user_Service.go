package user

type UserService struct {
	UserRepoI UserRepoInterface
}

func NewUserService(repo UserRepoInterface) *UserService {
	return &UserService{UserRepoI: repo}
}

func (userService *UserService) CheckUserEmailExist(email string) bool {
	_, err := userService.UserRepoI.CheckUserEmailExist(email)
	if err != nil {
		return false
	}
	return true
}

func (userService *UserService) CheckUserNameExist(username string) bool{
	_,err:=userService.UserRepoI.CheckUserNameExist(username)
	if err!=nil{
		return false
	}
	return true

}

func (userService *UserService) UserRegister(userInput *SignUpInput) *DBResponse {
	user, err := userService.UserRepoI.UserRegister(userInput)
	if err != nil {
		return nil
	}
	return user
}

func (userService UserService) UserLogin(userInput *SignInInput) *DBResponse {
	return userService.UserRepoI.UserLogin(userInput)
}
