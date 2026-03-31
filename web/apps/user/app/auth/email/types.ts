export type AuthView = "register" | "reset" | "login";

export type EmailAuthValues = Partial<
  API.UserLoginRequest & API.UserRegisterRequest & API.ResetPasswordRequest
>;
