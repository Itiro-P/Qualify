"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import type {
  IUserEditForm,
  IUserEditFormErrors,
} from "@/types/user/userEditForm";
import { validateEditForm } from "@/libs/validation";
import { userService } from "@/libs/services";
import {
  getSessionUser,
} from "@/libs/session";
import { FormInput, FormButton, FormPanel, Alert } from "@/components/ui";
import type { User } from "@/types/services/user";

export function EditUser() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [form, setForm] = useState<IUserEditForm>({
    name: "",
    surname: "",
    email: "",
    phone: "",
    timezone: "",
    country_name: "",
    country_code: "",
    country_state: "",
    city: "",
  });
  const [errors, setErrors] = useState<IUserEditFormErrors>({});
  const [submitError, setSubmitError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    async function loadUser() {
      const session = await getSessionUser();
      if (!session) {
        router.push("/user/register");
        return;
      }
      setUser(session);

      const nameParts = (session.name || "").split(" ");
      const firstName = nameParts[0] || "";
      const surname = nameParts.slice(1).join(" ") || "";

      setForm({
        name: firstName,
        surname,
        email: session.email || "",
        phone: session.phone || "",
        timezone: session.timezone || "",
        country_name: session.country_name || "",
        country_code: session.country_code || "",
        country_state: session.country_state || "",
        city: session.city || "",
      });
    }
    loadUser();
  }, [router]);

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
    setErrors((prev) => ({ ...prev, [name]: undefined }));
    setSuccess("");
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitError("");
    setSuccess("");

    const validationErrors = validateEditForm(form);
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length > 0) return;
    if (!user) return;

    setLoading(true);
    const fullName = `${form.name.trim()} ${form.surname.trim()}`;
    const updated = await userService.update(user.id, {
      name: fullName,
      email: form.email.trim(),
      phone: form.phone.trim(),
      city: form.city.trim(),
      country_code: form.country_code.trim(),
      country_name: form.country_name.trim(),
      country_state: form.country_state.trim(),
      timezone: form.timezone.trim(),
    });

    if (updated) {
      setSuccess("Dados atualizados com sucesso!");
    } else {
      setSubmitError("Erro ao atualizar. Tente novamente.");
    }
    setLoading(false);
  }

  if (!user) return null;

  return (
    <FormPanel
      title="Editar cadastro"
      description="Atualize suas informações pessoais."
    >
      {submitError && <Alert variant="error">{submitError}</Alert>}
      {success && <Alert variant="success">{success}</Alert>}

      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <div className="grid grid-cols-2 gap-4">
          <FormInput
            label="Nome"
            name="name"
            value={form.name}
            onChange={handleChange}
            error={errors.name}
            required
          />
          <FormInput
            label="Sobrenome"
            name="surname"
            value={form.surname}
            onChange={handleChange}
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
          error={errors.email}
          required
        />

        <FormInput
          label="Telefone"
          name="phone"
          type="tel"
          value={form.phone}
          onChange={handleChange}
          error={errors.phone}
        />

        <FormInput
          label="Fuso horário"
          name="timezone"
          value={form.timezone}
          onChange={handleChange}
          error={errors.timezone}
          required
        />

        <div className="grid grid-cols-2 gap-4">
          <FormInput
            label="País"
            name="country_name"
            value={form.country_name}
            onChange={handleChange}
            error={errors.country_name}
            required
          />
          <FormInput
            label="Código do país"
            name="country_code"
            value={form.country_code}
            onChange={handleChange}
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
            error={errors.country_state}
            required
          />
          <FormInput
            label="Cidade"
            name="city"
            value={form.city}
            onChange={handleChange}
            error={errors.city}
            required
          />
        </div>

        <FormButton
          type="submit"
          loading={loading}
          loadingText="Salvando..."
          className="mt-2"
        >
          Salvar alterações
        </FormButton>
      </form>
    </FormPanel>
  );
}
