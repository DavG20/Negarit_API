package user

type UserService struct {
	UserRepoI UserRepoInterface
}

func NewUserService(repo UserRepoInterface) *UserService {
	return &UserService{UserRepoI: repo}
}

func (userService *UserService) CheckUserEmail(email string) bool {
	_, err := userService.UserRepoI.CheckUserEmail(email)
	if err != nil {
		return false
	}
	return true
}

func (userService *UserService) UserRegister(userInput SignUpInput) *User {
	user, err := userService.UserRepoI.UserRegister(userInput)
	if err != nil {
		return nil
	}
	return user
}

func (userService UserService) UserLogin(email, password string) *User {
	return userService.UserRepoI.UserLogin(email, password)
}
