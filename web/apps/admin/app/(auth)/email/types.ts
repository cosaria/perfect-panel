import type {
  UserLoginRequest,
  UserRegisterRequest,
  ResetPasswordRequest,
} from "@/services/user-api/types.gen";

export type AuthView = "register" | "reset" | "login";

export type EmailAuthValues = Partial<
  UserLoginRequest & UserRegisterRequest & ResetPasswordRequest
>;
