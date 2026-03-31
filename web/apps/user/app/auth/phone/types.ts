export type AuthView = "register" | "reset" | "login";

export type PhoneAuthValues = Partial<
  API.TelephoneLoginRequest & API.TelephoneRegisterRequest & API.TelephoneResetPasswordRequest
>;
