package rest

type Handler struct {
    User *UserHandler
    Auth *AuthHandler
}

func ProvideHandler(user *UserHandler, auth *AuthHandler) *Handler {
    return &Handler{
        User: user,
        Auth: auth,
    }
}