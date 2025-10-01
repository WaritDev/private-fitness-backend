package rest

type Handler struct {
    User *UserHandler
}

func ProvideHandler(user *UserHandler) *Handler {
    return &Handler{
        User: user,
    }
}