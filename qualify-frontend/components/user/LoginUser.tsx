"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { FormInput, FormButton, FormPanel, Alert } from "@/components/ui";
import { userService } from "@/libs/services";
import { setStoredTokens } from "@/libs/session";

export function LoginUser() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<{
    email?: string;
    password?: string;
  }>({});

  function validate(): boolean {
    const errors: { email?: string; password?: string } = {};

    if (!email.trim()) {
      errors.email = "E-mail é obrigatório";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      errors.email = "E-mail inválido";
    }

    if (!password) {
      errors.password = "Senha é obrigatória";
    }

    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (!validate()) return;

    setLoading(true);
    const response = await userService.login({ email, password });

    if (response) {
      setStoredTokens({
        access_token: response.access_token,
        refresh_token: response.refresh_token,
      });
      router.push("/");
    } else {
      setError("E-mail ou senha incorretos.");
    }
    setLoading(false);
  }

  return (
    <FormPanel title="Entrar" description="Acesse sua conta na Qualify.">
      {error && <Alert variant="error">{error}</Alert>}

      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <FormInput
          label="E-mail"
          name="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={fieldErrors.email}
          placeholder="seu@email.com"
          required
        />

        <FormInput
          label="Senha"
          name="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          error={fieldErrors.password}
          placeholder="••••••••"
          required
        />

        <div className="flex justify-end">
          <Link
            href="/user/recover-password"
            className="text-xs text-accent hover:underline"
          >
            Esqueceu a senha?
          </Link>
        </div>

        <FormButton type="submit" loading={loading} loadingText="Entrando...">
          Entrar
        </FormButton>
      </form>

      <p className="text-center text-sm text-neutral-slate mt-6">
        Não tem uma conta?{" "}
        <Link
          href="/user/register"
          className="text-accent hover:underline font-medium"
        >
          Cadastre-se
        </Link>
      </p>
    </FormPanel>
  );
}
