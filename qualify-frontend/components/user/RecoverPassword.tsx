"use client";

import { useState } from "react";
import Link from "next/link";
import { FormInput, FormButton, FormPanel, Alert } from "@/components/ui";

export function RecoverPassword() {
  const [email, setEmail] = useState("");
  const [emailError, setEmailError] = useState("");
  const [submitted, setSubmitted] = useState(false);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setEmailError("");

    if (!email.trim()) {
      setEmailError("E-mail é obrigatório");
      return;
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setEmailError("E-mail inválido");
      return;
    }

    // TODO: integrar com endpoint de recuperação de senha quando disponível
    setSubmitted(true);
  }

  return (
    <FormPanel
      title="Recuperar senha"
      description="Informe seu e-mail para receber instruções de recuperação."
      maxWidth="max-w-md"
    >
      {submitted ? (
        <div className="text-center">
          <Alert variant="success">
            <p className="font-medium mb-1">E-mail enviado!</p>
            <p className="text-sm opacity-80">
              Se o e-mail <strong className="text-white">{email}</strong>{" "}
              estiver cadastrado, você receberá as instruções de recuperação em
              breve.
            </p>
          </Alert>
          <Link
            href="/user/register"
            className="inline-block mt-4 text-sm text-accent hover:underline"
          >
            Voltar ao cadastro
          </Link>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="flex flex-col gap-5">
          <FormInput
            label="E-mail"
            name="email"
            type="email"
            value={email}
            onChange={(e) => {
              setEmail(e.target.value);
              setEmailError("");
            }}
            placeholder="joao@exemplo.com"
            error={emailError}
            required
          />

          <FormButton type="submit" className="mt-2">
            Enviar instruções
          </FormButton>
        </form>
      )}

      <p className="text-center text-sm text-neutral-slate mt-6">
        Lembrou a senha?{" "}
        <Link href="#" className="text-accent hover:underline">
          Entrar
        </Link>
      </p>
    </FormPanel>
  );
}
