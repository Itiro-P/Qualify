import type { IUserRegisterForm, IUserRegisterFormErrors } from "@/types/user/userRegisterForm";
import type { IUserEditForm, IUserEditFormErrors } from "@/types/user/userEditForm";

function validateEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}

function validatePassword(password: string): string[] {
  const errors: string[] = [];
  if (password.length < 8) errors.push("Mínimo de 8 caracteres");
  if (!/[a-zA-Z]/.test(password)) errors.push("Pelo menos 1 letra");
  if (!/[0-9]/.test(password)) errors.push("Pelo menos 1 número");
  if (!/[A-Z]/.test(password)) errors.push("Pelo menos 1 letra maiúscula");
  return errors;
}

export function validateRegisterForm(form: IUserRegisterForm): IUserRegisterFormErrors {
  const errors: IUserRegisterFormErrors = {};

  if (!form.name.trim()) errors.name = "Nome é obrigatório";
  if (!form.surname.trim()) errors.surname = "Sobrenome é obrigatório";

  if (!form.email.trim()) {
    errors.email = "E-mail é obrigatório";
  } else if (!validateEmail(form.email)) {
    errors.email = "E-mail inválido";
  }

  if (!form.timezone.trim()) errors.timezone = "Fuso horário é obrigatório";
  if (!form.country_name.trim()) errors.country_name = "País é obrigatório";
  if (!form.country_code.trim()) errors.country_code = "Código do país é obrigatório";
  if (!form.country_state.trim()) errors.country_state = "Estado é obrigatório";
  if (!form.city.trim()) errors.city = "Cidade é obrigatória";

  if (!form.password) {
    errors.password = "Senha é obrigatória";
  } else {
    const pwdErrors = validatePassword(form.password);
    if (pwdErrors.length > 0) {
      errors.password = pwdErrors.join("; ");
    }
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = "Confirmação de senha é obrigatória";
  } else if (form.password !== form.confirmPassword) {
    errors.confirmPassword = "As senhas não coincidem";
  }

  return errors;
}

export function validateEditForm(form: IUserEditForm): IUserEditFormErrors {
  const errors: IUserEditFormErrors = {};

  if (!form.name.trim()) errors.name = "Nome é obrigatório";
  if (!form.surname.trim()) errors.surname = "Sobrenome é obrigatório";

  if (!form.email.trim()) {
    errors.email = "E-mail é obrigatório";
  } else if (!validateEmail(form.email)) {
    errors.email = "E-mail inválido";
  }

  if (!form.timezone.trim()) errors.timezone = "Fuso horário é obrigatório";
  if (!form.country_name.trim()) errors.country_name = "País é obrigatório";
  if (!form.country_code.trim()) errors.country_code = "Código do país é obrigatório";
  if (!form.country_state.trim()) errors.country_state = "Estado é obrigatório";
  if (!form.city.trim()) errors.city = "Cidade é obrigatória";

  return errors;
}
