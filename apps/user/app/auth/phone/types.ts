import type {
  TelephoneLoginRequest,
  TelephoneRegisterRequest,
  TelephoneResetPasswordRequest,
} from "@/services/user-api/types.gen";

export type AuthView = "register" | "reset" | "login";

export type PhoneAuthValues = Partial<
  TelephoneLoginRequest & TelephoneRegisterRequest & TelephoneResetPasswordRequest
>;
