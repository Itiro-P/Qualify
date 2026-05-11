"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import type {
  IUserRegisterForm,
  IUserRegisterFormErrors,
} from "@/types/user/userRegisterForm";
import { validateRegisterForm } from "@/libs/validation";
import { userService } from "@/libs/services";
import { setSessionUser } from "@/libs/session";
import type { ApiError } from "@/libs/api";
import { FormInput, FormButton, FormPanel, Alert } from "@/components/ui";

const initialForm: IUserRegisterForm = {
  name: "",
  surname: "",
  email: "",
  phone: "",
  timezone: "",
  country_name: "",
  country_code: "",
  country_state: "",
  city: "",
  password: "",
  confirmPassword: "",
};

export function RegisterUser() {
  const router = useRouter();
  const [form, setForm] = useState<IUserRegisterForm>(initialForm);
  const [errors, setErrors] = useState<IUserRegisterFormErrors>({});
  const [submitError, setSubmitError] = useState("");
  const [loading, setLoading] = useState(false);

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
    setErrors((prev) => ({ ...prev, [name]: undefined }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitError("");

    const validationErrors = validateRegisterForm(form);
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length > 0) return;

    setLoading(true);
    try {
      const fullName = `${form.name.trim()} ${form.surname.trim()}`;
      const response = await userService.register({
        name: fullName,
        email: form.email.trim(),
        password: form.password,
        phone: form.phone.trim(),
        city: form.city.trim(),
        country_code: form.country_code.trim(),
        country_name: form.country_name.trim(),
        country_state: form.country_state.trim(),
        timezone: form.timezone.trim(),
      });

      setSessionUser({
        id: response.user.id!,
        name: response.user.name!,
        email: response.user.email!,
        phone: response.user.phone,
        city: response.user.city,
        country_code: response.user.country_code,
        country_name: response.user.country_name,
        country_state: response.user.country_state,
        timezone: response.user.timezone,
      });

      router.push("/");
    } catch (err) {
      const apiErr = err as ApiError;
      setSubmitError(apiErr.message || "Erro ao cadastrar. Tente novamente.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <FormPanel
      title="Criar conta"
      description="Preencha os dados abaixo para se cadastrar na plataforma."
    >
      {submitError && <Alert variant="error">{submitError}</Alert>}

      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <div className="grid grid-cols-2 gap-4">
          <FormInput
            label="Nome"
            name="name"
            value={form.name}
            onChange={handleChange}
            placeholder="João"
            error={errors.name}
            required
          />
          <FormInput
            label="Sobrenome"
            name="surname"
            value={form.surname}
            onChange={handleChange}
            placeholder="Silva"
            error={errors.surname}
            required
          />
        </div>

        <FormInput
          label="E-mail"
          name="email"
          type="email"
          value={form.email}
          onChange={handleChange}
          placeholder="joao@exemplo.com"
          error={errors.email}
          required
        />

        <FormInput
          label="Telefone"
          name="phone"
          type="tel"
          value={form.phone}
          onChange={handleChange}
          placeholder="+55 11 99999-9999"
          error={errors.phone}
        />

        <FormInput
          label="Fuso horário"
          name="timezone"
          value={form.timezone}
          onChange={handleChange}
          placeholder="America/Sao_Paulo"
          error={errors.timezone}
          required
        />

        <div className="grid grid-cols-2 gap-4">
          <FormInput
            label="País"
            name="country_name"
            value={form.country_name}
            onChange={handleChange}
            placeholder="Brasil"
            error={errors.country_name}
            required
          />
          <FormInput
            label="Código do país"
            name="country_code"
            value={form.country_code}
            onChange={handleChange}
            placeholder="BR"
            error={errors.country_code}
            required
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <FormInput
            label="Estado"
            name="country_state"
            value={form.country_state}
            onChange={handleChange}
            placeholder="São Paulo"
            error={errors.country_state}
            required
          />
          <FormInput
            label="Cidade"
            name="city"
            value={form.city}
            onChange={handleChange}
            placeholder="São Paulo"
            error={errors.city}
            required
          />
        </div>

        <FormInput
          label="Senha"
          name="password"
          type="password"
          value={form.password}
          onChange={handleChange}
          placeholder="••••••••"
          error={errors.password}
          hint="Mínimo 8 caracteres, 1 letra, 1 número e 1 maiúscula."
          required
        />

        <FormInput
          label="Confirmar senha"
          name="confirmPassword"
          type="password"
          value={form.confirmPassword}
          onChange={handleChange}
          placeholder="••••••••"
          error={errors.confirmPassword}
          required
        />

        <FormButton
          type="submit"
          loading={loading}
          loadingText="Cadastrando..."
          className="mt-2"
        >
          Cadastrar
        </FormButton>
      </form>

      <p className="text-center text-sm text-neutral-slate mt-6">
        Já tem uma conta?{" "}
        <Link href="#" className="text-accent hover:underline">
          Entrar
        </Link>
      </p>
    </FormPanel>
  );
}
