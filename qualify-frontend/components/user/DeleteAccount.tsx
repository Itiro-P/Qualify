"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { userService } from "@/libs/services";
import { getSessionUser, clearSession, type SessionUser } from "@/libs/session";
import type { ApiError } from "@/libs/api";
import { FormInput, FormButton, FormPanel, Alert } from "@/components/ui";

export function DeleteAccount() {
  const router = useRouter();
  const [user, setUser] = useState<SessionUser | null>(null);
  const [confirmation, setConfirmation] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const session = getSessionUser();
    if (!session) {
      router.push("/User/register");
      return;
    }
    setUser(session);
  }, [router]);

  async function handleDelete(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (confirmation !== "EXCLUIR") {
      setError("Digite EXCLUIR para confirmar");
      return;
    }

    if (!user) return;

    setLoading(true);
    try {
      await userService.delete(user.id);
      clearSession();
      router.push("/");
    } catch (err) {
      const apiErr = err as ApiError;
      setError(apiErr.message || "Erro ao excluir conta. Tente novamente.");
    } finally {
      setLoading(false);
    }
  }

  if (!user) return null;

  return (
    <FormPanel title="Excluir conta" maxWidth="max-w-md">
      <p className="text-sm text-neutral-slate mb-4">
        Esta ação é <strong className="text-red-400">irreversível</strong>.
        Todos os seus dados serão permanentemente removidos.
      </p>

      <Alert variant="error">
        Você está prestes a excluir a conta de{" "}
        <strong className="text-white">{user.email}</strong>.
      </Alert>

      {error && <Alert variant="error">{error}</Alert>}

      <form onSubmit={handleDelete} className="flex flex-col gap-5">
        <FormInput
          label={`Digite EXCLUIR para confirmar`}
          name="confirmation"
          value={confirmation}
          onChange={(e) => {
            setConfirmation(e.target.value);
            setError("");
          }}
          placeholder="EXCLUIR"
        />

        <div className="flex gap-4">
          <FormButton
            type="button"
            variant="outline"
            onClick={() => router.back()}
            fullWidth
          >
            Cancelar
          </FormButton>
          <FormButton
            type="submit"
            variant="danger"
            loading={loading}
            loadingText="Excluindo..."
            disabled={confirmation !== "EXCLUIR"}
            fullWidth
          >
            Excluir conta
          </FormButton>
        </div>
      </form>
    </FormPanel>
  );
}
